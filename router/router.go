package router

import (
    "html/template"
    "net/http"
    "sort"
    "strconv"

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
        "seq": func(start, end int) []int {
            s := make([]int, end-start+1)
            for i := range s {
                s[i] = start + i
            }
            return s
        },
        "add":       func(a, b int) int { return a + b },
        "cellClass": cellClass,
    }

    // Page d'accueil (jeu)
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(
            template.New("index.html").Funcs(funcMap).ParseFiles("template/index.html"),
        )

        stats := game.GetScores()

        type Player struct {
            Name  string
            Icon  string
            Score int
            Class string
        }
        players := []Player{
            {Name: "Joueur Rouge", Icon: "🔴", Score: stats.Red, Class: "player-red"},
            {Name: "Joueur Jaune", Icon: "🟡", Score: stats.Yellow, Class: "player-yellow"},
        }
        sort.Slice(players, func(i, j int) bool {
            return players[i].Score > players[j].Score
        })

        data := map[string]interface{}{
            "Title":       "Puissance 4",
            "Message":     lastMessage,
            "Grid":        currentGame.Grid,
            "Current":     currentGame.Current,
            "Winner":      currentGame.Winner,
            "ScoreRed":    stats.Red,
            "ScoreYellow": stats.Yellow,
            "GamesPlayed": stats.Games,
            "Draws":       stats.Draws,
            "Players":     players,
        }

        tmpl.Execute(w, data)
    })

    // Page À propos
    mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("template/about.html"))
        tmpl.Execute(w, map[string]interface{}{"Title": "À propos"})
    })

    // Page Contact
    mux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("template/contact.html"))
        tmpl.Execute(w, map[string]interface{}{"Title": "Contact"})
    })

    // Page Tableau des scores
    mux.HandleFunc("/tableau", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(
            template.New("tableau.html").Funcs(funcMap).ParseFiles("template/tableau.html"),
        )

        stats := game.GetScores()
        type Player struct {
            Name  string
            Icon  string
            Score int
            Class string
        }
        players := []Player{
            {Name: "Joueur Rouge", Icon: "🔴", Score: stats.Red, Class: "player-red"},
            {Name: "Joueur Jaune", Icon: "🟡", Score: stats.Yellow, Class: "player-yellow"},
        }
        sort.Slice(players, func(i, j int) bool {
            return players[i].Score > players[j].Score
        })

        data := map[string]interface{}{
            "Title":   "Tableau des scores",
            "Players": players,
        }
        tmpl.Execute(w, data)
    })

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
            lastMessage = "🗑️ Scores et parties réinitialisés !"
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    // Fichiers statiques (CSS, images…)
    mux.Handle("/stylecss/", http.StripPrefix("/stylecss/", http.FileServer(http.Dir("stylecss"))))
    mux.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image"))))

    return mux
}


