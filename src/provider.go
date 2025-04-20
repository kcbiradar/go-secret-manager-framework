package secretsmanager

// GetAwsClient creates and returns an AWS Secrets Manager client
func GetAwsClient(config AwsSecretClientConfig) (SecretClient, error) {
	if config.Application == "" {
		return nil, NewConfigurationError("Application name is required")
	}

	if config.Environment == "" {
		return nil, NewConfigurationError("Environment is required")
	}

	if config.Region == "" {
		return nil, NewConfigurationError("AWS region is required")
	}

	if config.CacheTTL < 0 {
		return nil, NewConfigurationError("CacheTTL must be a positive number")
	}

	return NewAwsSecretClient(config)
}
