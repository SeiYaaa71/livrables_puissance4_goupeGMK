package main

import (
    "html/template"
    "net/http"
    "strconv"

    "power4/game"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var currentGame = game.NewGame()

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/play", playHandler)
    http.HandleFunc("/reset", resetHandler)
    http.HandleFunc("/resetall", resetAllHandler)

    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")

    red, yellow := game.GetScores()
    data := struct {
        Title       string
        Grid        [game.Rows][game.Cols]int
        Current     int
        Winner      int
        ScoreRed    int
        ScoreYellow int
    }{
        Title:       "Puissance 4",
        Grid:        currentGame.Grid,
        Current:     currentGame.Current,
        Winner:      currentGame.Winner,
        ScoreRed:    red,
        ScoreYellow: yellow,
    }

    tmpl.Execute(w, data)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
    col, _ := strconv.Atoi(r.URL.Query().Get("col"))
    currentGame.Play(col)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
    currentGame.Reset()
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func resetAllHandler(w http.ResponseWriter, r *http.Request) {
    currentGame.Reset()
    game.ResetScores()
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
