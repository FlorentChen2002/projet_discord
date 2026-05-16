import React,{useState , useEffect} from "react";
import { useLocation } from 'react-router-dom';
import axios from 'axios';
import "./profile.css";

// Composant d'affichage de profil utilisateur
// Affiche les informations de l'utilisateur, son statut (admin ou non) et l'historique de ses messages visibles selon les droits

function Profile({user}) {
    // state initial pour les messages, les informations de l'utilisateur affiché, et le statut d'admin
    const [messages,setMessages] = useState([]);
    const location = useLocation();
    const [users,setUsers] = useState(null);
    const [rang,setRang] = useState(users?.rang?.toString()==="admin");
    const showPrive = user?.rang?.toString()==="admin";
    
    // Fonction pour récupérer les informations d'un autre utilisateur en fonction de son ID
    const otherUser = async(userid) => {
        try {
            const response = await axios.get(`http://localhost:8000/api/user/${userid}`,{ withCredentials: true });
            console.log("Requête GET envoyée avec succès :", response.data);
            if (response.status==200){
                setUsers(response.data);
                setRang(response.data.rang.toString()==="admin");
                getMessage(response.data.id);
            }
        }catch(e){
            setUsers({id: '', pseudo: 'Non inscrit', mdp: '', date: 'None', rang: 'users'});
            setRang(false);
            getMessage(response.data.id);
            console.error("Erreur lors de l'envoi de la requête :", e);
        }
    }
    // Fonction pour récupérer les messages d'un utilisateur en fonction de son ID
    const getMessage = async(userid) =>{
        try {
            const response = await axios.get('http://localhost:8000/api/forum/getallthread',{ withCredentials: true });
            console.log("Requête GET envoyée avec succès :", response.data);
            if (response.status==200){
                let responseCopy=[...response.data];
                responseCopy=responseCopy.filter((responses) => responses.userid.toString()===userid.toString());
                setMessages(responseCopy);
            }
        }catch(e){
            console.error("Erreur lors de l'envoi de la requête :", e);
        }
    }
    // Effet de bord pour charger les données de l'utilisateur affiché au chargement du composant ou lors du changement d'ID utilisateur
    useEffect(() => {
        const loadData = async () => {
            if (location.state?.userid) {
                await otherUser(location.state.userid);
            }else{
                setUsers(user);
                await getMessage(user?.id);
            }
        };
        loadData();
    }, [location.state?.userid]);
    //rendu
    return (
        <div className="container-profile">
        <div className="profile-info">
            <h2>
            {users?.pseudo}
            {rang&&<span className="badge">Admin</span>}
            </h2>
            <p>Inscrit : {users?.date}</p>
            <p>Nombre de messages : {messages.length}</p>
        </div>
        <div className="messages">
            <h3>Historique des messages</h3>
            {messages.map((message) => (
            ((showPrive||!message.prive) && <div className="message" key={message.id}>{message.content}
                <div className="meta">{message.date}
                </div>
            </div>)
            ))}
        </div>
        </div>
    );
}

export default Profile;
