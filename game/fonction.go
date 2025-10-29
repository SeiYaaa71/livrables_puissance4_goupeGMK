package game

import (
    "math/rand"
)

const (
    Rows = 6
    Cols = 7
)

type Stats struct {
    Red    int
    Yellow int
    Games  int
    Draws  int
}

var GlobalStats Stats

type Game struct {
    Grid    [Rows][Cols]int
    Current int
    Winner  int
}

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

func NewGame() *Game {
    return &Game{Current: 1}
}

func (g *Game) switchPlayer() {
    if g.Current == 1 {
        g.Current = 2
    } else {
        g.Current = 1
    }
}

func (g *Game) Play(col int) (bool, string) {
    if col < 0 || col >= Cols || g.Winner != 0 {
        return false, "Partie terminée refaite nouvelle partie🔄"
    }

    for row := Rows - 1; row >= 0; row-- {
        if g.Grid[row][col] == 0 {
            g.Grid[row][col] = g.Current

            if g.checkWin(row, col) {
                g.Winner = g.Current
                GlobalStats.Games++
                if g.Winner == 1 {
                    GlobalStats.Red++
                } else {
                    GlobalStats.Yellow++
                }

                if g.Winner == 1 {
                    return true, "🎉 Le Joueur 1 🔴 a gagné ! 🏆"
                } else {
                    return true, "🎉 Le Joueur 2 🟡 a gagné ! 🏆"
                }
            }

            if g.isBoardFull() {
                GlobalStats.Games++
                GlobalStats.Draws++
                return true, "🤝 Match nul !"
            }

            g.switchPlayer()

            msg := encouragements[rand.Intn(len(encouragements))]
            return true, msg
        }
    }
    return false, "⚠️ Colonne pleine"
}

func (g *Game) checkWin(row, col int) bool {
    player := g.Grid[row][col]
    if player == 0 {
        return false
    }

    directions := [][2]int{
        {0, 1},
        {1, 0},
        {1, 1},
        {1, -1},
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

func (g *Game) isBoardFull() bool {
    for c := 0; c < Cols; c++ {
        if g.Grid[0][c] == 0 {
            return false
        }
    }
    return true
}

func (g *Game) Reset() {
    for r := 0; r < Rows; r++ {
        for c := 0; c < Cols; c++ {
            g.Grid[r][c] = 0
        }
    }
    g.Current = 1
    g.Winner = 0
}

func GetScores() Stats {
    return GlobalStats
}

func ResetScores() {
    GlobalStats = Stats{}
}
