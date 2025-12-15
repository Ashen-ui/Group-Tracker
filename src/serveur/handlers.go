package serveur

import (
	"group-tracker/src/modules"
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := modules.FetchAllArtists()
	if err != nil {
		http.Error(w, "Impossible de récupérer les artistes", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}
	data := struct{ Artists []modules.Artist }{Artists: artists}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur exécution template", http.StatusInternalServerError)
		return
	}
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/artists.html")
}
