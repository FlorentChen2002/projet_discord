## Projet web PC3R

**Auteur** : CHEN Florent (21101813)

**GitHub** : [Lien vers le projet web](https://github.com/FlorentChen2002/projet_discord)

**Discord** : [Serveur discord](https://discord.gg/aXWtgNb6kz)

**Site Web** : [lien du site web](https://projetdiscord-production.up.railway.app)

---

Le travail a été réalisé sur Windows en utilisant Ubuntu via WSL2.

## Architecture

Le projet est organisé en trois grands blocs:

- `rapport_florent_chen_21101813/` contient le rapport final du projet.

### Frontend local et en ligne

Les dossiers `client_localhost/` et `client_online/` correspondent à la partie frontend du projet: `client_localhost/` sert au développement local et `client_online/` à la mise en ligne.

- `app/` contient l'application React.
- `src/App.jsx` et `src/main.jsx` sont les points d'entrée de l'interface.
- `src/components/auth/` contient la connexion, l'inscription et le routage de l'authentification.
- `src/components/dashboard/` contient l'espace utilisateur après connexion.
	- `forum/` gère la liste des sujets et la création de sujet.
	- `thread/` gère l'affichage d'un sujet, des messages et des réponses.
	- `profile/` gère les profils utilisateurs et la liste des membres.
	- `deconnexion/` gère la déconnexion.
	- `recherche/` gère la recherche dans les sujets.

### Backend

Le dossier `serveur/` correspond à l'API Go du projet. Il est organisé en couches pour séparer les responsabilités:

- `main.go` lance le serveur et branche les routes.
- `handler/` gère les requêtes HTTP et les réponses.
- `middleware/` regroupe le CORS, l'authentification et le logging.
- `service/` contient la logique métier et la synchronisation avec Discord.
- `repository/` s'occupe de l'accès à MongoDB.

Cette séparation permet de garder le code plus lisible et plus simple à faire évoluer.

### Hébergement

Le projet est déployé sur plusieurs services:

- `client_online/` est hébergé sur Railway. [lien du site web](https://dashboard.render.com)
- `serveur/` est hébergé sur Render. [lien du site web](https://railway.com)
- La base de données MongoDB est hébergée sur MongoDB Atlas. [lien du site web](https://www.mongodb.com/products/platform/atlas-database)

### Lancement du projet

Pour lancer le projet en local, il est nécessaire d’ouvrir trois terminaux distincts afin de démarrer la base de données, le serveur backend, et le client React.

1. Dans le deuxième terminal, on lance le serveur :
```bash
cd serveur/
DISCORD_TOKEN="" go run main.go
```
2. Dans le troisième terminal, on lance le client :
```bash
cd /client_localhost/app
npm run dev
```
---

3. Le serveur MongoDB est lancé en ligne.
