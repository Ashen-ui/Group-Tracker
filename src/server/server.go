package server

import (
	"fmt"
	"log"
	"net/http"
)

// The server function starts on the web server
func Server() {
	mux := http.NewServeMux()

	//Static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/artists/", ArtistDetailHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/api/artists", apiArtistsHandler)
	mux.HandleFunc("/", indexHandler)

	//Starting the server
	fmt.Println("Serveur démarré sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
