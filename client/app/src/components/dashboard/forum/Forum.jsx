import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import Recherche from "../recherche/Recherche";
import axios from 'axios';
import "./styles.css";

// Composant d'affichage des sujets/topic du forum.
// Récupère tous les sujets et permet une recherche par titre.
// Les sujets privés sont visibles uniquement par les admins.
function Forum({ user }) {
    // État initial
    const [memoire, setMemoire] = useState([]);
    const [sujets, setSujets] = useState([]);
    const [recherche, setRecherche] = useState("");
    const rechercheRef = useRef("");
    const navigate = useNavigate();

    const filtreSujets = (data, valeur) => {
        const rechercheNormalisee = valeur.trim().toLowerCase();
        if (!rechercheNormalisee) {
            return data;
        }
        return data.filter((obj) =>
            obj.titre.toLowerCase().includes(rechercheNormalisee)
        );
    };

    // comportement
    const getAllSujet = async() =>{
        try {
            const response = await axios.get('http://localhost:8000/api/forum/sujet',{ withCredentials: true });
            console.log("Requête GET envoyée avec succès :", response.data);
            if (response.status==200){
                setMemoire(response.data);
                setSujets(filtreSujets(response.data, rechercheRef.current));
            }
        }catch(e){
            console.error("Erreur lors de l'envoi de la requête :", e);
        }
    }

    const cherche = (e,resultat) =>{
        e.preventDefault();
        rechercheRef.current = resultat;
        setRecherche(resultat);
        setSujets(filtreSujets(memoire, resultat));
    }

    useEffect(() => {
        rechercheRef.current = recherche;
    }, [recherche]);

    useEffect(() => {
        document.body.classList.add("forum");
        getAllSujet();

        const source = new EventSource("http://localhost:8000/api/forum/events", {
            withCredentials: true,
        });

        source.addEventListener("forum", (event) => {
            const update = JSON.parse(event.data);
            if (update.type === "sujet") {
                getAllSujet();
            }
        });

        source.onerror = () => {
            source.close();
        };

        return () => {
            document.body.classList.remove("forum");
            source.close();
        };
    }, []);

    // Rendu
    return (
        <div className="forum-container">
            <Recherche onRecherche={cherche} />
            <main>
                { sujets.map((sujet) => {
                    const showSujet = !sujet.prive || user.rang?.toString() === "admin";
                    if (!showSujet) return null;

                    return (
                        <div 
                            className="post-link" 
                            key={sujet.id}
                            onClick={() => navigate("/forum/sujet", { state: { sujet } })}
                            role="button"
                            tabIndex={0}
                        >
                            <div className="post-forum">
                                <h3>{sujet.titre}
                                    {sujet.prive&& <span className="badge-forum"> privée </span>}
                                </h3>
                                <p>
                                    - Auteur: {sujet.user_pseudo} - Posté: {sujet.date}
                                </p>
                            </div>
                        </div>
                    );
                })}
            </main>
        </div>
        
    );
}

export default Forum;