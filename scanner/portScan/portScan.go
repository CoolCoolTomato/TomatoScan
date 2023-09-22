package portScan

import (
	"sync/atomic"
	"tomatoscan/common"
	"tomatoscan/scanner/parse"
)

func PortScan(inputHosts []string, inputPorts []string, speed int) ([]string, []string, []string, string) {
	//  解析输入
	hosts, errorHosts = parse.ParseHosts(inputHosts)
	ports, errorPorts = parse.ParsePorts(inputPorts)
	//  去重
	hosts = common.RemoveDuplicates(hosts)
	ports = common.RemoveDuplicates(ports)
	//  初始化参数
	initArguments(speed)
	//  判断hosts或ports是否为空
	if len(hosts) <= 0 || len(ports) <= 1 {
		scanType = "NoHosts or NiPorts"
	}
	for _, host := range hosts {
		TcpConnectScan(host)
		atomic.AddInt64(&AccomplishHostsNum, 1)
	}
	scanType = "TcpConnectScan"
	return common.RemoveDuplicates(addressList), errorHosts, errorPorts, scanType
}
