package handler

import (
    "fmt"
    "net/http"
    "encoding/json"
    "serveur/service"
    "serveur/repository"
)

// Conteneur pour avoir accès à la base de données et au hub d'événements pour le forum
type Env_forum struct {
    Forum_reposite *service.Forum_service
    Forum_Events *service.ForumEventHub
}

// Structure représentant une requête pour créer un sujet ou un message dans le forum
type Sujet_requete struct {
    Id string`json:"id"`
    Titre string `json:"titre"`
    Description string `json:"description"`
    User_id string `json:"userid"`
    User_pseudo string `json:"userpseudo"`
    Private bool `json:"prive"`
}

// Structure représentant une requête pour créer un message dans le forum
type Message_requete struct {
    Id string`json:"id"`
    Sujet_id string `json:"sujetid"`
    Content string `json:"content"`
    Date string `json:"date"`
    User_id string `json:"userid"`
    User_pseudo string `json:"userpseudo"`
    Private bool `json:"prive"`
    Repond []Repond_requete`json:"repond"`
}

// Structure représentant une requête pour créer une réponse à un message dans le forum
type Repond_requete struct {
    Id string `json:"id"`
    Content string `json:"content"`
    User_id string `json:"userid"`
    User_pseudo string `json:"userpseudo"`
    Date string `json:"date"`
}

// Handler pour gérer les événements du forum et les requêtes liées aux sujets et messages du forum
// cela permet au site de se mettre à jour en temps réel lorsqu'un sujet ou un message est créé ou supprimé dans le forum
func (e *Env_forum) ForumEventsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("X-Accel-Buffering", "no")
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Erreur : Streaming non supporté", http.StatusInternalServerError)
        return
    }
    ch := e.Forum_Events.Subscribe()
    defer e.Forum_Events.Unsubscribe(ch)
    flusher.Flush()
    for {
        select {
        case <-r.Context().Done():
            return
        case event, ok := <-ch:
            if !ok {
                return
            }
            data, err := json.Marshal(event)
            if err != nil {
                continue
            }
            fmt.Fprintf(w, "event: forum\ndata: %s\n\n", data)
            flusher.Flush()
        }
    }
}

// Handler pour gérer les requêtes liées aux sujets du forum
// Pour obtenir tous les sujets
func (e *Env_forum) GetAllSujetHandler (w http.ResponseWriter, r *http.Request) {
    sujets, err := e.Forum_reposite.GetAllSujet(r.Context())
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Récupération des sujets échouée",})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(sujets)
}

// Handler pour gérer les requêtes liées aux créations de sujets dans le forum et dans le discord
func (e *Env_forum) CreateSujetHandler (w http.ResponseWriter, r *http.Request) {
    var req Sujet_requete
    w.Header().Set("Content-Type", "application/json")
    err:= json.NewDecoder(r.Body).Decode(&req)
    if err!=nil || req.Titre == "" || req.User_id == "" || req.User_pseudo == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Les champs sont vides",})
        return
    }
    etat, err := e.Forum_reposite.Create_sujet(r.Context(), "", req.Titre, req.Description, req.User_id, req.User_pseudo, req.Private, true)
    if err != nil || !etat{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Création du Sujet a échoué",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Création du sujet réussie"})
    if e.Forum_Events != nil {
        e.Forum_Events.Publish(service.ForumEvent{Type: "sujet", Action: "created"})
    }
}

// Handler pour gérer les requêtes liées aux messages du forum
// Pour obtenir tous les messages
func (e *Env_forum) GetAllMessageHandler (w http.ResponseWriter, r *http.Request) {
    messages, err := e.Forum_reposite.GetAllMessages(r.Context())
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Récupération des messages échouée",})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(messages)
}

// Handler pour gérer les requêtes liées à l'obtention des messages d'un sujet du forum
// Pour obtenir tous les messages d'un sujet
func (e *Env_forum) GetSujetMessageHandler (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    sujet_id := r.URL.Query().Get("sujetid")
    if sujet_id == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Les champs sont vides",})
        return
    }
    messages, err := e.Forum_reposite.GetSujetMessage(r.Context(), sujet_id)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Récupération des messages échouée",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(messages)
}

// Handler pour gérer les requêtes liées à la création de messages dans le forum et discord
func (e *Env_forum) CreateMessageHandler (w http.ResponseWriter, r *http.Request) {
    var req Message_requete
    w.Header().Set("Content-Type", "application/json")
    err:= json.NewDecoder(r.Body).Decode(&req)
    if err!=nil || req.Sujet_id == "" || req.User_id == "" || req.User_pseudo == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Les champs sont vides",})
        return
    }
    var reponsDB []repository.Repond_db
    for _, reponse := range req.Repond { // on reconvertit les réponses de la requête en réponses pour la base de données
        if reponse.User_id == ""{
            continue
        }
        nouvelle_reponse := repository.Repond_db{ Id: reponse.Id, Content: reponse.Content, User_id: reponse.User_id, User_pseudo: reponse.User_pseudo, Date: reponse.Date,}
        reponsDB = append(reponsDB, nouvelle_reponse)
    }
    etat, err := e.Forum_reposite.Create_messages(r.Context(), req.Sujet_id, "", req.Content, req.User_id, req.User_pseudo, req.Private, reponsDB, true)
    if err != nil || !etat{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Création du message a échoué",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Création du message a réussie"})
    if e.Forum_Events != nil {
        e.Forum_Events.Publish(service.ForumEvent{Type: "message", Action: "created", SujetID: req.Sujet_id})
    }
}

// Handler pour gérer les requêtes liées à la suppression de sujets et de messages du forum et du discord
func (e *Env_forum) DeleteSujetHandler (w http.ResponseWriter, r *http.Request) {
    var req Sujet_requete
    w.Header().Set("Content-Type", "application/json")
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Id == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Lecture du body a échoué",})
        return
    }
    etat, err := e.Forum_reposite.DeleteSujet(r.Context(), req.Id)
    if err != nil || !etat{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Suppresion du sujet a échoué",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Suppresion du sujet a réussie"})
    if e.Forum_Events != nil {
        e.Forum_Events.Publish(service.ForumEvent{Type: "sujet", Action: "deleted", SujetID: req.Id})
    }
}

// Handler pour gérer les requêtes liées à la suppression de sujets et de messages du forum et du discord
func (e *Env_forum) DeleteMessageHandler (w http.ResponseWriter, r *http.Request) {
    var req Message_requete
    w.Header().Set("Content-Type", "application/json")
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Id == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Lecture du body a échoué",})
        return
    }
    etat, err := e.Forum_reposite.DeleteMessage(r.Context(), req.Id)
    if err != nil || !etat{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Suppresion du message a échoué",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Suppresion du message a réussie"})
    if e.Forum_Events != nil {
        e.Forum_Events.Publish(service.ForumEvent{Type: "message", Action: "deleted"})
    }

}