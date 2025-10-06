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
}

// Fonction utilitaire
func RenderTemplate(w http.ResponseWriter, filename string, data PageData) {
    tmpl := template.Must(template.ParseFiles("template/" + filename))
    tmpl.Execute(w, data)
}

// About
func About(w http.ResponseWriter, r *http.Request) {
    data := PageData{
        Title:   "À propos",
        Message: "Ceci est la page À propos ✨",
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
        }
        RenderTemplate(w, "contact.html", data)
        return
    }
    data := PageData{
        Title:   "Contact",
        Message: "Envoie-nous un message 📩",
    }
    RenderTemplate(w, "contact.html", data)
}


