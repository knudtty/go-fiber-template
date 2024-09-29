package utils

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func getGoogleInfo(oauth2Token *oauth2.Token) (userInfo, error) {
	var ui userInfo
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return ui, err
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return ui, errors.New("id_token field missing from google oauth2Token")
	}

	// Parse and verify ID Token payload.
	var verifier = provider.Verifier(&oidc.Config{ClientID: os.Getenv("GOOGLE_CLIENT_ID")})
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return ui, fmt.Errorf("Unable to verify id_token: %s", rawIDToken)
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		PictureUrl    string `json:"picture"`
		Name          string `json:"name"`
	}
	if err = idToken.Claims(&claims); err != nil {
		return ui, err
	}

	ui.email = claims.Email
	ui.id = idToken.Subject
	ui.avatarURL = claims.PictureUrl
	ui.name = claims.Name

	return ui, nil
}
