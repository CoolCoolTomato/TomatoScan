package scanner

import "regexp"

// 端口类型正则匹配
func judgePort(port string) int {
	//  结束标志
	reEnd := regexp.MustCompile(`^-1$`)
	if reEnd.MatchString(port) {
		return -1
	}
	//  范围端口(2)
	re2 := regexp.MustCompile(`^\d+\-\d+$`)
	if re2.MatchString(port) {
		return 2
	}
	//  指定端口(1)
	re1 := regexp.MustCompile(`^\d+$`)
	if re1.MatchString(port) {
		return 1
	}
	//  err(0)
	return 0
}
