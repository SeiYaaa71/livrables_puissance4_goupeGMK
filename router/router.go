package router

import (
	"html/template"
	"net/http"
	"strconv"

	"power4/controller" // ✅ IMPORT AJOUTÉ
	"power4/game"
)

var currentGame = game.NewGame()
var lastMessage string // message du dernier coup

// cellClass traduit 0/1/2 en classes CSS
func cellClass(v int) string {
	switch v {
	case 1:
		return "player1"
	case 2:
		return "player2"
	default:
		return ""
	}
}

// New retourne un ServeMux avec toutes les routes configurées
func New() *http.ServeMux {
	mux := http.NewServeMux()

	// Fonctions utilitaires pour les templates
	funcMap := template.FuncMap{
		"cellClass": cellClass,
	}

	// Page d'accueil (jeu)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(
			template.New("index.html").Funcs(funcMap).ParseFiles("template/index.html"),
		)

		stats := game.GetScores()

		data := map[string]interface{}{
			"Title":   "Puissance 4",
			"Message": lastMessage,
			"Grid":    currentGame.Grid,
			"Current": currentGame.Current,
			"Winner":  currentGame.Winner,
			"Stats":   stats, // ✅ Envoi de la structure Stats complète
		}

		tmpl.Execute(w, data)
	})

	// === ROUTES VERS LE CONTRÔLEUR ===
	// Au lieu de gérer la logique ici, on appelle les fonctions du contrôleur
	
	// Page À propos
	mux.HandleFunc("/about", controller.About) // ✅ MODIFIÉ

	// Page Contact
	mux.HandleFunc("/contact", controller.Contact) // ✅ MODIFIÉ

	// Page Tableau des scores
	mux.HandleFunc("/tableau", controller.HandleTableau) // ✅ MODIFIÉ (Étape 3)
	
	// API pour sauvegarder le jeu
	mux.HandleFunc("/api/save-game", controller.HandleSaveGame) // ✅ AJOUTÉ (Étape 2)
	
	// === ROUTES DE JEU ===

	// Route pour jouer un coup
	mux.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
		colStr := r.URL.Query().Get("col")
		col, err := strconv.Atoi(colStr)
		if err == nil {
			_, msg := currentGame.Play(col)
			lastMessage = msg
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Reset (nouvelle partie, scores conservés)
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			currentGame.Reset()
			lastMessage = "🔄 Nouvelle partie lancée !"
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Reset complet (scores + partie)
	mux.HandleFunc("/resetall", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			currentGame.Reset()
			game.ResetScores()
			// ✅ On DOIT aussi effacer le fichier JSON
			controller.ClearHistoryFile() // On appelle la fonction de reset du contrôleur
			lastMessage = "🗑️ Scores et parties réinitialisés !"
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// === FICHIERS STATIQUES ===
	// ✅ CORRIGÉ: Chemins mixtes (certains à la racine, d'autres dans static)
	mux.Handle("/stylecss/", http.StripPrefix("/stylecss/", http.FileServer(http.Dir("stylecss"))))
	mux.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js")))) // ✅ Chemin corrigé pour JS

	return mux
}

