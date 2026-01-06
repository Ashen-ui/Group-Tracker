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
	tmpl, err := template.ParseFiles("./templates/artists.html")
	if err != nil {
		http.Error(w, "Erreur template", http.StatusInternalServerError)
		return
	}

	artists, err := modules.FetchAllArtists()
	if err != nil {
		return
	}

	recherche := r.URL.Query().Get("query")
	result := modules.Recherche(recherche, artists)
	println("------------------- Debut de la recherhce -------------------------------")
	println("il a trouver ", len(result), " résultat(s) pour la recherche de ", recherche)
	println(result)
	for i := range result {
		println("Artiste trouvé :", result[i].Groupe.Name)
		println("Champs correspondants :", result[i].Where)
		println("Nom corespondant a le recherche :", result[i].Name)
	}
	println("------------------- Fin de la recherche-------------------------------")

	data := struct{ Artists []modules.Resultat }{Artists: result}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur exécution template", http.StatusInternalServerError)
		return
	}
}
