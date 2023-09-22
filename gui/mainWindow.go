package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func MainWindow() {
	App := app.New()
	Window := App.NewWindow("TomatoScan")

	mainBox := container.NewMax()

	portscan := PortScanBox()
	livescan := LiveScanBox()

	tabs := container.NewAppTabs(
		container.NewTabItem("portscan", portscan),
		container.NewTabItem("livescan", livescan),
	)

	mainBox.Add(tabs)

	Window.SetContent(mainBox)

	Window.Resize(fyne.Size{
		Width:  1200,
		Height: 700,
	})
	Window.Show()
	//App.Settings().SetTheme(theme.LightTheme())
	App.Run()
}
