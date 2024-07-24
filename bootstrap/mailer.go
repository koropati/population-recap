package bootstrap

import (
	"context"
	"time"

	"github.com/koropati/population-recap/internal/mailer"
	"gopkg.in/gomail.v2"
)

func NewMailer(env *Config) (mailerData mailer.Mailer) {

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	host := env.SmtpHost
	port := env.SmtpPort
	user := env.SmtpUser
	pass := env.SmtpPass

	smtp := gomail.NewDialer(
		host,
		port,
		user,
		pass,
	)

	return mailer.New(smtp)
}
