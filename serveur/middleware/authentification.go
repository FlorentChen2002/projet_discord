package middleware

import (
    "net/http"
    "fmt"
    "context"
    "serveur/handler"
)

func Auth_middleware(next http.Handler) http.Handler {
    return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request){
        if (r.Method == "POST") || (r.Method == "PUT"){
            if r.ContentLength == 0  {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(w,`{"Erreur":"La requête est vide"}`)
                return
            }
        }
        next.ServeHTTP(w,r)
    })
}

func Forum_middleware(next http.Handler) http.Handler {
    return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request){
        if r.Method == "OPTIONS" {
            next.ServeHTTP(w, r)
            return
        }
        cookie, err := r.Cookie("session_token")
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        user_id, ok := handler.GetSession(cookie.Value)
        if !ok {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), "user_id", user_id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}