package pbft

import (
	"blockEmulator/params"
	"blockEmulator/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func RunClient(path string) {
	tx_cnt := CountTx(path)
	finished := make([]bool, tx_cnt)
	finished_cnt := 0
	all_finish := make(chan int)

	// 消息处理器
	var handle func([]byte) = func(data []byte) {
		cmd, content := splitMessage(data)
		switch command(cmd) {
		case cReply:
			ids := new([]int)
			err := json.Unmarshal(content, ids)
			if err != nil {
				log.Panic(err)
			}
			for _, id := range *ids {
				if !finished[id] {
					finished[id] = true
					finished_cnt += 1
					if finished_cnt == tx_cnt {
						all_finish <- 1
						break
					}
				}
			}
		default:
			log.Panic()
		}

	}

	// 消息监听器
	var listener func() = func() {
		listen, err := net.Listen("tcp", params.ClientAddr)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("客户端开启监听，地址：%s\n", params.ClientAddr)
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
			handle(b)
		}
	}

	go listener()

	<-all_finish
	for shardID, nodes := range params.NodeTable {
		for nodeID, addr := range nodes {
			fmt.Printf("客户端向分片%v的节点%v发送终止运行消息\n", shardID, nodeID)
			m := jointMessage(cStop, nil)
			utils.TcpDial(m, addr)
		}
	}

}

func CountTx(path string) int {
	file, err := os.Open(path)
	if err != nil {
		log.Panic()
	}
	defer file.Close()
	r := csv.NewReader(file)
	_, err = r.Read()
	if err != nil {
		log.Panic()
	}
	cnt := 0
	for {
		_, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}
		cnt += 1
	}
	return cnt
}
