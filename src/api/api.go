package api

import (
	"encoding/json"
	"io"
	"net/http"
)

// URL de base de l'API
const BaseURL = "https://groupietrackers.herokuapp.com/api"

// La structure Artist représente un artiste/groupe
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

// La structure Location représente les lieux de concert
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

// La structure Date représente les dates de concert
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

// La structure Relation représente la relation entre dates et lieux
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// La fonction GetArtists récupère tous les artistes depuis l'API
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

// La fonction GetLocations récupère toutes les locations depuis l'API
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

// La fonction GetDates récupère toutes les dates depuis l'API
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

// La fonction GetRelations récupère toutes les relations depuis l'API
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
