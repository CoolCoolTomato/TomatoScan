package gui

import (
	"bufio"
	"encoding/json"
	"fyne.io/fyne/v2/data/binding"
	"os"
)

// 从列表中删除元素
func delItemFromList(l []string, item string) []string {
	var res []string
	for _, s := range l {
		if s != item {
			res = append(res, s)
		}
	}
	return res
}

// 从绑定数据的列表中删除元素
func delItemFromBindStringList(l binding.StringList, item string) {
	oldList, _ := l.Get()
	newList := delItemFromList(oldList, item)
	_ = l.Set(newList)
}

// 导出到TXT文件
func exportToTXT(result []string, txtPath string) bool {
	file, err := os.OpenFile(txtPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return false
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	for _, r := range result {
		write.WriteString(r + "\n")
	}
	_ = write.Flush()
	return true
}

// 导出到JSON文件
func exportToJSON(result []string, jsonPath string) bool {
	jsonResult := make(map[string]interface{})
	jsonResult["result"] = result
	file, err1 := os.OpenFile(jsonPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err1 != nil {
		return false
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err2 := encoder.Encode(jsonResult)
	if err2 != nil {
		return false
	}
	return true
}
