package encryption

import (
	"bulletin-board-api/internal/constants"
	"bulletin-board-api/internal/secrets"
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ErrMissingKMSKeyID = errors.New("KMS_KEY_ID is missing")

func NewEncryptor(config *secrets.AddressSecretConfig, logger *zap.Logger) (Encrypter, error) {
	stage := viper.GetString("stage")

	if stage == constants.StageLocal {
		return NewAESEncrypter(config)
	}

	keyID := viper.GetString("KMS_KEY_ID")
	if keyID == "" {
		return nil, ErrMissingKMSKeyID
	}

	mySession := session.Must(session.NewSession())
	logger.Info("Using KMS for encryption", zap.String("KeyID", keyID))

	return NewKMSEncrypter(
		kms.New(mySession),
		keyID,
	), nil
}

var Provider = wire.NewSet(
	NewAESEncrypter,
	wire.Bind(new(Encrypter), new(*AESEncrypter)),
)
