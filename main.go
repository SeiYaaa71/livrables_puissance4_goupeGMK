package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "runtime"

    "power4/router"
)

const (
    green  = "\033[32m"
    yellow = "\033[33m"
    blue   = "\033[34m"
    red    = "\033[31m"
    reset  = "\033[0m"
)

// clearConsole efface la console selon l'OS
func clearConsole() {
    switch runtime.GOOS {
    case "windows":
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    default: // linux, mac, etc.
        cmd := exec.Command("clear")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

// loggingMiddleware affiche uniquement certaines requêtes HTTP
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 👉 On ignore tous les GET
        if r.Method == http.MethodGet {
            next.ServeHTTP(w, r)
            return
        }

        // Couleur selon méthode
        methodColor := reset
        switch r.Method {
        case "POST":
            methodColor = blue
        case "DELETE":
            methodColor = red
        }

        fmt.Printf("%s➡️  %s %s%s\n", methodColor, r.Method, r.URL.Path, reset)
        next.ServeHTTP(w, r)
    })
}

func main() {
    // Nettoyage de la console au démarrage
    clearConsole()

    mux := router.New()
    loggedMux := loggingMiddleware(mux)

    fmt.Printf("%s🚀 Serveur lancé !🚀%s\n", green, reset)
    fmt.Printf("%s🌐 http://localhost:8080 🌐%s\n", yellow, reset)

    log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
