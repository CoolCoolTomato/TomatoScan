package parse

import (
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// 主机类型正则匹配
func getHostType(host string) int {
	//  结束标志(-1)
	reEnd := regexp.MustCompile(`^-1$`)
	if reEnd.MatchString(host) {
		return -1
	}
	//  CIDR型主机(4)
	re4 := regexp.MustCompile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})/(\d{1,2})$`)
	if re4.MatchString(host) {
		return 4
	}
	//  单个正常主机(3)
	re3 := regexp.MustCompile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`)
	if re3.MatchString(host) {
		return 3
	}
	//  192.168.0.0-255型主机(3)
	re2 := regexp.MustCompile(`^(\d{1,3}(-\d{1,3})?)\.(\d{1,3}(-\d{1,3})?)\.(\d{1,3}(-\d{1,3})?)\.(\d{1,3}(-\d{1,3})?)$`)
	if re2.MatchString(host) {
		return 2
	}
	//  url型主机(1)
	re1 := regexp.MustCompile(`^([^.]+)(\.\S+)+$`)
	if re1.MatchString(host) {
		return 1
	}
	//  err(0)
	return 0
}

// CIDR型主机(4)
func parseHost4(host string) []string {
	var hostsList []string
	_, ipNet, err := net.ParseCIDR(host)
	if err != nil {
		return hostsList
	}
	ip := ipNet.IP
	mask, _ := ipNet.Mask.Size()
	for i := 0; i < int(math.Pow(2, float64(32-mask))); i++ {
		hostsList = append(hostsList, ip.String())
		ip[3]++
		if ip[3] == 0 {
			ip[2] += 1
			if ip[2] == 0 {
				ip[1] += 1
				if ip[1] == 0 {
					ip[0] += 1
				}
			}
		}
	}
	return hostsList
}

// 单个正常主机(3)
func parseHost3(host string) string {
	return host
}

// 192.168.0.0-255型主机(3)
func parseHost2(host string) []string {
	var hostsList []string
	ipBlock := strings.Split(host, ".")
	for i0, j0 := getRange(ipBlock[0]); i0 <= j0; i0++ {
		for i1, j1 := getRange(ipBlock[1]); i1 <= j1; i1++ {
			for i2, j2 := getRange(ipBlock[2]); i2 <= j2; i2++ {
				for i3, j3 := getRange(ipBlock[3]); i3 <= j3; i3++ {
					ip := strconv.Itoa(i0) + "." + strconv.Itoa(i1) + "." + strconv.Itoa(i2) + "." + strconv.Itoa(i3)
					hostsList = append(hostsList, ip)
				}
			}
		}
	}
	return hostsList
}

// url型主机(1)
func parseHost1(host string) string {
	return host
}

// ParseHosts 解析输入的所有主机
func ParseHosts(inputHosts []string) ([]string, []string) {
	var hosts []string
	var errorHosts []string
	for _, host := range inputHosts {
		switch getHostType(host) {
		case -1:
			hosts = append(hosts, host)
		case 0:
			errorHosts = append(errorHosts, host)
		case 1:
			hosts = append(hosts, parseHost1(host))
		case 2:
			hosts = append(hosts, parseHost2(host)...)
		case 3:
			hosts = append(hosts, parseHost3(host))
		case 4:
			hosts = append(hosts, parseHost4(host)...)
		default:
			continue
		}
	}
	return hosts, errorHosts
}

// 获取某个block的范围
func getRange(block string) (int, int) {
	re := regexp.MustCompile(`-`)
	if re.MatchString(block) {
		blockRange := strings.Split(block, "-")
		blockStart, _ := strconv.Atoi(blockRange[0])
		blockEnd, _ := strconv.Atoi(blockRange[1])
		return blockStart, blockEnd
	} else {
		blockRange, _ := strconv.Atoi(block)
		return blockRange, blockRange
	}
}
