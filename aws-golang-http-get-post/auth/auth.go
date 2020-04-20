package auth

import (
	"aws-golang-http-get-post/models"
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

//NewAuthenticator creates new authenticator a via Auth0 credentials from .env file
func NewAuthenticator() (*models.Authenticator, error) {
	ctx := context.Background()

	sess := session.New()
	svc := ssm.New(sess)

	domain, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/Domain"),
		},
	)
	clientID, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/ClientID"),
		},
	)
	clientSecret, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/ClientSecret"),
		},
	)
	redirectURL, err := svc.GetParameter(
		&ssm.GetParameterInput{
			Name: aws.String("/dev/RedirectURL"),
		},
	)

	provider, err := oidc.NewProvider(ctx, aws.StringValue(domain.Parameter.Value))
	if err != nil {
		log.Printf("Failed to get provider :%v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     aws.StringValue(clientID.Parameter.Value),
		ClientSecret: aws.StringValue(clientSecret.Parameter.Value),
		RedirectURL:  aws.StringValue(redirectURL.Parameter.Value),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &models.Authenticator{
		Ctx:      ctx,
		Provider: provider,
		Config:   conf,
	}, nil

}
