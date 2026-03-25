package encryption

import "context"

type Encrypter interface {
	Encrypt(ctx context.Context, value string) (string, error)
	Decrypt(ctx context.Context, value string) (string, error)
}
