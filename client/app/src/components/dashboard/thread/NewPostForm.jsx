// components/NewPostForm/index.jsx
import React, { useState } from "react";
import "./styles_module.css";
import axios from 'axios';

//Formulaire pour poster un nouveau commentaire ou une réponse dans un sujet de forum
//Envoie les données au serveur via une requête POST 

const NewPostForm = ({ user,sujet,commentaire, onAdd }) => {
    const [content, setContent] = useState("");
    //serveur
    //Gère la soumission du formulaire en envoyant une requête POST à l'API avec les données du nouveau commentaire ou sujet
    const handleSubmit = async(e) => {
        e.preventDefault();
        if (!content.trim()) return;
        try{
            console.log("Envoi de la requête POST avec les données :", commentaire)
            const repondData = commentaire ? [{
                id: commentaire.id,
                content: commentaire.content,
                userid: commentaire.userid,
                userpseudo: commentaire.userpseudo,
                date: commentaire.date
            }] : [];
            const response =await axios.post('http://localhost:8000/api/forum/postthread', {
                sujetid: sujet.id,
                content: content,
                userid: user.id,
                userpseudo: user.pseudo,
                date: new Date().toLocaleString('fr-FR'),
                prive:sujet.prive,
                repond: repondData
            },{ withCredentials: true });
            console.log("Requête POST envoyée avec succès :", response.data);
            setContent("");
            onAdd(response.data);
        }catch(e){
            console.error("Erreur lors de l'envoi de la requête :", e);
        }
    };
    // rendu
    return (
        <form onSubmit={handleSubmit} className="newPostForm-post">
        <textarea
            placeholder="Votre commentaire"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="textareaField-post"
            required
        />
        <button type="submit" className="submitButton-post">
            Envoyer
        </button>
        </form>
    );
};

export default NewPostForm;
