package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/ktr0731/go-fuzzyfinder"
)

type SecretContent struct {
	User     string
	Password string
}

func getUserSecretsForHost(host string) []types.SecretListEntry {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("configuration error: " + err.Error())
	}

	svc := secretsmanager.NewFromConfig(cfg)

	l_params := &secretsmanager.ListSecretsInput{
		Filters: []types.Filter{
			types.Filter{
				Key:    "tag-key",
				Values: []string{"birdie:database_name"},
			},
			types.Filter{
				Key:    "tag-value",
				Values: []string{host},
			},
		},
	}

	secrets, err := svc.ListSecrets(context.Background(), l_params)

	if err != nil {
		log.Fatal(err)
	}

	return secrets.SecretList
}

func getPasswordForUser(secretId string) (string, string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("configuration error: " + err.Error())
	}

	svc := secretsmanager.NewFromConfig(cfg)

	res, err := svc.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretId,
	})

	if err != nil {
		log.Fatal(err)
	}

	var secretContent SecretContent

	err = json.Unmarshal([]byte(*res.SecretString), &secretContent)

	if err != nil {
		log.Fatal(err)
	}

	return secretContent.User, secretContent.Password

}

func FuzzUsers(host string) (string, string) {
	userSecrets := getUserSecretsForHost(host)

	idx, err := fuzzyfinder.Find(
		userSecrets,
		func(i int) string {
			for _, tag := range userSecrets[i].Tags {
				if *tag.Key == "birdie:role_name" {
					return *tag.Value
				}
			}
			return "Unknown"
		},
		fuzzyfinder.WithHeader("Select a user to connect with:"),
	)

	if err != nil {
		log.Fatal(err)
	}

	return getPasswordForUser(*userSecrets[idx].Name)
}
