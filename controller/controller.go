package controller

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.Respo,seWriter, filename string, data map{string}string) {
	tmpl := template.Must(template.ParseFiles("views/" + filename))
	tmpl.Execute(w, data)
}
