package main

import (
	"context"
	"fmt"
	"log"

	secretsmanager "github.com/kcbiradar/go-secret-manager/src"
)

func main() {
	client, err := secretsmanager.GetAwsClient(secretsmanager.AwsSecretClientConfig{
		Application: "XXX",
		Environment: "XXX",
		Region:      "us-east-1",
		CacheTTL:    600,
	})

	if err != nil {
		log.Fatalf("Failed to create AWS Secrets Manager client: %v", err)
	}

	ctx := context.Background()

	secretValue, err := client.Get(ctx, "homewise_bot_api_key")
	if err != nil {
		if secretsmanager.IsSecretNotFoundError(err) {
			log.Printf("Secret not found: %v", err)
		} else {
			log.Fatalf("Error retrieving secret: %v", err)
		}
	} else {
		fmt.Printf("Secret value: %s\n", secretValue)
	}
}
