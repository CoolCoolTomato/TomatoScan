package liveScan

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
	"tomatoscan/common"
)

// 尝试连接目标主机
func icmpConnect(hostsChan chan string, resultsChan chan string) {
	//  从管道接收host
	for host := range hostsChan {
		atomic.AddInt64(&AccomplishHostsNum, 1)
		//  管道关闭条件
		if host == "-1" {
			close(hostsChan)
			break
		}
		//  尝试与host建立连接
		conn, err := net.DialTimeout("ip4:icmp", host, 5*time.Second)
		if err != nil || conn == nil {
			continue
		}
		//  与Timeout不同
		//  Deadline这是一个绝对的时间点，超过了就就报错啦
		//  而Timeout规定的是一个IO操作的最长时间
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			conn.Close()
			continue
		}
		//  构造数据包
		icmpByte := makePacket(host)
		//  发送数据包
		if _, err := conn.Write(icmpByte); err != nil {
			conn.Close()
			continue
		}
		//  接收数据包
		receive := make([]byte, 60)
		if _, err := conn.Read(receive); err != nil {
			conn.Close()
			continue
		}
		conn.Close()
		//  成功接收后向管道发送host
		resultsChan <- host
	}
}

// IcmpConnectScan 通过连接扫描主机
func IcmpConnectScan() bool {
	defer common.GC()
	var (
		//  发送主机的管道
		hostsChan = make(chan string, len(hosts))
		//  接收结果的管道
		resultsChan = make(chan string)
		//  icmpConnect等待组
		wg sync.WaitGroup
	)

	//  向管道发送主机
	go sendHosts(hostsChan)

	//  为icmpConnect创建goroutine池
	for i := 0; i < routineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//  从管道接收主机，再连接
			icmpConnect(hostsChan, resultsChan)
		}()
	}

	//  等待所有的主机都尝试连接完毕
	//  hostsChan关闭后，所有icmpConnect都退出阻塞
	go func() {
		//  所有icmpConnect都退出阻塞后，锁就开了
		wg.Wait()
		//  关闭结果管道
		close(resultsChan)
	}()

	//  接收结果
	for host := range resultsChan {
		liveHosts = append(liveHosts, host)
	}

	return true
}
