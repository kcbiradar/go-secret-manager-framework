package secretsmanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smithy "github.com/aws/smithy-go"
)

type AwsSecretClient struct {
	client *secretsmanager.Client
	config AwsSecretClientConfig
	cache  *SecretCache
}

// NewAwsSecretClient creates a new AWS Secrets Manager client
func NewAwsSecretClient(cfg AwsSecretClientConfig) (*AwsSecretClient, error) {
	// Set default values if not provided
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 600 // Default to 10 minutes
	}

	if cfg.Region == "" {
		cfg.Region = "us-east-1" // Default region
	}

	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(awsCfg)

	return &AwsSecretClient{
		client: client,
		config: cfg,
		cache:  NewSecretCache(cfg.CacheTTL),
	}, nil
}

// Get retrieves a secret by its key
func (c *AwsSecretClient) Get(ctx context.Context, secretKey string) (string, error) {
	formats := c.getAllSecretFormats(secretKey)
	log.Printf("Trying secret formats: %v", formats)

	var lastError error
	for _, format := range formats {
		value, err := c.fetchSecret(ctx, format)
		if err == nil {
			return value, nil
		}
		lastError = err
	}

	if lastError != nil {
		return "", lastError
	}
	return "", NewSecretNotFoundError(fmt.Sprintf("Secret %s not found in any format", secretKey))
}

// GetMultiple retrieves multiple secrets by their keys
func (c *AwsSecretClient) GetMultiple(ctx context.Context, secretKeys []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, key := range secretKeys {
		value, err := c.Get(ctx, key)
		if err != nil {
			if !IsSecretNotFoundError(err) {
				return nil, err
			}
			// If it's a "not found" error, we just continue to the next key
			continue
		}
		result[key] = value
	}

	return result, nil
}

// getAllSecretFormats returns all possible formats for a secret name
func (c *AwsSecretClient) getAllSecretFormats(secretKey string) []string {
	app := c.config.Application
	env := c.config.Environment

	if strings.Contains(secretKey, "/") {
		cleanKey := secretKey
		if strings.HasPrefix(cleanKey, "/") {
			cleanKey = cleanKey[1:]
		}

		parts := strings.Split(cleanKey, "/")
		var filteredParts []string
		for _, part := range parts {
			if part != "" {
				filteredParts = append(filteredParts, part)
			}
		}

		if len(filteredParts) >= 3 {
			actualKey := filteredParts[len(filteredParts)-1]
			return []string{
				secretKey,
				fmt.Sprintf("%s/%s/%s", app, env, actualKey),
				fmt.Sprintf("/%s/%s/%s", app, env, actualKey),
			}
		}
	}

	return []string{
		secretKey,
		fmt.Sprintf("%s/%s/%s", app, env, secretKey),
		fmt.Sprintf("/%s/%s/%s", app, env, secretKey),
	}
}

// fetchSecret retrieves a secret from AWS Secrets Manager or cache
func (c *AwsSecretClient) fetchSecret(ctx context.Context, secretName string) (string, error) {
	// Normalize to prevent double slashes
	normalizedSecretName := strings.Replace(secretName, "//", "/", -1)

	// Check cache first
	if cachedValue, found := c.cache.Get(normalizedSecretName); found {
		return cachedValue, nil
	}

	// If not in cache or expired, fetch from AWS
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(normalizedSecretName),
	}

	result, err := c.client.GetSecretValue(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "ResourceNotFoundException" {
				return "", NewSecretNotFoundError(fmt.Sprintf("The secret %s is not found in AWS secrets manager", normalizedSecretName))
			}
		}
		return "", err
	}

	if result.SecretString == nil {
		return "", NewSecretNotFoundError(fmt.Sprintf("Secret %s exists but has no string value", normalizedSecretName))
	}

	// Store in cache
	c.cache.Set(normalizedSecretName, *result.SecretString)

	return *result.SecretString, nil
}
