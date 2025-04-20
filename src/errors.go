package secretsmanager

// SecretNotFoundError is returned when a secret cannot be found
type SecretNotFoundError struct {
	Message string
}

func (e *SecretNotFoundError) Error() string {
	return e.Message
}

// NewSecretNotFoundError creates a new SecretNotFoundError
func NewSecretNotFoundError(message string) *SecretNotFoundError {
	return &SecretNotFoundError{Message: message}
}

// ConfigurationError is returned when there's an issue with the client configuration
type ConfigurationError struct {
	Message string
}

func (e *ConfigurationError) Error() string {
	return e.Message
}

// NewConfigurationError creates a new ConfigurationError
func NewConfigurationError(message string) *ConfigurationError {
	return &ConfigurationError{Message: message}
}

// IsSecretNotFoundError checks if an error is a SecretNotFoundError
func IsSecretNotFoundError(err error) bool {
	_, ok := err.(*SecretNotFoundError)
	return ok
}
