package secrets_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/require"
	"github.com/zuttrax/rdlf-components/pkg/secrets"
)

type AWSSecretStub struct {
	secret string
	err    error
}

func (a AWSSecretStub) GetSecretValueByInput(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	if a.err != nil {
		return nil, errors.New("")
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(a.secret),
	}, nil
}

func TestSecret_GetSecretByName(t *testing.T) {
	type fields struct {
		manager secrets.AWSManager
	}

	type args struct {
		secretName string
	}

	tt := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"case ok return secret",
			fields{
				manager: AWSSecretStub{
					secret: "test",
					err:    nil,
				},
			},
			args{
				secretName: "test",
			},
			"test",
			false,
		},
		{
			"case error getting secret",
			fields{
				manager: AWSSecretStub{
					secret: "",
					err:    errors.New("not found"),
				},
			},
			args{
				secretName: "test",
			},
			"",
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := secrets.Secret{
				AWSManager: tc.fields.manager,
			}

			got, err := s.GetSecretByName(tc.args.secretName)
			if tc.wantErr {
				require.Error(t, err)
				require.Equal(t, "", got)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
