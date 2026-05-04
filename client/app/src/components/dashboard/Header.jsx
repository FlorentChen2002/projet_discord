import React from "react";
import { Link,useNavigate } from "react-router-dom";



//Barre de navigation latérale avec liens vers les principales sections.
const Header = ({user}) => {
  //comportement
  const navigate = useNavigate();
  const submissionHandler = () =>{
    navigate("/postforum");
  }
  //affichage
  return (
    <div className="sidebar-gauche">
        <ul>
          <button className="submitButton" onClick={submissionHandler}> +    Sujet </button>
          <li>
            <Link to="/forum">Forum</Link>
          </li>
          <li>
            <Link to="/profile">Profile</Link>
          </li>
          <li>
            <Link to="/listuser">Utilisateurs</Link>
          </li>
          {(user.rang.toString()==="admin")&&<li>
            <Link to="/adminpanel">AdminPanel</Link>
          </li>}
        </ul>
    </div>
  );
};

export default Header;
