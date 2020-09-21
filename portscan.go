package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var max int

type Flag struct {
	Ips     string // ip列表
	Ports   string // 端口列表
	Threads int
}

func Init() *Flag {
	portscan_flag := Flag{}
	flag.StringVar(&portscan_flag.Ips, "i", "", "输入ip地址 eg:192.0.0.1-192.0.0.255")
	flag.StringVar(&portscan_flag.Ports, "p", "22,80,1433,3306,3389", "端口列表 eg:22,80,1433,3306,3389")
	flag.IntVar(&portscan_flag.Threads, "t", 300, "线程数量 默认50")
	return &portscan_flag
}

func StringIpToInt(ipstring string) int {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt int = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func IpIntToString(ipInt int) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[len-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

//验证语句
func verifyip(ips string) bool {

	if ips == "" {
		return true
	}
	return false
}
func checkips(ipstr string) (int, int) {
	ips_arr := strings.Split(ipstr, "-")
	if len(ips_arr) < 2 {
		return StringIpToInt(ips_arr[0]), StringIpToInt(ips_arr[0])
	}
	return StringIpToInt(ips_arr[0]), StringIpToInt(ips_arr[1])
}
func checkports(portstr string) []string {
	port_arr := strings.Split(portstr, ",")
	return port_arr
}
func testTCPConnection(ip string, port int, doneChannel chan bool) {
	_, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(port),
		time.Second*1)
	if err == nil {
		fmt.Printf("%s:%d: Open\n", ip, port)
	}
	max--
	//fmt.Println(max)
	doneChannel <- true
}

func main() {
	max = 0
	portscan_flag := &Flag{}
	portscan_flag = Init()
	flag.Parse()
	if verifyip(portscan_flag.Ips) {
		flag.Usage()
		return
	}
	activeThreads := 0
	doneChannel := make(chan bool)
	startip, endip := checkips(portscan_flag.Ips)
	ports := checkports(portscan_flag.Ports)
	//fmt.Println(ports)
	//fmt.Println(startip, endip)
	for i := startip; i <= endip; i++ {
		for p_number := 0; p_number < len(ports); p_number++ {
			//fmt.Println(ports[p_number])
			for {
				if max < portscan_flag.Threads {
					//fmt.Println("123")
					break
				}
			}
			port, _ := strconv.Atoi(ports[p_number])

			go testTCPConnection(IpIntToString(i), port, doneChannel)
			activeThreads++
			max++

		}
	}
	// Wait for all threads to finish
	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}
}
