package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"insighthub.uk/connectron/v2/ui"
)

func main() {
	// Initialize the application
	connectronApp := app.New()
	mainWindow := connectronApp.NewWindow("Connectron")

	// Set the window size
	mainWindow.Resize(fyne.Size{1200,1000})
	mainWindow.CenterOnScreen()

	// Create menu items
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New Game", func() { ui.SetupWindow(connectronApp) }),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Settings", func() {  }),
		),
	)

	// Set the menu
	mainWindow.SetMainMenu(menu)

	// Welcome message
	welcomeLabel := widget.NewLabel("Welcome to Connectron")
	welcomeLabel.Alignment = fyne.TextAlignCenter
	mainWindow.SetContent(container.NewCenter(welcomeLabel))

	mainWindow.ShowAndRun()
}
