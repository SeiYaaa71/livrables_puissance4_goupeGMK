package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "os/signal"
    "runtime"
    "syscall"
    "time"

    "power4/router"
)

const (
    green  = "\033[32m"
    yellow = "\033[33m"
    reset  = "\033[0m"
)

func clearConsole() {
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }
    cmd.Stdout = os.Stdout
    _ = cmd.Run()
}

func main() {
    clearConsole()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    mux := router.New()

    srv := &http.Server{
        Addr:         ":" + port,
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    fmt.Printf("%s🚀 Serveur lancé ! 🚀%s\n", green, reset)
    fmt.Printf("%s🌐 http://localhost:%s 🌐%s\n", yellow, port, reset)

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Erreur serveur: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    fmt.Println("\n🛑 Arrêt du serveur en cours...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Arrêt forcé: %v", err)
    }

    fmt.Println("✅ Serveur arrêté proprement")
}


