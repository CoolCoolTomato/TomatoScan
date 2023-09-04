package scanner

import (
	"fmt"
	"math"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// 计算端口数
func getPortsNum(ports []string) int {
	num := 0
	for i := 0; i < len(ports); i++ {
		port := ports[i]
		switch judgePort(port) {
		case 2:
			se := strings.Split(port, "-")
			start, _ := strconv.Atoi(se[0])
			end, _ := strconv.Atoi(se[1])
			num = num + end - start + 1
		case 1:
			num += 1
		}
	}
	return num
}

// tcp连接函数
func connectPort(host string, portsChan chan int, resultsChan chan int) {
	//  从管道中读取端口
	for p := range portsChan {
		//  关闭管道
		if p == -1 {
			close(portsChan)
			break
		}
		//  尝试TCP连接
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			continue
		}
		//  关闭连接
		conn.Close()
		//  发送开放端口
		resultsChan <- p
	}
}

// 迭代运行ip
func rangeIpRun(ipBlock []string, ind int, ports []string, speed int, openPorts *[]string) {
	if ind == 4 {
		ip := ipBlock[0] + "." + ipBlock[1] + "." + ipBlock[2] + "." + ipBlock[3]
		runRes := getOpenPorts(ip, ports, speed)
		for i := 0; i < len(runRes); i++ {
			*openPorts = append(*openPorts, ip+":"+strconv.Itoa(runRes[i]))
		}
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
				rangeIpRun(newIpBlock, ind+1, ports, speed, openPorts)
			}
		} else {
			rangeIpRun(ipBlock, ind+1, ports, speed, openPorts)
		}
	}
}

// 发送端口到管道
func sendPorts(ports []string, portsChan chan int) {
	for i := 0; i < len(ports); i++ {
		port := ports[i]
		switch judgePort(port) {
		case -1:
			p, _ := strconv.Atoi(port)
			portsChan <- p
		case 2:
			se := strings.Split(port, "-")
			start, _ := strconv.Atoi(se[0])
			end, _ := strconv.Atoi(se[1])
			for i := start; i <= end; i++ {
				portsChan <- i
			}
		case 1:
			p, _ := strconv.Atoi(port)
			portsChan <- p
		}
	}
}

// 端口存活探测函数，返回存活端口
func getOpenPorts(host string, ports []string, speed int) []int {
	//  端口管道
	portsChan := make(chan int, speed)
	//  结果管道
	resultsChan := make(chan int)
	//  goroutine计数
	var complete sync.WaitGroup
	complete.Add(speed)
	//  开放端口切片，储存结果
	var openPorts []int
	//  创建goroutine池
	for i := 0; i < speed; i++ {
		//  并发，建立TCP连接
		go func() {
			defer complete.Done()
			connectPort(host, portsChan, resultsChan)
		}()
	}
	//  向端口管道发送端口
	go sendPorts(ports, portsChan)
	//  结束判断
	go func() {
		//  complete为0时结束
		complete.Wait()
		close(resultsChan)
	}()
	//  从结果管道获取开放端口
	for r := range resultsChan {
		openPorts = append(openPorts, r)
	}
	//  将端口进行排序
	sort.Ints(openPorts)
	//  返回结果
	return openPorts
}

func PortScan(hosts []string, ports []string, speed int) []string {
	var results []string
	//  获取host，启动run函数
	for i := 0; i < len(hosts); i++ {
		host := hosts[i]
		//  判断host类型
		switch judgeHost(host) {
		//  CIDR类型
		case 4:
			_, ipNet, err := net.ParseCIDR(host)
			if err != nil {
				fmt.Println("err")
			}
			ip := ipNet.IP
			mask, _ := ipNet.Mask.Size()
			for i := 0; i < int(math.Pow(2, float64(32-mask))); i++ {
				openPorts := getOpenPorts(ip.String(), ports, speed)
				for j := 0; j < len(openPorts); j++ {
					results = append(results, ip.String()+":"+strconv.Itoa(openPorts[j]))
				}
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
			openPorts := getOpenPorts(host, ports, speed)
			for j := 0; j < len(openPorts); j++ {
				results = append(results, host+":"+strconv.Itoa(openPorts[j]))
			}
		//  用-表示范围的ip
		case 2:
			var openPorts []string
			ipBlock := strings.Split(host, ".")
			rangeIpRun(ipBlock, 0, ports, speed, &openPorts)
			for j := 0; j < len(openPorts); j++ {
				results = append(results, openPorts[j])
			}
		//  域名主机
		case 1:
			openPorts := getOpenPorts(host, ports, speed)
			for j := 0; j < len(openPorts); j++ {
				results = append(results, host+":"+strconv.Itoa(openPorts[j]))
			}
		}
	}
	return results
}
