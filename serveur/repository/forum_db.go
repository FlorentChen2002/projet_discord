package repository

import (
    "context"
    "time"
    "log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Forum_reposite struct {
    Sujet_Collection *mongo.Collection
    Thread_Collection *mongo.Collection
}

type Sujet_db struct {
    Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Titre string `bson:"titre" json:"titre"`
    Description string `bson:"description" json:"description"`
    Date string `bson:"date" json:"date"`
    User_id primitive.ObjectID `bson:"user_id" json:"userid"`
    User_pseudo string `bson:"user_pseudo" json:"userpseudo"`
    Prive bool `bson:"prive" json:"prive"`
}

type Thread_db struct {
    Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Sujet_id primitive.ObjectID `bson:"sujet_id" json:"sujetid"`
    Content string `bson:"content" json:"content"`
    Date string `bson:"date" json:"date"`
    User_id primitive.ObjectID `bson:"user_id" json:"userid"`
    User_pseudo string `bson:"user_pseudo" json:"userpseudo"`
    Prive bool `bson:"prive" json:"prive"`
    Repond []Repond_db `bson:"repond" json:"repond"`
}

type Repond_db struct {
    Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Content string `bson:"content" json:"content"`
    User_id primitive.ObjectID `bson:"user_id" json:"userid"`
    User_pseudo string `bson:"user_pseudo" json:"userpseudo"`
    Date string `bson:"date" json:"date"`
}

func (r *Forum_reposite) Create_sujet (ctx context.Context, titre string, description string, user_id string, user_pseudo string, private bool)(bool, error){
    obj_user_id, err := primitive.ObjectIDFromHex(user_id)
    if err != nil {
        return false, err
    }
    sujet := Sujet_db{Titre: titre, Description: description, Date: time.Now().Format("2006-01-02"), User_id: obj_user_id, User_pseudo: user_pseudo, Prive: private,}
    _, err = r.Sujet_Collection.InsertOne(ctx, sujet)
    if err != nil {
        return false, err
    }
    return true, nil
}

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

func (r *Forum_reposite) DeleteSujet (ctx context.Context, id string)(bool,error){
    objId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return false, err
    }
    _, err = r.Sujet_Collection.DeleteOne(ctx, bson.M{"_id": objId})
    if err != nil {
        return false, err
    }
    thread := bson.M{"sujet_id": objId} 
    
    _, err = r.Thread_Collection.DeleteMany(ctx, thread)
    if err != nil {
        return false, err
    }
    return true, nil
}

func (r *Forum_reposite) Create_thread (ctx context.Context, sujet_id string, content string, user_id string, user_pseudo string, private bool, repond []Repond_db)(bool, error){
    obj_sujet_id, err := primitive.ObjectIDFromHex(sujet_id)
    if err != nil{
        log.Printf("erreur InsertOne: %v", err)
        return false, err
    }
    obj_user_id, err := primitive.ObjectIDFromHex(user_id)
    if err != nil{
        log.Printf("erreur InsertOne: %v", err)
        return false, err
    }
    thread := Thread_db{Sujet_id: obj_sujet_id, Content: content, Date: time.Now().Format("2006-01-02"), User_id: obj_user_id, User_pseudo: user_pseudo, Prive: private, Repond: repond,}
    _, err = r.Thread_Collection.InsertOne(ctx, thread)
    if err != nil {
        log.Printf("erreur InsertOne: %v", err)
        return false, err
    }
    return true, nil
}

func (r *Forum_reposite) GetThread (ctx context.Context, sujet_id string)([]Thread_db,error){
    obj_sujet_id, err := primitive.ObjectIDFromHex(sujet_id)
    if err != nil {
        return nil, err
    }
    cursor, err := r.Thread_Collection.Find(ctx, bson.M{"sujet_id": obj_sujet_id})
    defer cursor.Close(ctx)
    if err != nil {
        return nil, err
    }
    var threads []Thread_db
    err = cursor.All(ctx, &threads);
    if err != nil {
        return nil, err
    }
    if threads == nil {
        threads = []Thread_db{}
    }
    return threads, nil
}

func (r *Forum_reposite) GetAllThread (ctx context.Context)([]Thread_db,error){
    var thread []Thread_db
    cursor, err := r.Thread_Collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    err = cursor.All(ctx, &thread)
    if err != nil {
        return nil, err
    }
    return thread, nil
}

func (r *Forum_reposite) DeleteThread (ctx context.Context, id string)(bool,error){
    objId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return false, err
    }
    resultat, err := r.Thread_Collection.DeleteOne(ctx, bson.M{"_id": objId})
    if err != nil {
        return false, err
    }
    if resultat.DeletedCount == 0 {
        return false, nil
    }
    return true, nil
}