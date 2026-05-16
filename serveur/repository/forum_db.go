package repository

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// Accès à la collection "sujet" et "thread" de la base de données MongDB, (sujet = thread principal, thread = message)
type Forum_reposite struct {
    Sujet_Collection *mongo.Collection
    Thread_Collection *mongo.Collection
}

// Structure représentant un sujet dans la base de données
type Sujet_db struct {
    Id string `bson:"_id,omitempty" json:"id"`
    Titre string `bson:"titre" json:"titre"`
    Description string `bson:"description" json:"description"`
    Date string `bson:"date" json:"date"`
    User_id string `bson:"user_id" json:"userid"`
    User_pseudo string `bson:"user_pseudo" json:"userpseudo"`
    Prive bool `bson:"prive" json:"prive"`
}

// Structure représentant un message dans la base de données
type Message_db struct {
    Id string `bson:"_id,omitempty" json:"id"`
    Sujet_id string `bson:"sujet_id" json:"sujetid"`
    Content string `bson:"content" json:"content"`
    Date string `bson:"date" json:"date"`
    User_id string `bson:"user_id" json:"userid"`
    User_pseudo string `bson:"user_pseudo" json:"userpseudo"`
    Prive bool `bson:"prive" json:"prive"`
    Repond []Repond_db `bson:"repond" json:"repond"`
}

// Structure représentant une réponse à un message dans la base de données
type Repond_db struct {
    Id string `bson:"_id,omitempty" json:"id"`
    Content string `bson:"content" json:"content"`
    User_id string `bson:"user_id" json:"userid"`
    User_pseudo string `bson:"user_pseudo" json:"userpseudo"`
    Date string `bson:"date" json:"date"`
}

// fonctions qui crée un sujet dans la base de données
func (r *Forum_reposite) Create_sujet (ctx context.Context, id string,titre string, description string, user_id string, user_pseudo string, private bool)(bool, error){
    sujet := Sujet_db{Id: id, Titre: titre, Description: description, Date: time.Now().Format("2006-01-02"), User_id: user_id, User_pseudo: user_pseudo, Prive: private,}
    _, err := r.Sujet_Collection.InsertOne(ctx, sujet)
    if err != nil {
        return false, err
    }
    return true, nil
}

// fonction qui permet d'obtenir tous les sujets de la base de données
func (r *Forum_reposite) GetAllSujet (ctx context.Context)([]Sujet_db,error){
    var sujet []Sujet_db
    cursor, err := r.Sujet_Collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    err = cursor.All(ctx, &sujet)
    if err != nil {
        return nil, err
    }
    return sujet, nil
}

// fonction qui permet d'obtenir un sujet de la base de données via son id
func (r *Forum_reposite) GetSujet (ctx context.Context, sujet_id string)(Sujet_db,error){
    var sujet Sujet_db
    err := r.Sujet_Collection.FindOne(ctx, bson.M{"_id": sujet_id}).Decode(&sujet)
    if err != nil {
        return Sujet_db{}, err
    }
    return sujet, nil
}

// fonction qui permet de supprimer un sujet de la base de données via son ID, ainsi que tous les messages associés
func (r *Forum_reposite) DeleteSujet (ctx context.Context, id string)(bool,error){
    _, err := r.Sujet_Collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        return false, err
    }
    thread := bson.M{"sujet_id": id} 
    _, err = r.Thread_Collection.DeleteMany(ctx, thread)
    if err != nil {
        return false, err
    }
    return true, nil
}

// fonction qui permet de créer un message dans la base de données
func (r *Forum_reposite) Create_message (ctx context.Context, msg_id string, sujet_id string, content string, user_id string, user_pseudo string, private bool, repond []Repond_db)(bool, error){
    msg := Message_db{Id: msg_id, Sujet_id: sujet_id, Content: content, Date: time.Now().Format("2006-01-02"), User_id: user_id, User_pseudo: user_pseudo, Prive: private, Repond: repond,}
    _, err := r.Thread_Collection.InsertOne(ctx, msg)
    if err != nil {
        return false, err
    }
    return true, nil
}

// fonction qui permet d'obtenir un message de la base de données via son id 
func (r *Forum_reposite) GetMessage (ctx context.Context, msg_id string)(Message_db,error){
    var msg Message_db
    err := r.Thread_Collection.FindOne(ctx, bson.M{"_id": msg_id}).Decode(&msg)
    if err != nil {
        return Message_db{}, err
    }
    return msg, nil
}

// fonction qui permet d'obtenir tous les messages d'un sujet de la base de données via l'id du sujet
func (r *Forum_reposite) GetSujetMessage (ctx context.Context, sujet_id string)([]Message_db,error){
    cursor, err := r.Thread_Collection.Find(ctx, bson.M{"sujet_id": sujet_id})
    defer cursor.Close(ctx)
    if err != nil {
        return nil, err
    }
    var messages []Message_db
    err = cursor.All(ctx, &messages);
    if err != nil {
        return nil, err
    }
    if messages == nil {
        messages = []Message_db{}
    }
    return messages, nil
}

// fonction qui permet d'obtenir tous les messages de la base de données
func (r *Forum_reposite) GetAllMessage (ctx context.Context)([]Message_db,error){
    var msg []Message_db
    cursor, err := r.Thread_Collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    err = cursor.All(ctx, &msg)
    if err != nil {
        return nil, err
    }
    return msg, nil
}

// fonction qui permet de supprimer un message via son id
func (r *Forum_reposite) DeleteMessage (ctx context.Context, id string)(bool,error){
    resultat, err := r.Thread_Collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        return false, err
    }
    if resultat.DeletedCount == 0 {
        return false, nil
    }
    return true, nil
}