import React, { useState } from 'react';
//import axios from 'axios';
import './search.css';

// Composant de barre de recherche générique.
// Prend une fonction onRecherche en prop, appelée lors de la soumission avec la requête utilisateur.
function Recherche({onRecherche}) {
  const [query, setQuery] = useState('');
  /*
  const [results, setResults] = useState([]);
  const handleSearch = async (e) => {
    e.preventDefault();
    const res = await axios.get(`http://localhost:8000/api/search/topic?q=${query}`);
    setResults(res.data);
  };*/

  return (
      <form onSubmit={(e) => onRecherche(e,query)} className="search-form horizontal">
        <input
          type="text"
          placeholder="Recherche ..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button type="submit">Rechercher</button>
      </form>
  );
}

export default Recherche;

