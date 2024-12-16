package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	lineLengthValue := widget.NewLabel("4") // Label to display the current slider value
	lineLengthSlider := widget.NewSlider(4, 10)
	lineLengthSlider.OnChanged = func(value float64) {
		lineLengthValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Number of Players
	playerCountLabel := widget.NewLabel("Number of Players (0-10):")
	playerCountValue := widget.NewLabel("0") // Label to display the current slider value
	playerCountSlider := widget.NewSlider(0, 10)
	playerCountSlider.OnChanged = func(value float64) {
		playerCountValue.SetText(fmt.Sprintf("%d", int(value)))
	}

	// AI Inclusion Checkbox
	aiCheckbox := widget.NewCheck("Include AI for Missing Players", nil)

	// Alliance Option Checkbox
	allianceCheckbox := widget.NewCheck("Enable Player Alliances", nil)

	// Buttons for Actions
	startGameButton := widget.NewButton("Start Game", func() {
		startGameSetup(
			int(gridWidthSlider.Value),
			int(gridHeightSlider.Value),
			int(lineLengthSlider.Value),
			int(playerCountSlider.Value),
			aiCheckbox.Checked,
			allianceCheckbox.Checked,
		)
	})

	cancelButton := widget.NewButton("Cancel", func() {
		setupWindow.Close()
	})

	// Layout
	content := container.NewVBox(
		gridWidthLabel,
		gridWidthSlider, gridWidthValue,
		gridHeightLabel,
		gridHeightSlider, gridHeightValue,
		lineLengthLabel,
		lineLengthSlider, lineLengthValue,
		playerCountLabel,
		playerCountSlider, playerCountValue,
		aiCheckbox,
		allianceCheckbox,
		container.NewHBox(startGameButton, cancelButton),
	)

	setupWindow.SetContent(content)
	setupWindow.Show()
}

func startGameSetup(gridWidth, gridHeight, lineLength, playerCount int, includeAI, enableAlliances bool) {
	// Validate inputs and prepare game settings
	// This is a placeholder for further development
	fmt.Printf("Grid Width: %d\n", gridWidth)
	fmt.Printf("Grid Height: %d\n", gridHeight)
	fmt.Printf("Line Length to Win: %d\n", lineLength)
	fmt.Printf("Player Count: %d\n", playerCount)
	fmt.Printf("Include AI: %v\n", includeAI)
	fmt.Printf("Enable Alliances: %v\n", enableAlliances)
}
