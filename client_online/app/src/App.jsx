/**
 * =========================================================================================================================
 * Composant React : TestAxios
 * Envoie une requête POST vers un serveur Express via Axios avec un champ texte
 * =========================================================================================================================
 */

import axios from 'axios'; // Importe la bibliothèque Axios pour faire des requêtes HTTP
import React,{ useState,useEffect } from 'react'; // Importe le hook useState de React pour gérer l'état local
import AuthRouter from "./components/auth/AuthRouter";
import DashBoard from "./components/dashboard/DashBoard";
import Deconnexion from "./components/dashboard/deconnexion/Deconnexion";
import { BrowserRouter } from "react-router-dom";

function App() {
    const [user,setUser]=useState("");
    const [isLoading, setIsLoading] = useState(true);
    // Fonction pour vérifier l'authentification de l'utilisateur en envoyant une requête GET à l'API
    const checkAuth = async () => {
        try {
        const response = await axios.get('https://projet-discord.onrender.com/api/user/me', {
            withCredentials: true
        });
        if(response.data){
            console.log("Requête GEt envoyée avec succès :", response.data);
            setUser(response.data);
            setIsLoading(true);
        }
        } catch (error) {
        setIsLoading(false);
        }
    };

    useEffect(() => { checkAuth(); }, []); // lance une fois pour vérifier si il est connecté

    const refreshAuth = () => {
        setIsLoading(true);
        checkAuth();
    };
    if (isLoading){
        return (
        <div className="App">
            <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
                <DashBoard users={user}/>
            </BrowserRouter>
        </div>
        );
    }


    return (
        <div className="App">
            <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
                <AuthRouter refreshAuth={refreshAuth} />
            </BrowserRouter>
        </div>
    );
}

export default App;