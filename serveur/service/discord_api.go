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
    Channel_ID string
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
}

//l'ensemble des threads d'un channel
type DiscordThreadsResponse struct {
    Threads []DiscordThread `json:"threads"` // La liste des sous-parties
}

//structure de thread
type DiscordThread struct {
    ID string `json:"id"`
    Name string `json:"name"`
    ParentID string `json:"parent_id"`
    OwnerID string `json:"owner_id"`
    Type int `json:"type"`
    Message DiscordMessage `json:"message"`
}

//structure de message ( récupération )
type DiscordMessage struct {
    ID string `json:"id"`
    Content string `json:"content"`
    Auteur DiscordUser `json:"author"`
    Timestamp string `json:"timestamp"`
    ReferenceMessage *DiscordMessageReference `json:"referenced_message,omitempty"`
    MessageReference *DiscordMessageReference `json:"message_reference,omitempty"`
}

// Structure des utilisateurs
type DiscordUser struct {
    ID string `json:"id"`
    Username string `json:"username"`
}

// Structure de réponse d'un message de discord
type DiscordMessageReference struct {
    MessageID string `json:"message_id"`
    Content string `json:"content"`
    Auteur DiscordUser `json:"author"`
    Date string `json:"timestamp"`
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

//Fonction qui retourne l'id du salon rechercher
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

// fonction qui crée un thread dans un channel du forum
func (s *Discord_service) CreateThreadDiscord(ctx context.Context, channelID string, titre string, message string, pseudo string) (DiscordThread, error) {
    url := fmt.Sprintf("%s/channels/%s/threads", s.Discord_API, channelID)
    message_auteur := fmt.Sprintf("Posté par **%s** :\n\n%s", pseudo, message)
    newMsg := DiscordThread{
        Name: titre,
        Message: DiscordMessage{
            Content: message_auteur,
        },
    }
    jsonData, err := json.Marshal(newMsg)
    if err != nil {
        return DiscordThread{}, err
    }
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return DiscordThread{}, err
    }
    req.Header.Set("Authorization", "Bot "+s.Token)
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return DiscordThread{}, err
    }
    defer resp.Body.Close()
    var result DiscordThread
    err = json.NewDecoder(resp.Body).Decode(&result)
    return result, err
}

// fonction qui poste un message dans un thread
func (s *Discord_service) PostMessageDiscord(ctx context.Context, threadID string, text string, authorPseudo string, repond_id string) (DiscordMessage, error) {
    url := fmt.Sprintf("%s/channels/%s/messages", s.Discord_API, threadID)
    message_auteur := fmt.Sprintf("**%s** : %s", authorPseudo, text) // format du message afficher sur le discord
    newMsg := DiscordMessage{
        Content: message_auteur,
    }
    if repond_id != "" {// si on répond à un message
        newMsg.MessageReference = &DiscordMessageReference{MessageID: repond_id}
    }
    jsonData, err := json.Marshal(newMsg)
    if err != nil {
        return DiscordMessage{}, err
    }
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return DiscordMessage{}, err
    }
    req.Header.Set("Authorization", "Bot "+s.Token)
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return DiscordMessage{}, err
    }
    defer resp.Body.Close()
    var result DiscordMessage
    err = json.NewDecoder(resp.Body).Decode(&result)
    return result, nil
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
    for _, t := range reponse.Threads {// on le transforme en liste
        if t.ParentID == channelID { 
            allThreads = append(allThreads, t)
        }
    }
    return allThreads, nil
}

// fonction qui obtient tous les messages du thread en question
func (s *Discord_service) GetNewMessages(ctx context.Context, threadID string, lastID string) ([]DiscordMessage, error) {
    url := fmt.Sprintf("%s/channels/%s/messages", s.Discord_API, threadID)
    if lastID != "" {
        url = fmt.Sprintf("%s?after=%s", url, lastID) // Format de comment il va être afficher sur son discord
    }
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bot "+s.Token)
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

// foncrion qui permet de supprimer un message dans Discord
func (s *Discord_service) DeleteMessageDiscord(ctx context.Context, channel_id string, message_id string) error {
    url := fmt.Sprintf("%s/channels/%s/messages/%s", s.Discord_API, channel_id, message_id)
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bot "+s.Token)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}

// fonction qui permet de supprimer un message dans Discord
func (s *Discord_service) DeleteThreadDiscord(ctx context.Context, thread_id string) error {
    url := fmt.Sprintf("%s/channels/%s", s.Discord_API, thread_id)
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bot "+s.Token)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
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