package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

func LiveScanWarp(liveScan func(hosts []string, speed int, count int) []string) *fyne.Container {
	ipBox, ipList := IpInput()

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

	pingLabel := widget.NewRichTextFromMarkdown("## PingNum")
	pingList := []string{"1", "2", "3"}
	pingNumInput := widget.NewSelectEntry(pingList)
	pingNumInput.SetText("2")

	speedBox := container.NewGridWithColumns(5, statusLabel, speedLabel, speedInput, pingLabel, pingNumInput)

	startPortScan := widget.NewButton("Run", func() {
		ips, _ := ipList.Get()
		speed, _ := strconv.Atoi(speedInput.Text)
		pingNum, _ := strconv.Atoi(pingNumInput.Text)
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
		_ = resultList.Set(liveScan(append(ips, "-1"), speed, pingNum))
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
	box.Add(ipBox)
	box.Add(resultBox)
	return box
}
