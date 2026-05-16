import { useNavigate } from "react-router-dom";
import React, { useState } from "react";
import NewPostForm from "./NewPostForm";
import DeleteForum from "./deleteForum";
import "./styles.css";

//Affiche un message d'utilisateur avec ses métadonnées, les réponses associées,
//et des actions comme répondre, supprimer ou accéder au profil
function Commentaire({ sujet, commentaire, user, onDelete, onAdd }) {
    // state initial pour l'affichage du formulaire de réponse et du composant de suppression
    const [showForm, setShowForm] = useState(false);
    const [showDeleteComposant, setShowDeleteComposant] = useState(false);
    //const [showRepond, setShowRepond] = useState(false);
    const navigate = useNavigate();
    const showRepond = commentaire?.repond?.[0]?.content;
    const showDelete = user.id.toString()===commentaire.userid.toString() || user.rang.toString()==="admin";
    const userid=commentaire.userid;
    //console.log(user._id.toString(),commentaire.user_id.toString(), showDelete);
    // rendu
    return (
        <div className="commentaire-body">
        <div className="commentaire-meta">
            <span className="author">{commentaire.userpseudo}</span>
            <span className="meta">{commentaire.date}</span>
        </div>
        <p>{commentaire.content}</p>
        {showRepond && (
            <div className="reply-block">
            <div className="commentaire-meta">
                <span className="author">{commentaire.repond[0].userpseudo}</span>
                <span className="meta">{commentaire.repond[0].date}</span>
            </div>
            <p>{commentaire.repond[0].content}</p>
            </div>
        )}
        <div className="commentaire-actions">
            <span onClick={() => setShowForm(!showForm)}>💬 Répondre</span>
            <span onClick={() => navigate("/profile", { state: { userid } })}>👤 Profil</span>
            {showDelete && <span onClick={() => setShowDeleteComposant(true)} >🗑️ Supprimer</span>}
            {showDeleteComposant && <DeleteForum commentaire={commentaire} query="thread" onDelete={() => { onDelete(commentaire.id); setShowDeleteComposant(false);}}/>}
        </div>
        {showForm && <NewPostForm user={user} sujet={sujet} commentaire={commentaire} onAdd={onAdd}/>}
        </div>
    );
}

export default Commentaire;
