package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"insighthub.uk/connectron/v2/settings"
	"insighthub.uk/connectron/v2/ui"
)

const (
	EasyAI   = iota
	MediumAI
	HardAI
)

func main() {
	// Initialize the application
	connectronApp := app.New()
	connectronApp.Settings().SetTheme(theme.LightTheme())
	mainWindow := connectronApp.NewWindow("Connectron")

	// Set the window size
	mainWindow.Resize(fyne.Size{1200,1000})
	mainWindow.CenterOnScreen()

	// Create menu items
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			//fyne.NewMenuItem("New Game", func() { ui.SetupWindow(connectronApp) }),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Settings", func() { settings.ShowSettingsWindow(connectronApp, "1.0.0") }),
		),
	)

	// Set the menu
	mainWindow.SetMainMenu(menu)

	// Grid Width Selection
	gridWidthLabel := widget.NewLabel("Grid width (6-100):")
	gridWidthValue := widget.NewLabel("6") // Label to display the current slider value
	gridWidthSlider := widget.NewSlider(6, 100)
	gridWidthSlider.OnChanged = func(value float64) {
		gridWidthValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Grid Height Selection
	gridHeightLabel := widget.NewLabel("Grid height (6-100):")
	gridHeightValue := widget.NewLabel("6") // Label to display the current slider value
	gridHeightSlider := widget.NewSlider(6, 100)
	gridHeightSlider.OnChanged = func(value float64) {
		gridHeightValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Line Length to Win
	lineLengthLabel := widget.NewLabel("Line Length to Win (4-10):")
	lineLengthValue := widget.NewLabel("4") // Default line length
	lineLengthSlider := widget.NewSlider(4, 10)
	lineLengthSlider.OnChanged = func(value float64) {
		lineLengthValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Number of Players
	playerCountLabel := widget.NewLabel("Number of Players (0-10):")
	playerCountValue := widget.NewLabel("0") // Label to display the current slider value
	playerCountSlider := widget.NewSlider(0, 10)

	// Player Dropdowns Container
	playerDropdownsContainer := container.NewVBox()
	playerTypes := make([]int, 10) // Track player types

	// Update player dropdowns based on the number of players
	updatePlayerDropdowns := func(count int) {
		playerDropdownsContainer.RemoveAll() // Clear existing dropdowns
		for i := 0; i < count; i++ {
			options := []string{"Easy AI", "Medium AI", "Hard AI", "Person"}
			dropdown := widget.NewSelect(options, func(selected string) {
				fmt.Printf("Player %d type set to: %s\n", i+1, selected)
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
			playerDropdownsContainer.Add(dropdown)
		}
		playerDropdownsContainer.Refresh() // Refresh the container to show new dropdowns
	}

	playerCountSlider.OnChanged = func(value float64) {
		playerCountValue.SetText(fmt.Sprintf("%d", int(value)))
		updatePlayerDropdowns(int(value))
	}

	// AI/Player Configuration
	aiForMissingCheckbox := widget.NewCheck("AI for Missing Players", nil)
	bestOfValue := widget.NewLabel("Best of: 1")

	// Define the options for the radio buttons
	options := []string{"1", "3", "5", "7"}

	// Create a radio group with the options
	bestOfRadioGroup := widget.NewRadioGroup(options, func(selected string) {
		bestOfValue.SetText(fmt.Sprintf("Best of: %s", selected))
	})

	// Set the default selection (optional)
	bestOfRadioGroup.SetSelected("1") // Set default to "3"

	// Special Rule Options
	cornerBonusCheckbox := widget.NewCheck("Enable Corner Bonus", nil)
	solitaireRuleCheckbox := widget.NewCheck("Enable Solitaire Destruction", nil)
	bombCounterCheckbox := widget.NewCheck("Enable Bomb Counter", nil)
	overflowRuleCheckbox := widget.NewCheck("Enable Overflow Rule", nil)

	startGameButton := widget.NewButton("Start Game", func() {
		bestofConverted, _ := strconv.Atoi(bestOfRadioGroup.Selected)
		startGameSetup(int(gridWidthSlider.Value), int(gridHeightSlider.Value), int(lineLengthSlider.Value), int(playerCountSlider.Value), false, playerTypes, bestofConverted, cornerBonusCheckbox.Checked, solitaireRuleCheckbox.Checked, bombCounterCheckbox.Checked, overflowRuleCheckbox.Checked, aiForMissingCheckbox.Checked)
	})

	leftcontent := container.NewVBox(
		gridWidthLabel,
		gridWidthSlider, gridWidthValue,
		gridHeightLabel,
		gridHeightSlider, gridHeightValue,
		lineLengthLabel,
		lineLengthSlider, lineLengthValue,
		playerCountLabel,
		playerCountSlider, playerCountValue,
		aiForMissingCheckbox,
		bestOfValue,
		bestOfRadioGroup,
		cornerBonusCheckbox,
		solitaireRuleCheckbox,
		bombCounterCheckbox,
		overflowRuleCheckbox,
		container.NewHBox(startGameButton),
	)

	//temp data
	playerData := [][]string{
		{"Name", "Score", "Min Speed", "Max Speed", "UUID"}, // Header row
		{"Alice", "150", "1.5", "3.0", "UUID-001"},
		{"Bob", "200", "2.0", "4.5", "UUID-002"},
		{"Charlie", "120", "1.0", "3.5", "UUID-003"},
		{"Diana", "180", "1.8", "4.0", "UUID-004"},
		{"Eve", "220", "2.2", "5.0", "UUID-005"},
	}

	setupgameContainer := container.NewHBox(leftcontent, playerDropdownsContainer)
	content := container.NewBorder(nil,nil,nil,setupgameContainer,ui.CreateLeaderboard(playerData))
	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}

func startGameSetup(gridWidth, gridHeight, lineLength, playerCount int, enableAlliances bool, playerTypes []int, bestOf int, cornerBonus, solitaireRule, bombCounter, overflowRule, aiForMissing bool) {
	// Create a new Game instance with the setup parameters
	game := ui.NewGame(gridWidth, gridHeight, playerCount, lineLength, 0, bestOf, playerTypes, aiForMissing, cornerBonus, solitaireRule, bombCounter, overflowRule)

	// Pass the game instance to MainGameWindow and display the window
	ui.MainGameWindow(game, fyne.CurrentApp())
}