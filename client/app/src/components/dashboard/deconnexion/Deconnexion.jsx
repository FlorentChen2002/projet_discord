import { useNavigate } from "react-router-dom";
import axios from "axios";
import "./deconnexion.css";

// Composant de déconnexion utilisateur
// Envoie une requête pour déconnecter l'utilisateur et redirige vers la page d'accueil
function Deconnexion() {
    const navigate = useNavigate();
    
    const logout = async () => {
        try {
            const response = await axios.post('http://localhost:8000/api/user/deconnexion', {}, {
                withCredentials: true
            });
            console.log("Requête POST envoyée avec succès :", response.data);
            if (response.status === 200) {
                navigate("/");
                window.location.reload();
            }
        } catch (e) {
            console.error("Erreur lors de la déconnexion :", e);
        }
    };
    // rendu
    return (
        <button className="submitButton-post" onClick={logout}>
            Déconnexion
        </button>
    );
}

export default Deconnexion;
