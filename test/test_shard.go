package test

import (
	"blockEmulator/params"
	"blockEmulator/pbft"
	"blockEmulator/shard"
	"fmt"
	"log"

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
		pbft.RunClient(testFile)
		return
	}
	if shard_num == 1 {
		shardID = "S0"
	}
	if shardID == "" || nodeID == "" || shardID != "SC" && nodeID == "N0" && testFile == "" {
		fmt.Println(nodeID)

		log.Panic("参数不正确！")
	}
	// 修改全局变量 Config，之后其他地方会调用
	config := params.Config
	config.NodeID = nodeID
	config.ShardID = shardID
	config.Malicious_num = int(malicious_num)
	config.Shard_num = int(shard_num)

	if config.NodeID == "N0" {
		config.Path = testFile
	}

	if _, ok := params.NodeTable[shardID][nodeID]; ok && shardID != "SC" {
		node = shard.NewNode()
	} else if shardID == "SC" {
		node = shard.NewCenterNode()
	} else {
		log.Fatal("无此节点编号！")
	}

	<-node.P.Stop
	fmt.Printf("节点收到终止节点消息，停止运行")
}

// func Test() {
// 	shard.StaleGenReconfigProposal()
// 	// NewGenReconfigProposal()
// }
