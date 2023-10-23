package application

import (
	"context"
	log "github.com/sirupsen/logrus"
	"prem-gateway/auth/internal/core/domain"
	"sync"
	"time"
)

// ApiKeyService defines the interface for managing API keys.
type ApiKeyService interface {
	// CreateApiKey creates a new API key.
	CreateApiKey(ctx context.Context, key CreateApiKeyReq) (string, error)
	// AllowRequest checks if a given API key is allowed to access a specified path.
	AllowRequest(apiKey string, path string) bool
	// GetServiceApiKey fetches the API key for a specific service.
	GetServiceApiKey(ctx context.Context, service string) (string, error)
}

// NewApiKeyService constructs a new instance of the ApiKeyService.
func NewApiKeyService(
	ctx context.Context,
	apiKeyRepository domain.ApiKeyRepository,
) (ApiKeyService, error) {
	keysDb := make(map[string]apiKeyInfo)
	keys, err := apiKeyRepository.GetAllApiKeys(ctx) // Fetch all existing API keys from the repository.
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		keysDb[key.ID] = newApiKeyInfo(key) // Populate the in-memory cache with the fetched keys.
	}

	return &apiKeyService{
		apiKeyRepository: apiKeyRepository,
		keysDb:           keysDb,
	}, nil
}

// apiKeyService is the concrete implementation of the ApiKeyService interface.
type apiKeyService struct {
	apiKeyRepository domain.ApiKeyRepository // Repository for interacting with the API key datastore.

	keysMtx sync.RWMutex          // Mutex for safe concurrent access to the `keysDb` map.
	keysDb  map[string]apiKeyInfo // In-memory cache of API keys for fast lookup.

	rootKeyMtx sync.RWMutex // Mutex for safe concurrent access to the `rootApiKey` field.
	rootApiKey string       // The root API key with unrestricted access.
}

// CreateApiKey creates a new API key and saves it in the datastore.
func (a *apiKeyService) CreateApiKey(
	ctx context.Context, key CreateApiKeyReq,
) (string, error) {
	if key.IsRootKey {
		if a.rootKeyExists() {
			return "", ErrRootKeyExists // Ensure only one root key exists.
		}
	}

	// Construct a new API key domain object.
	apiKey, err := domain.NewApiKey(
		key.IsRootKey,
		key.Services,
		domain.RateLimit{
			RequestsPerRange: key.RequestsPerRange,
			RangeInSeconds:   key.RangeInSeconds,
		},
	)
	if err != nil {
		return "", err
	}

	// Save the new API key to the repository.
	if err = a.apiKeyRepository.CreateApiKey(ctx, *apiKey); err != nil {
		return "", err
	}

	a.insertKey(newApiKeyInfo(*apiKey)) // Cache the new key for quick lookup.

	log.Debugf("Created new API key %s", apiKey.ID)

	return apiKey.ID, nil
}

// AllowRequest checks if a given API key is allowed to access a specified path.
func (a *apiKeyService) AllowRequest(apiKey string, path string) bool {
	// Check if the key exists and if it's allowed to access the given path.
	key, exists := a.getKey(apiKey)
	if !exists || !key.canAccessServicePath(path) {
		log.Debugf("Api key %s is not allowed to access %s", apiKey, path)

		return false
	}

	if key.isRootKey {
		return true
	}

	// Check if the key has exceeded its rate limit.
	isRateLimited, aki := key.isRateLimited(time.Now())
	a.updateKey(aki)

	if isRateLimited {
		log.Debugf("Api key %s has exceeded its rate limit", apiKey)
	}
	log.Debugf("Api key %s is allowed to access %s", apiKey, path)

	return !isRateLimited
}

// GetServiceApiKey retrieves the API key associated with a specific service.
func (a *apiKeyService) GetServiceApiKey(
	ctx context.Context, service string,
) (string, error) {
	apiKey, err := a.apiKeyRepository.GetServiceApiKey(ctx, service)
	if err != nil {
		return "", err
	}

	return apiKey.ID, nil
}

// Below are helper methods for the apiKeyService.

func (a *apiKeyService) insertKey(key apiKeyInfo) {
	a.keysMtx.Lock()
	defer a.keysMtx.Unlock()

	a.keysDb[key.id] = key
}

func (a *apiKeyService) getKey(key string) (apiKeyInfo, bool) {
	a.keysMtx.RLock()
	defer a.keysMtx.RUnlock()

	keyInfo, exists := a.keysDb[key]
	return keyInfo, exists
}

func (a *apiKeyService) isRootKey(key string) bool {
	a.rootKeyMtx.RLock()
	defer a.rootKeyMtx.RUnlock()

	return a.rootApiKey == key
}

func (a *apiKeyService) rootKeyExists() bool {
	a.rootKeyMtx.RLock()
	defer a.rootKeyMtx.RUnlock()

	return a.rootApiKey != ""
}

func (a *apiKeyService) updateKey(key apiKeyInfo) {
	a.keysMtx.Lock()
	defer a.keysMtx.Unlock()

	a.keysDb[key.id] = key
}

// apiKeyInfo represents detailed information about an API key.
type apiKeyInfo struct {
	id                  string              // Unique identifier of the API key.
	allowedEndpoints    map[string]struct{} // List of services or paths the API key has access to.
	firstRequestInRange *time.Time          // Timestamp of the first request made within the current rate limit range.
	requestsPerRange    int                 // Max number of requests allowed within the rate limit range.
	rangeInSeconds      int                 // Duration of the rate limit range in seconds.
	requestCount        int                 // Number of requests made within the current rate limit range.
	isRootKey           bool                // Flag indicating if the API key is a root key with unrestricted access.
}

// Construct a new apiKeyInfo from a domain API key.
func newApiKeyInfo(apiKey domain.ApiKey) apiKeyInfo {
	allowedEndpoints := make(map[string]struct{})
	for _, endpoint := range apiKey.Services {
		allowedEndpoints[endpoint] = struct{}{}
	}

	return apiKeyInfo{
		id:                  apiKey.ID,
		allowedEndpoints:    allowedEndpoints,
		firstRequestInRange: nil,
		requestsPerRange:    apiKey.RateLimit.RequestsPerRange,
		rangeInSeconds:      apiKey.RateLimit.RangeInSeconds,
		requestCount:        0,
		isRootKey:           apiKey.IsRoot,
	}
}

// Check if the API key is allowed to access a given service or path.
func (ak apiKeyInfo) canAccessServicePath(servicePath string) bool {
	_, exists := ak.allowedEndpoints[servicePath]

	return exists
}

// Determine if the API key has exceeded its rate limit.
func (ak apiKeyInfo) isRateLimited(now time.Time) (bool, apiKeyInfo) {
	aki := ak

	// If it's the first request in the rate limit range
	if ak.firstRequestInRange == nil {
		aki.firstRequestInRange = &now
	}

	// Check if the current request is outside the rate limit range, if so reset the count and timestamp
	if now.Sub(*aki.firstRequestInRange).Seconds() >= float64(aki.rangeInSeconds) {
		aki.firstRequestInRange = &now
		aki.requestCount = 1

		return false, aki
	}

	aki.requestCount++

	// If the key has already reached its rate limit, just return true without incrementing the count
	if aki.requestCount > aki.requestsPerRange {
		return true, aki
	}

	return false, aki
}
