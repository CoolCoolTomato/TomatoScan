package liveScan

// 初始化参数
func initArguments(speed int) {
	//  所有主机数
	AllHostsNum = int64(len(hosts))
	//  已完成主机数
	AccomplishHostsNum = 0
	//  设置goroutine数
	if speed > len(hosts) {
		routineNum = len(hosts)
	} else {
		routineNum = speed
	}
}

// 向管道发送主机
func sendHosts(hostsChan chan string) {
	for _, host := range hosts {
		hostsChan <- host
	}
}

// 制作数据包
func makePacket(host string) []byte {
	msg := make([]byte, 40)
	id0, id1 := genIdentifier(host)
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4], msg[5] = id0, id1
	msg[6], msg[7] = genSequence(1)
	check := checkSum(msg[0:40])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)
	return msg
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	answer := uint16(^sum)
	return answer
}

func genSequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genIdentifier(host string) (byte, byte) {
	return host[0], host[1]
}
