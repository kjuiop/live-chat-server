package domain

import "fmt"

const (
	InternalRedisError = "internal redis error occur"
)

const (
	NoError                int = 0
	ErrParsing                 = 4001
	ErrNotFoundChatRoom        = 4002
	ErrNotConnectSocket        = 4003
	ErrEmptyParam              = 4004
	ErrRedisHMSETError         = 5001
	ErrRedisExistError         = 5002
	ErrRedisHMDELError         = 5003
	ErrInternalServerError     = 5004
)

var codeToMessage = map[int]string{
	NoError:                "ok",
	ErrParsing:             "invalid request body",
	ErrNotFoundChatRoom:    "not found chat room",
	ErrNotConnectSocket:    "not connect socket",
	ErrEmptyParam:          "invalid params",
	ErrRedisHMSETError:     InternalRedisError,
	ErrRedisExistError:     InternalRedisError,
	ErrRedisHMDELError:     InternalRedisError,
	ErrInternalServerError: "internal server error",
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
