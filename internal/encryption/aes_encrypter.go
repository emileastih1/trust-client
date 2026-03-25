package encryption

import (
	"bulletin-board-api/internal/secrets"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/xdg-go/pbkdf2"
)

const (
	saltSize       = 8
	keySize        = 32
	iterationCount = 4096
	nonceSize      = 12
)

type AESEncrypter struct {
	encodedSecret []byte
}

var ErrShortMessage = errors.New("encrypted message too short")

func deriveKey(address, salt []byte) []byte {
	return pbkdf2.Key(address, salt, iterationCount, keySize, sha256.New)
}

func NewAESEncrypterWithSecret(secret []byte) *AESEncrypter {
	return &AESEncrypter{
		encodedSecret: secret,
	}
}

func NewAESEncrypter(config *secrets.AddressSecretConfig) (*AESEncrypter, error) {
	secret, err := base64.StdEncoding.DecodeString(config.EncodedSecret)

	return &AESEncrypter{
		encodedSecret: secret,
	}, err
}

func (e *AESEncrypter) Encrypt(_ context.Context, value string) (string, error) {
	// create a random salt
	salt := make([]byte, saltSize)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := deriveKey(e.encodedSecret, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(value), nil)

	// combine salt, nonce and ciphertext
	//nolint:gocritic,makezero
	result := append(salt, append(nonce, ciphertext...)...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func (e *AESEncrypter) Decrypt(_ context.Context, value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	if len(data) < saltSize+nonceSize {
		return "", ErrShortMessage
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := deriveKey(e.encodedSecret, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decodedValue, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decodedValue), err
}
