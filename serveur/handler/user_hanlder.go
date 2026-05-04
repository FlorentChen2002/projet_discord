package handler

import (
    "sync"
    "net/http"
    "encoding/json"
    "serveur/service"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Login_requete struct {
    Pseudo string `json:"pseudo"`
    Mdp string `json:"mdp"`
}

type Inscription_requete struct {
    Pseudo string `json:"pseudo"`
    Mdp string `json:"mdp"`
    Rang string `json:"rang"`
}

type Env struct {
    User_reposite *service.User_service
}

var (
    Sessions = make(map[string]string)
    Sessions_mux sync.RWMutex
)

func SetSession(session_id, user_id string) {
    Sessions_mux.Lock()
    Sessions[session_id] = user_id
    defer Sessions_mux.Unlock()
}

func GetSession(session_id string) (string, bool) {
    Sessions_mux.RLock()
    user_id, tmp := Sessions[session_id]
    defer Sessions_mux.RUnlock()
    return user_id, tmp
}

func DeleteSession(session_id string) {
    Sessions_mux.Lock()
    delete(Sessions, session_id)
    defer Sessions_mux.Unlock()
}

func (e *Env) LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req Login_requete
    w.Header().Set("Content-Type", "application/json")
    err:= json.NewDecoder(r.Body).Decode(&req)
    //fmt.Printf("pseudo='%s',mdp='%s'\n",req.Pseudo, req.Mdp)
    if err!=nil || req.Pseudo == "" || req.Mdp == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Le pseudo ou le mdp est vide",})
        return
    }
    user_id, err := e.User_reposite.Connexion(r.Context(), req.Pseudo, req.Mdp)
    if err != nil || user_id =="" {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Pseudo incorrecte",})
        return
    }
    session_id := primitive.NewObjectID().Hex()
    SetSession(session_id, user_id)
    http.SetCookie(w, &http.Cookie{ Name: "session_token", Value: session_id, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: 3600,})
    w.WriteHeader(http.StatusOK)
    retour := map[string]interface{}{"id": user_id, "message": "Connexion réussie", "status": 200,}
    json.NewEncoder(w).Encode(retour)
}

func (e *Env) InscriptionHandler(w http.ResponseWriter, r *http.Request) {
    var req Inscription_requete
    w.Header().Set("Content-Type", "application/json")
    err:= json.NewDecoder(r.Body).Decode(&req)
    if err!=nil || req.Pseudo == "" || req.Mdp == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Le pseudo ou le mdp est vide",})
        return
    }
    err = e.User_reposite.Inscription(r.Context(), req.Pseudo, req.Mdp, req.Rang)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Le mot de passe est incorrecte",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Inscription réussie"})
}

func (e *Env) GetMeHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    w.Header().Set("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    user_id, tmp := GetSession(cookie.Value)
    if !tmp {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    user, err := e.User_reposite.GetId(r.Context(), user_id)
    if err != nil {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Id inconnu"})
        return
    }
    json.NewEncoder(w).Encode(user)
}

func (e *Env) LogoutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    cookie, err := r.Cookie("session_token")
    if err != nil {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Déjà déconnecté"})
        return
    }
    DeleteSession(cookie.Value)
    http.SetCookie(w, &http.Cookie{ Name: "session_token", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode,})
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Déconnexion réussie"})
}

func (e *Env) GetAllUserHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    user, err := e.User_reposite.GetAllUsers(r.Context())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "récupération de tous les utilsateurs échouer"})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}

func (e *Env) GetUserByIdHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    userID := r.PathValue("id")
    user, err := e.User_reposite.GetId(r.Context(), userID)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Id inconnu"})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}