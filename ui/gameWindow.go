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
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
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
	EnableAlliances bool
	Winners		   []int
	GridHistory    [][][]int
	BombCounters   []bool
	Alliances	   [][]string
}


func NewGame(gridWidth, gridHeight, players, winLength, roundCounter, bestOf int, playerTypes []int, aiForMissing, cornerBonus, solitaireRule, bombCounter, overflowRule, enableAlliances bool, alliances [][]string) *Game {
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
		RoundCount:     roundCounter,
		CornerBonus:    cornerBonus,
		SolitaireRule:  solitaireRule,
		BombCounter:    bombCounter,
		OverflowRule:   overflowRule,
		AIForMissing:   aiForMissing,
		EnableAlliances: enableAlliances,
		BombCounters:   make([]bool, players),
		Alliances: 	alliances,
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

    // Helper function to check if two players are in the same alliance
	inSameAlliance := func(player1, player2 int) bool {
		if player1 == -1 || player2 == -1 {
			return false // Blank circles are not in any alliance
		}
		for _, alliance := range g.Alliances {
			inAlliance1 := false
			inAlliance2 := false
			for _, p := range alliance {
				if p == fmt.Sprintf("Player-%d", player1+1) {
					inAlliance1 = true
				}
				if p == fmt.Sprintf("Player-%d", player2+1) {
					inAlliance2 = true
				}
			}
			if inAlliance1 && inAlliance2 {
				return true
			}
		}
		return false
	}

    for _, dir := range directions {
        count := 1
        for _, sign := range []int{-1, 1} {
            r, c := row, column
            for {
                r += dir[0] * sign
                c += dir[1] * sign
                if r < 0 || r >= len(g.Grid) || c < 0 || c >= len(g.Grid[0]) {
                    break
                }
				fmt.Printf("In same alliance: %d, %d, %t\n", player, g.Grid[r][c], inSameAlliance(player, g.Grid[r][c])) // Debugging statement
                if g.Grid[r][c] != player && (!g.EnableAlliances || !inSameAlliance(player, g.Grid[r][c])) {
                    break
                }
                count++
                fmt.Printf("Count: %d, Player: %d, Position: (%d, %d)\n", count, player, r, c) // Debugging statement
                // Apply corner bonus
                if g.CornerBonus {
                    if (r == 0 || r == len(g.Grid)-1) && (c == 0 || c == len(g.Grid[0])-1) {
                        if g.WinLength >= 7 {
                            count += 2 // 3 counters for win length 7 or more
                        } else {
                            count++ // 2 counters for win length less than 7
                        }
                    }
                }
            }
        }
        if count >= g.WinLength {
            fmt.Printf("Winning condition met. Count: %d, WinLength: %d\n", count, g.WinLength) // Debugging statement
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

func (g *Game) CheckCornerBonus(row, col int) {
	fmt.Println(g.CornerBonus)
	if !g.CornerBonus {
		return // Exit if the corner bonus is not enabled
	}
    if (row == 0 || row == len(g.Grid)-1) && (col == 0 || col == len(g.Grid[0])-1) {
        bonus := 2
        if g.WinLength >= 7 {
            bonus = 3
        }
        player := g.Grid[row][col]
        if player >= 0 {
            fmt.Printf("Player %d gets a corner bonus of %d points!\n", player, bonus)
        }
    }
}

func (g *Game) CheckSolitaire() {
	if !g.SolitaireRule {
		fmt.Println("Solitaire rule is not enabled.")
		return // Exit if the solitaire rule is not enabled
	}

	for row := 0; row < len(g.Grid); row++ {
		for col := 0; col < len(g.Grid[0]); col++ {
			player := g.Grid[row][col]
			if player == -1 {
				continue // Skip empty cells
			}

			// Check if all neighbors belong to the same player
			neighborPlayer := -1
			surroundedBySamePlayer := true
			for _, dir := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
				r, c := row+dir[0], col+dir[1]
				if r >= 0 && r < len(g.Grid) && c >= 0 && c < len(g.Grid[0]) {
					neighbor := g.Grid[r][c]
					if neighbor == -1 || (neighborPlayer != -1 && neighbor != neighborPlayer) {
						surroundedBySamePlayer = false
						break
					}
					neighborPlayer = neighbor
				} else {
					surroundedBySamePlayer = false // Out-of-bound neighbors do not count
					break
				}
			}

			if surroundedBySamePlayer {
				// Remove the solitaire counter
				for r := row; r > 0; r-- {
					g.Grid[r][col] = g.Grid[r-1][col]
				}
				g.Grid[0][col] = -1 // Set the top cell to empty

				// Reset the loop to re-check the updated grid
				row = 0
				col = -1 // This ensures the outer loop resets correctly after modification
				break
			}
		}
	}
}

func (g *Game) UseBombCounter(row, col int) {
    if !g.BombCounter {
		return // Exit if the bomb counter is not enabled
    }
	for _, dir := range [][2]int{{0, 0}, {0, 1}, {1, 0}, {0, -1}, {-1, 0}, {1, 1}, {1, -1}, {-1, -1}, {-1, 1}} {
		r, c := row+dir[0], col+dir[1]
		if r >= 0 && r < len(g.Grid) && c >= 0 && c < len(g.Grid[0]) {
			g.Grid[r][c] = -1
		}
	}
}

func (g *Game) CheckOverflow(column int) {
	if g.OverflowRule && len(g.Grid) >= 6 {
		full := true
		for _, cell := range g.Grid {
			if cell[column] == -1 {
				full = false
				break
			}
		}
		if full {
			// Drop a counter in the left adjacent column if possible
			if column > 0 {
				row, success := g.DropCounter(column - 1)
				if success {
					// Set the color of the newly dropped counter
					g.Grid[row][column-1] = g.CurrentTurn
				}
			}
			// Drop a counter in the right adjacent column if possible
			if column < len(g.Grid[0])-1 {
				row, success := g.DropCounter(column + 1)
				if success {
					// Set the color of the newly dropped counter
					g.Grid[row][column+1] = g.CurrentTurn
				}
			}
		}
	}
} 

// Main game window
func MainGameWindow(gw *Game, connectronApp fyne.App) {
    gameWindow := connectronApp.NewWindow("Connectron - Game")
    infoLabel := widget.NewLabel("Game Start!")
    //gameWindow.SetFullScreen(true)

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

        // Apply special rules
        gw.CheckCornerBonus(row, column)
        gw.CheckSolitaire()
        gw.CheckOverflow(column)

        // Update the UI for the newly added counters
        for i := 0; i < len(gw.Grid); i++ {
            for j := 0; j < len(gw.Grid[0]); j++ {
                cell := gridContainer.Objects[i*len(gw.Grid[0])+j].(*canvas.Circle)
                if gw.Grid[i][j] != -1 {
                    cell.FillColor = gw.Colors[gw.Grid[i][j]]
                } else {
                    cell.FillColor = color.RGBA{240, 240, 240, 255} // Default color for empty cells
                }
                cell.Refresh()
            }
        }

        // Check for win
        if gw.CheckWin(row, column) {
            infoLabel.SetText(fmt.Sprintf("Player %d Wins!", gw.CurrentTurn+1))
            
            // Record the winner
            gw.Winners = append(gw.Winners, gw.CurrentTurn+1)
            gw.GridHistory = append(gw.GridHistory, copyGrid(gw.Grid)) // Store current grid state
            fmt.Println("A player won: Round count", gw.RoundCount, "Best of", gw.BestOf)
            if gw.RoundCount+1 < gw.BestOf {
                gw.RoundCount++
                // Start a new game
                nextGame := NewGame(len(gw.Grid[0]), len(gw.Grid), gw.Players, gw.WinLength, gw.RoundCount, gw.BestOf, gw.PlayerTypes, gw.AIForMissing, gw.CornerBonus, gw.SolitaireRule, gw.BombCounter, gw.OverflowRule, gw.EnableAlliances, gw.Alliances)
                MainGameWindow(nextGame, connectronApp)
                gameWindow.Close()
            } else {
                // Show results window
                ShowResultsWindow(gw, connectronApp)
                gameWindow.Close()
            }
            return true
        }

        // Check for draw
        if gw.IsFull() {
            fmt.Println("Draw: Round count", gw.RoundCount, "Best of", gw.BestOf)
            infoLabel.SetText("The game is a draw!")
            gw.Winners = append(gw.Winners, 0) // 0 indicates a draw
            gw.GridHistory = append(gw.GridHistory, copyGrid(gw.Grid))

            if gw.RoundCount+1 < gw.BestOf {
                gw.RoundCount++
                gameWindow.Close()
                // Start a new game
                nextGame := NewGame(len(gw.Grid[0]), len(gw.Grid), gw.Players, gw.WinLength, gw.RoundCount, gw.BestOf, gw.PlayerTypes, gw.AIForMissing, gw.CornerBonus, gw.SolitaireRule, gw.BombCounter, gw.OverflowRule, gw.EnableAlliances, gw.Alliances)
                MainGameWindow(nextGame, connectronApp)
            } else {
                gameWindow.Close()
                // Show results window
                ShowResultsWindow(gw, connectronApp)
            }
            return true
        }

        gw.CurrentTurn = (gw.CurrentTurn + 1) % gw.Players
        infoLabel.SetText(fmt.Sprintf("Player %d's Turn", gw.CurrentTurn+1))

        // AI move handling
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
	
    bombButton := widget.NewButton("Use Bomb Counter", func() {
        if gw.BombCounters[gw.CurrentTurn] {
            infoLabel.SetText("You have already used your bomb counter!")
            return
        }
        col, err := strconv.Atoi(columnEntry.Text)
        if err != nil || col < 1 || col > len(gw.Grid[0]) {
            infoLabel.SetText("Invalid column number!")
            return
        }
        row, _ := gw.DropCounter(col - 1)
        if row == -1 {
            infoLabel.SetText("Column is full!")
            return
        }
        gw.UseBombCounter(row, col-1)
        gw.BombCounters[gw.CurrentTurn] = true
        infoLabel.SetText("Bomb counter used!")
        // Update the UI for the newly added counters
        for i := 0; i < len(gw.Grid); i++ {
            for j := 0; j < len(gw.Grid[0]); j++ {
                cell := gridContainer.Objects[i*len(gw.Grid[0])+j].(*canvas.Circle)
                if gw.Grid[i][j] != -1 {
                    cell.FillColor = gw.Colors[gw.Grid[i][j]]
                } else {
                    cell.FillColor = color.RGBA{240, 240, 240, 255} // Default color for empty cells
                }
                cell.Refresh()
            }
        }
    })

	if !gw.BombCounter {
		bombButton.Disable()
	}

    content := container.NewBorder(
        container.NewVBox(infoLabel, columnEntry, dropButton, bombButton),
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

func updateLeaderboard(gw *Game) {
	filePath := filepath.Join("files", "leaderboard.csv")
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening leaderboard file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading leaderboard file:", err)
		return
	}

	// Create a map to store player stats
	playerStats := make(map[string][]string)
	if len(records) > 1 {
		for _, record := range records[1:] { // Skip header
			playerStats[record[0]] = record
		}
	}

	// Update player stats based on game results
	for _, winner := range gw.Winners {
		if winner == 0 {
			// Draw case
			for _, record := range playerStats {
				drawCount, err := strconv.Atoi(record[5])
				if err == nil {
					record[5] = strconv.Itoa(drawCount + 1) // Increment Draw count
				}
			}
		} else {
			player := fmt.Sprintf("Player-%d", winner)
			if record, exists := playerStats[player]; exists {
				playedCount, err := strconv.Atoi(record[3])
				if err == nil {
					record[3] = strconv.Itoa(playedCount + 1) // Increment Played count
				}
				wonCount, err := strconv.Atoi(record[4])
				if err == nil {
					record[4] = strconv.Itoa(wonCount + 1) // Increment Won count
				}
			} else {
				playerStats[player] = []string{player, "0", "UUID", "1", "1", "0", "0"}
			}
		}
	}

	// Write updated stats back to the file
	file.Seek(0, 0)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write(records[0])

	// Write updated records
	for _, record := range playerStats {
		writer.Write(record)
	}
}

func ShowResultsWindow(gw *Game, connectronApp fyne.App) {
	updateLeaderboard(gw)

	resultsWindow := connectronApp.NewWindow("Series Results")
	resultsText := "Series Results:\n\n"
	for i, winner := range gw.Winners {
		if winner == 0 {
			resultsText += fmt.Sprintf("Game %d: Draw\n", i+1)
		} else {
			resultsText += fmt.Sprintf("Game %d: Player %d Wins\n", i+1, winner)
		}
	}

	// Sort players by their wins
	playerWins := make(map[int]int)
	for _, winner := range gw.Winners {
		if winner != 0 {
			playerWins[winner]++
		}
	}

	type playerResult struct {
		Player int
		Wins   int
	}

	var results []playerResult
	for player, wins := range playerWins {
		results = append(results, playerResult{Player: player, Wins: wins})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Wins > results[j].Wins
	})

	resultsText += "\nFinal Standings:\n"
	for i, result := range results {
		resultsText += fmt.Sprintf("%d. Player %d with %d wins\n", i+1, result.Player, result.Wins)
	}

	resultsLabel := widget.NewLabel(resultsText)
	closeButton := widget.NewButton("Close", func() {
		resultsWindow.Close()
	})

	resultsWindow.SetContent(container.NewVBox(resultsLabel, closeButton))
	resultsWindow.Resize(fyne.NewSize(400, 300))
	resultsWindow.Show()
}


// Utility function to copy the grid
func copyGrid(grid [][]int) [][]int {
	newGrid := make([][]int, len(grid))
	for i := range grid {
		newGrid[i] = make([]int, len(grid[i]))
		copy(newGrid[i], grid[i])
	}
	return newGrid
}

func init() {
	rand.Seed(time.Now().UnixNano())
}