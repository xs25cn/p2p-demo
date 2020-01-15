package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	cName = flag.String("name", "duzhenxun", "客户端名称（注意2台客户端名不要用一样的哦）")
	cPort = flag.Int("cPort", 12525, "客户端要占用本机的端口")
	sPort = flag.Int("sPort", 2525, "服务端开放接口")
	sIp   = flag.String("sIp", "xs25.cn", "服务端地址（如小手25服务器）")
)

func main() {

	flag.Parse()

	cAddr := &net.UDPAddr{IP: net.IPv4zero, Port: *cPort}
	conn, err := net.DialUDP("udp", cAddr, &net.UDPAddr{IP: net.ParseIP(*sIp), Port: *sPort})

	if err != nil {
		log.Println(err)
		return
	}
	//向服务端写消息，如果服务端会应答
	if _, err = conn.Write([]byte("Hi~小手25,My name is " + *cName + " ,我想问一下我需要和哪个地址建立连接？")); err != nil {
		log.Println(err)
		return
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		log.Println(err)
		return
	}
	conn.Close()
	//服务端发来的另一个客户端信息
	dstAddrArr := strings.Split(string(data[:n]), ":")
	port, _ := strconv.Atoi(dstAddrArr[1])
	dstAddr := &net.UDPAddr{IP: net.ParseIP(dstAddrArr[0]), Port: port}
	log.Printf("本地:%s,对方:%s,中转服务端:%s", cAddr, dstAddr, remoteAddr)

	//对方地址获取到后与对方进行连接
	p2pConn(cAddr, dstAddr)
}

func p2pConn(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) {
	time.Sleep(3 * time.Second)
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if _, err = conn.Write([]byte("p2p...")); err != nil {
		log.Println(err)
		return
	}
	//给对方发数据
	go func() {
		for {
			time.Sleep(time.Second * 1)
			s := fmt.Sprintf("你好，我是：%s,我来自：%s", *cName, srcAddr)
			if _, err = conn.Write([]byte(s)); err != nil {
				log.Println(err)
			}
		}
	}()

	//输出对方发来的数据
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("收到对方数据->", data[:n])
		}
	}
}
