package repository

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Accès à la collection "users" de la base de données MongDB 
type User_reposite struct {
    Collection *mongo.Collection
}

// Structure représentant un utilisateur dans la base de données
type User_db struct {
    Id string `bson:"_id,omitempty" json:"id"`
    Pseudo string `bson:"pseudo" json:"pseudo"`
    Mdp string `bson:"mdp" json:"mdp"`
    Date string `bson:"date" json:"date"`
    Rang string `bson:"rang" json:"rang"`
}

// fonctions qui permet d'accéder à la base de donnée pour obtenir les données de l'utilsateur via son pseudo
func (r *User_reposite) GetUser(ctx context.Context, pseudo string)(*User_db, error){
    var user User_db
    err:= r.Collection.FindOne(ctx,bson.M{"pseudo":pseudo}).Decode(&user)
    if err!=nil {
        return nil, err
    }
    return &user, nil
}

// fonction qui permet d'accéder à la base de donnée pour obtenir les données de l'utilisateur via son id
func (r *User_reposite) GetID(ctx context.Context, id string)(*User_db, error){
    var user User_db
    err := r.Collection.FindOne(ctx,bson.M{"_id":id}).Decode(&user)
    if err!=nil {
        return nil, err
    }
    return &user, nil
}

// fonction qui vérifie si un utilisateur existe dans la base de données via son pseudo
func (r *User_reposite) Exists(ctx context.Context, pseudo string)bool {
    count, err := r.Collection.CountDocuments(ctx, bson.M{"pseudo":pseudo})
    if err!=nil {
        return false
    }
    return count >0
}

// fonction qui permet de créer un utilisateur dans la base de données
func (r *User_reposite) CreateUser(ctx context.Context, pseudo string, mdp string, rang string) error{
    user := User_db{ Id: primitive.NewObjectID().Hex(), Pseudo: pseudo, Mdp: mdp, Date: time.Now().Format("2006-01-02"), Rang: rang, }
    _, err := r.Collection.InsertOne(ctx, user)
    return err
}

// fonction qui permet d'obtenir tous les utilisateurs de la base de données
func (r *User_reposite) GetAllUser(ctx context.Context)([]User_db, error){
    var users []User_db
    cursor, err := r.Collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    err = cursor.All(ctx, &users)
    if err != nil {
        return nil, err
    }
    return users, nil
}