package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	port := flag.Int("p", 2525, "服务端开放接口")
	flag.Parse()
	listen, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: *port})
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("服务端开放端口: <%s> \n", listen.LocalAddr().String())

	data := make([]byte, 1024)
	clients := make([]net.UDPAddr, 0, 2) //这里咱们测试只放2个客户端，要是够2个咱就让他们互相通信
	for {
		n, remoteAddr, err := listen.ReadFromUDP(data)
		fmt.Println("读取监听的udp端口数据", n, remoteAddr, err)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%s -> %s\n", remoteAddr.String(), data[:n])
		clients = append(clients, *remoteAddr)
		//如果读取的到2个客户端数据就开始互连
		if len(clients) == 2 {
			fmt.Println(clients)
			log.Printf("开始P2P打洞 %s <--> %s 的连接\n服务器写客户端断开后，测试一下客户端之间是否还在通信中...\n", clients[0].String(), clients[1].String())
			//分别告诉2个客户端对方的地址与可通信的端口
			c1, err := listen.WriteToUDP([]byte(clients[1].String()), &clients[0])
			if err != nil {
				log.Println("c1", c1, err)
				return
			}

			c2, err := listen.WriteToUDP([]byte(clients[0].String()), &clients[1])
			if err != nil {
				log.Println("c2", c2, err)
				return
			}
			log.Println("服务器务与客户端通信完后，客户端那收到信息后主动断开~")
			time.Sleep(time.Second * 2)

			log.Println("清空P2P连接信息，重新等待新的连接产生...")
			clients = make([]net.UDPAddr, 0, 2) //清空client

		}
	}
}
