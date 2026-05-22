package handler

import (
    "sync"
    "net/http"
    "encoding/json"
    "serveur/service"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Structure représentant une requête de connexion, obtenue à partir du corps de la requête HTTP
type Login_requete struct {
    Pseudo string `json:"pseudo"`
    Mdp string `json:"mdp"`
}

// Structure représentant une requête d'inscription, obtenue à partir du corps de la requête HTTP
type Inscription_requete struct {
    Pseudo string `json:"pseudo"`
    Mdp string `json:"mdp"`
    Rang string `json:"rang"`
}

// Conteneur pour avoir accès à la base de données
type Env struct {
    User_reposite *service.User_service
}

// variables globales pour la gestion des sessions utilisateur
var (
    Sessions = make(map[string]string)
    Sessions_mux sync.RWMutex
)
// fonctions ajouter un nouveau utilisateur dans la session
func SetSession(session_id, user_id string) {
    Sessions_mux.Lock()
    Sessions[session_id] = user_id
    defer Sessions_mux.Unlock()
}

// fonction qui permet d'obtenir l'id de l'utilisateur à partir de son session_id
func GetSession(session_id string) (string, bool) {
    Sessions_mux.RLock()
    user_id, tmp := Sessions[session_id]
    defer Sessions_mux.RUnlock()
    return user_id, tmp
}

// fonction qui permet de supprimer une session à partir de son session_id ( utilisé pour la déconnexion )
func DeleteSession(session_id string) {
    Sessions_mux.Lock()
    delete(Sessions, session_id)
    defer Sessions_mux.Unlock()
}

// fonction qui permet de se connecter à un compte utilisateur
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
    http.SetCookie(w, &http.Cookie{ Name: "session_token", Value: session_id, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: true, MaxAge: 3600,})
    w.WriteHeader(http.StatusOK)
    retour := map[string]interface{}{"id": user_id, "message": "Connexion réussie", "status": 200,}
    json.NewEncoder(w).Encode(retour)
}

// fonction qui permet de s'inscrire à un compte utilisateur
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

// fonction qui permet d'obtenir les données d'un utilisateur via son id
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

// fonction qui permet de faire une deconnexion d'un utilisateur en supprimant l'utilisateur de la session
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

// fonction qui permet d'obtenir tous les utilsateurs
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

// fonction qui permet d'obtenir les données d'un utilisateur via son id
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