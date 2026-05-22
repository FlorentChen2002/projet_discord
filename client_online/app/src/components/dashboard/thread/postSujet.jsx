import { useNavigate } from "react-router-dom";
import React, { useEffect, useState } from "react";
import DeleteForum from "./deleteForum";
import "./styles.css";


//Affiche les détails d’un sujet de forum avec options pour répondre
//accéder au profil de l’auteur, et supprimer si autorisé
function PostSujet({ user, sujet }) {
    const [showDeleteComposant, setShowDeleteComposant] = useState(false);
    const showDelete = user.id.toString()===sujet.userid.toString() || user.rang.toString()==="admin";
    const navigate = useNavigate();
    const userid =sujet.user_id;
    //Effet de bord pour faire défiler la page vers le bas
    useEffect(() => {
        window.scrollTo({
        top: document.body.scrollHeight,
        behavior: "smooth",
        });
    }, []);
    //Gère le clic sur le bouton de réponse en faisant défiler la page vers le bas
    const handleClick = () => {
        window.scrollTo(0, document.body.scrollHeight);
    };
    // rendu
    return (
        <div className="main-content">
        <h2>{sujet.titre}</h2>
        <p className="meta">
            {sujet.userpseudo} · {sujet.date}
        </p>
        <p className="article-text">{sujet.description}</p>
        <div className="commentaire-actions">
            <span onClick={handleClick}>💬 Répondre</span>
            <span onClick={() => navigate("/profile", { state: { userid } })}>👤 profile</span>
            {showDelete && <span onClick={() => setShowDeleteComposant(true)} >🗑️ Supprimer</span>}
            {showDeleteComposant && <DeleteForum commentaire={sujet} query="sujet"/>}
        </div>
        </div>
    );
}

export default PostSujet;
