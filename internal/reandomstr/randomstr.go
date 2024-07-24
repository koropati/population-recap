package randomstr

import (
	"math/rand"
	"time"
)

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers          = "0123456789"
	specialChars     = "!@#$%^&*()-_=+[]{}|;:,.<>/?"
)

type RandomStr interface {
	GenerateRandomString(length int) string
}

type randomStr struct {
	includeLowercase bool
	includeUppercase bool
	includeNumeric   bool
	includeSpecial   bool
}

func New(includeLowercase, includeUppercase, includeNumeric, includeSpecial bool) RandomStr {
	return &randomStr{
		includeLowercase: includeLowercase,
		includeUppercase: includeUppercase,
		includeNumeric:   includeNumeric,
		includeSpecial:   includeSpecial,
	}
}

func (r *randomStr) GenerateRandomString(length int) string {
	var characterSet string

	if r.includeLowercase {
		characterSet += lowercaseLetters
	}
	if r.includeUppercase {
		characterSet += uppercaseLetters
	}
	if r.includeNumeric {
		characterSet += numbers
	}
	if r.includeSpecial {
		characterSet += specialChars
	}

	if characterSet == "" {
		return ""
	}

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	output := make([]byte, length)
	for i := range output {
		output[i] = characterSet[rnd.Intn(len(characterSet))]
	}

	return string(output)
}
