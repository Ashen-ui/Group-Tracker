package server

import (
	"group-tracker/src/api"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Récupération des artistes depuis l'API
	artists, err := api.GetArtists()
	if err != nil {
		log.Printf("Error while getting artists: %v", err)
		http.Error(w, "Error while getting data", http.StatusInternalServerError)
		return
	}
	// Parsing du template
	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		log.Printf("Error while parsing template: %v", err)
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}
	// Exécution du template avec les données
	err = tmpl.Execute(w, artists)
	if err != nil {
		log.Printf("Error while executing template: %v", err)
		http.Error(w, "Display error", http.StatusInternalServerError)
		return
	}
	log.Println("Index page succesfully loaded")
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("aboutHandler: request received: %s %s", r.Method, r.URL.Path)

	// Contenu par défaut pour la page À propos
	readmeContent := `Group Tracker

Welcome to the group tracker, a web app that lets you discover 
different artists and their concerts

Functionalities:
- Search by name
- Complete artist list display
- Full details about artists (members, creation dates, albums)
- Concert and location information
- Date and location relations

Used Tech:
- Go (Golang) for the backend
- HTML/CSS for the frontend
- API REST to get the data

Developed with the passion for music!`

	//Template parsing
	tmpl, err := template.ParseFiles("./template/about.html")
	if err != nil {
		log.Printf("Error while parsing template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	//Template execution with data
	err = tmpl.Execute(w, AboutPageData{Readme: readmeContent})
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Display Error", http.StatusInternalServerError)
		return
	}
	log.Println("Page was successfuly loaded")
}

type SearchQuery struct {
	Name, Genre, Type, Location string
}
type PageData struct {
	Artists []api.Artist
	Query   SearchQuery
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("request received: %s %s", r.Method, r.URL.Path)

	query := SearchQuery{
		Name:     r.URL.Query().Get("name"),
		Genre:    r.URL.Query().Get("genre"),
		Type:     r.URL.Query().Get("type"),
		Location: r.URL.Query().Get("location"),
	}

	log.Printf("Search: name='%s' genre='%s' type='%s' location='%s'", query.Name, query.Genre, query.Type, query.Location)

	//Template parsing
	tmpl, err := template.ParseFiles("./template/search.html")
	if err != nil {
		log.Printf("Error while parsing template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	artists, err := api.FilterSearch(query.Name, query.Genre, query.Type, query.Location)
	if err != nil {
		log.Printf("Error while searching: %v", err)
		artists = []api.Artist{}
	}

	err = tmpl.Execute(w, PageData{Artists: artists, Query: query})
	if err != nil {
		log.Printf("Error while executing the template: %v", err)
		http.Error(w, "Display Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Page was successfuly loaded: %d results", len(artists))
}

type AboutPageData struct {
	Readme string
}
type ArtistData struct {
	Artist    api.Artist
	Locations []string
	Dates     []string
	Relations map[string][]string
}

func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== ArtistDetailHandler called for: %s ===", r.URL.Path)
	//Uses ID from URL /artists/{id}
	//Example: /artists/1 -> path = "1"
	path := r.URL.Path
	if !strings.HasPrefix(path, "/artists/") {
		log.Printf("Error: path doesn't start with /artists/: %s", path)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idStr := strings.TrimPrefix(path, "/artists/")
	idStr = strings.TrimSuffix(idStr, "/")
	idStr = strings.TrimSpace(idStr)
	if idStr == "" {
		log.Printf("Error: empty ID in the path %s", path)
		http.Error(w, "ID missing", http.StatusBadRequest)
		return
	}
	log.Printf("ID extract: '%s' from path: %s", idStr, path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting ID: '%s' isn't a number (erreur: %v)", idStr, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	log.Printf("ID successfuly converted: %d", id)
	//Getting artist
	artists, err := api.GetArtists()
	if err != nil {
		log.Printf("Error while getting artists: %v", err)
		http.Error(w, "Error while getting data", http.StatusInternalServerError)
		return
	}
	var artist *api.Artist
	for i := range artists {
		if artists[i].ID == id {
			artist = &artists[i]
			break
		}
	}
	if artist == nil {
		log.Printf("Artist not found with ID: %d", id)
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}
	//Getting locations
	locations, err := api.GetLocations()
	if err != nil {
		log.Printf("Error while getting locations: %v", err)
		http.Error(w, "Error while getting data", http.StatusInternalServerError)
		return
	}
	var artistLocations []string
	for _, loc := range locations {
		if loc.ID == id {
			artistLocations = loc.Locations
			break
		}
	}
	//Getting dates
	dates, err := api.GetDates()
	if err != nil {
		log.Printf("Error while getting dates: %v", err)
		http.Error(w, "Error while getting data", http.StatusInternalServerError)
		return
	}

	var artistDates []string
	for _, date := range dates {
		if date.ID == id {
			artistDates = date.Dates
			break
		}
	}
	//Get the relations
	relations, err := api.GetRelations()
	if err != nil {
		log.Printf("Error while getting relations: %v", err)
		http.Error(w, "Error while getting data", http.StatusInternalServerError)
		return
	}
	var artistRelations map[string][]string
	for _, rel := range relations {
		if rel.ID == id {
			artistRelations = rel.DatesLocations
			break
		}
	}
	//Preparing data for the template
	data := ArtistData{
		Artist:    *artist,
		Locations: artistLocations,
		Dates:     artistDates,
		Relations: artistRelations,
	}
	//Parsing the template
	tmpl, err := template.ParseFiles("./template/artist.html")
	if err != nil {
		log.Printf("Error while parsing template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	//Template execution with data
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error while executing template: %v", err)
		http.Error(w, "Display error", http.StatusInternalServerError)
		return
	}
	log.Printf("Artist details found: %s (ID: %d)", artist.Name, id)
}

func apiArtistsHandler(w http.ResponseWriter, r *http.Request) {
	//Set headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//Header check
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Get from external API
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error: failed to fetch from API"}`))
		return
	}
	defer resp.Body.Close()

	//Post to client
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error: failed to read response"}`))
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
