package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	tr "Projet/test3"
)

// Équivalent de app.use(cors(...))
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Un simple logger pour voir les requêtes
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requête reçue : %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/user", tr.LoginHandler)
	siteHandler := corsMiddleware(mux)
	siteHandler = loggingMiddleware(siteHandler)
	server := &http.Server{
		Addr: ":8000",
		Handler: siteHandler,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Serveur Go démarré sur http://localhost:8000")
	log.Fatal(server.ListenAndServe())
}