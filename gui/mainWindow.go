package gui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"tomatoscan/scanner"
)

func MainWindow() {
	App := app.New()
	Window := App.NewWindow("TomatoScan")

	mainBox := container.NewMax()

	portscan := PortScanWarp(scanner.PortScan)
	livescan := LiveScanWarp(scanner.LiveScan)

	tabs := container.NewAppTabs(
		container.NewTabItem("portscan", portscan),
		container.NewTabItem("livescan", livescan),
	)

	mainBox.Add(tabs)

	Window.SetContent(mainBox)

	Window.Show()
	App.Run()
}
