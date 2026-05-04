import React, { useEffect, useState } from 'react';
import axios from 'axios';
import "./styleadmin.css";


//Composant AdminUsers
//Interface d’administration pour gérer les utilisateurs :
//promotion/dégradation admin, validation, suppression et recherche par pseudo.
function AdminUsers() {
  const [users, setUsers] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    const res = await axios.get("http://localhost:8000/api/admin/rocher", { withCredentials: true });
    setUsers(res.data);
    console.log("Requête GET envoyée avec succès :", res.data);
  };

  const handleAction = async (id, action) => {
    let url = `http://localhost:8000/api/admin/${action}/${id}`;
    try {
      if (action === "refuse") {
        const res = await axios.delete(url, { withCredentials: true });
        console.log("Requête DELETE envoyée avec succès :", res.data);
      } else {
        const res = await axios.post(url, {}, { withCredentials: true });
        console.log("Requête POST envoyée avec succès :", res.data);
      }
      fetchUsers(); // refresh list
    } catch (err) {
      console.error("Erreur API:", err);
    }
  };

  const filteredUsers = users.filter(user =>
    user.pseudo.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="admin-container">
      <h2 className="title">Gestion des utilisateurs</h2>

      <input
        type="text"
        placeholder="🔍 Rechercher un utilisateur"
        className="search-input"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
      />

      <ul className="user-list">
        {filteredUsers.map(user => {
          const userClass = user.rang === "admin"
            ? "admin"
            : user.status
              ? "valide"
              : "non-valide";

          return (
            <li key={user._id} className={`user-item ${userClass}`}>
              <span className="user-info">
                <strong>{user.pseudo}</strong> — <em>{user.rang}</em> — Status: {user.status ? "✔️ Validé" : "❌ En attente"}
              </span>
              <div className="actions">
                {!(user.rang.toString()==="admin")&&<button
                  className={`btn-promote ${user.status ? "active" : "disabled"}`}
                  disabled={!user.status}
                  onClick={() => handleAction(user._id, "promote")}
                >
                  Promouvoir admin
                </button>}

                <button className="btn-demote" onClick={() => handleAction(user._id, "demote")}>
                  Retirer admin
                </button>

                {!user.status && (
                  <button className="btn-validate" onClick={() => handleAction(user._id, "validate")}>
                    Valider
                  </button>
                )}

                <button className="btn-refuse" onClick={() => handleAction(user._id, "refuse")}>
                  Supprimer
                </button>
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

export default AdminUsers;

