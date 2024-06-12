package pbft

import (
	"encoding/binary"
	"fmt"
	"io"
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

// 节点使用的tcp监听
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
		go p.handleConn(conn)
		// time.Sleep(10 * time.Millisecond)
		// b, err := ioutil.ReadAll(conn)
		// if err != nil {
		// 	log.Println(err)
		// 	continue
		// }
		// conn.Close()
		// p.handleRequest(b)
	}
}

// // 读取连接数据, \n 做结束符
// func (p *Pbft) handleConn(conn net.Conn) {
// 	defer conn.Close()
// 	reader := bufio.NewReader(conn)
// 	for {
// 		b, err := reader.ReadBytes('\n')
// 		switch err {
// 		case nil:
// 			p.handleRequest(b)
// 		case io.EOF:
// 			log.Println("client closed the connection by terminating the process")
// 			return
// 		default:
// 			log.Printf("error: %v\n", err)
// 			return
// 		}
// 	}
// }

// 读取连接数据，利用长度前缀防粘包
func (p *Pbft) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		// 创建一个字节切片来存储消息长度前缀
		lengthPrefix := make([]byte, 4)

		// 读取消息长度前缀
		if _, err := conn.Read(lengthPrefix); err != nil {
			log.Fatal("Error reading message length prefix:", err.Error())
		}

		// 将消息长度前缀解析为一个无符号整数
		length := binary.BigEndian.Uint32(lengthPrefix)

		// 创建一个字节切片来存储消息内容
		message := make([]byte, length)

		// 读取消息内容
		if n, err := io.ReadFull(conn, message); err != nil {
			if err == io.ErrUnexpectedEOF {
				log.Fatal("Error reading message:", err.Error(), n)
			}else {
				log.Fatal("Error reading message:", err.Error())
			}
			
		}

		p.handleRequest(message)
	}
}
