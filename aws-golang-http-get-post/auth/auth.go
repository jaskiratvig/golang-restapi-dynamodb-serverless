package auth

import (
	"aws-golang-http-get-post/models"
	"context"
	"log"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

//NewAuthenticator creates new authenticator a via Auth0 credentials from .env file
func NewAuthenticator() (*models.Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://jasvig.auth0.com/")
	if err != nil {
		log.Printf("Failed to get provider :%v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     "kf9yX2qaBa7J5tV1PtL5SpcdZ2GXHEo9",
		ClientSecret: "q21XgkalMEr1WL81tx62ZgErY2-ZJTS5o3xs_ntm8uiwMiib5Z7N5SL92TLux2_1",
		RedirectURL:  "https://bhvn5rgkmd.execute-api.us-east-1.amazonaws.com/dev/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &models.Authenticator{
		Ctx:      ctx,
		Provider: provider,
		Config:   conf,
	}, nil

}
