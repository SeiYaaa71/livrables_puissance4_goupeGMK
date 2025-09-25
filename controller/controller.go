package controller

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.Respo,seWriter, filename string, data map{string}string) {
	tmpl := template.Must(template.ParseFiles("views/" + filename))
	tmpl.Execute(w, data)
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Welcome to Connect Four",
	}
	renderTemplate(w, "home.html", data)	
}

func Play(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Play Connect Four",
	}
	renderTemplate(w, "play.html", data)	
}

func Contact(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Contact Us",
	}
	renderTemplate(w, "contact.html", data)	
}
