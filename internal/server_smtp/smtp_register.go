package serversmtp

import (
	"fmt"

	mail "github.com/xhit/go-simple-mail/v2"
)

func SendConfremRegister(toEmail, code string) error {
	smtpClient, err := getSMTPClient()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	fromEmail := getEnvOrDefault("SMTP_EMAIL", "misuraaleksej60@gmail.com")
	email.SetFrom(fromEmail).
		AddTo(toEmail).
		SetSubject("Код подтверждения регистрации").
		SetBody(mail.TextHTML, fmt.Sprintf("<h2>Ваш код: <b>%s</b></h2>", code))
	if email.Error != nil {
		return email.Error
	}
	fmt.Println(err)
	return email.Send(smtpClient)
}
