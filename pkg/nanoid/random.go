package nanoid

import gonanoid "github.com/jaevor/go-nanoid"

const (
	alphabet = "1234567890qwertyupasdfghkzxcvbnm"
)

var (
	randomShortGen, _     = gonanoid.CustomASCII(alphabet, 10)
	randomGen, _          = gonanoid.CustomASCII(alphabet, 16)
	randomLongGen, _      = gonanoid.CustomASCII(alphabet, 24)
	randomExtraLongGen, _ = gonanoid.CustomASCII(alphabet, 32)
)

func RandomID() string {
	return randomGen()
}

func RandomShortID() string {
	return randomShortGen()
}

func RandomLongID() string {
	return randomLongGen()
}

func RandomCryptoID() string {
	return randomExtraLongGen()
}
