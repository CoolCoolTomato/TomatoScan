package portScan

func initArguments(speed int) {
	//  所有主机数
	AllHostsNum = int64(len(hosts))
	//  已完成主机数
	AccomplishHostsNum = 0
	//  设置goroutine数
	if speed > len(ports) {
		routineNum = len(ports)
	} else {
		routineNum = speed
	}
	//  清空addressList
	addressList = addressList[0:0]
}
