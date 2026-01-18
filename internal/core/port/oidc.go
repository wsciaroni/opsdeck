package port

import (
	"context"
)

// UserInfo represents basic user information retrieved from an OIDC provider.
type UserInfo struct {
	Email     string
	Name      string
	AvatarURL string
}

// OIDCProvider defines the interface for interacting with an OIDC provider.
type OIDCProvider interface {
	AuthCodeURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*UserInfo, error)
}
