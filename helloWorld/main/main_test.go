package main

import (
	"testing"
)

func TestCreateArtist(t *testing.T) {
	var s []string
	s = append(s, "Take You Down")
	s = append(s, "Crashing")
	s = append(s, "Takeaway")
	
	a := createArtist("Illenium", s, "Future Bass", true)

	if a.Name != "Illenium" {
		t.Errorf("createArtist(Illenium) failed, expected %s, got %s", "Illenium", a.Name)
	} else {
		t.Logf("createArtist(Illenium) success, expected %s, got %s", "Illenium", a.Name)
	}
}

func TestGetArtist(t *testing.T) {
	var s []string
	s = append(s, "Take You Down")
	s = append(s, "Crashing")
	s = append(s, "Takeaway")
	
	_ = createArtist("Illenium", s, "Future Bass", true)
	a := getArtist("Illenium")

	if a.Name != "Illenium" {
		t.Errorf("getArtist(Illenium) failed, expected %s, got %s", "Illenium", a.Name)
	} else {
		t.Logf("getArtist(Illenium) success, expected %s, got %s", "Illenium", a.Name)
	}
}

func TestGetAllArtists(t *testing.T) {
	var s []string
	s = append(s, "Take You Down")
	s = append(s, "Crashing")
	s = append(s, "Takeaway")
	
	_ = createArtist("Illenium", s, "Future Bass", true)

	if len(artists) == 0 {
		t.Errorf("getAllArtists(Illenium) failed, expected %s, got %s", "More than one artist", artists)
	} else {
		t.Logf("getAllArtists(Illenium) success, expected %s, got %s", "More than one artist", artists)
	}
}

func TestEditArtist(t *testing.T) {
	var s []string
	s = append(s, "Take You Down")
	s = append(s, "Crashing")
	s = append(s, "Takeaway")
	
	a := createArtist("Illenium", s, "Future Bass", true)

	a = artists[0]
	a.Domestic = true
	a = a.editArtist(a)

	if a.Domestic != true {
		t.Errorf("editArtist(Illenium) failed, expected %v, got %v", true, a.Domestic)
	} else {
		t.Logf("editArtist(Illenium) success, expected %v, got %v", true, a.Domestic)
	}
}