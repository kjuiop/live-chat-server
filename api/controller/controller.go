package controller

func IsInternalServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}
