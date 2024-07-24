package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv                      string `mapstructure:"APP_ENV"`
	ServerAddress               string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout              int    `mapstructure:"CONTEXT_TIMEOUT"`
	AppName                     string `mapstructure:"APP_NAME"`
	AppFeUrl                    string `mapstructure:"APP_FE_URL"`
	DBHost                      string `mapstructure:"DB_HOST"`
	DBPort                      string `mapstructure:"DB_PORT"`
	DBUser                      string `mapstructure:"DB_USER"`
	DBPass                      string `mapstructure:"DB_PASS"`
	DBName                      string `mapstructure:"DB_NAME"`
	DefaultPageNumber           int64  `mapstructure:"DEFAULT_PAGE_NUMBER"`
	DefaultPageSize             int64  `mapstructure:"DEFAULT_PAGE_SIZE"`
	AccessTokenExpiryHour       int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour      int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	VerificationEmailExpiryHour int    `mapstructure:"VERIFICATION_EMAIL_EXPIRY_HOUR"`
	ForgotTokenExpiryHour       int    `mapstructure:"FORGOT_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret           string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret          string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AmqpHost                    string `mapstructure:"AMQP_HOST"`
	AmqpPort                    string `mapstructure:"AMQP_PORT"`
	AmqpUser                    string `mapstructure:"AMQP_USER"`
	AmqpPass                    string `mapstructure:"AMQP_PASS"`
	AmqpReconRetry              int    `mapstructure:"AMQP_RECON_RETRY"`
	AmqpReconInterval           int    `mapstructure:"AMQP_RECON_INTERVAL"`
	AmqpQueueEx                 string `mapstructure:"AMQP_QUEUE_EX"`
	AmqpDebug                   bool   `mapstructure:"AMQP_DEBUG"`
	AmqpConsumerLimit           int    `mapstructure:"AMQP_CONSUMER_LIMIT"`
	AmqpWorkerLimit             int    `mapstructure:"AMQP_WORKER_LIMIT"`
	SmtpHost                    string `mapstructure:"SMTP_HOST"`
	SmtpPort                    int    `mapstructure:"SMTP_PORT"`
	SmtpUser                    string `mapstructure:"SMTP_USER"`
	SmtpPass                    string `mapstructure:"SMTP_PASS"`
	SmtpSenderMail              string `mapstructure:"SMTP_SENDER_EMAIL"`
	SmtpEncryption              string `mapstructure:"SMTP_ENCRYPTION"`
	CasbinModelPath             string `mapstructure:"CASBIN_MODEL_PATH"`
	CasbinPolicyPath            string `mapstructure:"CASBIN_POLICY_PATH"`
	SecretKey                   string `mapstructure:"SECRET_KEY"`
	SessionKey                  string `mapstructure:"SESSION_KEY"`
	TelegramBotToken            string `mapstructure:"TELEGRAM_BOT_TOKEN"`
}

func NewConfig() *Config {
	env := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
