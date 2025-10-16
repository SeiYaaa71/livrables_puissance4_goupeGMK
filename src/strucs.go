package game

import "fmt"

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

// Crée une nouvelle partie
func NewGame() *Game {
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

// Joue un coup dans la colonne donnée
// Retourne (succès, message)
func (g *Game) Play(col int) (bool, string) {
    if col < 0 || col >= Cols || g.Winner != 0 {
        return false, "Coup invalide"
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
            return true, ""
        }
    }
    return false, "Colonne pleine"
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
