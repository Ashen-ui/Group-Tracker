package server

import (
	"fmt"
	"log"
	"net/http"
)

// La fonction Server démarre le serveur web
func Server() {
	mux := http.NewServeMux()

	// Fichiers statiques
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/artists/", ArtistDetailHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/", indexHandler)

	// Démarrage du serveur
	fmt.Println("Serveur démarré sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
