package emailservice

type EmailServiceInterface interface {
	Send(data *SendData) error
}
