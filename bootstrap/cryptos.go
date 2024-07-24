package bootstrap

import (
	"sync"

	"github.com/koropati/population-recap/internal/cryptos"
)

var (
	myCryptoOnce sync.Once
	myCrypto     cryptos.Cryptos
)

// NewRedisClient creates a new Redis client connection
func NewCryptos(config *Config) cryptos.Cryptos {
	myCryptoOnce.Do(func() {
		myCrypto = cryptos.New(config.SecretKey)
	})

	return myCrypto
}
