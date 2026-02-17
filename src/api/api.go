package api

import (
	"encoding/json"
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

// SearchBar
func SearchBar(name string) ([]Artist, error) {
	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}
	lowCaseName := strings.ToLower(strings.TrimSpace(name))
	if lowCaseName == "" {
		return artists, nil //return all artists if empty
	}
	var matchingArtists []Artist
	for _, artist := range artists {
		artistlowCaseName := strings.ToLower(artist.Name)
		if strings.Contains(artistlowCaseName, lowCaseName) {
			matchingArtists = append(matchingArtists, artist)
			continue //To avoid adding twice if already member
		}
		for _, member := range artist.Members {
			memberLower := strings.ToLower(member)
			if strings.Contains(memberLower, lowCaseName) {
				matchingArtists = append(matchingArtists, artist)
				break //Add once per artist
			}
		}
	}
	return matchingArtists, nil
}

var genreMapping = map[string][]string{
	"Rock": {
		"Queen", "Pink Floyd", "Scorpions", "ACDC", "Pearl Jam", "Genesis", "Phil Collins", "Led Zeppelin",
		"The Jimi Hendrix Experience", "Bee Gees", "Deep Purple", "Aerosmith", "Dire Straits",
		"The Rolling Stones", "U2", "Guns N' Roses", "Eagles", "Linkin Park", "Red Hot Chili Peppers",
		"Green Day", "Metallica", "Coldplay", "Foo Fighters", "Arctic Monkeys", "Fall Out Boy",
	},
	"Hip-Hop/Rap": {
		"XXXTentacion", "Mac Miller", "Joyner Lucas", "Kendrick Lamar", "Juice Wrld", "Logic",
		"Post Malone", "Travis Scott", "J. Cole", "Mobb Deep", "NWA", "Eminem",
	},
	"Pop": {
		"Katy Perry", "Rihanna", "Thirty Seconds to Mars", "Imagine Dragons", "Maroon 5",
		"Twenty One Pilots", "The Chainsmokers",
	},
	"Electronic/Alternative": {
		"SOJA", "Gorillaz", "R3HAB", "Muse", "Nickelback",
	},
}

func ArtistGenreSearch(artistName string) string {
	lowCaseName := strings.ToLower(artistName)
	for genre, bands := range genreMapping {
		for _, band := range bands {
			if strings.Contains(lowCaseName, strings.ToLower(band)) {
				return genre
			}
		}
	}
	return "Other"
}

func FilterSearch(name, genre, artistType, location string) ([]Artist, error) {
	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}

	locations, err := GetLocations()
	if err != nil {
		return nil, err
	}

	lowCaseName := strings.ToLower(strings.TrimSpace(name))
	lowCaseLocation := strings.ToLower(strings.TrimSpace(location))
	genre = strings.TrimSpace(genre)

	var results []Artist

	for _, i := range artists {
		if lowCaseName != "" {
			matched := false
			if strings.Contains(strings.ToLower(i.Name), lowCaseName) {
				matched = true
			} else {
				for _, j := range i.Members {
					if strings.Contains(strings.ToLower(j), lowCaseName) {
						matched = true
						break
					}
				}
			}
			if !matched {
				continue
			}
		}

		//Genre
		if genre != "" {
			if ArtistGenreSearch(i.Name) != genre {
				continue
			}
		}

		//Type
		if artistType == "solo" && len(i.Members) != 1 {
			continue
		}
		if artistType == "group" && len(i.Members) == 1 {
			continue
		}

		//Location
		if lowCaseLocation != "" {
			found := false
			for _, loc := range locations {
				if loc.ID != i.ID {
					continue
				}
				for _, l := range loc.Locations {
					if strings.Contains(strings.ToLower(l), lowCaseLocation) {
						found = true
						break
					}
				}
				break
			}
			if !found {
				continue
			}
		}

		results = append(results, i)
	}

	return results, nil
}
