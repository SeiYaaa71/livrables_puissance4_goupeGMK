package router

import (
	"net/http"
	"power4/controller"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", controller.HomeHandler)
	mux.HandleFunc("/play", controller.PlayHandler)
	mux.HandleFunc("/reset", controller.ResetHandler)
	mux.HandleFunc("/status", controller.StatusHandler)
	mux.HandleFunc("/contact", controller.ContactHandler)
	return mux
}
