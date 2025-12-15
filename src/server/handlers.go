package server

import (
	"group-tracker/src/api"
	"html/template"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex

func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

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
}
