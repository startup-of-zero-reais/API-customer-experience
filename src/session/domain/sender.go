package domain

type Sender interface {
	SendSMS(message, phone string) error
}
