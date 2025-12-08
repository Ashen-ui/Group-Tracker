package serveur

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Artist struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func fetchAllArtists() ([]Artist, error) {
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := fetchAllArtists()
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
	data := struct{ Artists []Artist }{Artists: artists}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur exécution template", http.StatusInternalServerError)
		return
	}
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/artists.html")
}
