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

func main() {
    // Nettoyage de la console au démarrage
    clearConsole()

    mux := router.New()

    green := "\033[32m"
    yellow := "\033[33m"
    reset := "\033[0m"

    fmt.Printf("%s🚀 Serveur lancé !🚀%s\n", green, reset)
    fmt.Printf("%s🌐 http://localhost:8080 🌐%s\n", yellow, reset)

    log.Fatal(http.ListenAndServe(":8080", mux))
}



