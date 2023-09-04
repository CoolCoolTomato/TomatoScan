package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func myNewLine() *canvas.Line {
	l := canvas.NewLine(color.RGBA{
		R: 230,
		G: 230,
		B: 230,
		A: 255,
	})
	l.StrokeWidth = 5
	return l
}

func delItemFromList(item string, l []string) []string {
	var res []string
	for _, s := range l {
		if s != item {
			res = append(res, s)
		}
	}
	return res
}

func delItemFromBindStringList(item string, l binding.StringList) {
	oldList, _ := l.Get()
	newList := delItemFromList(item, oldList)
	_ = l.Set(newList)
}

func IpInput() (*fyne.Container, binding.StringList) {
	ipList := binding.BindStringList(&[]string{})
	var delList []string

	ipCheckList := widget.NewListWithData(ipList, func() fyne.CanvasObject {
		ipCheck := widget.NewCheck("", func(b bool) {})
		return ipCheck
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		ip, _ := item.(binding.String).Get()
		object.(*widget.Check).Text = ip
		object.(*widget.Check).OnChanged = func(b bool) {
			if b {
				delList = append(delList, ip)
			} else {
				delList = delItemFromList(ip, delList)
			}
		}
		object.(*widget.Check).Refresh()
	})

	ipEntry := widget.NewEntry()
	ipSubmitButton := widget.NewButton("Submit", func() {
		ip := ipEntry.Text
		if ip == "" {
			return
		}
		_ = ipList.Append(ip)
		ipEntry.Text = ""
		ipEntry.Refresh()
		ipCheckList.Refresh()
	})
	ipDelButton := widget.NewButton("Delete", func() {
		for _, ip := range delList {
			delItemFromBindStringList(ip, ipList)
		}
		delList = []string{}
	})
	ipHandle := container.NewVBox()
	ipHandle.Add(widget.NewRichTextFromMarkdown("## Input your hosts"))
	ipHandle.Add(myNewLine())
	ipHandle.Add(ipEntry)
	ipHandle.Add(ipSubmitButton)
	ipHandle.Add(ipDelButton)

	ipInputBox := container.NewBorder(
		myNewLine(),
		myNewLine(),
		myNewLine(),
		container.NewBorder(nil, nil, myNewLine(), myNewLine(), ipHandle))
	ipInputBox.Add(ipCheckList)
	return ipInputBox, ipList
}

func PortInput() (*fyne.Container, binding.StringList) {
	portList := binding.BindStringList(&[]string{})
	var delList []string

	portCheckList := widget.NewListWithData(portList, func() fyne.CanvasObject {
		portCheck := widget.NewCheck("", func(b bool) {})
		return portCheck
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		port, _ := item.(binding.String).Get()
		object.(*widget.Check).Text = port
		object.(*widget.Check).OnChanged = func(b bool) {
			if b {
				delList = append(delList, port)
			} else {
				delList = delItemFromList(port, delList)
			}
		}
		object.(*widget.Check).Refresh()
	})

	portEntry := widget.NewEntry()
	portSubmitButton := widget.NewButton("Submit", func() {
		port := portEntry.Text
		if port == "" {
			return
		}
		_ = portList.Append(port)
		portEntry.Text = ""
		portEntry.Refresh()
		portCheckList.Refresh()
	})
	portDelButton := widget.NewButton("Delete", func() {
		for _, port := range delList {
			delItemFromBindStringList(port, portList)
		}
		delList = []string{}
	})
	portHandle := container.NewVBox()
	portHandle.Add(widget.NewRichTextFromMarkdown("## Input your ports"))
	portHandle.Add(myNewLine())
	portHandle.Add(portEntry)
	portHandle.Add(portSubmitButton)
	portHandle.Add(portDelButton)

	portInputBox := container.NewBorder(
		myNewLine(),
		myNewLine(),
		myNewLine(),
		container.NewBorder(nil, nil, myNewLine(), myNewLine(), portHandle))
	portInputBox.Add(portCheckList)
	return portInputBox, portList
}
