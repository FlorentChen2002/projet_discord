package repository

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Accès à la collection "messages" de la base de données MongDB (pour stocker le dernier message de chaque thread)
type Discord_reposite struct {
    Message_Collection *mongo.Collection
}

// Structure représentant le dernier message d'un thread dans la base de données
type LastMessage_db struct {
    Id string `bson:"_id,omitempty" json:"id"`
    Msg_id string `bson:"msg_id" json:"msgid"`
    Thread_id string `bson:"thread_id" json:"threadid"`
}

// fonction qui permet de stocker le dernier message d'un thread dans la base de données
func (r *Discord_reposite) SetLastMessage(ctx context.Context, msg_id string, thread_id string) (error) {
    filter := bson.M{"thread_id": thread_id}
    update := bson.M{"$set": bson.M{"msg_id": msg_id, "thread_id": thread_id},}
    opts := options.Update().SetUpsert(true)
    _, err := r.Message_Collection.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        return err
    }
    return nil
}

// fonction qui permet de récupérer le dernier message d'un thread dans la base de données
func (r *Discord_reposite) GetLastMessage (ctx context.Context, thread_id string)(LastMessage_db,error){
    var result LastMessage_db
    err := r.Message_Collection.FindOne(ctx, bson.M{"thread_id": thread_id}).Decode(&result)
    if err != nil {
        return LastMessage_db{}, err
    }
    return result, nil
}

// fonction qui permet de supprimer le dernier message d'un thread dans la base de données
func (r *Discord_reposite) DeleteLastMessage (ctx context.Context, id string)(bool,error){
    _, err := r.Message_Collection.DeleteOne(ctx, bson.M{"thread_id": id})
    if err != nil {
        return false, err
    }
    return true, nil
}