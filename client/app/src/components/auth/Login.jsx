import React, { useState,useEffect } from "react";
import axios from 'axios';
import { Link,useNavigate } from "react-router-dom";
import "./Login.css";

// Composant de connexion utilisateur.
// Gère l'envoi du pseudo/mot de passe et redirige après authentification réussie.
function Login({ refreshAuth }) {
  //state
  const [pseudo, setPseudo] = useState("");
  const [password, setPassword] = useState("");
  const [passOK, setPassOK] = useState(true);
  const navigate = useNavigate();
  //comportement
  const getPseudo = (evt) => {
    setPseudo(evt.target.value);
  };
  const getPassword = (evt) => {
    setPassword(evt.target.value);
  };

  useEffect(() => {
    document.body.classList.add("login");
    return () => {
      document.body.classList.remove("login");
    };
  },[]);

  const submissionHandler = async (event) => {
    event.preventDefault();

    try {
      const response = await axios.post('http://localhost:8000/api/user/connexion', {
        Pseudo: pseudo, Mdp: password
      },{ withCredentials: true });
      console.log("Requête POST envoyée avec succès :", response.data);
      setPassOK(true);
      if (response.data.status==200){
  
        const userId = response.data.id;
        console.log("Connexion réussie, user ID :", userId);
        refreshAuth();
        navigate("/")
      }
    } catch (error) {
      setPassOK(false);
      console.error("Erreur lors de l'envoi de la requête :", error);
    }
  };

  //affichage
  return (
    <form className="login-box">
      <header className="header-login">Se connecter</header>
      {passOK ? (
        <p></p>
      ) : (
        <p style={{ color: "red" }}>
          Veuillez vérifier votre mot de passe et votre nom de compte et réessayer !
        </p>
      )}
      <div className="input-box">
        <input
          type="text"
          className="input-field"
          placeholder="Pseudo"
          value={pseudo}
          onChange={getPseudo}
        />
      </div>
      <div className="input-box">
        <input
          type="password"
          className="input-field"
          placeholder="Mot de passe"
          value={password}
          onChange={getPassword}
        />
      </div>
      <div className="input-submit">
      <button className="submit-btn" onClick={submissionHandler}></button>
        <label>Connexion</label>
      </div>
      <div className="sign-up-link">
        <p>
          Vous n'avez pas de compte ? <Link to="/register">Créer un compte</Link>
        </p>
      </div>
    </form>
  );
}

export default Login;
