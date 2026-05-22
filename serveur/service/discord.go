package service

import (
    "time"
    "sync"
    "context"
    "serveur/repository"
)

// Conteneur pour avoir accès à la base de données, à l'api discord, au service (forum, discord) et au hub d'événements du forum
// ajout de 2 cannaus pour le systeme des signaux Esterel (pauseSyncro et resumeSyncro) pour éviter les problèmes de synchronisation lors de la suppression d'un sujet
// on utilise un mutex pour éviter les problèmes de synchronisation lors de la création d'un thread et l'ajout de messages dans le forum
type Env_discord struct {
    Discord_repo *repository.Discord_reposite
    Discord_api *Discord_service
    Forum_service *Forum_service
    Forum_Events *ForumEventHub
    mu sync.Map
    pauseSyncro chan string
    resumeSyncro chan string
}

//constructeur pour le conteneur Env_discord
func NewEnvDiscord(discordRepo *repository.Discord_reposite, discordAPI *Discord_service, forumService *Forum_service, forumEvents *ForumEventHub) *Env_discord {
    return &Env_discord{
        Discord_repo: discordRepo,
        Discord_api: discordAPI,
        Forum_service: forumService,
        Forum_Events: forumEvents,
        pauseSyncro: make(chan string, 1),
        resumeSyncro: make(chan string, 1),
    }
}

// fonction qui permet de créer un nouveau verrou par sujet_id
func (e *Env_discord) getMutex(sujet_id string) *sync.RWMutex {
    mu, _ := e.mu.LoadOrStore(sujet_id, &sync.RWMutex{})
    return mu.(*sync.RWMutex)
}

// fonction qui permet de créer un sujet dans Discord et de stocker le dernier message du thread dans la base de données
func (e *Env_discord) SetSujetService(ctx context.Context, titre string, description string, user_id string, user_pseudo string)(string, error) {
    //mutex ?
    thread, err := e.Discord_api.CreateThreadDiscord(ctx, e.Discord_api.Channel_ID, titre, description, user_pseudo)
    if err != nil {
        return "", err
    }
    mu := e.getMutex(thread.ID)
    mu.Lock()
    defer mu.Unlock()
    err = e.Discord_repo.SetLastMessage(ctx, thread.Message.ID,thread.ID)
    if err != nil {
        return "", err
    }
    return thread.ID, nil
}

// fonction qui permet de créer un message dans Discord et de stocker le dernier message du thread dans la base de données
func (e *Env_discord) SetMessagesService(ctx context.Context, sujet_id string, content string, user_pseudo string, repond_id string)(string,error){
    mu := e.getMutex(sujet_id)
    mu.Lock()
    defer mu.Unlock()
    msg, err := e.Discord_api.PostMessageDiscord(ctx, sujet_id, content, user_pseudo, repond_id)
    if err != nil {
        return "", err
    }
    err = e.Discord_repo.SetLastMessage(ctx, msg.ID, sujet_id)
    if err != nil {
        return "", err
    }
    return msg.ID, nil
}

// fonction qui permet de récupérer les nouveaux messages d'un thread Discord et de les ajouter dans le forum
// si le thread n'existe pas dans le forum, on le crée
// si le thread existe déjà, on ajoute les nouveaux messages dans le forum
// on publie un événement dans le hub d'événements du forum à chaque fois qu'on ajoute un message dans le forum
func (e *Env_discord) GetMessagesService(ctx context.Context, sujet_id string, name string)(error) {
    mu := e.getMutex(sujet_id)
    mu.Lock()
    defer mu.Unlock()
    // obtenir le dernier message du thread dans la base de données pour ne récupérer que les nouveaux messages depuis ce message
    last_msg, _ :=  e.Discord_repo.GetLastMessage(ctx, sujet_id)
    liste_msg, err := e.Discord_api.GetNewMessages(ctx, sujet_id, last_msg.Msg_id)
    if err != nil{
        return err
    }
    if len(liste_msg) == 0 { // pas de nouveau message, donc on ne fait rien
        return nil
    }
    _, err = e.Forum_service.GetSujet(ctx, sujet_id) // vérifier si le thread existe déjà dans le forum
    tmp := false
    if err != nil {
        tmp = true
    }
    for i := len(liste_msg) - 1; i >= 0; i-- {
        if tmp { // si elle n'existe pas alors on la crée
            e.Forum_service.Create_sujet(ctx, sujet_id, name, liste_msg[i].Content, liste_msg[i].Auteur.ID, liste_msg[i].Auteur.Username, false, false)
            tmp = false
            if e.Forum_Events != nil {
                e.Forum_Events.Publish(ForumEvent{Type: "sujet", Action: "updated", SujetID: sujet_id})
            }
        } else { // sinon on ajoute les nouveaux messages dans le forum
            var repond []repository.Repond_db
            if liste_msg[i].MessageReference != nil { // gestion du cas si on répond à un message sur Discord
                repond_tmp, err := e.Forum_service.GetMessage(ctx, liste_msg[i].MessageReference.MessageID)
                if err != nil {
                    return err
                }
                repond = []repository.Repond_db{{ // reconstitution de la structure Repond_db via une obtention du message référencé dans le forum pour pouvoir le concorder avec le msg sur le site
                        Id: liste_msg[i].MessageReference.MessageID,
                        Content: repond_tmp.Content,
                        User_pseudo: repond_tmp.User_pseudo,
                        Date: repond_tmp.Date,
                    },
                }
            }
            if sujet_id == liste_msg[i].ID {
                continue
            }
            e.Forum_service.Create_messages(ctx, sujet_id, liste_msg[i].ID, liste_msg[i].Content, liste_msg[i].Auteur.ID, liste_msg[i].Auteur.Username, false, repond, false)
        }
    }
    // puis sauvegarde du dernier message lu
    err = e.Discord_repo.SetLastMessage(ctx, liste_msg[0].ID, sujet_id)
    if err != nil {
        return err
    }
    if e.Forum_Events != nil { // signaler au abonnée pour dire qu'il y a un nouveau message dans le thread
        e.Forum_Events.Publish(ForumEvent{Type: "message", Action: "created", ThreadID: sujet_id})
    }
    return nil
}

