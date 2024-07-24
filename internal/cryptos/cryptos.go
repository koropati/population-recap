package cryptos

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
)

type Cryptos interface {
	Encrypt(inputString string) (output string, err error)
	Decrypt(inputString string) (output string, err error)
}

type cryptos struct {
	block cipher.Block
}

func New(secretKey string) Cryptos {

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		log.Printf("[Error][Initial Cipher] E: %v", err)
	}

	return &cryptos{
		block: block,
	}
}

func (c *cryptos) Encrypt(inputString string) (output string, err error) {
	plaintextBytes := []byte(inputString)
	cipherText := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(c.block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintextBytes)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (c *cryptos) Decrypt(inputString string) (output string, err error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(inputString)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return "", errors.New("invalid ciphertext size")
	}

	iv := decodedCiphertext[:aes.BlockSize]
	plaintext := decodedCiphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(c.block, iv)
	stream.XORKeyStream(plaintext, plaintext)
	return string(plaintext), nil
}
