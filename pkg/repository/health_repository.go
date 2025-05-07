package repository

func (r RepositoryImpl) HealthCheck() error {
	db, err := r.DB.DB()
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}
