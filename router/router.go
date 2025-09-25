package router

import (
    "html/template"
    "net/http"
    "strconv"

    "power4/game" // adapte au nom de module dans go.mod
)

var currentGame = game.NewGame()

func New() *http.ServeMux {
    mux := http.NewServeMux()

    // Fonction utilitaire seq
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

    // Page "À propos"
    mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("template/about.html"))
        tmpl.Execute(w, map[string]string{
            "Title":   "À propos",
            "Message": "Ceci est un projet Puissance 4 en Go.",
        })
    })

    // Page "Contact"
    mux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.ParseFiles("template/contact.html"))
        tmpl.Execute(w, map[string]string{
            "Title":   "Contact",
            "Message": "Envoyez-nous un message.",
        })
    })

    // Fichiers statiques
    mux.Handle("/stylecss/", http.StripPrefix("/stylecss/", http.FileServer(http.Dir("stylecss"))))

    return mux
}

