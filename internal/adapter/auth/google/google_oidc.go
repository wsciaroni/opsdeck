package google

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/wsciaroni/opsdeck/internal/core/port"
	"golang.org/x/oauth2"
)

// Ensure GoogleProvider implements port.OIDCProvider
var _ port.OIDCProvider = (*GoogleProvider)(nil)

type GoogleProvider struct {
	provider    *oidc.Provider
	oauthConfig *oauth2.Config
}

type claims struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewGoogleProvider(ctx context.Context, clientID, clientSecret, callbackURL string) (*GoogleProvider, error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create google oidc provider: %w", err)
	}

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &GoogleProvider{
		provider:    provider,
		oauthConfig: oauthConfig,
	}, nil
}

func (g *GoogleProvider) AuthCodeURL(state string) string {
	return g.oauthConfig.AuthCodeURL(state)
}

func (g *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*port.UserInfo, error) {
	oauth2Token, err := g.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	verifier := g.provider.Verifier(&oidc.Config{ClientID: g.oauthConfig.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %w", err)
	}

	var c claims
	if err := idToken.Claims(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &port.UserInfo{
		Email:     c.Email,
		Name:      c.Name,
		AvatarURL: c.Picture,
	}, nil
}
