package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	port   := flag.Int("p", 2525, "服务端开放接口")
	flag.Parse()
	listen, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: *port})
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("服务端开放端口: <%s> \n", listen.LocalAddr().String())

	data := make([]byte, 1024)
	clients:=make([]net.UDPAddr,0,2) //这里咱们测试只放2个客户端，要是够2个咱就让他们互相通信
	for {
		n,remoteAddr,err:=listen.ReadFromUDP(data)
		fmt.Println("读取监听的udp端口数据",n,remoteAddr,err)

		if err!=nil{
			log.Print(err)
		}
		log.Printf("%s -> %s\n",remoteAddr.String(),data[:n])
		clients = append(clients,*remoteAddr)
		//如果读取的到2个客户端数据就开始互连
		if len(clients)==2{
			log.Printf("开始P2P打洞 %s <--> %s 的连接\n服务器程序在10秒后退出，测试一下客户端是否还在通信中...\n",clients[0].String(),clients[1].String())

			//分别告诉2个客户端对方的地址与可通信的端口
			listen.WriteToUDP([]byte(clients[0].String()),&clients[1])
			listen.WriteToUDP([]byte(clients[1].String()),&clients[0])

			time.Sleep(time.Second*10)
			log.Println("服务器已退出，请观察一下你的2台客户端是否还在继续保持通信~")
		}
	}
}
