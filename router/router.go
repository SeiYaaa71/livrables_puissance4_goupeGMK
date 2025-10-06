package router

import (
    "html/template"
    "net/http"
    "strconv"
    "power4/game" // correspond au module défini dans go.mod
)

var currentGame = game.NewGame()

// New retourne un ServeMux avec toutes les routes configurées
func New() *http.ServeMux {
    mux := http.NewServeMux()

    // Fonction utilitaire seq pour générer les colonnes (0 → 6)
    funcMap := template.FuncMap{
        "seq": func(start, end int) []int {
            s := make([]int, end-start+1)
            for i := range s {
                s[i] = start + i
            }
            return s
        },
    }

    // Page d'accueil
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(
            template.New("index.html").Funcs(funcMap).ParseFiles("template/index.html"),
        )
        data := map[string]interface{}{
            "Title":   "Puissance 4",
            "Message": "Bienvenue sur le jeu !",
            "Grid":    currentGame.Grid,
            "Current": currentGame.Current,
            "Winner":  currentGame.Winner,
        }
        tmpl.Execute(w, data)
    })

    // Route pour jouer un coup
    mux.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
        colStr := r.URL.Query().Get("col")
        col, err := strconv.Atoi(colStr)
        if err == nil {
            currentGame.Play(col)
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    // Route pour reset
    mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            currentGame.Reset()
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    // Fichiers statiques (CSS, images…)
    mux.Handle("/stylecss/", http.StripPrefix("/stylecss/", http.FileServer(http.Dir("stylecss"))))

    return mux
}



