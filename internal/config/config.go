package config

import (
	"bytes"
	"errors"

	"encoding/pem"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/spf13/viper"
)

var (
	ErrUnrecognizedBlockType = errors.New("unrecognized block type in UnpackCertificateBundle")
	ErrCertificateNotPresent = errors.New("certificate not present in bundle")
	ErrKeyNotPresent         = errors.New("key not present in bundle")
	ErrFailedToSetEnv        = errors.New("failed to set BBAPI ENV")
)

type DatabaseConfig struct {
	DatabaseURL string
}

type SecretsManagerConfig struct {
	Region string
}

func InitializeSecretsManagerConfig(logger *zap.Logger) (config SecretsManagerConfig, err error) {
	if err := viper.UnmarshalKey("secret-manager", &config); err != nil {
		return config, fmt.Errorf("can not find secret manager config %w", err)
	}
	logger.Sugar().Infow("initialized secret manager config",
		zap.String("region", config.Region))
	return config, nil
}

const (
	certificateType     = "CERTIFICATE"
	rsaPrivateKeyType   = "RSA PRIVATE KEY"
	ecdsaPrivateKeyType = "EC PRIVATE KEY"
	pkcs8PrivateKeyType = "PRIVATE KEY"
)

func UnpackCertificateBundle(bundle string) (cert string, key string, err error) {
	var certificateBytes bytes.Buffer
	var privKeyBytes bytes.Buffer

	for block, rest := pem.Decode([]byte(bundle)); block != nil; block, rest = pem.Decode(rest) {
		switch block.Type {
		case certificateType:
			if err := pem.Encode(&certificateBytes, block); err != nil {
				return "", "", fmt.Errorf("error encoding certificate in UnpackCertificateBundle: %w", err)
			}
		case rsaPrivateKeyType, ecdsaPrivateKeyType, pkcs8PrivateKeyType:
			if err := pem.Encode(&privKeyBytes, block); err != nil {
				return "", "", fmt.Errorf("error encoding private key in UnpackCertificateBundle: %w", err)
			}
		default:
			return "", "", fmt.Errorf("%w: %s", ErrUnrecognizedBlockType, block.Type)
		}
	}

	if certificateBytes.Len() == 0 {
		return "", "", ErrCertificateNotPresent
	}
	if privKeyBytes.Len() == 0 {
		return "", "", ErrKeyNotPresent
	}

	return certificateBytes.String(), privKeyBytes.String(), nil
}

func ProcessCertificateBundle() error {
	certBundle := viper.GetString("CERT_BUNDLE")
	if certBundle == "" {
		return nil
	}

	serverCert, serverKey, err := UnpackCertificateBundle(certBundle)
	if err != nil {
		return fmt.Errorf("failed to unpack certificate bundle: %w", err)
	}

	if err := os.Setenv("SERVER_CERT", serverCert); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSetEnv, err)
	}
	if err := os.Setenv("SERVER_KEY", serverKey); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSetEnv, err)
	}

	return nil
}
