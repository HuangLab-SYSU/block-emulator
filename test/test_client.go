package test

import (
	"blockEmulator/pbft"
	"log"
	"os"
)

func Test_client() {
	if len(os.Args) != 2 {
		log.Panic("输入的参数不对！")
	}
	pbft.RunClient(os.Args[1])
}