// fonction qui permet de supprimer un message dans Discord et si c'est le dernier message du thread alors on le initialise vide
func (e *Env_discord) DeleteMessageService(ctx context.Context, sujet_id string, msg_id string){
    mu := e.getMutex(sujet_id)
    mu.Lock()
    defer mu.Unlock()
    last_msg, _ :=  e.Discord_repo.GetLastMessage(ctx, sujet_id)
    if msg_id == last_msg.Msg_id {
        e.Discord_repo.SetLastMessage(ctx, "", sujet_id)
    }
    e.Discord_api.DeleteMessageDiscord(ctx, sujet_id, msg_id)
}

// fonction qui permet de supprimer le thread et tous les messages associés dans Discord
// et de supprimer la collection concernant
func (e *Env_discord) DeleteSujetService(ctx context.Context, sujet_id string){
    if e.pauseSyncro != nil {
        e.pauseSyncro <- sujet_id
    }
    mu := e.getMutex(sujet_id)
    mu.Lock()
    defer mu.Unlock()
    e.Discord_api.DeleteThreadDiscord(ctx, sujet_id)
    e.Discord_repo.DeleteLastMessage(ctx, sujet_id)
    if e.resumeSyncro != nil {
        e.resumeSyncro <- sujet_id
    }
}

// Lancement de la synchronisation entre Discord et le forum
// on récupère tous les threads Discord
// on utilise un ticker pour limiter la fréquence de récupération des messages et pour éviter la saturation du serveur Discord
func (e *Env_discord) StartSyncro(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    var threads []DiscordThread
    index := 0
    // on utilise une map pour stocker les sujets en pause de synchronisation pour éviter les problèmes de synchronisation lors de la suppression d'un sujet
    paused := make(map[string]bool)
    for {
        select {
        case <-ctx.Done():
            return
        case sujet_id := <-e.pauseSyncro: // si on reçoit un signal de pause pour un sujet, on le met en pause dans la map
            paused[sujet_id] = true 
            index = 0  //on force le refresh la liste des threads pour éviter de continuer sur une liste de threads obsolète
            threads = nil
        case sujet_id := <-e.resumeSyncro: // si on reçoit un signal de reprise pour un sujet, on le supprime de la map pour reprendre la synchronisation
            delete(paused, sujet_id)
        case <-ticker.C: // à chaque tick, on récupère les messages de 10 threads Discord pour les ajouter dans le forum
            if index == 0 {
                threads, _ = e.Discord_api.GetAllThreadDiscord(ctx, e.Discord_api.Channel_ID)
            }
            if len(threads) == 0 {
                continue
            }
            end := index + 10
            if end > len(threads) {
                end = len(threads)
            }
            for _, t := range threads[index:end] {
                id := t.ID
                name := t.Name
                if paused[id] { // si on rencontre le sujet en question alors on le saute 
                    continue
                }
                go func(threadID, threadName string) {
                    e.GetMessagesService(ctx, threadID, threadName)
                }(id, name)
            }
            index += 10
            if index >= len(threads) {
                index = 0
            }
        }
    }
}

// 2 Problèmes :
// Plus il y a d'utilisateurs, plus il y a de threads, plus le temps de synchronisation est long
// Latence entre la création d'un thread/message et son apparition sur le forum

// Problème de fonctionnalité : 
// on peut supprimer un message sur discord mais elle ne sera pas supprimé sur le forum ( le cas inverse marche)
// quand on répond à un message sur le forum, la réponse n'apparaît pas sur discord et idem pour la création d'un message sur le forum