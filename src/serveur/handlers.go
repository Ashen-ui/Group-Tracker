package serveur

import (
	"net/http"
	"sync"
)

var mu sync.Mutex

func indexHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	http.ServeFile(w, r, "./templates/index.html")

}
