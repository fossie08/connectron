package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"insighthub.uk/connectron/v2/ui"
)

const (
	EasyAI   = iota
	MediumAI
	HardAI
)

var alliances = map[string][]string{}
var unassigned []string

func main() {
	// Initialize the application
	connectronApp := app.New()
	connectronApp.Settings().SetTheme(theme.LightTheme())
	mainWindow := connectronApp.NewWindow("Connectron")

	// Set the window size
	mainWindow.Resize(fyne.Size{1200, 1000})
	mainWindow.CenterOnScreen()

	// Create menu items
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File"),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Settings", func() { ShowSettingsWindow(connectronApp) }),
		),
	)

	// Set the menu
	mainWindow.SetMainMenu(menu)

	// Grid Width Selection
	gridWidthLabel := widget.NewLabel("Grid Width (6-100):")
	gridWidthValue := widget.NewLabel("6")
	gridWidthSlider := widget.NewSlider(6, 100)
	gridWidthSlider.OnChanged = func(value float64) {
		gridWidthValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Grid Height Selection
	gridHeightLabel := widget.NewLabel("Grid Height (6-100):")
	gridHeightValue := widget.NewLabel("6")
	gridHeightSlider := widget.NewSlider(6, 100)
	gridHeightSlider.OnChanged = func(value float64) {
		gridHeightValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Line Length to Win
	lineLengthLabel := widget.NewLabel("Line Length to Win (4-10):")
	lineLengthValue := widget.NewLabel("4")
	lineLengthSlider := widget.NewSlider(4, 10)
	lineLengthSlider.OnChanged = func(value float64) {
		lineLengthValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Number of Players
	playerCountLabel := widget.NewLabel("Number of Players (1-10):")
	playerCountValue := widget.NewLabel("1")
	playerCountSlider := widget.NewSlider(1, 10)

	// Player Dropdowns Container
	playerDropdownsContainer := container.NewVBox()
	playerTypes := make([]int, 10)

	updatePlayerDropdowns := func(count int) {
		playerDropdownsContainer.RemoveAll()
		for i := 0; i < count; i++ {
			options := []string{"Easy AI", "Medium AI", "Hard AI", "Person"}
			dropdown := widget.NewSelect(options, func(selected string) {
				switch selected {
				case "Easy AI":
					playerTypes[i] = EasyAI
				case "Medium AI":
					playerTypes[i] = MediumAI
				case "Hard AI":
					playerTypes[i] = HardAI
				case "Person":
					playerTypes[i] = -1
				}
			})
			dropdown.SetSelected("Person")
			playerDropdownsContainer.Add(container.NewHBox(widget.NewLabel(fmt.Sprintf("Player %d:", i+1)), dropdown))
		}
		playerDropdownsContainer.Refresh()
	}

	updatePlayerDropdowns(int(playerCountSlider.Value))
	playerCountSlider.OnChanged = func(value float64) {
		playerCountValue.SetText(fmt.Sprintf("%d", int(value)))
		updatePlayerDropdowns(int(value))
	}

	// Missing player AI Configuration
	aiForMissingCheckbox := widget.NewCheck("AI for Missing Players", nil)

	// Special Rule Options
	cornerBonusCheckbox := widget.NewCheck("Enable Corner Bonus", nil)
	solitaireRuleCheckbox := widget.NewCheck("Enable Solitaire Destruction", nil)
	bombCounterCheckbox := widget.NewCheck("Enable Bomb Counter", nil)
	overflowRuleCheckbox := widget.NewCheck("Enable Overflow Rule", nil)

	// Alliance Rule
	allianceRuleCheckbox := widget.NewCheck("Enable Alliances Rule", nil)
	allianceSetupButton := widget.NewButton("Configure Alliances", func() {
		showAlliancesWindow(connectronApp, playerCountSlider)
	})

	// Layout containers
	gridSettings := container.NewVBox(
		gridWidthLabel, gridWidthSlider, gridWidthValue,
		gridHeightLabel, gridHeightSlider, gridHeightValue,
		lineLengthLabel, lineLengthSlider, lineLengthValue,
	)

	playerSettings := container.NewVBox(
		playerCountLabel, playerCountSlider, playerCountValue,
		playerDropdownsContainer,
		aiForMissingCheckbox,
	)

	ruleSettings := container.NewVBox(
		cornerBonusCheckbox,
		solitaireRuleCheckbox,
		bombCounterCheckbox,
		overflowRuleCheckbox,
		allianceRuleCheckbox,
		allianceSetupButton,
	)

	leftPane := container.NewVBox(
		widget.NewAccordion(
			widget.NewAccordionItem("Grid Settings", gridSettings),
			widget.NewAccordionItem("Player Settings", playerSettings),
			widget.NewAccordionItem("Rules", ruleSettings),
		),
	)

	// Start Game Button
	startGameButton := widget.NewButton("Start Game", func() {
		bestOfConverted, _ := strconv.Atoi("1")
		startGameSetup(int(gridWidthSlider.Value), int(gridHeightSlider.Value), int(lineLengthSlider.Value), int(playerCountSlider.Value), allianceRuleCheckbox.Checked, playerTypes, bestOfConverted, cornerBonusCheckbox.Checked, solitaireRuleCheckbox.Checked, bombCounterCheckbox.Checked, overflowRuleCheckbox.Checked, aiForMissingCheckbox.Checked)
	})

	// Main Tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Setup Game", leftPane),
	)

	mainWindow.SetContent(container.NewBorder(nil, startGameButton, nil, nil, tabs))
	mainWindow.ShowAndRun()
}

// ShowSettingsWindow creates a simple settings window
func ShowSettingsWindow(a fyne.App) {
	win := a.NewWindow("Settings")
	win.SetContent(container.NewVBox(
		widget.NewLabel("Settings Window"),
		widget.NewButton("Close", func() {
			win.Close()
		}),
	))
	win.Resize(fyne.NewSize(400, 300))
	win.Show()
}

// startGameSetup initiates the game setup based on selected settings
func startGameSetup(gridWidth, gridHeight, lineLength, playerCount int, enableAlliances bool, playerTypes []int, bestOf int, cornerBonus, solitaireRule, bombCounter, overflowRule, aiForMissing bool) {
	// Create and configure the game instance here (this part is a placeholder)
	game := ui.NewGame(gridWidth, gridHeight, playerCount, lineLength, 0, bestOf, playerTypes, aiForMissing, cornerBonus, solitaireRule, bombCounter, overflowRule)

	// Display the main game window
	ui.MainGameWindow(game, fyne.CurrentApp())
}

func showAlliancesWindow(a fyne.App, playerCountSlider *widget.Slider) {
	win := a.NewWindow("Configure Alliances")

	// Generate players dynamically based on the playerCountSlider value
	playerCount := int(playerCountSlider.Value)
	var players []string

	// Create players but only add them if they haven't been generated before
	for i := 0; i < playerCount; i++ {
		playerName := fmt.Sprintf("Player-%d", i+1)

		// Check if the player already exists in the alliances or unassigned lists
		if !playerExists(playerName) {
			players = append(players, playerName)
		}
	}

	// Update unassigned players only with new players
	unassigned = append([]string{}, players...)

	// Create the alliance manager window with the player count slider
	allianceManagerWindow := CreateAllianceManagerWindow(playerCountSlider)

	// Set the content of the window to the alliance manager
	win.SetContent(allianceManagerWindow)
	win.Resize(fyne.NewSize(600, 400))
	win.Show()
}

// Helper function to check if a player already exists
func playerExists(playerName string) bool {
	// Check in alliances map
	for _, alliedPlayers := range alliances {
		for _, alliedPlayer := range alliedPlayers {
			if alliedPlayer == playerName {
				return true
			}
		}
	}

	// Check in unassigned players list
	for _, unassignedPlayer := range unassigned {
		if unassignedPlayer == playerName {
			return true
		}
	}

	return false
}


// CreateAllianceManagerWindow creates the alliance manager UI with dynamic player assignment
func CreateAllianceManagerWindow(playerCountSlider *widget.Slider) fyne.CanvasObject {
	// Variables to track the selected item in lists
	var selectedUnassignedItem int
	var selectedAllianceItem int

	// Unassigned players list
	unassignedList := widget.NewList(
		func() int { return len(unassigned) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(unassigned[id])
		},
	)

	// Set the OnSelected function to track the selected item in unassignedList
	unassignedList.OnSelected = func(id widget.ListItemID) {
		selectedUnassignedItem = id
	}

	// Container to hold alliance lists
	alliancesContainer := container.NewVBox()

	// Function to refresh all lists
	refreshLists := func() {
		unassignedList.Refresh()
		alliancesContainer.Refresh()
	}
	

	// Load existing alliances into the window
	for allianceName, playersInAlliance := range alliances {
		// Create a list for each existing alliance
		allianceList := widget.NewList(
			func() int { return len(playersInAlliance) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(id widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(playersInAlliance[id])
			},
		)

		// Track the selected item in allianceList
		allianceList.OnSelected = func(id widget.ListItemID) {
			selectedAllianceItem = id
		}

		// Drag-and-drop buttons for the alliance
		moveToAlliance := widget.NewButton("→ Assign", func() {
			if selectedUnassignedItem >= 0 && selectedUnassignedItem < len(unassigned) {
				player := unassigned[selectedUnassignedItem]
				unassigned = append(unassigned[:selectedUnassignedItem], unassigned[selectedUnassignedItem+1:]...)
				alliances[allianceName] = append(alliances[allianceName], player)
				refreshLists()
			}
		})

		moveToUnassigned := widget.NewButton("← Unassign", func() {
			if selectedAllianceItem >= 0 && selectedAllianceItem < len(alliances[allianceName]) {
				player := alliances[allianceName][selectedAllianceItem]
				alliances[allianceName] = append(alliances[allianceName][:selectedAllianceItem], alliances[allianceName][selectedAllianceItem+1:]...)
				unassigned = append(unassigned, player)
				refreshLists()
			}
		})

		// Add the new alliance to the container
		alliancesContainer.Add(container.NewVBox(
			widget.NewLabel(allianceName),
			container.NewVBox(
				allianceList,
				container.NewHBox(moveToAlliance, moveToUnassigned),
			),
		))
	}


	// Create a new alliance
	newAllianceButton := widget.NewButton("Add Alliance", func() {
		allianceName := fmt.Sprintf("Alliance-%d", len(alliances)+1)
		alliances[allianceName] = []string{}

		// Create a list for the new alliance
		allianceList := widget.NewList(
			func() int { return len(alliances[allianceName]) },
			func() fyne.CanvasObject { return widget.NewLabel("") },
			func(id widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(alliances[allianceName][id])
			},
		)

		// Track the selected item in allianceList
		allianceList.OnSelected = func(id widget.ListItemID) {
			selectedAllianceItem = id
		}

		// Drag-and-drop buttons for the alliance
		moveToAlliance := widget.NewButton("→ Assign", func() {
			if selectedUnassignedItem >= 0 && selectedUnassignedItem < len(unassigned) {
				player := unassigned[selectedUnassignedItem]
				unassigned = append(unassigned[:selectedUnassignedItem], unassigned[selectedUnassignedItem+1:]...)
				alliances[allianceName] = append(alliances[allianceName], player)
				refreshLists()
			}
		})

		moveToUnassigned := widget.NewButton("← Unassign", func() {
			if selectedAllianceItem >= 0 && selectedAllianceItem < len(alliances[allianceName]) {
				player := alliances[allianceName][selectedAllianceItem]
				alliances[allianceName] = append(alliances[allianceName][:selectedAllianceItem], alliances[allianceName][selectedAllianceItem+1:]...)
				unassigned = append(unassigned, player)
				refreshLists()
			}
		})

		// Add the new alliance to the container
		alliancesContainer.Add(container.NewVBox(
			widget.NewLabel(allianceName),
			container.NewVBox(
				allianceList,
				container.NewHBox(moveToAlliance, moveToUnassigned),
			),
		))
	})

	// Confirm button to save the alliances into a 2D array
	confirmButton := widget.NewButton("Confirm Alliances", func() {
		// Save alliances to a 2D array (flattening alliances map into a 2D array)
		var allianceArray [][]string
		for _, playersInAlliance := range alliances {
			allianceArray = append(allianceArray, playersInAlliance)
		}

		// Optionally, save `alliances` globally to persist across windows
		fmt.Println("Alliances confirmed:", allianceArray)
	})

	// Layout for the entire alliance manager window
	content := container.NewBorder(
		container.NewVBox(newAllianceButton, confirmButton), // Include the confirm button
		nil,
		container.NewVBox(widget.NewLabel("Unassigned Players"), unassignedList),
		nil,
		alliancesContainer,
	)

	return content
}
