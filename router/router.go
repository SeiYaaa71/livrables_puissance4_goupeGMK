package router

import (
    "html/template"
    "net/http"
    "strconv"

    "power4/controller"
    "power4/game"
)

var currentGame = game.NewGame()
var lastMessage string

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

func New() *http.ServeMux {
    mux := http.NewServeMux()

    funcMap := template.FuncMap{
        "cellClass": cellClass,
    }

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(
            template.New("index.html").Funcs(funcMap).ParseFiles("template/index.html"),
        )

        stats := game.GetScores()

        data := map[string]interface{}{
            "Title":     "Puissance 4",
            "Message":   lastMessage,
            "Grid":      currentGame.Grid,
            "Current":   currentGame.Current,
            "Winner":    currentGame.Winner,
            "Stats":     stats,
        }

        tmpl.Execute(w, data)
    })

    mux.HandleFunc("/about", controller.About)

    mux.HandleFunc("/contact", controller.Contact)

    mux.HandleFunc("/tableau", controller.HandleTableau)
    
    mux.HandleFunc("/api/save-game", controller.HandleSaveGame)

    mux.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
        colStr := r.URL.Query().Get("col")
        col, err := strconv.Atoi(colStr)
        if err == nil {
            _, msg := currentGame.Play(col)
            lastMessage = msg
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            currentGame.Reset()
            lastMessage = "🔄 Nouvelle partie lancée !"
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    mux.HandleFunc("/resetall", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            currentGame.Reset()
            game.ResetScores()
            controller.ClearHistoryFile()
            lastMessage = "🗑️ Scores et parties réinitialisés !"
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    })

    mux.Handle("/stylecss/", http.StripPrefix("/stylecss/", http.FileServer(http.Dir("stylecss"))))
    mux.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image"))))
    mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))

    return mux
}