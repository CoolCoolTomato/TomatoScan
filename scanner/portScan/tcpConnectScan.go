package portScan

import (
	"fmt"
	"net"
	"sync"
	"tomatoscan/common"
)

// tcp连接函数
func tcpConnect(host string, portsChan chan string, resultsChan chan string) {
	//  从管道中读取端口
	for port := range portsChan {
		//  关闭管道
		if port == "-1" {
			close(portsChan)
			break
		}
		//  尝试TCP连接
		address := fmt.Sprintf("%s:%s", host, port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			continue
		}
		//  关闭连接
		conn.Close()
		//  发送开放端口
		resultsChan <- address
	}
}

// 向管道发送主机
func sendPorts(portsChan chan string) {
	for _, port := range ports {
		portsChan <- port
	}
}

func TcpConnectScan(host string) bool {
	defer common.GC()
	var (
		//  发送端口的管道
		portsChan = make(chan string, len(ports))
		//  接收结果的管道
		resultsChan = make(chan string)
		//  tcpConnect等待组
		wg sync.WaitGroup
	)

	//  向管道发送端口
	go sendPorts(portsChan)

	//  为tcpConnect创建goroutine池
	for i := 0; i < routineNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tcpConnect(host, portsChan, resultsChan)
		}()
	}

	//  等待所有的端口都尝试连接完毕
	//  portsChan关闭后，所有tcpConnect都退出阻塞
	go func() {
		//  所有tcpConnect都退出阻塞后，锁就开了
		wg.Wait()
		//  关闭结果管道
		close(resultsChan)
	}()

	//  接收结果
	for address := range resultsChan {
		addressList = append(addressList, address)
	}

	return true
}
