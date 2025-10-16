package game

import (
    "fmt"
    "math/rand"
    "time"
)

const (
    Rows = 6
    Cols = 7
)

// Stats globales (cumulées sur toutes les parties)
type Stats struct {
    Red    int
    Yellow int
    Games  int
    Draws  int
}

var GlobalStats Stats

type Game struct {
    Grid    [Rows][Cols]int // 0 = vide, 1 = joueur rouge, 2 = joueur jaune
    Current int             // joueur courant
    Winner  int             // 0 = pas de gagnant, 1 = rouge, 2 = jaune
}


// Liste de messages d’encouragement
var encouragements = []string{
    "💡 Belle tentative, continue comme ça !",
    "🔥 La partie s’échauffe, ne lâche rien !",
    "🎯 Stratégie intéressante, à toi de jouer !",
    "💪 Tu peux le faire, concentre-toi !",
    "⚡ Beau coup, la tension monte !",
    "🚀 Tu prends de la vitesse, continue !",
    "🌟 Impressionnant, quel sens du jeu !",
    "🧠 Belle réflexion, ça se voit que tu anticipes !",
    "🏹 Tu vises juste, garde le cap !",
    "🎶 Le rythme est bon, ne t’arrête pas !",
    "🔥 Tu mets la pression, bien joué !",
    "💥 Coup puissant, ça change la partie !",
    "🌈 Quelle créativité, bravo !",
    "🕹️ Tu joues comme un pro !",
    "⚔️ La bataille est serrée, tiens bon !",
    "🏆 Tu te rapproches de la victoire !",
    "🎉 Super mouvement, ça va payer !",
    "🌀 Tu crées la surprise, excellent !",
    "🧩 Ton coup s’emboîte parfaitement !",
    "🌍 Toute la salle retient son souffle !",
    "✨ Tu brilles sur ce coup !",
    "📈 Ta stratégie monte en puissance !",
    "💎 Coup précieux, bien trouvé !",
    "🔮 On dirait que tu vois l’avenir !",
}


// Nouvelle partie
func NewGame() *Game {
    rand.Seed(time.Now().UnixNano())
    return &Game{Current: 1}
}

// Change de joueur
func (g *Game) switchPlayer() {
    if g.Current == 1 {
        g.Current = 2
    } else {
        g.Current = 1
    }
}

// Joue un coup et retourne un message
func (g *Game) Play(col int) (bool, string) {
    if col < 0 || col >= Cols || g.Winner != 0 {
        return false, "❌ Coup invalide"
    }

    for row := Rows - 1; row >= 0; row-- {
        if g.Grid[row][col] == 0 {
            g.Grid[row][col] = g.Current

            // Vérifie victoire
            if g.checkWin(row, col) {
                g.Winner = g.Current
                GlobalStats.Games++
                if g.Winner == 1 {
                    GlobalStats.Red++
                } else {
                    GlobalStats.Yellow++
                }
                return true, fmt.Sprintf("🎉 Joueur %d a gagné ! 🏆", g.Winner)
            }

            // Vérifie égalité
            if g.isBoardFull() {
                GlobalStats.Games++
                GlobalStats.Draws++
                return true, "🤝 Match nul !"
            }

            // Sinon, on change de joueur
            g.switchPlayer()

            // Tirer un message d’encouragement aléatoire
            msg := encouragements[rand.Intn(len(encouragements))]
            return true, msg
        }
    }
    return false, "⚠️ Colonne pleine"
}

// Vérifie si le coup joué est gagnant
func (g *Game) checkWin(row, col int) bool {
    player := g.Grid[row][col]
    if player == 0 {
        return false
    }

    directions := [][2]int{
        {0, 1},  // horizontal
        {1, 0},  // vertical
        {1, 1},  // diagonale ↘
        {1, -1}, // diagonale ↙
    }

    for _, d := range directions {
        count := 1
        count += g.countDir(row, col, d[0], d[1], player)
        count += g.countDir(row, col, -d[0], -d[1], player)
        if count >= 4 {
            return true
        }
    }
    return false
}

// Compte les pions alignés dans une direction donnée
func (g *Game) countDir(r, c, dr, dc, player int) int {
    count := 0
    for {
        r += dr
        c += dc
        if r < 0 || r >= Rows || c < 0 || c >= Cols {
            break
        }
        if g.Grid[r][c] != player {
            break
        }
        count++
    }
    return count
}

// Vérifie si la grille est pleine
func (g *Game) isBoardFull() bool {
    for c := 0; c < Cols; c++ {
        if g.Grid[0][c] == 0 {
            return false
        }
    }
    return true
}

// Réinitialise la grille mais garde les scores
func (g *Game) Reset() {
    for r := 0; r < Rows; r++ {
        for c := 0; c < Cols; c++ {
            g.Grid[r][c] = 0
        }
    }
    g.Current = 1
    g.Winner = 0
}

// Retourne les scores globaux
func GetScores() Stats {
    return GlobalStats
}

// Réinitialise les scores
func ResetScores() {
    GlobalStats = Stats{}
}



