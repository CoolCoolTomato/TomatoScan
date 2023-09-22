package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"tomatoscan/scanner/liveScan"
)

func LiveScanBox() *fyne.Container {
	//  输入框和其绑定的列表
	ipBox, ipList := inputList("Input Hosts")
	//  装输入框的容器
	inputBox := container.NewGridWithRows(1)
	inputBox.Add(ipBox)
	//  输出框和其绑定的列表
	resultDisplay, resultList := outputList()
	//  状态栏
	statusLabel := widget.NewLabel("")
	//  速度输入框
	speedLabel := widget.NewRichTextFromMarkdown("## Speed")
	speedList := []string{"1000", "2000", "3000", "4000"}
	speedInput := widget.NewSelectEntry(speedList)
	speedInput.SetText("3000")
	speedBox := container.NewGridWithColumns(3, statusLabel, speedLabel, speedInput)
	//  运行进度条
	progress := widget.NewProgressBar()
	if liveScan.AllHostsNum == 0 {
		progress.Value = float64(0)
	} else {
		progress.Value = float64(liveScan.AccomplishHostsNum) / float64(liveScan.AllHostsNum)
	}
	//  结束标志
	var endFlag bool
	//  开始按钮
	startPortScan := widget.NewButton("Run", func() {
		endFlag = false
		statusLabel.SetText("Scanning...")
		ips, _ := ipList.Get()
		//  判断参数是否正确
		if len(ips) == 0 {
			statusLabel.SetText("Empty Arguments")
			return
		}
		speed, err := strconv.Atoi(speedInput.Text)
		if err != nil {
			statusLabel.SetText("Invalid Arguments")
			return
		}
		endFlag = false
		//  开始扫描
		go func() {
			addressList, _, _ := liveScan.LiveScan(append(ips, "-1"), speed)
			resultList.Set(addressList)
			endFlag = true
		}()
		for true {
			if endFlag {
				break
			}
			if liveScan.AllHostsNum == 0 {
				progress.Value = float64(0)
			} else {
				progress.Value = float64(liveScan.AccomplishHostsNum) / float64(liveScan.AllHostsNum)
			}
			progress.Refresh()
		}
		statusLabel.SetText("Accomplish")
	})
	//  装按钮的容器
	buttonBox := container.NewGridWithColumns(1, startPortScan)
	//  导出类型
	exportTypes := []string{"TXT", "JSON"}
	exportTypeInput := widget.NewSelectEntry(exportTypes)
	exportTypeInput.SetText("TXT")
	//  导出按钮
	exportButton := widget.NewButton("Export", func() {
		exportType := exportTypeInput.Text
		result, _ := resultList.Get()
		switch exportType {
		case "TXT":
			f1ag := exportToTXT(result, "result.txt")
			if f1ag {
				statusLabel.SetText("Export succeed")
			} else {
				statusLabel.SetText("Export error")
			}
		case "JSON":
			f1ag := exportToJSON(result, "result.json")
			if f1ag {
				statusLabel.SetText("Export succeed")
			} else {
				statusLabel.SetText("Export error")
			}
		}
	})
	//  装导出类型和导出按钮的容器
	exportBox := container.NewGridWithColumns(2, exportButton, exportTypeInput)

	//  控制台
	resultHandle := container.NewVBox()
	resultHandle.Add(speedBox)
	resultHandle.Add(myNewLine())
	resultHandle.Add(progress)
	resultHandle.Add(myNewLine())
	resultHandle.Add(buttonBox)
	resultHandle.Add(myNewLine())
	resultHandle.Add(exportBox)

	//  装结果和控制台的容器
	resultBox := container.NewBorder(
		myNewLine(),
		container.NewBorder(myNewLine(), myNewLine(), myNewLine(), myNewLine(), resultHandle),
		myNewLine(),
		myNewLine(),
		resultDisplay)

	//  总容器
	box := container.NewGridWithColumns(2)
	box.Add(inputBox)
	box.Add(resultBox)
	return box
}
