package pbft

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

// //客户端使用的tcp监听
// func clientTcpListen() {
// 	listen, err := net.Listen("tcp", clientAddr)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	defer listen.Close()

// 	for {
// 		conn, err := listen.Accept()
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		b, err := ioutil.ReadAll(conn)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		fmt.Println(string(b))
// 	}

// }

//节点使用的tcp监听
func (p *Pbft) TcpListen() {
	listen, err := net.Listen("tcp", p.Node.addr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("节点开启监听，地址：%s\n", p.Node.addr)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		p.handleRequest(b)
	}

}
