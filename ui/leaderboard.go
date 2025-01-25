package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sort"
	"strconv"
)

// CreateLeaderboard creates a leaderboard UI from the provided 2D array of player data
func CreateLeaderboard(playerData [][]string) fyne.CanvasObject {
	// Handle the case where no player data is provided
	if len(playerData) == 0 {
		return container.NewVBox(widget.NewLabel("No data available"))
	}

	// Ensure all rows have the same number of columns
	numRows := len(playerData)
	numCols := len(playerData[0])

	for _, row := range playerData {
		if len(row) != numCols {
			return container.NewVBox(widget.NewLabel("Error: Inconsistent data"))
		}
	}

	// Create the table to display the leaderboard data
	list := widget.NewTable(
		func() (int, int) { return numRows, numCols },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(playerData[id.Row][id.Col])
			if id.Row == 0 {
				// Bold headers
				o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
			}
		},
	)

	// Sorting buttons
	sortByScoreButton := widget.NewButton("Sort by Score", func() {
		sort.Slice(playerData[1:], func(i, j int) bool {
			scoreI, _ := strconv.Atoi(playerData[i+1][1]) // +1 to skip header
			scoreJ, _ := strconv.Atoi(playerData[j+1][1])
			return scoreI > scoreJ // Sort by score (descending)
		})
		list.Refresh() // Refresh the list with the updated playerData
	})

	sortByNameButton := widget.NewButton("Sort by Name", func() {
		sort.Slice(playerData[1:], func(i, j int) bool {
			return playerData[i+1][0] < playerData[j+1][0] // Sort by name (alphabetical)
		})
		list.Refresh() // Refresh the list with the updated playerData
	})

	sortByUUIDButton := widget.NewButton("Sort by UUID", func() {
		sort.Slice(playerData[1:], func(i, j int) bool {
			return playerData[i+1][4] < playerData[j+1][4] // Sort by UUID (alphabetical)
		})
		list.Refresh() // Refresh the list with the updated playerData
	})

	// Create a toolbar
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MediaReplayIcon(), func() {
			list.Refresh() // Refresh the table
		}),
		widget.NewToolbarSeparator(),
	)

	// Display the table and controls
	content := container.NewBorder(
		container.NewVBox(toolbar, sortByNameButton, sortByScoreButton, sortByUUIDButton),
		nil, nil, nil,
		list,
	)

	return content
}


// calculateColumnWidths calculates the maximum width for each column based on the content
func calculateColumnWidths(data [][]string) []float32 {
	widths := make([]float32, len(data[0]))

	for _, row := range data {
		for col, cell := range row {
			// Measure the width of the cell text
			width := float32(len(cell) * 10) // Approximate width (10 pixels per character)
			if width > widths[col] {
				widths[col] = width
			}
		}
	}

	return widths
}