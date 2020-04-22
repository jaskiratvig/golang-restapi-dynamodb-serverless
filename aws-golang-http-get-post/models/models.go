package models

import (
	"context"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// Artist is our self-made struct to process JSON request from Client
type Artist struct {
	ArtistID    string   `json:"ArtistID"`
	Name        string   `json:"Name"`
	Songs       []string `json:"Songs"`
	Subcategory string   `json:"Subcategory"`
	Domestic    bool     `json:"Domestic"`
}

//Authenticator is our struct to create a new authenticator for Auth0
type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

// SessionData is our struct for saving the state of our authentication request
type SessionData struct {
	ClientID string                 `json:"ClientID"`
	State    string                 `json:"State"`
	Profile  map[string]interface{} `json:"Profile"`
	Token    string                 `json:"Token"`
}

func (art Artist) editArtist(artist Artist) Artist {
	art.ArtistID = artist.ArtistID
	art.Name = artist.Name
	art.Songs = artist.Songs
	art.Subcategory = artist.Subcategory
	art.Domestic = artist.Domestic
	return art
}
