package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

func PortScanWarp(portScan func(hosts []string, ports []string, speed int) []string) *fyne.Container {
	ipBox, ipList := IpInput()
	portBox, portList := PortInput()

	inputBox := container.NewGridWithRows(2)
	inputBox.Add(ipBox)
	inputBox.Add(portBox)

	resultList := binding.BindStringList(&[]string{})

	resultDisplay := widget.NewListWithData(resultList, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		object.(*widget.Label).Bind(item.(binding.String))
	})

	statusLabel := widget.NewLabel("")

	speedLabel := widget.NewRichTextFromMarkdown("## Speed")
	speedList := []string{"2000", "3000", "4000", "5000", "6000", "7000", "8000"}
	speedInput := widget.NewSelectEntry(speedList)
	speedInput.SetText("5000")
	speedBox := container.NewGridWithColumns(3, statusLabel, speedLabel, speedInput)

	startPortScan := widget.NewButton("Run", func() {
		ips, _ := ipList.Get()
		ports, _ := portList.Get()
		speed, _ := strconv.Atoi(speedInput.Text)
		fl4g := 0
		go func() {
			for true {
				if fl4g == 1 {
					break
				}
				statusLabel.SetText("Scanning")
				time.Sleep(time.Millisecond * 500)
				statusLabel.SetText("Scanning.")
				time.Sleep(time.Millisecond * 500)
				statusLabel.SetText("Scanning..")
				time.Sleep(time.Millisecond * 500)
				statusLabel.SetText("Scanning...")
				time.Sleep(time.Millisecond * 500)
			}
			statusLabel.SetText("Accomplish")
		}()
		_ = resultList.Set(portScan(ips, append(ports, "-1"), speed))
		fl4g = 1
	})

	exportTypes := []string{"TXT", "JSON"}
	exportTypeInput := widget.NewSelectEntry(exportTypes)
	exportTypeInput.SetText("TXT")
	exportButton := widget.NewButton("Export", func() {
		exportType := exportTypeInput.Text
		result, _ := resultList.Get()
		switch exportType {
		case "TXT":
			f1ag := exportToTXT(result)
			if f1ag {
				statusLabel.SetText("Export succeed")
			} else {
				statusLabel.SetText("Export error")
			}
		case "JSON":
			f1ag := exportToJSON(result)
			if f1ag {
				statusLabel.SetText("Export succeed")
			} else {
				statusLabel.SetText("Export error")
			}
		}
	})
	exportBox := container.NewGridWithColumns(2, exportButton, exportTypeInput)

	resultHandle := container.NewVBox()
	resultHandle.Add(speedBox)
	resultHandle.Add(myNewLine())
	resultHandle.Add(startPortScan)
	resultHandle.Add(myNewLine())
	resultHandle.Add(exportBox)

	resultBox := container.NewBorder(
		myNewLine(),
		container.NewBorder(myNewLine(), myNewLine(), myNewLine(), myNewLine(), resultHandle),
		myNewLine(),
		myNewLine(),
		resultDisplay)

	box := container.NewGridWithColumns(2)
	box.Add(inputBox)
	box.Add(resultBox)
	return box
}
