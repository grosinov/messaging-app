package helpers

const DefaultMessagesLimit = "100"

var (
	MessageTypes = map[string]bool{
		"text":  true,
		"image": true,
		"video": true}
)
