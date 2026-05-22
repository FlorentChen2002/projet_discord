package middleware

import (
    "net/http"
    "log"
)

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if origin == "http://localhost:5173" || origin == "https://projet-discord.onrender.com" || origin == "https://projetdiscord-production.up.railway.app" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Un simple logger pour voir les requêtes
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Requête reçue : %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func Direction_middleware( h http.Handler, auto_direction ...func(http.Handler) http.Handler) http.Handler {
    for i := len(auto_direction) - 1; i >= 0; i-- {
        h = auto_direction[i](h)
    }
    return h
}