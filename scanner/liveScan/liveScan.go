package liveScan

import (
	"tomatoscan/common"
	"tomatoscan/scanner/parse"
)

func LiveScan(inputHosts []string, speed int) ([]string, []string, string) {
	//  解析输入
	hosts, errorHosts = parse.ParseHosts(inputHosts)
	//  去重
	hosts = common.RemoveDuplicates(hosts)
	//  初始化参数
	initArguments(speed)
	//  判断hosts是否为空
	if len(hosts) <= 1 {
		scanType = "NoHosts"
		return hosts, errorHosts, scanType
	}
	//  开始扫描
	if IcmpListenScan() {
		scanType = "IcmpListenScan"
	} else if IcmpConnectScan() {
		scanType = "IcmpConnectScan"
	} else {
		scanType = "NoScan"
	}
	//  返回结果
	return common.RemoveDuplicates(liveHosts), errorHosts, scanType
}
