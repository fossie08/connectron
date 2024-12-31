package ui

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Game struct {
    Grid           [][]int
    Players        int
    WinLength      int
    CurrentTurn    int
    Colors         []color.RGBA
    PlayerTypes    []int // -1 for human 0+ for ai levels
    BestOf         int
    RoundCount     int
    CornerBonus    bool
    SolitaireRule  bool
    BombCounter    bool
    OverflowRule   bool
    AIForMissing   bool
}


func NewGame(gridWidth, gridHeight, players, winLength, bestOf int, playerTypes []int, aiForMissing, cornerBonus, solitaireRule, bombCounter, overflowRule bool) *Game {
    grid := make([][]int, gridHeight)
    for i := range grid {
        grid[i] = make([]int, gridWidth)
        for j := range grid[i] {
            grid[i][j] = -1
        }
    }

    playerColors := []color.RGBA{
        {255, 0, 0, 255},   // Red
        {0, 255, 0, 255},   // Green
        {0, 0, 255, 255},   // Blue
        {255, 255, 0, 255}, // Yellow
        {255, 0, 255, 255}, // Magenta
        {0, 255, 255, 255}, // Cyan
        {128, 0, 128, 255}, // Purple
        {255, 165, 0, 255}, // Orange
        {128, 128, 128, 255}, // Gray
        {0, 128, 128, 255}, // Teal
    }

    return &Game{
        Grid:           grid,
        Players:        players,
        WinLength:      winLength,
        CurrentTurn:    0,
        Colors:         playerColors[:players],
        PlayerTypes:    playerTypes,
        BestOf:         bestOf,
        RoundCount:     0,
        CornerBonus:    cornerBonus,
        SolitaireRule:  solitaireRule,
        BombCounter:    bombCounter,
        OverflowRule:   overflowRule,
        AIForMissing:   aiForMissing,
    }
}

func (g *Game) CheckCornerBonus(row, col int) {
    if (row == 0 || row == len(g.Grid)-1) && (col == 0 || col == len(g.Grid[0])-1) && g.CornerBonus {
        if g.Grid[row][col] != -1 {
            // Apply the corner bonus (2 counters, 3 if win length is 7 or more)
            //bonus := 2
            if g.WinLength >= 7 {
                //bonus = 3
            }
            // Apply bonus logic (for simplicity, assume we apply the bonus to the current player)
            // You can expand this logic depending on your needs
        }
    }
}


func (g *Game) DropCounter(column int) (int, bool) {
	if column < 0 || column >= len(g.Grid[0]) {
		return -1, false // Invalid column
	}
	for i := len(g.Grid) - 1; i >= 0; i-- {
		if g.Grid[i][column] == -1 {
			g.Grid[i][column] = g.CurrentTurn
			return i, true
		}
	}
	return -1, false // Column is full
}

func (g *Game) CheckWin(row, column int) bool {
	player := g.Grid[row][column]
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, dir := range directions {
		count := 1
		for _, sign := range []int{-1, 1} {
			r, c := row, column
			for {
				r += dir[0] * sign
				c += dir[1] * sign
				if r < 0 || r >= len(g.Grid) || c < 0 || c >= len(g.Grid[0]) || g.Grid[r][c] != player {
					break
				}
				count++
			}
		}
		if count >= g.WinLength {
			return true
		}
	}
	return false
}

func (g *Game) IsFull() bool {
	for _, row := range g.Grid {
		for _, cell := range row {
			if cell == -1 {
				return false
			}
		}
	}
	return true
}

// AI Strategies
func (g *Game) GetAIColumn(aiType int) (int, int) {
	switch aiType {
	case 1: // EasyAI
		return g.easyAI()
	case 2: // MediumAI
		return g.mediumAI()
	case 3: // HardAI
		return g.hardAI()
	default:
		return g.easyAI()
	}
}

// EasyAI - Random move
func (g *Game) easyAI() (int, int) {
	for {
		column := rand.Intn(len(g.Grid[0]))
		if row, success := g.DropCounter(column); success {
			return column, row
		}
	}
}

// MediumAI - Block or Win strategy
func (g *Game) mediumAI() (int, int) {
	// Try to win or block the player
	for col := 0; col < len(g.Grid[0]); col++ {
		// Check if we can win or block in this column
		if row, success := g.DropCounter(col); success {
			if g.CheckWin(row, col) {
				return col, row // If AI wins
			}
			g.Grid[row][col] = -1 // Undo move

			// Block player's winning move
			g.Grid[row][col] = 2 // Assuming 2 is the opponent's player
			if g.CheckWin(row, col) {
				g.Grid[row][col] = -1 // Undo block
				return col, row
			}
			g.Grid[row][col] = -1 // Undo move
		}
	}

	// If no win or block, return a random move
	return g.easyAI()
}

