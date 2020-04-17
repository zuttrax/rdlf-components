package secrets

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
	DecryptError  = errors.New("aws decryption failure")
	InternalError = errors.New("aws internal error")
	NotFoundError = errors.New("secret not found")
)

type AWSManager interface {
	GetSecretValueByInput(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

type AWSSecret struct {
	manager *secretsmanager.SecretsManager
}

func (a AWSSecret) GetSecretValueByInput(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return a.manager.GetSecretValue(input)
}

type SecretsManager interface {
	GetSecretByName(string) (string, error)
}

type Secret struct {
	AWSManager AWSManager
}

func NewSecrets() Secret {
	return Secret{
		AWSManager: AWSSecret{
			manager: initializeSecretsManager(),
		},
	}
}

func initializeSecretsManager() *secretsmanager.SecretsManager {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.ApSoutheast1RegionID),
	}))
	secretsManager := secretsmanager.New(sess)

	return secretsManager
}

func (s Secret) GetSecretByName(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSPREVIOUS"),
	}

	resp, err := s.AWSManager.GetSecretValueByInput(input)
	if err != nil {
		return "", handError(err)
	}

	value := *resp.SecretString

	return value, nil
}

func handError(err error) error {
	return fmt.Errorf("%w: %v", InternalError, err)
}
