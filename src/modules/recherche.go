package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// exemple
//	artists{
//	    "id": 1,
//	    "image": "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
//	    "name": "Queen",
//	    "members": [
//	      "Freddie Mercury",
//	      "Brian May",
//	      "John Daecon",
//	      "Roger Meddows-Taylor",
//	      "Mike Grose",
//	      "Barry Mitchell",
//	      "Doug Fogie"
//	    ],
//	    "creationDate": 1970,
//	    "firstAlbum": "14-12-1973",
//	    "locations": "https://groupietrackers.herokuapp.com/api/locations/1",
//	    "concertDates": "https://groupietrackers.herokuapp.com/api/dates/1",
//	    "relations": "https://groupietrackers.herokuapp.com/api/relation/1"
//	  },

// Permet d'avoir en locale la liste de tout les Groupes (en gros de récuperer tout les données de l'api)
func FetchAllArtists() ([]Artist, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status non-OK: %d", resp.StatusCode)
	}
	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}
	return artists, nil
}

type Resultat struct {
	Groupe Artist //le groupe trouvé
	Where  string //le champ où le mot clef a été trouvé savoir si c est le nom du groupe ou un des membres
	Name   string // le "nom" trouver
}

// mot clef = se que tu met dans la bare de recherche ex: "queen"
//
// artists = la liste de tout les artistes récupérés par FetchAllArtists
func Recherche(motclef string, artists []Artist) []Resultat {
	var rep []Resultat
	var name string
	motclef = strings.ToLower(strings.TrimSpace(motclef))
	if motclef == "" {
		return rep
	}
	for _, artist := range artists {
		where := ""

		if strings.Contains(strings.ToLower(artist.Name), motclef) {
			where = "name"
			name = artist.Name
		}

		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), motclef) {
				where = "members"
				name = member
				break
			}
		}

		if len(where) > 0 {
			rep = append(rep, Resultat{Groupe: artist, Where: where, Name: name})
		}
	}
	return rep
}