// HardAI - Minimax with Alpha-Beta pruning
func (g *Game) hardAI() (int, int) {
	bestMove := g.minimax(4, -10000, 10000, true)
	return bestMove[0], bestMove[1]
}

// Minimax algorithm with alpha-beta pruning
func (g *Game) minimax(depth, alpha, beta int, isMaximizing bool) [2]int {
	var bestMove [2]int

	if depth == 0 {
		return bestMove // Base case: return score
	}

	for col := 0; col < len(g.Grid[0]); col++ {
		if row, success := g.DropCounter(col); success {
			// Calculate score of this move
			score := g.evaluateBoard(isMaximizing)

			// Maximize or minimize the score based on the current player
			if isMaximizing {
				if score > alpha {
					alpha = score
					bestMove = [2]int{col, row}
				}
			} else {
				if score < beta {
					beta = score
					bestMove = [2]int{col, row}
				}
			}

			// Undo move
			g.Grid[row][col] = -1
		}
	}

	return bestMove
}

func (g *Game) evaluateBoard(isMaximizing bool) int {
	// Evaluate the current board. Positive score for AI, negative for player.
	if isMaximizing {
		// AI's perspective
		return 1
	}
	// Player's perspective
	return -1
}

// Main game window
func MainGameWindow(gw *Game, connectronApp fyne.App) {
	gameWindow := connectronApp.NewWindow("Connectron - Game")
	infoLabel := widget.NewLabel("Game Start!")
	gameWindow.SetFullScreen(true)

	gridContainer := container.NewGridWithColumns(len(gw.Grid[0]))
	for j := 0; j < len(gw.Grid[0]); j++ {
		for i := 0; i < len(gw.Grid); i++ {
			cell := canvas.NewCircle(color.RGBA{240, 240, 240, 255})
			gridContainer.Add(cell)
		}
	}

	var processTurn func(column int, row int) bool
	processTurn = func(column int, row int) bool {
		if row == -1 {
			row, _ = gw.DropCounter(column)
		}
		if row == -1 {
			infoLabel.SetText("Column is full!")
			return false
		}

		cell := gridContainer.Objects[row*len(gw.Grid[0])+column].(*canvas.Circle)
		cell.FillColor = gw.Colors[gw.CurrentTurn]
		cell.Refresh()

		if gw.CheckWin(row, column) {
			infoLabel.SetText(fmt.Sprintf("Player %d Wins!", gw.CurrentTurn+1))
			return true
		}

		if gw.IsFull() {
			infoLabel.SetText("The game is a draw!")
			return true
		}

		gw.CurrentTurn = (gw.CurrentTurn + 1) % gw.Players
		infoLabel.SetText(fmt.Sprintf("Player %d's Turn", gw.CurrentTurn+1))

		if gw.PlayerTypes[gw.CurrentTurn] != -1 {
			aiColumn, aiRow := gw.GetAIColumn(gw.PlayerTypes[gw.CurrentTurn])
			time.AfterFunc(10*time.Millisecond, func() {
				processTurn(aiColumn, aiRow)
			})
		}
		return false
	}

	columnEntry := widget.NewEntry()
	columnEntry.SetPlaceHolder("Enter Column")
	dropButton := widget.NewButton("Drop", func() {
		if gw.PlayerTypes[gw.CurrentTurn] == -1 {
			col, err := strconv.Atoi(columnEntry.Text)
			if err != nil || col < 1 || col > len(gw.Grid[0]) {
				infoLabel.SetText("Invalid column number!")
				return
			}
			if processTurn(col-1, -1) {
				columnEntry.SetText("")
			}
		} else {
			infoLabel.SetText("It's not your turn!")
		}
	})

	content := container.NewBorder(
		container.NewVBox(infoLabel, columnEntry, dropButton),
		nil, nil, nil, gridContainer,
	)

	gameWindow.Resize(fyne.NewSize(800, 600))
	gameWindow.SetContent(content)
	gameWindow.Show()

	if gw.PlayerTypes[gw.CurrentTurn] != -1 {
		aiColumn, aiRow := gw.GetAIColumn(gw.PlayerTypes[gw.CurrentTurn])
		time.AfterFunc(100*time.Millisecond, func() {
			processTurn(aiColumn, aiRow)
		})
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
