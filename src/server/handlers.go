package server

import (
	"group-tracker/src/api"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

// Bloque les accès concurrents à la base de données
var mu sync.Mutex

func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("Requête reçue: %s %s", r.Method, r.URL.Path)

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

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requête reçue: %s %s", r.Method, r.URL.Path)

	b, err := os.ReadFile("./README.md")
	if err != nil {
		log.Printf("Erreur lors de la lecture du README: %v", err)
		http.Error(w, "README.md introuvable", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./template/about.html")
	if err != nil {
		log.Printf("Erreur lors du parsing du template about: %v", err)
		http.Error(w, "Template about.html introuvable", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, AboutPageData{Readme: string(b)})
	if err != nil {
		log.Printf("Erreur lors de l'exécution du template about: %v", err)
		http.Error(w, "Erreur d'affichage", http.StatusInternalServerError)
		return
	}

	log.Println("Page À propos servie avec succès")
}
