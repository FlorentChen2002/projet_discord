package api

import (
    "encoding/json"
    "net/http"
    // "votre_projet/service" // Importer vos services
)

// Au lieu de init(db), on passe le service utilisateur
type UserAPI struct {
    Service *service.UserService
}

// Handler pour le Login (Équivalent de router.post("/login"))
func (api *UserAPI) LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }

    // 1. Lecture du Body (Équivalent de req.body)
    var credentials struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        sendError(w, 400, "Requête invalide")
        return
    }

    // 2. Logique métier (Vérification mot de passe)
    userID, err := api.Service.CheckPassword(credentials.Login, credentials.Password)
    if err != nil || userID == "" {
        sendError(w, 403, "Login et/ou mot de passe invalide(s)")
        return
    }

    // 3. Gestion de la Session (Concept Page 40 du cours)
    // En Go, on utilise souvent un cookie sécurisé pour stocker l'ID
    setSessionCookie(w, userID)

    // 4. Réponse JSON (Concept Page 3 du cours)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": 200,
        "message": "Login accepté",
        "id": userID,
    })
}