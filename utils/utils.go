package utils

import (
	"blockEmulator/params"
	"log"
	"net"
	"strconv"
)

func Addr2Shard(senderAddr string) int {
	// 只取地址后五位已绝对够用
	senderAddr = senderAddr[len(senderAddr)-5:]
	num, err := strconv.ParseInt(senderAddr, 16, 32)
	// num, err := strconv.ParseInt(senderAddr, 10, 32)
	if err != nil {
		log.Panic()
	}
	return int(num) % params.Config.Shard_num
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//使用tcp发送消息
func TcpDial(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect error", err)
		return
	}

	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
