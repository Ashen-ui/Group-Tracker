package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const BaseURL = "https://groupietrackers.herokuapp.com/api"

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func GetArtists() ([]Artist, error) {
	resp, err := http.Get(BaseURL + "/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func GetLocations() ([]Location, error) {
	resp, err := http.Get(BaseURL + "/locations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Index []Location `json:"index"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Index, nil
}

func GetDates() ([]Date, error) {
	resp, err := http.Get(BaseURL + "/dates")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Index []Date `json:"index"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Index, nil
}

func GetRelations() ([]Relation, error) {
	resp, err := http.Get(BaseURL + "/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Index []Relation `json:"index"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Index, nil
}

// Barre de recherche
func SearchBar(name string) (Artist, error) {
	artists, err := GetArtists()
	if err != nil {
		return Artist{}, err
	}
	nameLower := strings.ToLower(strings.TrimSpace(name))
	if nameLower == "" {
		return Artist{}, errors.New("nom de recherche vide")
	}
	for _, artist := range artists {
		artistNameLower := strings.ToLower(artist.Name)
		if artistNameLower == nameLower || strings.Contains(artistNameLower, nameLower) {
			return artist, nil
		}
	}
	return Artist{}, errors.New("artist not found")
}
