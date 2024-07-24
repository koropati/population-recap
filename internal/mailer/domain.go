package mailer

type ForgotPasswordData struct {
	AppName             string
	Name                string
	Email               string
	From                string
	ForgotPasswordToken string
	UrlRedirect         string
	UrlVerification     string
}

type VerificationEmailData struct {
	AppName           string
	Name              string
	Email             string
	From              string
	VerificationToken string
	UrlRedirect       string
	UrlVerification   string
}

type Notification struct {
	AppName string
	Name    string
	Email   string
	From    string
	Title   string
	Subject string
	Message string
}
