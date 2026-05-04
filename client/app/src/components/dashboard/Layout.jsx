import { Outlet } from "react-router-dom";
import React, { useEffect } from "react";
import Header from "./Header";
import Deconnexion from "./deconnexion/Deconnexion";
import "./forum/styles.css";



//Structure principale de la page avec header, barre de navigation latérale et contenu dynamique via Outlet.
const Layout = ({user}) => {
    useEffect(() => {
        document.body.classList.add("forum");
        return () => {
        document.body.classList.remove("forum");
        };
    }, []);
    return (
        <div>
            <header className="header">
                <h1>Forum</h1>
                <div className="nav-links">
                    <a href="#">Vidéo</a>
                    <a href="https://github.com/FlorentChen2002/Omega">Github  </a>
                    <Deconnexion/>
                </div>
            </header>
            <div className="container">
                <Header user={user}/>
                <div className="content">
                    <Outlet />
                </div>
            </div>

        </div>
    );
};

export default Layout;