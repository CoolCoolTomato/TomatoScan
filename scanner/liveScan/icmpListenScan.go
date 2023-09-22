package liveScan

import (
	"golang.org/x/net/icmp"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"tomatoscan/common"
)

// 向主机发送数据
func icmpWrite(hostsChan chan string, conn *icmp.PacketConn) {
	//  从管道接收host
	for host := range hostsChan {
		atomic.AddInt64(&AccomplishHostsNum, 1)
		//  管道关闭条件
		if host == "-1" {
			close(hostsChan)
			break
		}
		//  将host解析成ip，返回一个*IPAddr数据
		addr, _ := net.ResolveIPAddr("ip", host)
		//  构造数据包
		icmpByte := makePacket(host)
		//  向目标主机发送数据包
		conn.WriteTo(icmpByte, addr)
	}
}

// 接收存活主机发送的信号
func icmpReceive(conn *icmp.PacketConn, resultsChan chan string, lastReceiveTime *time.Time) {
	for {
		//  从conn中读取数据
		//  数据会被读取到msg中
		//  ReadFrom函数会返回数据长度和来源host
		msg := make([]byte, 100)
		_, liveHost, _ := conn.ReadFrom(msg)
		//  如果收到host，就发到结果管道
		if liveHost != nil {
			resultsChan <- liveHost.String()
			//  获取最后一次接收到信息的时间
			*lastReceiveTime = time.Now()
		}
	}
}

// IcmpListenScan 通过监听扫描主机
func IcmpListenScan() bool {
	defer common.GC()
	var (
		//  获取最后一次接收到信息的时间
		lastReceiveTime time.Time
		//  发送主机的管道
		hostsChan = make(chan string, len(hosts))
		//  接收结果的管道
		resultsChan = make(chan string)
		//  icmpWrite等待组
		wg sync.WaitGroup
	)

	//  监听本地网络
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false
	}
	defer conn.Close()

	//  向管道发送主机
	go sendHosts(hostsChan)

	//  为icmpWrite创建一个goroutine池
	for i := 0; i < routineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//  从管道接收主机，再发包
			icmpWrite(hostsChan, conn)
		}()
	}

	//  接收响应包，将存活的主机写入resultsChan
	go icmpReceive(conn, resultsChan, &lastReceiveTime)

	//  所有的包发完，并且3秒内没有收到响应，关闭结果管道
	//  hostsChan关闭后，所有icmpWrite都退出阻塞
	go func() {
		//  所有icmpWrite都退出阻塞后，锁就开了
		wg.Wait()
		//  上一次收到响应到现在所经历的时间，大于三秒则break
		for {
			if time.Now().Sub(lastReceiveTime) > 3*time.Second {
				break
			}
		}
		//  关闭结果管道
		close(resultsChan)
	}()

	//  接收结果
	for host := range resultsChan {
		liveHosts = append(liveHosts, host)
	}

	return true
}
