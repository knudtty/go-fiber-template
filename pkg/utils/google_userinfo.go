package utils

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type googleUser struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func getGoogleInfo(oauth2Token *oauth2.Token) (string, string, error) {
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return "", "", err
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", "", errors.New("id_token field missing from google oauth2Token")
	}

	// Parse and verify ID Token payload.
	var verifier = provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_CLIENT_ID")})
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return "", "", fmt.Errorf("Unable to verify id_token: %s", rawIDToken)
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		PictureUrl    string `json:"picture"`
	}

	if err = idToken.Claims(&claims); err != nil {
		return "", "", err
	}
	return idToken.Subject, claims.Email, nil
}
