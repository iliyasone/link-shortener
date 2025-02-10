package generator

import (
	"crypto/rand"
	"math/big"
)

const DefaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const DefaultURLLength = 10

type Generator struct {
	Charset   string
	URLLength int
}

func NewGenerator() *Generator {
	return &Generator{
		Charset:   DefaultCharset,
		URLLength: DefaultURLLength,
	}
}

// Generator is responsible for generating random short URLs based on a configurable character set and length.
func (g *Generator) Generate() (string, error) {
	result := make([]byte, g.URLLength)
	for i := 0; i < g.URLLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(g.Charset))))
		if err != nil {
			return "", err
		}
		result[i] = g.Charset[n.Int64()]
	}
	return string(result), nil
}
