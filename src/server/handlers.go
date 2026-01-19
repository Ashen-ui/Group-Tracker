package server

import (
	"group-tracker/src/api"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Bloque les accès concurrents à la base de données
var mu sync.Mutex

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Ne pas traiter les requêtes vers /artists/, /about ou /search
	if strings.HasPrefix(r.URL.Path, "/artists/") || r.URL.Path == "/about" || r.URL.Path == "/search" {
		log.Printf("indexHandler: requête ignorée: %s", r.URL.Path)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	log.Printf("indexHandler: Requête reçue: %s %s", r.Method, r.URL.Path)
	// Récupération des artistes depuis l'API
	artists, err := api.GetArtists()
	if err != nil {
		log.Printf("Erreur lors de la récupération des artistes: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		return
	}
	// Parsing du template
	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		log.Printf("Erreur lors du parsing du template: %v", err)
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}
	// Exécution du template avec les données
	err = tmpl.Execute(w, artists)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
		return
	}
	log.Println("Page servie avec succès")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("aboutHandler: Requête reçue: %s %s", r.Method, r.URL.Path)

	// Contenu par défaut pour la page À propos
	readmeContent := `Groupie Tracker

Bienvenue sur Groupie Tracker, une application web permettant de découvrir 
les artistes et leurs concerts.

Fonctionnalités:
- Recherche d'artistes par nom
- Affichage de la liste complète des artistes
- Détails complets pour chaque artiste (membres, dates de création, albums)
- Informations sur les concerts et locations
- Relations entre dates et locations

Technologies utilisées:
- Go (Golang) pour le backend
- HTML/CSS pour le frontend
- API REST pour la récupération des données

Développé avec passion pour les fans de musique!`

	// Parsing du template
	tmpl, err := template.ParseFiles("./template/about.html")
	if err != nil {
		log.Printf("Erreur lors du parsing du template: %v", err)
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}
	// Exécution du template avec les données
	err = tmpl.Execute(w, AboutPageData{Readme: readmeContent})
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
		return
	}
	log.Println("Page servie avec succès")
}
func search_bar_handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("Requête reçue: %s %s", r.Method, r.URL.Path)

	searchName := r.URL.Query().Get("name")
	log.Printf("Recherche de: '%s'", searchName)
	// Parsing du template
	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		log.Printf("Erreur lors du parsing du template: %v", err)
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}
	// Si aucun nom n'est fourni, afficher tous les artistes
	if searchName == "" {
		artists, err := api.GetArtists()
		if err != nil {
			log.Printf("Erreur lors de la récupération des artistes: %v", err)
			http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, artists)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template: %v", err)
			http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
			return
		}
		return
	}
	// Recherche de l'artiste
	artist, err := api.SearchBar(searchName)
	if err != nil {
		log.Printf("Artiste non trouvé: '%s'", searchName)
		// Afficher une liste vide au lieu d'une erreur
		artists := []api.Artist{}
		err = tmpl.Execute(w, artists)
		if err != nil {
			log.Printf("Erreur lors de l'exécution du template: %v", err)
			http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
			return
		}
		return
	}
	// Convertir l'artiste unique en slice pour le template
	artists := []api.Artist{artist}
	// Exécution du template avec les données
	err = tmpl.Execute(w, artists)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
		return
	}

	log.Printf("Recherche servie avec succès: %s", artist.Name)
}

type AboutPageData struct {
	Readme string
}
type ArtistDetailData struct {
	Artist    api.Artist
	Locations []string
	Dates     []string
	Relations map[string][]string
}

func ArtistDetailHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("=== ArtistDetailHandler appelé pour: %s ===", r.URL.Path)
	// Extraction de l'ID depuis l'URL /artists/{id}
	// Exemple: /artists/1 -> path = "1"
	path := r.URL.Path
	if !strings.HasPrefix(path, "/artists/") {
		log.Printf("Erreur: le path ne commence pas par /artists/: %s", path)
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}
	idStr := strings.TrimPrefix(path, "/artists/")
	idStr = strings.TrimSuffix(idStr, "/")
	idStr = strings.TrimSpace(idStr)
	if idStr == "" {
		log.Printf("Erreur: ID vide dans le path: %s", path)
		http.Error(w, "ID manquant", http.StatusBadRequest)
		return
	}
	log.Printf("ID extrait: '%s' depuis path: %s", idStr, path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Erreur conversion ID: '%s' n'est pas un nombre (erreur: %v)", idStr, err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}
	log.Printf("ID converti avec succès: %d", id)
	// Récupération de l'artiste
	artists, err := api.GetArtists()
	if err != nil {
		log.Printf("Erreur lors de la récupération des artistes: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
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
		log.Printf("Artiste non trouvé avec l'ID: %d", id)
		http.Error(w, "Artiste non trouvé", http.StatusNotFound)
		return
	}
	// Récupération des locations
	locations, err := api.GetLocations()
	if err != nil {
		log.Printf("Erreur lors de la récupération des locations: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		return
	}
	var artistLocations []string
	for _, loc := range locations {
		if loc.ID == id {
			artistLocations = loc.Locations
			break
		}
	}
	// Récupération des dates
	dates, err := api.GetDates()
	if err != nil {
		log.Printf("Erreur lors de la récupération des dates: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		return
	}

	var artistDates []string
	for _, date := range dates {
		if date.ID == id {
			artistDates = date.Dates
			break
		}
	}
	// Récupération des relations
	relations, err := api.GetRelations()
	if err != nil {
		log.Printf("Erreur lors de la récupération des relations: %v", err)
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		return
	}
	var artistRelations map[string][]string
	for _, rel := range relations {
		if rel.ID == id {
			artistRelations = rel.DatesLocations
			break
		}
	}
	// Préparation des données pour le template
	data := ArtistDetailData{
		Artist:    *artist,
		Locations: artistLocations,
		Dates:     artistDates,
		Relations: artistRelations,
	}
	// Parsing du template
	tmpl, err := template.ParseFiles("./template/artist.html")
	if err != nil {
		log.Printf("Erreur lors du parsing du template: %v", err)
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}
	// Exécution du template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
		return
	}
	log.Printf("Détails de l'artiste servis: %s (ID: %d)", artist.Name, id)
}
