package server

import (
	"fmt"
	"log"
	"net/http"
)

// La fonction Server démarre le serveur web
func Server() {
	// Fichiers statiques
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/about", AboutHandler)

	// Handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", search_bar_handler)
	// Démarrage du serveur
	fmt.Println("Serveur démarré sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
