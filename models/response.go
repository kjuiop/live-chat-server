package models

type SuccessRes struct {
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

type FailRes struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
