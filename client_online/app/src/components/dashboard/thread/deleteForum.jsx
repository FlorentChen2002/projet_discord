import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import axios from 'axios';

//Supprime un commentaire ou un sujet du forum via une requête API
function DeleteForum({commentaire,query, onDelete}){
    //comportement
    const called = useRef(false);
    const navigate = useNavigate();
    //Envoie une requête DELETE à l'API pour supprimer le sujet ou le commentaire
    const supprimer = async() =>{
        try {
            const response = await axios.delete(`https://projet-discord.onrender.com/api/forum/delete/${query}`, {
                data: { id: commentaire.id },
                withCredentials: true
            });
            console.log("Requête DELETE envoyée avec succès :", response.data);
            if (response.status==200){
                if (query==="sujet"){
                    navigate("/forum");
                }else{
                    onDelete();
                }
            }
        }catch(e){
            console.error("Erreur lors de l'envoi de la requête :", e);
        }
    }
    //Effet de bord pour appeler la fonction de suppression une seule fois au montage du composant
    useEffect(() => {
        if (called.current) return;
        called.current = true;
        supprimer();
    }, []);
    return null;
}

export default DeleteForum;