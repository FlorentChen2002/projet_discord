package repository

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User_reposite struct {
    Collection *mongo.Collection
}

type User_db struct {
    Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Pseudo string `bson:"pseudo" json:"pseudo"`
    Mdp string `bson:"mdp" json:"mdp"`
    Date string `bson:"date" json:"date"`
    Rang string `bson:"rang" json:"rang"`
}

func (r *User_reposite) GetUser(ctx context.Context, pseudo string)(*User_db, error){
    var user User_db
    err:= r.Collection.FindOne(ctx,bson.M{"pseudo":pseudo}).Decode(&user)
    if err!=nil {
        return nil, err
    }
    return &user, nil
}

func (r *User_reposite) GetID(ctx context.Context, id string)(*User_db, error){
    objId, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
    var user User_db
    err = r.Collection.FindOne(ctx,bson.M{"_id":objId}).Decode(&user)
    if err!=nil {
        return nil, err
    }
    return &user, nil
}

func (r *User_reposite) Exists(ctx context.Context, pseudo string)bool {
    count, err := r.Collection.CountDocuments(ctx, bson.M{"pseudo":pseudo})
    if err!=nil {
        return false
    }
    return count >0
}

func (r *User_reposite) CreateUser(ctx context.Context, pseudo string, mdp string, rang string) error{
    user := User_db{ Pseudo: pseudo, Mdp: mdp, Date: time.Now().Format("2006-01-02"), Rang: rang, }
    _, err := r.Collection.InsertOne(ctx, user)
    return err
}

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