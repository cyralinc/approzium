package approzium

import (
	sdk "github.com/cyralinc/approzium/sdk/go/approzium"
)

// This is here to improve the UX for Go developers. It enables them to simply
// run "go get github.com/cyralinc/approzium" to get the Go SDK, which is normal
// for them, and for their imports to simply be "github.com/cyralinc/approzium".
// However, the meat of our SDK does live further down in our project.
func NewAuthClient(grpcAddress string, config *Config) (*AuthClient, error) {
	authClient, err := sdk.NewAuthClient(grpcAddress, config.Config)
	if err != nil {
		return nil, err
	}
	return &AuthClient{authClient}, nil
}

type AuthClient struct {
	*sdk.AuthClient
}

type Config struct {
	*sdk.Config
}
