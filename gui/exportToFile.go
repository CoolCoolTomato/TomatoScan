package gui

import (
	"bufio"
	"encoding/json"
	"os"
)

func exportToTXT(result []string) bool {
	txtPath := "result.txt"
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

func exportToJSON(result []string) bool {
	jsonPath := "result.json"
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
