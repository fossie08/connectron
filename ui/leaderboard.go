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
func CreateLeaderboard(playerData [][]string) *fyne.Container {
	// Create a widget to show leaderboard data
	list := widget.NewTable(
		func() (int, int) { return len(playerData), 5 },
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

	// Creating a toolbar
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MediaReplayIcon(), func() {
			list.Refresh()
		}),
		widget.NewToolbarSeparator(),
	)

	// Setting column widths
	list.SetColumnWidth(0, 140)
	list.SetColumnWidth(1, 140)
	list.SetColumnWidth(2, 140)
	list.SetColumnWidth(3, 140)
	list.SetColumnWidth(4, 280)

	// Display list and sorting buttons
	content := container.NewBorder(
		container.NewHBox(toolbar, sortByNameButton, sortByScoreButton, sortByUUIDButton),
		nil, nil, nil,
		list,
	)

	return content
}