package game

const (
    Rows = 6
    Cols = 7
)

// Scores globaux (cumulés sur toutes les parties)
var scoreRed int
var scoreYellow int
var gamesPlayed int // compteur de parties jouées
var draws int       // compteur d'égalités

type Game struct {
    Grid    [Rows][Cols]int // 0 = vide, 1 = joueur rouge, 2 = joueur jaune
    Current int             // joueur courant
    Winner  int             // 0 = pas de gagnant, 1 = rouge, 2 = jaune
}

// Crée une nouvelle partie
func NewGame() *Game {
    return &Game{Current: 1}
}

// Joue un coup dans une colonne
func (g *Game) Play(col int) bool {
    if col < 0 || col >= Cols || g.Winner != 0 {
        return false
    }

    for row := Rows - 1; row >= 0; row-- {
        if g.Grid[row][col] == 0 {
            g.Grid[row][col] = g.Current
            if g.checkWin(row, col) {
                g.Winner = g.Current
                gamesPlayed++ // ✅ incrémente le compteur de parties
                if g.Winner == 1 {
                    scoreRed++
                } else if g.Winner == 2 {
                    scoreYellow++
                }
            } else if g.isBoardFull() {
                // ✅ si la grille est pleine et pas de gagnant → égalité
                g.Winner = 0
                gamesPlayed++
                draws++
            } else {
                // changement de joueur
                if g.Current == 1 {
                    g.Current = 2
                } else {
                    g.Current = 1
                }
            }
            return true
        }
    }
    return false
}

// Vérifie si le coup est gagnant
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

// Compte les pions alignés dans une direction
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
func GetScores() (int, int, int, int) {
    return scoreRed, scoreYellow, gamesPlayed, draws
}

// Réinitialise les scores
func ResetScores() {
    scoreRed = 0
    scoreYellow = 0
    gamesPlayed = 0
    draws = 0
}

