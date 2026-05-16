package service

import (
    "fmt"
    "time"
    "context"
    "errors"
    "serveur/repository"
)

// Conteneur pour avoir accès à la base de données
type User_service struct {
    User_reposite *repository.User_reposite
}

// fonction qui permet de se connecter à un compte utilisateur
// vérifie si le pseudo existe et si le mot de passe est correct ( en vérifiant que le pseudo et le mot de passe se concordent )
// retourne l'id de l'utilisateur si la connexion est réussie
func (r *User_service) Connexion(ctx context.Context, pseudo string, mdp string) (string, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    user, err := r.User_reposite.GetUser(ctx, pseudo)
    if err != nil {
        fmt.Printf("Erreur : l'utilisateur non trouvé\n")
        return "", errors.New("l'utilisateur non trouvé")
    }
    if user.Pseudo == pseudo && user.Mdp == mdp {
        return user.Id, nil
    }
    fmt.Printf("Erreur : le mot de passe est faux\n")
    return "", errors.New("le mot de passe est faux")
}

// fonction qui permet de s'inscrire à un compte utilisateur
// vérifie si le pseudo existe déjà dans la base de données
// si le pseudo n'existe pas, on crée un nouvel utilisateur
func (r *User_service) Inscription(ctx context.Context, pseudo string, mdp string, rang string) error {
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    if r.User_reposite.Exists(ctx, pseudo) {
        fmt.Printf("Erreur : l'utilisateur existe déjà\n")
        return errors.New("l'utilisateur existe déjà")
    }
    return r.User_reposite.CreateUser(ctx, pseudo, mdp, rang)
}

// fonction qui permet d'obtenir les données d'un utilisateur via son id
func (r *User_service) GetId(ctx context.Context, id string) (repository.User_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    user, err := r.User_reposite.GetID(ctx, id)
    if err != nil {
        return repository.User_db{}, errors.New("Erreur : l'id non trouvé")
    }
    return *user,nil
}

// fonction qui permet d'avoir tous les utilsateurs
func (r *User_service) GetAllUsers(ctx context.Context)([]repository.User_db, error){
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    user, err := r.User_reposite.GetAllUser(ctx)
    if err != nil {
        return nil, errors.New("Erreur : liste d'utilisateur non trouvé")
    }
    return user,nil
}