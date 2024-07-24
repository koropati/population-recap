package mailer

import (
	"html/template"
	"strings"

	"gopkg.in/gomail.v2"
)

type mailer struct {
	dialer *gomail.Dialer
}

type Mailer interface {
	SendForgotPassword(data ForgotPasswordData) (err error)
	SendVerificationEmail(data VerificationEmailData) (err error)
	SendNotification(data Notification) (err error)
}

func New(dialer *gomail.Dialer) Mailer {
	return &mailer{
		dialer: dialer,
	}
}

func (m *mailer) SendForgotPassword(data ForgotPasswordData) (err error) {
	// Membaca template HTML dari file

	htmlTemplate, err := template.ParseFiles("./internal/mailer/html/forgot_password.html")
	if err != nil {
		return err
	}

	// Membuat buffer untuk menyimpan hasil render template
	var emailBodyBuilder = new(strings.Builder)

	// Menerapkan data ke dalam template HTML
	err = htmlTemplate.Execute(emailBodyBuilder, data)
	if err != nil {
		return err
	}

	// Konfigurasi pengaturan email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", data.From)
	mailer.SetHeader("To", data.Email)
	mailer.SetHeader("Subject", data.AppName+" - Forgot Password")

	// Isi pesan email dengan HTML yang telah dibuat
	mailer.SetBody("text/html", emailBodyBuilder.String())

	// Kirim email
	err = m.dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func (m *mailer) SendVerificationEmail(data VerificationEmailData) (err error) {
	// Membaca template HTML dari file
	// Mendapatkan path dari file saat ini

	htmlTemplate, err := template.ParseFiles("./internal/mailer/html/verification_email.html")
	if err != nil {
		return err
	}

	// Membuat buffer untuk menyimpan hasil render template
	var emailBodyBuilder = new(strings.Builder)

	// Menerapkan data ke dalam template HTML
	err = htmlTemplate.Execute(emailBodyBuilder, data)
	if err != nil {
		return err
	}

	// Konfigurasi pengaturan email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", data.From)
	mailer.SetHeader("To", data.Email)
	mailer.SetHeader("Subject", data.AppName+" - Email Verification")

	// Isi pesan email dengan HTML yang telah dibuat
	mailer.SetBody("text/html", emailBodyBuilder.String())

	// Kirim email
	err = m.dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func (m *mailer) SendNotification(data Notification) (err error) {
	// Membaca template HTML dari file

	htmlTemplate, err := template.ParseFiles("./internal/mailer/html/notification.html")
	if err != nil {
		return err
	}

	// Membuat buffer untuk menyimpan hasil render template
	var emailBodyBuilder = new(strings.Builder)

	// Menerapkan data ke dalam template HTML
	err = htmlTemplate.Execute(emailBodyBuilder, data)
	if err != nil {
		return err
	}

	// Konfigurasi pengaturan email
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", data.From)
	mailer.SetHeader("To", data.Email)
	mailer.SetHeader("Subject", data.AppName+" - "+data.Subject)

	// Isi pesan email dengan HTML yang telah dibuat
	mailer.SetBody("text/html", emailBodyBuilder.String())

	// Kirim email
	err = m.dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
