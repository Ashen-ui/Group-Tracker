package main

import (
	"fmt"
	"group-tracker/src/modules"
	"group-tracker/src/serveur"
	"io/ioutil"
	"net/http"
)

func main() {
	modules.Test() // metre cette ligne en commentaire si tu veux pas tester la fonction recherche
	urlbase := "https://groupietrackers.herokuapp.com/api/artists"
	// slice de urlbase pour chaque id possible

	response, err := http.Get(urlbase + "/1")
	if err != nil {

		fmt.Println("Erreur lors de la requête HTTP :", err)
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse :", err)
		return
	}

	fmt.Println(string(body))
	serveur.Serveur()
}
