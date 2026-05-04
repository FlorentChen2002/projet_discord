package handler

import (
    "net/http"
    "encoding/json"
    "fmt"
    "serveur/service"
    "serveur/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Env_forum struct {
    Forum_reposite *service.Forum_service
    Forum_Events *service.ForumEventHub
}

type Sujet_requete struct {
    Id string`json:"id"`
    Titre string `json:"titre"`
    Description string `json:"description"`
    User_id string `json:"userid"`
    User_pseudo string `json:"userpseudo"`
    Private bool `json:"prive"`
}
type Thread_requete struct {
    Id string`json:"id"`
    Sujet_id string `json:"sujetid"`
    Content string `json:"content"`
    Date string `json:"date"`
    User_id string `json:"userid"`
    User_pseudo string `json:"userpseudo"`
    Private bool `json:"prive"`
    Repond []Repond_requete`json:"repond"`
}

type Repond_requete struct {
    Content string `json:"content"`
    User_id string `json:"userid"`
    User_pseudo string `json:"userpseudo"`
    Date string `json:"date"`
}

func (e *Env_forum) ForumEventsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("X-Accel-Buffering", "no")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming unsupported", http.StatusInternalServerError)
        return
    }

    ch := e.Forum_Events.Subscribe()
    defer e.Forum_Events.Unsubscribe(ch)

    fmt.Fprint(w, ": connected\n\n")
    flusher.Flush()

    for {
        select {
        case <-r.Context().Done():
            return
        case event, ok := <-ch:
            if !ok {
                return
            }
            payload, err := json.Marshal(event)
            if err != nil {
                continue
            }
            fmt.Fprintf(w, "event: forum\n")
            fmt.Fprintf(w, "data: %s\n\n", payload)
            flusher.Flush()
        }
    }
}

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

func (e *Env_forum) CreateSujetHandler (w http.ResponseWriter, r *http.Request) {
    var req Sujet_requete
    w.Header().Set("Content-Type", "application/json")
    err:= json.NewDecoder(r.Body).Decode(&req)
    if err!=nil || req.Titre == "" || req.User_id == "" || req.User_pseudo == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Les champs sont vides",})
        return
    }
    etat, err := e.Forum_reposite.Create_sujet(r.Context(), req.Titre, req.Description, req.User_id, req.User_pseudo, req.Private)
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

func (e *Env_forum) GetAllThreadHandler (w http.ResponseWriter, r *http.Request) {
    threads, err := e.Forum_reposite.GetAllThread(r.Context())
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Récupération des threads échouée",})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(threads)
}

func (e *Env_forum) GetThreadHandler (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    sujet_id := r.URL.Query().Get("sujetid")
    if sujet_id == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Les champs sont vides",})
        return
    }
    threads, err := e.Forum_reposite.GetThread(r.Context(), sujet_id)
    if err!=nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Récupération des threads échouée",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(threads)
}

func (e *Env_forum) CreateThreadHandler (w http.ResponseWriter, r *http.Request) {
    var req Thread_requete
    w.Header().Set("Content-Type", "application/json")
    err:= json.NewDecoder(r.Body).Decode(&req)
    if err!=nil || req.Sujet_id == "" || req.User_id == "" || req.User_pseudo == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Les champs sont vides",})
        return
    }
    var reponsDB []repository.Repond_db
    for _, reponse := range req.Repond {
        if reponse.User_id == ""{
            continue
        }
        objUserID, err := primitive.ObjectIDFromHex(reponse.User_id)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"Erreur": "Transformation du id en objet id échoué",})
            return
        }
        nouvelle_reponse := repository.Repond_db{ Content: reponse.Content, User_id: objUserID, User_pseudo: reponse.User_pseudo, Date: reponse.Date,}
        reponsDB = append(reponsDB, nouvelle_reponse)
    }
    etat, err := e.Forum_reposite.Create_thread(r.Context(), req.Sujet_id, req.Content, req.User_id, req.User_pseudo, req.Private, reponsDB)
    if err != nil || !etat{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Création du thread a échoué",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Création du thread a réussie"})
    if e.Forum_Events != nil {
        e.Forum_Events.Publish(service.ForumEvent{Type: "thread", Action: "created", SujetID: req.Sujet_id})
    }
}

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

func (e *Env_forum) DeleteThreadHandler (w http.ResponseWriter, r *http.Request) {
    var req Thread_requete
    w.Header().Set("Content-Type", "application/json")
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.Id == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Lecture du body a échoué",})
        return
    }
    etat, err := e.Forum_reposite.DeleteThread(r.Context(), req.Id)
    if err != nil || !etat{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"Erreur": "Suppresion du thread a échoué",})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Suppresion du thread a réussie"})
    if e.Forum_Events != nil {
        e.Forum_Events.Publish(service.ForumEvent{Type: "thread", Action: "deleted"})
    }

}