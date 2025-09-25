package router

import (
    "html/template"
    "net/http"
)

func New() *http.ServeMux {
    mux := http.NewServeMux()

    // Fonction utilitaire seq pour générer des boucles dans les templates
    funcMap := template.FuncMap{
        "seq": func(start, end int) []int {
            s := make([]int, end-start+1)
            for i := range s {
                s[i] = start + i
            }
            return s
        },
    }

    // Page d'accueil (Puissance 4)
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(
            template.New("index.html").Funcs(funcMap).ParseFiles("template/index.html"),
        )
        tmpl.Execute(w, map[string]string{
            "Title":   "Puissance 4",
            "Message": "Bienvenue sur le jeu !",
        })
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

    // Fichiers statiques (CSS, images…)
    mux.Handle("/stylecss/", http.StripPrefix("/stylecss/", http.FileServer(http.Dir("stylecss"))))

    return mux
}
