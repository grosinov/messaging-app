package auth

import (
	"context"
	"fmt"
	"github.com/challenge/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// ValidateUser checks for a token and validates it
// before allowing the method to execute
func ValidateUser(db *gorm.DB) func(_ http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Println("Invalid token: token not start with bearer")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return JwtSecret, nil
			})

			if err != nil || !token.Valid {
				log.Println(err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Invalid token: claims not map claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userIDraw, ok := claims["user_id"].(float64)
			if !ok {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			_, ok = claims["exp"].(float64)
			if !ok {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userID := uint64(userIDraw)

			var user models.User
			if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
				log.Println("User not found")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
