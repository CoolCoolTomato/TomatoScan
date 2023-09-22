package parse

import (
	"regexp"
	"strconv"
	"strings"
)

// 端口类型正则匹配
func getPortType(port string) int {
	//  结束标志(-1)
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

// 范围端口
func parsePort2(port string) (int, int) {
	portRange := strings.Split(port, "-")
	portStart, _ := strconv.Atoi(portRange[0])
	portEnd, _ := strconv.Atoi(portRange[1])
	return portStart, portEnd
}

// 单个端口
func parsePort1(port string) string {
	return port
}

func ParsePorts(inputPorts []string) ([]string, []string) {
	var ports []string
	var errorPorts []string
	for _, port := range inputPorts {
		switch getPortType(port) {
		case -1:
			ports = append(ports, port)
		case 0:
			errorPorts = append(errorPorts, port)
		case 1:
			ports = append(ports, parsePort1(port))
		case 2:
			for i, j := parsePort2(port); i < j; i++ {
				ports = append(ports, strconv.Itoa(i))
			}
		}
	}
	return ports, errorPorts
}
