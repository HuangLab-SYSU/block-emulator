package test

import (
	"blockEmulator/params"
	"blockEmulator/pbft"
	"blockEmulator/shard"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	node          *shard.Node
	shard_num     int
	shardID       string
	malicious_num int
	nodeID        string
	testFile      string
	isClient      bool
	// requestlog    *csv.Writer
	EndTime int64
)

func Test_shard() {
	flag.IntVarP(&shard_num, "shard_num", "S", 1, "indicate that how many shards are deployed")
	flag.StringVarP(&shardID, "shardID", "s", "", "id of the shard to which this node belongs, for example, S0")
	flag.IntVarP(&malicious_num, "malicious_num", "f", 1, "indicate the maximum of malicious nodes in one shard")
	flag.StringVarP(&nodeID, "nodeID", "n", "", "id of this node, for example, N0")
	flag.StringVarP(&testFile, "testFile", "t", "", "path of the input test file")
	flag.BoolVarP(&isClient, "client", "c", false, "whether this node is a client")

	flag.Parse()

	if isClient {
		if testFile == "" {
			log.Panic("参数不正确！")
		}
		// 修改全局变量 Config，之后其他地方会调用
		config := params.Config
		config.NodeID = nodeID
		config.ShardID = shardID
		config.Malicious_num = int(malicious_num)
		config.Shard_num = int(shard_num)
		pbft.RunClient(testFile)
		return
	}
	if shard_num == 1 {
		shardID = "S0"
	}
	if shardID == "" || nodeID == "" || testFile == "" {
		log.Panic("参数不正确！")
	}

	// 修改全局变量 Config，之后其他地方会调用
	config := params.Config
	config.NodeID = nodeID
	config.ShardID = shardID
	config.Malicious_num = int(malicious_num)
	config.Shard_num = int(shard_num)
	config.Path = testFile
	// for i := 0; i < 184379; i++ {
	// 	params.Init_addrs = append(params.Init_addrs, utils.Int2hexString(i))
	// }

	file, err := os.Open(config.Path)
	if err != nil {
		log.Panic()
	}
	// defer file.Close()

	r := csv.NewReader(file)
	_, err = r.Read()
	if err != nil {
		log.Panic()
	}


	// 初始化读取所有账户
	isExist := make(map[string]bool)
	for i:=0; i<1000000; i++{
	// for i:=0; i<500000; i++{
		row, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}
		senderstr, recipientstr := row[1][2:], row[2][2:]
		if path=="0to999999_BlockTransaction.csv" || path=="300W.csv"  || path=="100W.csv"  || path=="20W.csv"  || path=="50W.csv" || path=="200W.csv" {
			if row[5] != "None" || row[6] == "1" || row[7] == "1" || len(row[4][2:]) != 40 || len(row[3][2:]) != 40 || row[4]==row[3] {
				continue
			}
			senderstr, recipientstr = row[3][2:], row[4][2:]
		}

		if !isExist[senderstr] {
			isExist[senderstr] = true
			params.Init_addrs = append(params.Init_addrs, senderstr)
		}
		if !isExist[recipientstr] {
			isExist[recipientstr] = true
			params.Init_addrs = append(params.Init_addrs, recipientstr)
		}
	}
	isExist = nil

	// if config.NodeID == "N0" {
	// 	// config.Path = testFile
	// 	// csvFile, err := os.Create("./log/" + shardID + "_requesttime.csv")
	// 	// if err != nil {
	// 	// 	log.Panic(err)
	// 	// }
	// 	// requestlog = csv.NewWriter(csvFile)
	// 	// requestlog.Write([]string{"txid", "waiting_time", "is_ctx", "1st_queueing_time", "2nd_queueing_time"})
	// 	// requestlog.Flush()
	// }

	if _, ok := params.NodeTable[shardID][nodeID]; ok {
		node = shard.NewNode()
	} else {
		log.Fatal("无此节点编号！")
	}

	file.Close()


	<-node.P.Stop
	fmt.Printf("节点收到终止节点消息，停止运行\n")

	// node.P.Node.CurChain.StatusTrie.PrintState()
	// fmt.Println(account.Account2Shard)
	// fmt.Println(account.AccountInOwnShard)
	// for _,v := range node.P.Node.CurChain.Tx_pool.Queue {
	// 	v.PrintTx()
	// }
}
