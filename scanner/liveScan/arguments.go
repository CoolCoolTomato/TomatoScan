package liveScan

var (
	// AllHostsNum 所有目标
	AllHostsNum int64
	// AccomplishHostsNum 已完成目标
	AccomplishHostsNum int64
	//  目标主机放在这里!!!
	hosts []string
	//  错误输入
	errorHosts []string
	//  存活主机放在这里!!!
	liveHosts []string
	//  goroutine数
	routineNum int
	//  扫描种类
	scanType string
)
