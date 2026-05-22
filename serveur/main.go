package main

import (
    "fmt"
    "os"
    "serveur/handler"
    "serveur/middleware"
    "serveur/service"
    "serveur/repository"
    "net/http"
    "time"
    "log"
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// variable globale pour le token du bot discord et l'url de l'api discord
const Discord_API = "https://discord.com/api/v10"
var Token = os.Getenv("DISCORD_TOKEN")

// fonction pour configurer le bot discord et récupérer les ids du serveur et du channel
func setupDiscordBot(ctx context.Context, s *service.Discord_service) error {
    guild_id, err := s.GetGuildIdDiscord(ctx)
    if err != nil {
        return err
    }
    s.Guild_ID = guild_id
    forumName := "forum"
    channel_id, err := s.GetChannelIdDiscord(ctx, forumName)
    if err != nil {
        fmt.Println("Création du channel")
        _, err := s.CreateChannelDiscord(ctx, forumName)
        if err != nil {
            return err
        }
        channel_id, _ = s.GetChannelIdDiscord(ctx, forumName)
    }
    s.Channel_ID = channel_id
    return nil
}

// fonction pour appliquer les middlewares à un handler 
func middleware_direction(h http.HandlerFunc, fonction_middleware func(http.Handler) http.Handler) http.Handler {
    return middleware.Direction_middleware(h, middleware.CorsMiddleware, middleware.LoggingMiddleware, fonction_middleware)
}

// fonction principale pour démarrer le serveur
func main() {
    // connexion à la base de données MongoDB
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://florentchen94_db_user:fPvofzjDMKKRVeLK@cluster0.qytxmrg.mongodb.net/?appName=Cluster0"))
    if err != nil {
        log.Fatal(err)
        return
    }
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal("Erreur : le serveur mongodb n'est pas lancé", err)
        return
    }
    // initialisation du conteneur pour des repositories, services et handlers pour les utilisateurs, le forum et le discord
    db := client.Database("projet_forum")
    // Conteneur pour les utilisateurs
    user_reposite := &repository.User_reposite{ Collection : db.Collection("user"),}
    user_service := &service.User_service{ User_reposite: user_reposite,}
    env := &handler.Env{ User_reposite: user_service,}
    // Conteneur pour le discord
    discord := &service.Discord_service{Token: Token,Discord_API: Discord_API,}
    discord_reposite := &repository.Discord_reposite{ Message_Collection: db.Collection("DiscordMessages"),}
    err = setupDiscordBot(ctx, discord)
    if err != nil {
        log.Fatalf("Erreur : échec de la connexion du bot %v", err)
    }
    // Conteneur pour le forum
    forum_reposite := &repository.Forum_reposite{ Sujet_Collection : db.Collection("ForumDB"), Thread_Collection : db.Collection("ThreadDB"),}
    forum_service := &service.Forum_service{ Forum_reposite: forum_reposite,}
    forum_events := service.NewForumEventHub()
    env_discord := &service.Env_discord{ Discord_repo: discord_reposite, Discord_api: discord, Forum_service: forum_service, Forum_Events: forum_events,}
    forum_service.Env_discord = env_discord
    env_forum := &handler.Env_forum{ Forum_reposite: forum_service, Forum_Events: forum_events,}

    // Lancement de la synchronisation entre le forum et le discord dans une goroutine
    go env_discord.StartSyncro(ctx, 5*time.Second)
    mux := http.NewServeMux()

    // Routes pour les utilisateurs
    mux.Handle("/api/user/connexion", middleware_direction(http.HandlerFunc(env.LoginHandler), middleware.Auth_middleware)) //a mettre une fonction
    mux.Handle("/api/user/inscription", middleware_direction(http.HandlerFunc(env.InscriptionHandler), middleware.Auth_middleware))
    mux.Handle("/api/user/me", middleware_direction(http.HandlerFunc(env.GetMeHandler), middleware.Auth_middleware))
    mux.Handle("/api/user/{id}", middleware_direction(http.HandlerFunc(env.GetUserByIdHandler), middleware.Auth_middleware))
    mux.Handle("/api/user/deconnexion", middleware_direction(http.HandlerFunc(env.LogoutHandler), middleware.Auth_middleware))
    mux.Handle("/api/user/alluser", middleware_direction(http.HandlerFunc(env.GetAllUserHandler), middleware.Auth_middleware))
    
    // Routes pour le forum
    mux.Handle("/api/forum/sujet", middleware_direction(http.HandlerFunc(env_forum.GetAllSujetHandler), middleware.Forum_middleware))
    mux.Handle("/api/forum/postforum", middleware_direction(http.HandlerFunc(env_forum.CreateSujetHandler), middleware.Forum_middleware))
    mux.Handle("/api/forum/getallthread", middleware_direction(http.HandlerFunc(env_forum.GetAllMessageHandler), middleware.Forum_middleware))
    mux.Handle("/api/forum/getthread", middleware_direction(http.HandlerFunc(env_forum.GetSujetMessageHandler), middleware.Forum_middleware))
    mux.Handle("/api/forum/postthread", middleware_direction(http.HandlerFunc(env_forum.CreateMessageHandler), middleware.Forum_middleware))
    mux.Handle("/api/forum/delete/sujet", middleware_direction(http.HandlerFunc(env_forum.DeleteSujetHandler), middleware.Forum_middleware))
    mux.Handle("/api/forum/delete/thread", middleware_direction(http.HandlerFunc(env_forum.DeleteMessageHandler), middleware.Forum_middleware))
    mux.Handle("/api/forum/events", middleware_direction(http.HandlerFunc(env_forum.ForumEventsHandler), middleware.Forum_middleware))
    // Start du serveur
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    server := &http.Server{
        Addr: ":" + port,
        Handler: mux,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 0,
        IdleTimeout: 60 * time.Second,
    }
    fmt.Println("Serveur Go démarré sur le port :" + port)
    log.Fatal(server.ListenAndServe())
}