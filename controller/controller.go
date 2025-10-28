package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort" // AJOUTÉ pour le tri
	"time"

	"power4/game" // C'est le paquet "game"
	src "power4/src"  // ✅ CORRIGÉ: On importe "power4/src" et on l'appelle "src"
)

// --- Structure de données pour les pages ---

// PageData pour les pages de jeu, about, contact
type PageData struct {
	Title   string
	Message string
	Grid    [game.Rows][game.Cols]int
	Current int
	Winner  int
	Stats   game.Stats
}

// TableauData pour la page de l'historique
type TableauData struct {
	Title   string
	History []src.GameHistoryEntry // ✅ CORRIGÉ: Utilise "src."
}

// Fonction utilitaire pour rendre un template (GÉNÉRALISÉE)
// Elle accepte n'importe quel type de données (data interface{})
func RenderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	// Ajout de la fonction "cellClass" (si elle n'est pas déjà dans votre router)
	funcMap := template.FuncMap{
		"cellClass": func(v int) string {
			switch v {
			case 1:
				return "player1"
			case 2:
				return "player2"
			default:
				return ""
			}
		},
	}

	// S'assure de parser le bon nom de template
	baseName := "template/" + filename

	tmpl, err := template.New(filename).Funcs(funcMap).ParseFiles(baseName)

	if err != nil {
		log.Printf("Erreur parsing template %s: %v", filename, err)
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Erreur exécution template %s: %v", filename, err)
		http.Error(w, "Erreur exécution template : "+err.Error(), http.StatusInternalServerError)
	}
}

// --- Handlers pour les pages statiques ---

// About
func About(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "À propos",
		Message: "Ceci est la page À propos ✨",
		Stats:   game.GetScores(),
	}
	RenderTemplate(w, "about.html", data)
}

// Contact
func Contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		msg := r.FormValue("msg")
		data := PageData{
			Title:   "Contact",
			Message: "Merci " + name + " pour ton message : " + msg,
			Stats:   game.GetScores(),
		}
		RenderTemplate(w, "contact.html", data)
		return
	}
	data := PageData{
		Title:   "Contact",
		Message: "Envoie-nous un message 📩",
		Stats:   game.GetScores(),
	}
	RenderTemplate(w, "contact.html", data)
}

// --- ÉTAPE 3 (LECTURE) ---

// HandleTableau lit le JSON et l'affiche
func HandleTableau(w http.ResponseWriter, r *http.Request) {
	history, err := readHistoryFile()
	if err != nil {
		log.Printf("Erreur lecture fichier pour tableau: %v", err)
		http.Error(w, "Could not read history", http.StatusInternalServerError)
		return
	}

	// Inverse la liste pour voir les plus récents en premier
	sort.SliceStable(history, func(i, j int) bool {
		return i > j // Simple inversion d'index
	})

	data := TableauData{
		Title:   "Tableau des scores",
		History: history, // ✅ CORRIGÉ: Utilise la structure de données
	}
	RenderTemplate(w, "tableau.html", data)
}

// --- ÉTAPE 2 (SAUVEGARDE) ---

// Chemin vers notre fichier de sauvegarde
const historyFilePath = "historique.json"

// HandleSaveGame reçoit les données du jeu (fetch) et les sauvegarde dans le JSON
func HandleSaveGame(w http.ResponseWriter, r *http.Request) {
	// 1. Décode le JSON envoyé par le client (index.html)
	var gameData struct {
		Winner string `json:"winner"`
		Loser  string `json:"loser"`
		IsDraw bool   `json:"isDraw"`
	}

	if err := json.NewDecoder(r.Body).Decode(&gameData); err != nil {
		log.Printf("Erreur décodage JSON: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Crée une nouvelle entrée d'historique
	newEntry := src.GameHistoryEntry{ // ✅ CORRIGÉ: Utilise "src."
		Winner:  gameData.Winner,
		Loser:   gameData.Loser,
		IsDraw:  gameData.IsDraw,
		DateStr: time.Now().Format("02/01/2006 15:04"),
	}

	// 3. Lit le fichier JSON existant
	history, err := readHistoryFile()
	if err != nil {
		log.Printf("Erreur lecture fichier: %v", err)
		http.Error(w, "Could not read history", http.StatusInternalServerError)
		return
	}

	// 4. Ajoute la nouvelle entrée
	history = append(history, newEntry)

	// 5. Écrit le fichier JSON mis à jour
	if err := writeHistoryFile(history); err != nil {
		log.Printf("Erreur écriture fichier: %v", err)
		http.Error(w, "Could not save history", http.StatusInternalServerError)
		return
	}

	// 6. Répond au client que tout s'est bien passé
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	log.Println("Partie sauvegardée:", newEntry)
}

// ClearHistoryFile vide le fichier d'historique (appelé par /resetall)
func ClearHistoryFile() {
	log.Println("Effacement du fichier d'historique...")
	// Crée un fichier avec une liste vide "[]"
	if err := ioutil.WriteFile(historyFilePath, []byte("[]"), 0644); err != nil {
		log.Printf("Erreur lors de l'effacement du fichier d'historique: %v", err)
	}
}

// --- Fonctions utilitaires ---

// readHistoryFile lit le fichier JSON et le retourne
func readHistoryFile() ([]src.GameHistoryEntry, error) { // ✅ CORRIGÉ: Utilise "src."
	var history []src.GameHistoryEntry // ✅ CORRIGÉ: Utilise "src."

	// Assure que le fichier existe
	if _, err := os.Stat(historyFilePath); os.IsNotExist(err) {
		log.Println("Fichier historique non trouvé, création...")
		// Crée un fichier avec une liste vide "[]"
		if err := ioutil.WriteFile(historyFilePath, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("could not create history file: %w", err)
		}
	}

	// Lit le fichier
	file, err := ioutil.ReadFile(historyFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read history file: %w", err)
	}

	// Décode le JSON
	if err := json.Unmarshal(file, &history); err != nil {
		// Si le fichier est vide ou corrompu, on repart d'une liste vide
		if len(file) == 0 {
			return []src.GameHistoryEntry{}, nil // ✅ CORRIGÉ: Utilise "src."
		}
		log.Printf("Fichier JSON corrompu ou vide, réinitialisation. Erreur: %v", err)
		return []src.GameHistoryEntry{}, nil // ✅ CORRIGÉ: Utilise "src."
	}

	return history, nil
}

// writeHistoryFile écrit la liste complète dans le fichier JSON
func writeHistoryFile(history []src.GameHistoryEntry) error { // ✅ CORRIGÉ: Utilise "src."
	// Convertit la liste en JSON (avec indentation pour la lisibilité)
	data, err := json.MarshalIndent(history, "", "  ") // 2 espaces pour indentation
	if err != nil {
		return fmt.Errorf("could not marshal history: %w", err)
	}

	// Écrit dans le fichier
	if err := ioutil.WriteFile(historyFilePath, data, 0644); err != nil {
		return fmt.Errorf("could not write history file: %w", err)
	}

	return nil
}

