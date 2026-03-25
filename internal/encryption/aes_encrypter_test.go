package encryption_test

import (
	"bulletin-board-api/internal/encryption"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	encrypter := encryption.NewAESEncrypterWithSecret([]byte("secret"))

	vaspUUID := "vaspUUID"
	encrypted, err := encrypter.Encrypt(context.Background(), vaspUUID)
	assert.NoError(t, err)

	decrypted, err := encrypter.Decrypt(context.Background(), encrypted)
	assert.NoError(t, err)

	assert.Equal(t, vaspUUID, decrypted)
}
