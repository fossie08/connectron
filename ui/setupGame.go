package ui

import (
	"fmt"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	EasyAI   = iota
	MediumAI
	HardAI
)

func SetupWindow(connectronApp fyne.App) {
	setupWindow := connectronApp.NewWindow("Connectron - Setup Game")
	setupWindow.Resize(fyne.NewSize(600, 400))

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

	// Update player dropdowns based on the number of players
	updatePlayerDropdowns := func(count int) {
		playerDropdownsContainer.RemoveAll() // Clear existing dropdowns
		for i := 0; i < count; i++ {
			options := []string{"Easy AI", "Medium AI", "Hard AI", "Person"}
			dropdown := widget.NewSelect(options, func(selected string) {
				// Handle selection change if needed
			})
			playerDropdownsContainer.Add(dropdown)
		}
		playerDropdownsContainer.Refresh() // Refresh the container to show new dropdowns
	}

	playerCountSlider.OnChanged = func(value float64) {
		playerCountValue.SetText(fmt.Sprintf("%d", int(value)))
		updatePlayerDropdowns(int(value))
	}

	// Player Type Dropdown
	playerTypes := make([]int, 4)
	for i := 0; i < 4; i++ {
		playerTypes[i] = EasyAI // Default AI type
	}

	// AI/Player Configuration
	aiForMissingCheckbox := widget.NewCheck("AI for Missing Players", nil)
	bestOfValue := widget.NewLabel("Best of (3, 5, etc.):")
	bestOfDropdown := widget.NewSelect([]string{"3", "5", "7"}, func(selected string) {
		bestOfValue.SetText(selected)
	})

	// Special Rule Options
	cornerBonusCheckbox := widget.NewCheck("Enable Corner Bonus", nil)
	solitaireRuleCheckbox := widget.NewCheck("Enable Solitaire Destruction", nil)
	bombCounterCheckbox := widget.NewCheck("Enable Bomb Counter", nil)
	overflowRuleCheckbox := widget.NewCheck("Enable Overflow Rule", nil)

	startGameButton := widget.NewButton("Start Game", func() {
		bestOf := 3 // Default value
		if bestOfDropdown.Selected != "" {
			bestOf, _ = strconv.Atoi(bestOfDropdown.Selected)
		}

		startGameSetup(int(gridWidthSlider.Value), int(gridHeightSlider.Value), int(lineLengthSlider.Value), int(playerCountSlider.Value), false, playerTypes, bestOf, cornerBonusCheckbox.Checked, solitaireRuleCheckbox.Checked, bombCounterCheckbox.Checked, overflowRuleCheckbox.Checked, aiForMissingCheckbox.Checked)
		setupWindow.Close()
	})

	cancelButton := widget.NewButton("Cancel", func() {
		setupWindow.Close()
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
		bestOfDropdown,
		cornerBonusCheckbox,
		solitaireRuleCheckbox,
		bombCounterCheckbox,
		overflowRuleCheckbox,
		container.NewHBox(startGameButton, cancelButton),
	)

	content := container.NewHBox(leftcontent, playerDropdownsContainer)
	setupWindow.SetContent(content)
	setupWindow.Show()
}

func startGameSetup(gridWidth, gridHeight, lineLength, playerCount int, enableAlliances bool, playerTypes []int, bestOf int, cornerBonus, solitaireRule, bombCounter, overflowRule, aiForMissing bool) {
    // Create a new Game instance with the setup parameters
    game := NewGame(gridWidth, gridHeight, playerCount, lineLength, bestOf, playerTypes, aiForMissing, cornerBonus, solitaireRule, bombCounter, overflowRule)

    // Pass the game instance to MainGameWindow and display the window
    MainGameWindow(game, fyne.CurrentApp())
}
