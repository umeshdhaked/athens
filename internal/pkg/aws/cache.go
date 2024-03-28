package aws

import (
	"fmt"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	secretCache, _ = secretcache.New()
)

// GetSecretKeyValue retrieves the value associated with a key from the cache
func GetSecretKeyValue(key string) (string, error) {
	// Retrieve the secret value from the cache
	secretValue, err := secretCache.GetSecretString(key)
	if err != nil {
		return "", fmt.Errorf("error retrieving secret value: %v", err)
	}
	return secretValue, nil
}
