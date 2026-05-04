package service

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

// Structure du service Discord qui contient les informations de la connexion avec le discord
type Discord_service struct {
    Token string
    Discord_API string
    Guild_ID string
}

// Structure de réponse pour obtenir les serveurs
type DiscordGuild struct {
    ID string `json:"id"`
    Name string `json:"name"`
}

// Struucture de réponse pour obtenir les channels
type DiscordChannel struct {
    ID string `json:"id"`
    Type int `json:"type"`
    Name string `json:"name"`
    LastMessageID string `json:"last_message_id,omitempty"`
}

//l'ensemble des channels
type DiscordThreadsResponse struct {
    Threads []DiscordThread `json:"threads"` // La liste des sous-parties
}

//structure de channel
type DiscordThread struct {
    ID string `json:"id"`
    Name string `json:"name"`
    ParentID string `json:"parent_id"` // L'ID du forum parent
    OwnerID string `json:"owner_id"`  // L'auteur
    Type int `json:"type"`
}

//structure de message ( récupération )
type DiscordMessage struct {
    ID string `json:"id"`
    Content string `json:"content"`
    Author DiscordUser `json:"author"`
    Timestamp string `json:"timestamp"`
}

// Structure des utilisateurs
type DiscordUser struct {
    ID string `json:"id"`
    Username string `json:"username"`
}

// Structure pour la création d'un thread avec un message
type CreateMsgThread struct {
    Name string `json:"name"`
    Message MessageContent `json:"message"`
}

// Structure pour le contenu d'un message ( envoyer )
type MessageContent struct {
    Content string `json:"content"`
}

// Fonction qui retourne le premier serveur que le bot est connecté
// Il fait une recherche sur l'ensemble de serveur discord que le bot est connecté et il retourne son id
func (s *Discord_service) GetGuildIdDiscord(ctx context.Context) (string, error) {
    //préparation de la requete
    url := fmt.Sprintf("%s/users/@me/guilds", s.Discord_API)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bot "+ s.Token)
    //on crée un client http et on envoie  une requete
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    //on transforme le text json en tableau
    var guilds []DiscordGuild
    err = json.NewDecoder(resp.Body).Decode(&guilds)
    if err != nil {
        return "", err
    }
    if len(guilds) > 0 {
        fmt.Println("Connexion au serveur : ", guilds[0].Name, guilds[0].ID)
        return guilds[0].ID, nil
    }
    return "", fmt.Errorf("le bot n'a pas trouvé de serveur ")
}

//Fonction qui retourne l'id du salon rechercher402
func (s *Discord_service) GetChannelIdDiscord(ctx context.Context, forum string) (string, error) {
    url := fmt.Sprintf("%s/guilds/%s/channels", s.Discord_API, s.Guild_ID)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bot "+ s.Token)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil{
        return "", err
    }
    defer resp.Body.Close()
    var channels []DiscordChannel
    json.NewDecoder(resp.Body).Decode(&channels)
    for _, ch := range channels {
        if ch.Name == forum && ch.Type == 15 {
            fmt.Printf("Channel forum existant : %s !\n",ch.ID)
            return ch.ID,nil
        }
    }
    return "",fmt.Errorf("le bot n'a pas trouvé de channel ")
}

//Fonction qui crée un channel forum sur le discord
func (s *Discord_service) CreateChannelDiscord(ctx context.Context, name string)(int, error){
    url := fmt.Sprintf("%s/guilds/%s/channels", s.Discord_API, s.Guild_ID)
    newChannel := DiscordChannel{ Name: name, Type: 15,}
    jsonData, err := json.Marshal(newChannel)
    if err != nil{
        return -1, err
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil{
        return -1, err
    }
    req.Header.Set("Authorization", "Bot "+ s.Token)
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Erreur lors de la création :", err)
        return -1,err
    }
    defer resp.Body.Close()
    return 1, nil
}

// fonction qui crée un thread dans un channel forum
func (s *Discord_service) CreateThreadDiscord(ctx context.Context, channelID string, titre string, message string, pseudo string) (string, error) {
    url := fmt.Sprintf("%s/channels/%s/threads", s.Discord_API, channelID)
    // Dans ton service
    message_auteur := fmt.Sprintf("Posté par **%s** :\n\n%s", pseudo, message)
    newMsg := CreateMsgThread{
        Name: titre,
        Message: MessageContent{
            Content: message_auteur,
        },
    }
    jsonData, err := json.Marshal(newMsg)
    if err != nil {
        return "", err
    }
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bot "+s.Token)
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Erreur %s : la création du thread a échoué", resp.StatusCode)
    }
    var result DiscordThread
    err = json.NewDecoder(resp.Body).Decode(&result)
    return result.ID, err
}

// fonction qui poste un message dans un thread
func (s *Discord_service) PostMessageDiscord(ctx context.Context, threadID string, text string, authorPseudo string) error {
    url := fmt.Sprintf("%s/channels/%s/messages", s.Discord_API, threadID)
    message_auteur := fmt.Sprintf("**%s** : %s", authorPseudo, text)
    newMsg := MessageContent{
        Content: message_auteur,
    }
    jsonData, err := json.Marshal(newMsg)
    if err != nil {
        return err
    }
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bot "+s.Token)
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("Erruer %s : l'envoie du message a échoué", resp.StatusCode)
    }
    return nil
}

// fonction qui retourne un tableau de tous les threads
func (s *Discord_service) GetAllThreadDiscord(ctx context.Context, channelID string) ([]DiscordThread,error) {
    url := fmt.Sprintf("%s/guilds/%s/threads/active", s.Discord_API, s.Guild_ID)
    client := &http.Client{}
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil{
        return nil, err
    }
    req.Header.Set("Authorization", "Bot "+ s.Token)
    resp, err := client.Do(req)
    if err != nil{
        return nil, err
    }
    defer resp.Body.Close()
    var allThreads []DiscordThread
    var reponse DiscordThreadsResponse
    json.NewDecoder(resp.Body).Decode(&reponse)
    for _, t := range reponse.Threads {
        if t.ParentID == channelID { 
            fmt.Printf("Thread trouvé : %s %s\n", t.Name, t.ID)
            allThreads = append(allThreads, t)
        }
    }
    return allThreads, nil
}

// fonction qui obtient tous les messages du thread en question
func (s *Discord_service) GetThreadMessageDiscord(ctx context.Context, threadID string) ([]DiscordMessage, error) {
    url := fmt.Sprintf("%s/channels/%s/messages", s.Discord_API, threadID)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bot "+ s.Token)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var messages []DiscordMessage
    err = json.NewDecoder(resp.Body).Decode(&messages)
    return messages, err
}
/*
func main() {
    fmt.Println("Démarrage du serveur !")
    Guild_ID = getGuild_ID()
    if Guild_ID == "" {
        fmt.Println("Erreur : aucun serveur n'a été détécté !")
        return
    }
    forum_name := "forum"
    channelid := getChannelID(forum_name)
    if channelid == "" {
        tmp := createChannel(forum_name);
        if tmp == -1 {
            fmt.Println("Erreur : impossible de créer un channel forum !")
            return
        }
        channelid = getChannelID(forum_name)
    }
    getAllThreads(channelid)
    getThreadMessages("1486667609759023135")
}
*/