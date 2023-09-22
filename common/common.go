package common

import (
	"runtime"
	"runtime/debug"
)

func GC() {
	runtime.GC()
	debug.FreeOSMemory()
}

// RemoveDuplicates 去重
func RemoveDuplicates(slice []string) []string {
	var newSlice []string
	temp := map[string]struct{}{}
	for _, s := range slice {
		if _, ok := temp[s]; !ok {
			temp[s] = struct{}{}
			newSlice = append(newSlice, s)
		}
	}
	return newSlice
}
