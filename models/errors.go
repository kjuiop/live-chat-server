package models

import "fmt"

const (
	NoError            int = 0
	ErrParsing             = 4001
	ErrRedisHMSETError     = 5001
)

var codeToMessage = map[int]string{
	NoError:            "ok",
	ErrParsing:         "invalid request body",
	ErrRedisHMSETError: "internal error occur",
}

func GetCustomErrMessage(code int, error string) string {
	message, exists := codeToMessage[code]
	if !exists {
		return "Unknown error"
	}

	return fmt.Sprintf("%s, err : %s", message, error)
}

func GetCustomMessage(code int) string {
	message, exists := codeToMessage[code]
	if !exists {
		return "Unknown error"
	}

	return message
}
