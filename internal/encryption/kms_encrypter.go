package encryption

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/service/kms"

	kmsapi "github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

type KMSEncrypter struct {
	kmsClient kmsapi.KMSAPI
	keyID     string
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func NewKMSEncrypter(kmsClient kmsapi.KMSAPI, keyID string) *KMSEncrypter {
	return &KMSEncrypter{
		kmsClient: kmsClient,
		keyID:     keyID,
	}
}

func (e *KMSEncrypter) Encrypt(_ context.Context, value string) (string, error) {
	res, err := e.kmsClient.Encrypt(&kms.EncryptInput{
		KeyId:     &e.keyID,
		Plaintext: []byte(value),
	})
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(res.CiphertextBlob), nil
}

func (e *KMSEncrypter) Decrypt(_ context.Context, value string) (string, error) {
	cipherTextBlob, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	res, err := e.kmsClient.Decrypt(&kms.DecryptInput{
		CiphertextBlob: cipherTextBlob,
	})
	if err != nil {
		return "", err
	}

	return string(res.Plaintext), nil
}
