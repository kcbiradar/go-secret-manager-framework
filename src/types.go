package secretsmanager

import "context"

// SecretClient defines the interface for secret management
type SecretClient interface {
	Get(ctx context.Context, secretKey string) (string, error)
	GetMultiple(ctx context.Context, secretKeys []string) (map[string]string, error)
}

// AwsSecretClientConfig holds configuration for the AWS Secrets Manager client
type AwsSecretClientConfig struct {
	Application string
	Environment string
	Region      string
	CacheTTL    int64 // in seconds
}
