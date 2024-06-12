package test

import (
	"blockEmulator/chain"
	"blockEmulator/params"
	"fmt"

	flag "github.com/spf13/pflag"
)

func Test_node() {
	flag.IntVarP(&shard_num, "shard_num", "S", 1, "indicate that how many shards are deployed")
	flag.StringVarP(&shardID, "shardID", "s", "", "id of the shard to which this node belongs, for example, S0")
	flag.IntVarP(&malicious_num, "malicious_num", "f", 1, "indicate the maximum of malicious nodes in one shard")
	flag.StringVarP(&nodeID, "nodeID", "n", "", "id of this node, for example, N0")
	flag.StringVarP(&testFile, "testFile", "t", "", "path of the input test file")

	flag.Parse()

	// 修改全局变量 Config
	config := params.Config
	config.NodeID = nodeID
	config.ShardID = shardID
	config.Malicious_num = int(malicious_num)
	config.Shard_num = int(shard_num)

	bc, _ := chain.NewBlockChain(config)

	curBlock := bc.CurrentBlock
	fmt.Printf("curBlock: \n")
	curBlock.PrintBlock()

	// stateTree := bc.StatusTrie
	// fmt.Printf("stateTree: \n")
	// stateTree.PrintState()

}
