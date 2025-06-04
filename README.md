# IoT Dashboard

Tableau de bord IoT pour la découverte, la visualisation et le contrôle des objets connectés sur un réseau local.

## Fonctionnalités principales
- Scan automatique et manuel du réseau local pour détecter les objets connectés (via Nmap)
- Affichage dynamique des objets (nom, type, description, statut, actions)
- Ajout d'un mock caméra locale (jamais supprimé automatiquement)
- Suppression et modification (nom, description) de chaque objet
- Actions dynamiques par objet (ex : voir la caméra)
- Notification en cas de changement réseau
- UI moderne et responsive

## Prérequis
- **Go** (1.18 ou supérieur)
- **Nmap** installé et accessible dans le PATH ([télécharger ici](https://nmap.org/download.html))

## Installation
1. **Clone le dépôt**
   ```powershell
   git clone https://github.com/Zephied/IoT_Dashboard
   cd IoT-Dashboard
   ```
2. **Installe Nmap**
   - Télécharge et installe Nmap depuis [nmap.org/download.html](https://nmap.org/download.html)
   - Ajoute le dossier d'installation (contenant `nmap.exe`) à ta variable d'environnement `PATH`
   - Vérifie dans un terminal :
     ```powershell
     nmap --version
     ```

3. **Lance le backend Go**
   ```powershell
   cd backend
   go run main.go
   ```
   Le serveur démarre sur [http://localhost:8080](http://localhost:8080)

4. **Accède à l'interface web**
   - Ouvre un navigateur sur [http://localhost:8080](http://localhost:8080)

## Utilisation
- **Scanner le réseau** : Clique sur "Scanner le réseau" pour détecter les objets connectés (scan manuel). Un scan automatique est aussi lancé toutes les 30 secondes.
- **Supprimer/Modifier** : Utilise les boutons sur chaque carte pour supprimer ou modifier le nom/description d'un objet.
- **Actions** : Clique sur les boutons d'action pour interagir avec l'objet (ex : voir la caméra pour le mock).
- **Notification** : Si des changements réseau sont détectés, une notification s'affiche pour inviter à rafraîchir la page.

## Structure du projet
```
backend/
  main.go           # Point d'entrée Go (API + serveur fichiers statiques)
  api/              # Handlers API, actions, mock
  db/               # Accès base SQLite
  models/           # Modèles Go
  scanner/          # Scan réseau (Nmap)
frontend/
  index.html        # Interface utilisateur
  script.js         # Logique frontend
  style.css         # Styles
```

## Notes
- Le mock caméra locale n'est jamais supprimé automatiquement et permet de tester l'affichage du flux webcam local.
- Les devices hors ligne restent visibles (statut rouge) jusqu'à suppression manuelle.
- Le projet fonctionne sous Windows, Linux ou Mac (tant que Go et Nmap sont installés).

## Dépannage
- **Erreur 500 lors du scan** : Vérifie que Nmap est bien installé et accessible dans le PATH.
- **404 sur les routes API** : Vérifie que tu lances bien le backend avec le bon routeur (voir main.go).
- **Problème de webcam** : Autorise l'accès à la caméra dans ton navigateur.

---

Développé pour la gestion locale d'objets connectés et la démo IoT.