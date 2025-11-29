package serversmtp

import (
	"fmt"

	mail "github.com/xhit/go-simple-mail/v2"
)

// функция для отправки кода подтверждения регистрации на email
func SendConfremRegister(toEmail, code string) error {
	smtpClient, err := getSMTPClient()
	if err != nil {
		return err
	}
	defer smtpClient.Close()

	email := mail.NewMSG()
	fromEmail := getEnvOrDefault("SMTP_EMAIL", "misuraaleksej60@gmail.com")
	email.SetFrom(fromEmail).
		AddTo(toEmail).
		SetSubject("Код подтверждения регистрации").
		SetBody(mail.TextHTML, generateEmailHTML(code))

	if email.Error != nil {
		return email.Error
	}

	return email.Send(smtpClient)
}

// generateEmailHTML создает простой HTML для письма
func generateEmailHTML(code string) string {
	return fmt.Sprintf(`
<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h2 style="color: #333; text-align: center;">Код подтверждения регистрации</h2>

    <p>Для завершения регистрации введите этот код в приложении:</p>

    <div style="background: #f5f5f5; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0;">
        <div style="font-size: 32px; font-weight: bold; color: #2c5aa0; letter-spacing: 5px;">
            %s
        </div>
    </div>

    <p style="color: #666; font-size: 14px;">
        Код действителен в течение 5 минут.<br>
        Если вы не запрашивали регистрацию, проигнорируйте это письмо.
    </p>
</div>
`, code)
}
