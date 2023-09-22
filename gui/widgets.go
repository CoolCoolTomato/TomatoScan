package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// 创建线条
func myNewLine() *canvas.Line {
	l := canvas.NewLine(color.RGBA{
		R: 128,
		G: 128,
		B: 128,
		A: 170,
	})
	l.StrokeWidth = 4
	return l
}

// 输入列表
func inputList(title string) (*fyne.Container, binding.StringList) {
	//  数据绑定：字符串切片
	bindList := binding.BindStringList(&[]string{})
	//  要删除的数据放这里
	var delList []string

	//  选择框列表，用来存放数据
	//  该选择框与bindList绑定
	checkList := widget.NewListWithData(bindList, func() fyne.CanvasObject {
		//  返回一个选择框
		checkItem := widget.NewCheck("", func(b bool) {})
		return checkItem
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		//  从item中获取数据
		itemText, _ := item.(binding.String).Get()
		//  将数据与选择框的文本绑定
		object.(*widget.Check).Text = itemText
		//  选择框被选中时将数据加入delList
		object.(*widget.Check).OnChanged = func(b bool) {
			if b {
				delList = append(delList, itemText)
			} else {
				delList = delItemFromList(delList, itemText)
			}
		}
		//  刷新选择框
		object.(*widget.Check).Refresh()
	})

	//  输入框
	dataInput := widget.NewEntry()
	//  提交按钮
	dataSubmitButton := widget.NewButton("Submit", func() {
		//  数据为空则不操作
		data := dataInput.Text
		if data == "" {
			return
		}
		//  向bindList添加数据
		bindList.Append(data)
		dataInput.SetText("")
		checkList.Refresh()
	})
	//  删除按钮
	dataDeleteButton := widget.NewButton("Delete", func() {
		//  删除bindList对应的数据
		for _, data := range delList {
			delItemFromBindStringList(bindList, data)
		}
		//  清空delList
		delList = []string{}
	})

	//  控制台
	dataHandle := container.NewVBox()
	//  控制台标题
	dataHandle.Add(widget.NewRichTextFromMarkdown("## " + title))
	dataHandle.Add(myNewLine())
	//  数据输入框
	dataHandle.Add(dataInput)
	//  提交按钮
	dataHandle.Add(dataSubmitButton)
	//  删除按钮
	dataHandle.Add(dataDeleteButton)
	dataHandle.Add(myNewLine())

	//  将所有组件放到一个容器中
	inputBox := container.NewBorder(
		myNewLine(),
		myNewLine(),
		myNewLine(),
		container.NewBorder(nil, nil, myNewLine(), myNewLine(), dataHandle),
		checkList,
	)
	//  返回容器和数据列表
	return inputBox, bindList
}

// 输出列表
func outputList() (*widget.List, binding.StringList) {
	//  结果列表
	resultList := binding.BindStringList(&[]string{})
	//  结果显示
	resultDisplay := widget.NewListWithData(resultList, func() fyne.CanvasObject {
		//  返回一个标签
		return widget.NewLabel("")
	}, func(item binding.DataItem, object fyne.CanvasObject) {
		//  将数据与标签绑定
		object.(*widget.Label).Bind(item.(binding.String))
	})
	//  返回组件和数据列表
	return resultDisplay, resultList
}
