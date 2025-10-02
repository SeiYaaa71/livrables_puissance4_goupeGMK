package main

import (
    "fmt"
    "log"
    "net/http"

    "power4/router"
)

func main() {
    mux := router.New()
    green := "\033[32m"
    yellow := "\033[33m"
    reset := "\033[0m"
    fmt.Printf("%s🚀 Serveur lancé !🚀%s\n", green, reset)
    fmt.Printf("%s🌐 http://localhost:8080 🌐%s\n", yellow, reset)

    log.Fatal(http.ListenAndServe(":8080", mux))
}


