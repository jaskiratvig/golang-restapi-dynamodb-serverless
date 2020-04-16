package models

// Artist is our self-made struct to process JSON request from Client
type Artist struct {
	ArtistID    string   `json:"ArtistID"`
	Name        string   `json:"Name"`
	Songs       []string `json:"Songs"`
	Subcategory string   `json:"Subcategory"`
	Domestic    bool     `json:"Domestic"`
}

func (art Artist) editArtist(artist Artist) Artist {
	art.ArtistID = artist.ArtistID
	art.Name = artist.Name
	art.Songs = artist.Songs
	art.Subcategory = artist.Subcategory
	art.Domestic = artist.Domestic
	return art
}
