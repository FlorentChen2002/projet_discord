 package service

import (
    "fmt"
    "time"
    "context"
    "serveur/repository"
    "go.mongodb.org/mongo-driver/mongo"
)

type Forum_service struct {
    Forum_reposite *repository.Forum_reposite
}

func (s *Forum_service) Create_sujet (ctx context.Context, titre string, description string, user_id string, user_pseudo string, private bool) (bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    sujet, err := s.Forum_reposite.Create_sujet(ctx, titre, description, user_id, user_pseudo, private)
    if err != nil || !sujet {
        fmt.Printf("Erreur : le sujet non trouvé\n")
        return false, err
    }
    return sujet, nil
}

func (s *Forum_service) Create_thread (ctx context.Context,  sujet_id string, content string, user_id string, user_pseudo string, private bool, repond []repository.Repond_db) (bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    thread, err := s.Forum_reposite.Create_thread( ctx, sujet_id, content, user_id, user_pseudo, private, repond)
    if err != nil || !thread {
        fmt.Printf("Erreur : le thread non trouvé\n")
        return false, err
    }
    return thread, nil
}

func (s *Forum_service) GetAllSujet (ctx context.Context)([]repository.Sujet_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    sujet, err := s.Forum_reposite.GetAllSujet(ctx)
    if err != nil{
        if err == mongo.ErrNoDocuments {
            fmt.Printf("Info : *aucun sujet n'a été trouver pour le sujet\n")
            return []repository.Sujet_db{}, nil
        }
        fmt.Printf("Erreur : les sujets non trouvé\n")
        return nil, err
    }
    if sujet == nil {
        sujet = []repository.Sujet_db{} 
    }
    return sujet, nil
}

func (s *Forum_service) GetAllThread (ctx context.Context)([]repository.Thread_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    thread, err := s.Forum_reposite.GetAllThread(ctx)
    if err != nil{
        if err == mongo.ErrNoDocuments {
            fmt.Printf("Info : *aucun thread n'a été trouver\n")
            return []repository.Thread_db{}, nil
        }
        fmt.Printf("Erreur : les threads non trouvé\n")
        return []repository.Thread_db{}, err
    }
    if thread == nil {
        thread = []repository.Thread_db{} 
    }
    return thread, nil
}

func (s *Forum_service) GetThread (ctx context.Context, id string)([]repository.Thread_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    thread, err := s.Forum_reposite.GetThread(ctx, id)
    if err != nil{
        if err == mongo.ErrNoDocuments {
            fmt.Printf("Info : *aucun thread n'a été trouver pour le sujet : %s\n", id)
            return []repository.Thread_db{}, nil
        }
        fmt.Printf("Erreur : le thread non trouvé\n")
        return []repository.Thread_db{}, err
    }
    return thread, nil
}

func (s *Forum_service) DeleteThread (ctx context.Context, id string)(bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    thread, err := s.Forum_reposite.DeleteThread(ctx, id)
    if err != nil{
        fmt.Printf("Erreur : on n'a pas réussi à supprimer le thread\n")
        return false, err
    }
    return thread, nil
}

func (s *Forum_service) DeleteSujet (ctx context.Context, id string)(bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    sujet, err := s.Forum_reposite.DeleteSujet(ctx, id)
    if err != nil{
        fmt.Printf("Erreur : on n'a pas réussi à supprimer le sujet + thread\n")
        return false, err
    }
    return sujet, nil
}