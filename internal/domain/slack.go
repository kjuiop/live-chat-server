package domain

type WebhookRes struct {
	Text string `json:"text"`
}

func NewWebhookRes(message string) *WebhookRes {
	return &WebhookRes{
		Text: message,
	}
}
