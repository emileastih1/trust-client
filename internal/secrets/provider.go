package secrets

import (
	"bulletin-board-api/internal/constants"
	"crypto/rand"
	"errors"
	"fmt"

	serviceConfig "bulletin-board-api/internal/config"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var errSecretStringEmpty = errors.New("secret string is empty")

type AddressSecretConfig struct {
	Name          string
	ARN           string
	VersionID     string
	EncodedSecret string
}

func GenerateSecureLen256Salt() (string, error) {
	bytes := make([]byte, 256)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("geneate secure salt failed: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

func InitializeSecrets(logger *zap.Logger, svc *secretsmanager.SecretsManager, _ serviceConfig.SecretsManagerConfig) (*AddressSecretConfig, error) {
	if viper.GetString("stage") == constants.StageLocal {
		return &AddressSecretConfig{
			Name:          "localname",
			ARN:           "localarn",
			VersionID:     "localversion",
			EncodedSecret: viper.GetString(constants.AddressSecret),
		}, nil
	}

	getInput := &secretsmanager.GetSecretValueInput{
		SecretId:     &constants.AddressSecret,
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	getResult, err := svc.GetSecretValue(getInput)
	if err != nil {
		return nil, fmt.Errorf("initialize secrets object failed: %w", err)
	}
	if getResult.SecretString == nil {
		return nil, fmt.Errorf("get secret failed: %w", errSecretStringEmpty)
	}
	logger.Sugar().Info("initialized address secret",
		zap.Int("secret_lenth", len(*getResult.SecretString)))

	return &AddressSecretConfig{
		Name:          *getResult.Name,
		ARN:           *getResult.ARN,
		VersionID:     *getResult.VersionId,
		EncodedSecret: *getResult.SecretString,
	}, nil
}

func InitializeSecretManagerClient(config serviceConfig.SecretsManagerConfig) (*secretsmanager.SecretsManager, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("initialize aws session failed: %w", err)
	}
	return secretsmanager.New(awsSession,
		aws.NewConfig().WithRegion(config.Region)), nil
}

var Provider = wire.NewSet(
	InitializeSecrets, InitializeSecretManagerClient,
)
