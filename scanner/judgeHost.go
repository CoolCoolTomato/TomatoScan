package scanner

import "regexp"

// 主机类型正则匹配
func judgeHost(host string) int {
	//  结束标志
	reEnd := regexp.MustCompile(`^-1$`)
	if reEnd.MatchString(host) {
		return -1
	}
	//  CIDR型主机(4)
	re4 := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}$`)
	if re4.MatchString(host) {
		return 4
	}
	//  单个正常主机(3)
	re3 := regexp.MustCompile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`)
	if re3.MatchString(host) {
		return 3
	}
	//  192.168.0.0-255型主机(3)
	re2 := regexp.MustCompile(`^(\d{1,3}(\-\d{1,3})?)\.(\d{1,3}(\-\d{1,3})?)\.(\d{1,3}(\-\d{1,3})?)\.(\d{1,3}(\-\d{1,3})?)$`)
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
