 package service

import (
    "fmt"
    "time"
    "context"
    "serveur/repository"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Conteneur pour avoir accès à la base de données et à l'api discord
type Forum_service struct {
    Forum_reposite *repository.Forum_reposite
    Env_discord *Env_discord
}

// fonction qui permet de créer un sujet dans le forum
// si api_discord est à vrai alors on crée un thread dans Discord et on utilise son id comme id du sujet
// sinon on génère un id aléatoire pour le sujet
func (s *Forum_service) Create_sujet (ctx context.Context, id string ,titre string, description string, user_id string, user_pseudo string, private bool, api_discord bool) (bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    if api_discord { // Pour le cas si on crée un sujet via le site
        tmp, err := s.Env_discord.SetSujetService(ctx, titre, description, user_id, user_pseudo)
        if err != nil {
            return false, err
        }
        id = tmp
    }
    // si l'id est vide, on génère un id aléatoire pour le sujet
    if id == "" {
        id = primitive.NewObjectID().Hex()
    }
    sujet, err := s.Forum_reposite.Create_sujet(ctx, id, titre, description, user_id, user_pseudo, private)
    if mongo.IsDuplicateKeyError(err) {
        return true, nil
    }
    if err != nil || !sujet {
        fmt.Printf("Erreur : Création du sujet échouée\n")
        return false, err
    }
    return sujet, nil
}

// fonction qui permet de créer un message dans le forum
// si api_discord est à vrai alors on crée un message dans Discord et on utilise son id comme id du message
// sinon on génère un id aléatoire pour le message
func (s *Forum_service) Create_messages (ctx context.Context, sujet_id string, msg_id string, content string, user_id string, user_pseudo string, private bool, repond []repository.Repond_db, api_discord bool) (bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    if api_discord { // Pour le cas si on crée un message via le site
        repond_tmp :=""
        if len(repond) > 0 { // si on répond à un message
            repond_tmp = repond[0].Id
        }
        tmp, err := s.Env_discord.SetMessagesService(ctx, sujet_id, content, user_pseudo, repond_tmp)
        if err != nil {
            return false, err
        }
        msg_id = tmp
    }
    if msg_id == "" {
        msg_id = primitive.NewObjectID().Hex()
    }
    msg, err := s.Forum_reposite.Create_message( ctx, msg_id, sujet_id, content, user_id, user_pseudo, private, repond)
    if mongo.IsDuplicateKeyError(err) {
        return true, nil
    }
    if err != nil || !msg {
        fmt.Printf("Erreur : Création du message échouée\n")
        return false, err
    }
    return msg, nil
}

// fonction qui permet d'obtenir un sujet de la base de données via son id
func (s *Forum_service) GetSujet (ctx context.Context, sujet_id string)(repository.Sujet_db,error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    sujet, err := s.Forum_reposite.GetSujet(ctx, sujet_id)
    if err != nil{
        return repository.Sujet_db{}, err
    }
    return sujet, nil
}

// fonction qui permet d'obtenir tous les sujets de la base de données
// on retourne une liste de sujet et une erreur si il y en a une
func (s *Forum_service) GetAllSujet (ctx context.Context)([]repository.Sujet_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    sujet, err := s.Forum_reposite.GetAllSujet(ctx)
    if err != nil{
        fmt.Printf("Erreur : les sujets non trouvé\n")
        return nil, err
    }
    if sujet == nil {
        sujet = []repository.Sujet_db{} 
    }
    return sujet, nil
}

// fonction qui permet d'obtenir tous les messages de la base de données
// on retourne une liste de message et une erreur si il y en a une
func (s *Forum_service) GetAllMessages (ctx context.Context)([]repository.Message_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    messages, err := s.Forum_reposite.GetAllMessage(ctx)
    if err != nil{
        fmt.Printf("Erreur : les messages non trouvé\n")
        return []repository.Message_db{}, err
    }
    if messages == nil {
        messages = []repository.Message_db{} 
    }
    return messages, nil
}

// fonction qui permet d'obtenir un message de la base de données via son id
func (s *Forum_service) GetMessage (ctx context.Context, id string)(repository.Message_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    messages, err := s.Forum_reposite.GetMessage(ctx, id)
    if err != nil{
        fmt.Printf("Erreur : les messages non trouvé %s\n", id)
        return repository.Message_db{}, err
    }
    return messages, nil
}

// fonction qui permet d'obtenir tous les messages d'un sujet de la base de données via l'id du sujet
func (s *Forum_service) GetSujetMessage (ctx context.Context, id string)([]repository.Message_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    messages, err := s.Forum_reposite.GetSujetMessage(ctx, id)
    if err != nil{
        fmt.Printf("Erreur : les messages du sujet non trouvé\n")
        return []repository.Message_db{}, err
    }
    return messages, nil
}

// fonction qui permet de supprimer un message via son id
func (s *Forum_service) DeleteMessage (ctx context.Context, id string)(bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    msg, err := s.Forum_reposite.GetMessage(ctx, id)
    if err != nil {
        fmt.Printf("Erreur : message non trouvé\n")
        return false, err
    }
    s.Env_discord.DeleteMessageService(ctx, msg.Sujet_id, id)
    message, err := s.Forum_reposite.DeleteMessage(ctx, id)
    if err != nil{
        fmt.Printf("Erreur : on n'a pas réussi à supprimer le message\n")
        return false, err
    }
    return message, nil
}

// fonction qui permet de supprimer un sujet via son id
func (s *Forum_service) DeleteSujet (ctx context.Context, id string)(bool, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    s.Env_discord.DeleteSujetService(ctx, id)
    sujet, err := s.Forum_reposite.DeleteSujet(ctx, id)
    if err != nil{
        fmt.Printf("Erreur : on n'a pas réussi à supprimer le sujet + thread\n")
        return false, err
    }
    return sujet, nil
}