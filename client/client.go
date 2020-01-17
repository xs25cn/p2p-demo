package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	cName  = flag.String("name", "duzhenxun", "客户端名称（注意2台客户端名不要用一样的哦）")
	cPort  = flag.Int("p", 25251, "客户端要占用本机的端口")
	server = flag.String("s", "39.106.231.36:2525", "服务端地址")
)

func main() {
	flag.Parse()
	cAddr := &net.UDPAddr{IP: net.IPv4zero, Port: *cPort}
	//从服务器获取对方客户端地址
	dstAddr, err := getDstAddr(cAddr, *server, *cName)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("对方地址",dstAddr.String())

	//与对方地址连接
	conn, err := net.DialUDP("udp", cAddr, dstAddr)
	if err != nil {
		fmt.Println("对方进行连接失败！！！" + err.Error())
	}

	//向另一方发送一条udp消息(对方的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方就可以发过来udp消息
	if _, err = conn.Write([]byte("connect...")); err != nil {
		log.Println("第一次发送失败", err)
	}
	log.Println("与对方客户端打洞成功....")
	//time.Sleep(2 * time.Second)
	//给对方每过5秒发一次心跳
	go func() {
		for {
			_, err = conn.Write([]byte("ping"))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	//输出对方发来的数据
	for {
		data := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println(err,"读取信息失败")
		} else {
			log.Println("--->", addr.String(), string(data[:n]))
		}
	}
}


func getDstAddr(cAddr *net.UDPAddr, server string, cName string) (*net.UDPAddr, error) {
	sAddr := strings.Split(server, ":")
	sPort, _ := strconv.Atoi(sAddr[1])
	conn, err := net.DialUDP("udp", cAddr, &net.UDPAddr{IP: net.ParseIP(sAddr[0]), Port: sPort})
	defer conn.Close()

	if err != nil {
		return nil, errors.New("连接服务器失败" + err.Error())
	}
	//向服务端写消息，如果服务端会应答
	_, err = conn.Write([]byte(cName))
	if err != nil {
		return nil, errors.New("向服务端发送消息失败！" + err.Error())
	}
	log.Println("与服务端连接并发送消息成功，等待服务器通知另一个客户端的ip地址与端口....")

	data := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(data)

	if err != nil {
		return nil, errors.New("读取服务器消息失败！" + err.Error())
	}
	log.Println("关闭与服务端的连接....")

	//分析服务端发来的另一个客户端信息
	dstAddrArr := strings.Split(string(data[:n]), ":")
	port, _ := strconv.Atoi(dstAddrArr[1])
	dstAddr := &net.UDPAddr{IP: net.ParseIP(dstAddrArr[0]), Port: port}

	return dstAddr, nil
}
