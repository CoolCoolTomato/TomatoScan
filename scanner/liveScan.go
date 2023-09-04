package scanner

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 封装 icmp 报头
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func checkSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index]) << 8
	}
	rt = uint16(sum) + uint16(sum>>16)
	return ^rt
}

func connectHost(hostsChan chan string, resultsChan chan string, count int) {
	for host := range hostsChan {
		if host == "-1" {
			close(hostsChan)
			break
		}
		var (
			//  得到本机的IP地址结构
			lAddr = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
			//  解析域名得到 IP 地址结构
			rAddr, _ = net.ResolveIPAddr("ip", host)
			//  icmp
			icmp ICMP
			//  起始字节
			originBytes = make([]byte, 10)
		)
		//  返回一个 ip socket
		conn, err := net.DialIP("ip4:icmp", &lAddr, rAddr)
		//  连接失败
		if err != nil {
			continue
		}
		//  自动关闭连接
		defer conn.Close()
		// 初始化 icmp 报文
		icmp = ICMP{8, 0, 0, 0, 0}

		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, icmp)
		binary.Write(&buffer, binary.BigEndian, originBytes)
		b := buffer.Bytes()
		binary.BigEndian.PutUint16(b[2:], checkSum(b))
		recv := make([]byte, 1024)

		for i := 0; i < count; i++ {
			//  发包
			_, err1 := conn.Write(buffer.Bytes())
			//  发送失败
			if err1 != nil {
				continue
			}
			//  设置超时
			conn.SetReadDeadline(time.Now().Add(time.Second))
			//  收包
			it, err2 := conn.Read(recv)
			//  接收失败
			if err2 != nil {
				continue
			}
			//  接收成功
			if it != 0 {
				resultsChan <- host
				break
			}
		}
	}
}

func rangeIpSend(ipBlock []string, ind int, hostsChan chan string) {
	if ind == 4 {
		ip := ipBlock[0] + "." + ipBlock[1] + "." + ipBlock[2] + "." + ipBlock[3]
		hostsChan <- ip
		return
	} else {
		ii := ipBlock[ind]
		re := regexp.MustCompile(`\-`)
		if re.MatchString(ii) {
			ise := strings.Split(ii, "-")
			iStart, _ := strconv.Atoi(ise[0])
			iEnd, _ := strconv.Atoi(ise[1])
			for i := iStart; i <= iEnd; i++ {
				newIpBlock := make([]string, 4)
				copy(newIpBlock, ipBlock)
				newIpBlock[ind] = strconv.Itoa(i)
				rangeIpSend(newIpBlock, ind+1, hostsChan)
			}
		} else {
			rangeIpSend(ipBlock, ind+1, hostsChan)
		}
	}
}

func sendHosts(hosts []string, hostsChan chan string) {
	for i := 0; i < len(hosts); i++ {
		host := hosts[i]
		switch judgeHost(host) {
		case -1:
			hostsChan <- host
		//  CIDR类型
		case 4:
			_, ipNet, err := net.ParseCIDR(host)
			if err != nil {
				fmt.Println("err")
			}
			ip := ipNet.IP
			mask, _ := ipNet.Mask.Size()
			for i := 0; i < int(math.Pow(2, float64(32-mask))); i++ {
				hostsChan <- ip.String()
				ip[3]++
				if ip[3] == 0 {
					ip[2] += 1
				}
				if ip[2] == 0 {
					ip[1] += 1
				}
				if ip[1] == 0 {
					ip[0] += 1
				}
			}
		//  正常ip
		case 3:
			hostsChan <- host
		//  用-表示范围的ip
		case 2:
			ipBlock := strings.Split(host, ".")
			rangeIpSend(ipBlock, 0, hostsChan)
		//  域名主机
		case 1:
			hostsChan <- host
		}
	}
}

func getLiveHosts(hosts []string, speed int, count int) []string {
	//  主机管道
	hostsChan := make(chan string, speed)
	//  结果管道
	resultsChan := make(chan string)
	//  goroutine计数
	var complete sync.WaitGroup
	complete.Add(speed)
	//  存活主机切片，储存结果
	var liveHosts []string
	for i := 0; i < speed; i++ {
		go func() {
			defer complete.Done()
			connectHost(hostsChan, resultsChan, count)
		}()
	}
	//  向管道发送主机
	go sendHosts(hosts, hostsChan)
	go func() {
		//  complete为0时结束
		complete.Wait()
		close(resultsChan)
	}()
	//  从结果管道获取存活主机
	for r := range resultsChan {
		liveHosts = append(liveHosts, r)
	}
	//  返回结果
	return liveHosts
}

func LiveScan(hosts []string, speed int, count int) []string {
	var results []string
	results = append(results, getLiveHosts(hosts, speed, count)...)
	//  结果排序
	var ips []net.IP
	for _, ip := range results {
		ips = append(ips, net.ParseIP(ip))
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	var sortResults []string
	for _, ip := range ips {
		sortResults = append(sortResults, ip.String())
	}
	return sortResults
}
