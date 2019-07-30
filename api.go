package infura

import "errors"

// Provider is an implementation of the IPFS provider for Infura.
type Provider struct {
	//	apiKey    string
	//	apiSecret string
}

// NewProvider creates a new Infura provider.
//func NewProvider(apiKey string, apiSecret string) (*Provider, error) {
func NewProvider() (*Provider, error) {

	provider := &Provider{}

	// Try a ping to ensure the service and credentials look good
	alive, err := provider.Ping()
	if err != nil {
		return nil, err
	}
	if !alive {
		return nil, errors.New("service unavailable")
	}
	return provider, nil
}
