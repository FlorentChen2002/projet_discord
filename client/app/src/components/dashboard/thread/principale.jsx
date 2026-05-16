import React, { useEffect, useState } from "react";
import { useLocation } from 'react-router-dom';
import NewPostForm from "./NewPostForm";
import PostSujet from "./postSujet";
import PostList from "./postList";
import axios from 'axios';

//Affiche un sujet de forum, sa liste de commentaires (threads) 
//et un formulaire pour poster un nouveau commentaire
//Charge les commentaires depuis l'API au chargement
function Principale({user}) {
    //
    const [threads,setThreads] = useState([]);
    const location = useLocation();
    const sujet = location.state?.sujet;
    //Fonction pour récupérer les commentaires d'un sujet de forum via une requête GET à l'API
    const getAllThread = async() =>{
        if (!sujet?.id) {
        return;
        }
        try {
            const response = await axios.get('http://localhost:8000/api/forum/getthread', {
                params: { sujetid: sujet.id },
                withCredentials: true
            });
            console.log("Requête GET envoyée avec succès :", response.data);
            if (response.status==200){
                setThreads(response.data);
            }
        }catch(e){
            console.error("Erreur lors de l'envoi de la requête :", e);
        }
    }
    //useEffect liste de commentaire
    //Effet de bord pour ajouter une classe CSS au body et charge les commentaires du sujet et configurer une connexion SSE pour recevoir les mises à jour
    useEffect(() => {
        document.body.classList.add("sujet");
        getAllThread();
        const source = new EventSource("http://localhost:8000/api/forum/events", {
            withCredentials: true,
        });
        source.addEventListener("forum", (event) => {
            const update = JSON.parse(event.data);
            console.log("Mise à jour reçue :", update);
            if (update.type === "message") {
                getAllThread();
            }
        });
        source.onerror = () => {
            source.close();
        };
        return () => {
            document.body.classList.remove("sujet");
            source.close();
        };
    }, [sujet?.id]);
    //rendu
    return (
        <main className="containers contents">
        <section>
            <PostSujet sujet={sujet} user={user}/>
            <PostList user={user} commentaires={threads} sujet={sujet} onAdd={getAllThread}/>
            <div className="commentaire">
            <NewPostForm user={user} sujet={sujet} commentaire={[]} onAdd={getAllThread} />
            </div>
        </section>
        </main>
    );
}

export default Principale;
