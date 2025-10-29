package controller

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "sort"
    "time"

    "power4/game"
    src "power4/src"
)

type PageData struct {
    Title   string
    Message string
    Grid    [game.Rows][game.Cols]int
    Current int
    Winner  int
    Stats   game.Stats
}

type TableauData struct {
    Title   string
    History []src.GameHistoryEntry
}

func RenderTemplate(w http.ResponseWriter, filename string, data interface{}) {
    funcMap := template.FuncMap{
        "cellClass": func(v int) string {
            switch v {
            case 1:
                return "player1"
            case 2:
                return "player2"
            default:
                return ""
            }
        },
    }

    baseName := "template/" + filename

    tmpl, err := template.New(filename).Funcs(funcMap).ParseFiles(baseName)

    if err != nil {
        log.Printf("Erreur parsing template %s: %v", filename, err)
        http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
        return
    }
    if err := tmpl.Execute(w, data); err != nil {
        log.Printf("Erreur exécution template %s: %v", filename, err)
        http.Error(w, "Erreur exécution template : "+err.Error(), http.StatusInternalServerError)
    }
}

func About(w http.ResponseWriter, r *http.Request) {
    data := PageData{
        Title:   "À propos",
        Message: "Ceci est la page À propos ✨",
        Stats:   game.GetScores(),
    }
    RenderTemplate(w, "about.html", data)
}

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

func HandleTableau(w http.ResponseWriter, r *http.Request) {
    history, err := readHistoryFile()
    if err != nil {
        log.Printf("Erreur lecture fichier pour tableau: %v", err)
        http.Error(w, "Could not read history", http.StatusInternalServerError)
        return
    }

    sort.SliceStable(history, func(i, j int) bool {
        return i > j
    })

    data := TableauData{
        Title:   "Tableau des scores",
        History: history,
    }
    RenderTemplate(w, "tableau.html", data)
}

const historyFilePath = "historique.json"

func HandleSaveGame(w http.ResponseWriter, r *http.Request) {
    var gameData struct {
        Winner string `json:"winner"`
        Loser  string `json:"loser"`
        IsDraw bool   `json:"isDraw"`
    }

    if err := json.NewDecoder(r.Body).Decode(&gameData); err != nil {
        log.Printf("Erreur décodage JSON: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    newEntry := src.GameHistoryEntry{
        Winner:  gameData.Winner,
        Loser:   gameData.Loser,
        IsDraw:  gameData.IsDraw,
        DateStr: time.Now().Format("02/01/2006 15:04"),
    }

    history, err := readHistoryFile()
    if err != nil {
        log.Printf("Erreur lecture fichier: %v", err)
        http.Error(w, "Could not read history", http.StatusInternalServerError)
        return
    }

    history = append(history, newEntry)

    if err := writeHistoryFile(history); err != nil {
        log.Printf("Erreur écriture fichier: %v", err)
        http.Error(w, "Could not save history", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
    log.Println("Partie sauvegardée:", newEntry)
}

func ClearHistoryFile() {
    log.Println("Effacement du fichier d'historique...")
    if err := ioutil.WriteFile(historyFilePath, []byte("[]"), 0644); err != nil {
        log.Printf("Erreur lors de l'effacement du fichier d'historique: %v", err)
    }
}

func readHistoryFile() ([]src.GameHistoryEntry, error) {
    var history []src.GameHistoryEntry

    if _, err := os.Stat(historyFilePath); os.IsNotExist(err) {
        log.Println("Fichier historique non trouvé, création...")
        if err := ioutil.WriteFile(historyFilePath, []byte("[]"), 0644); err != nil {
            return nil, fmt.Errorf("could not create history file: %w", err)
        }
    }

    file, err := ioutil.ReadFile(historyFilePath)
    if err != nil {
        return nil, fmt.Errorf("could not read history file: %w", err)
    }

    if err := json.Unmarshal(file, &history); err != nil {
        if len(file) == 0 {
            return []src.GameHistoryEntry{}, nil
        }
        log.Printf("Fichier JSON corrompu ou vide, réinitialisation. Erreur: %v", err)
        return []src.GameHistoryEntry{}, nil
    }

    return history, nil
}

func writeHistoryFile(history []src.GameHistoryEntry) error {
    data, err := json.MarshalIndent(history, "", "  ")
    if err != nil {
        return fmt.Errorf("could not marshal history: %w", err)
    }

    if err := ioutil.WriteFile(historyFilePath, data, 0644); err != nil {
        return fmt.Errorf("could not write history file: %w", err)
    }

    return nil
}
