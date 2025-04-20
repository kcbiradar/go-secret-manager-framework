package main

import (
	"context"
	"fmt"
	"log"

	secretsmanager "github.com/kcbiradar/go-secret-manager/src"
)

func main() {
	// Create a new secrets client
	client, err := secretsmanager.GetAwsClient(secretsmanager.AwsSecretClientConfig{
		Application: "common",
		Environment: "staging",
		Region:      "us-east-2",
		CacheTTL:    600, // 10 minutes
	})

	if err != nil {
		log.Fatalf("Failed to create AWS Secrets Manager client: %v", err)
	}

	ctx := context.Background()

	// Example: Get a single secret
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
