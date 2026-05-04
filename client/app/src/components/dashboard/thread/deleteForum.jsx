import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import axios from 'axios';

//Supprime un commentaire ou un sujet du forum via une requête API au montage.
function DeleteForum({commentaire,query, onDelete}){
    const called = useRef(false);
    const navigate = useNavigate();

    const supprimer = async() =>{
        try {
            const response = await axios.delete(`http://localhost:8000/api/forum/delete/${query}`, {
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

    useEffect(() => {
        if (called.current) return;
        called.current = true;
        supprimer();
    }, []);

    return null;
}

export default DeleteForum;