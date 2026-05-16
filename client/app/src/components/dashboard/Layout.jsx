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
                <h1>Forum X Discord</h1>
                <div className="nav-links">
                    <a href="https://discord.gg/aXWtgNb6kz">Discord</a>
                    <a href="https://github.com/FlorentChen2002/projet_discord">Github  </a>
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