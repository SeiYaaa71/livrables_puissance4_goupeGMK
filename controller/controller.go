package controller

import (
    "html/template"
    "net/http"
    "power4/game"
)

// Structure pour passer toutes les infos au template
type PageData struct {
    Title   string
    Message string
    Grid    [game.Rows][game.Cols]int
    Current int
    Winner  int
    Stats   game.Stats // ✅ Ajout des scores globaux
}

// Fonction utilitaire pour rendre un template
func RenderTemplate(w http.ResponseWriter, filename string, data PageData) {
    tmpl, err := template.ParseFiles("template/" + filename)
    if err != nil {
        http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
        return
    }
    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Erreur exécution template : "+err.Error(), http.StatusInternalServerError)
    }
}

// About
func About(w http.ResponseWriter, r *http.Request) {
    data := PageData{
        Title:   "À propos",
        Message: "Ceci est la page À propos ✨",
        Stats:   game.GetScores(), // ✅ On passe aussi les scores
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


