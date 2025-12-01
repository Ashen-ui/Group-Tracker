package serveur

import (
	"fmt"
	"log"
	"net/http"
)

func Serveur() {
	// Fichiers statiques
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers
	http.HandleFunc("/", indexHandler)

	// Démarrage du serveur
	fmt.Println("Serveur démarré sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
