import React, { useState, useEffect } from "react";
import Commentaire from "./commentaire";
import "./styles.css";

//Affiche la liste des commentaires d'un sujet de forum.
const PostList = ({ sujet,commentaires: initialCommentaires, user, onAdd }) => {
  const [liste, setListe] = useState(initialCommentaires);
  useEffect(() => {
    setListe(initialCommentaires);
  }, [initialCommentaires]);
  const handleSuppressionLocale = (idSupprime) => {
    setListe((prev) => prev.filter((c) => c.id !== idSupprime));
  };
  return (
    <div className="commentaire">
      <h1>Commentaire :</h1>
      {liste?.map((commentaire) => (
          <Commentaire
          key={commentaire.id}
          sujet={sujet}
          commentaire={commentaire}
          user={user}
          onDelete={handleSuppressionLocale}
          onAdd={onAdd}
        />
      ))}
    </div>
  );
};

export default PostList;
