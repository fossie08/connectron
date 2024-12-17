package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Game struct {
	Grid        [][]int        // 2D grid to represent the game board
	Players     int            // Number of players
	WinLength   int            // Line length needed to win
	CurrentTurn int            // Current player's turn
	Colors      []color.RGBA   // Player colors
}

// NewGame initializes the game with a grid and players
func NewGame(gridWidth, gridHeight, players, winLength int) *Game {
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
		Grid:        grid,
		Players:     players,
		WinLength:   winLength,
		CurrentTurn: 0,
		Colors:      playerColors[:players],
	}
}

// DropCounter method to drop a counter into a column
func (g *Game) DropCounter(column int) (row int, success bool) {
	for i := len(g.Grid) - 1; i >= 0; i-- {
		if g.Grid[i][column] == -1 {
			g.Grid[i][column] = g.CurrentTurn
			return i, true
		}
	}
	return -1, false
}

// CheckWin checks if the current player has won after placing a counter
func (g *Game) CheckWin(row, column int) bool {
	player := g.Grid[row][column]
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}} // Horizontal, Vertical, Diagonal
	for _, dir := range directions {
		count := 1
		// Check in both directions
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

// Show the main game window
func MainGameWindow(gw *Game, connectronApp fyne.App) {
	gameWindow := connectronApp.NewWindow("Connectron - Game")
	gameWindow.Resize(fyne.NewSize(800, 600))

	// Info panel at the top
	infoLabel := widget.NewLabel(fmt.Sprintf("Players: %d | Win Line: %d", gw.Players, gw.WinLength))

	// Create a container to hold the column buttons (changed names)
	buttonContainer := container.NewHBox()

	// Create a container to hold the column canvases (representing the grid)
	gridContainer := container.NewWithoutLayout()

	// Dimensions for the grid (we will rotate the whole thing)
	tileWidth := float32(100)
	tileHeight := float32(400 / len(gw.Grid)) // Adjust based on the number of rows

	// Add buttons and canvases for each column (rotate the grid)
	for j := 0; j < len(gw.Grid[0]); j++ {
		// Create column button (updated text)
		columnButton := widget.NewButton(fmt.Sprintf("Drop in Column %d", j+1), func(col int) func() {
			return func() {
				if row, success := gw.DropCounter(col); success {
					// Update the canvas with the color of the current player's turn
					cell := gridContainer.Objects[row*len(gw.Grid[0])+col].(*canvas.Circle)
					cell.FillColor = gw.Colors[gw.CurrentTurn]
					cell.Refresh()

					// Check for a win condition
					if gw.CheckWin(row, col) {
						infoLabel.SetText(fmt.Sprintf("Player %d Wins!", gw.CurrentTurn+1))
					} else {
						gw.CurrentTurn = (gw.CurrentTurn + 1) % gw.Players
						infoLabel.SetText(fmt.Sprintf("Player %d's Turn", gw.CurrentTurn+1))
					}
				}
			}
		}(j))
		buttonContainer.Add(columnButton)

		// Create column canvases (colored circles representing cells in each column)
		for i := 0; i < len(gw.Grid); i++ {
			cell := canvas.NewCircle(color.RGBA{240, 240, 240, 255})
			cell.Resize(fyne.NewSize(tileWidth, tileHeight))
			cell.Move(fyne.NewPos(tileWidth*float32(i), tileHeight*float32(j))) // Swap i and j for rotation
			gridContainer.Add(cell)
		}
	}

	// Layout for the window (rotate the whole layout 90 degrees clockwise)
	content := container.NewVBox(
		infoLabel,
		buttonContainer,
		gridContainer,
	)

	// Set the content of the window
	gameWindow.SetContent(content)

	// Show the window
	gameWindow.Show()
}
