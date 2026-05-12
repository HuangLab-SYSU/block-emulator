package main

import (
	"blockEmulator/build"
	"fmt"

	"github.com/spf13/pflag"
)

var (
	shardNum int
	nodeNum  int
	shardID  int
	nodeID   int
	modID    int
	isClient bool
	isGen    bool
)

var modeNames = map[int]string{
	0: "CLPA_Broker",
	1: "CLPA",
	2: "Broker",
	3: "Relay",
	4: "Broker_b2e",
}

func modeName(id int) string {
	if name, ok := modeNames[id]; ok {
		return name
	}
	return "Unknown"
}

func main() {
	pflag.IntVarP(&shardNum, "shardNum", "S", 2, "indicate that how many shards are deployed")
	pflag.IntVarP(&nodeNum, "nodeNum", "N", 4, "indicate how many nodes of each shard are deployed")
	pflag.IntVarP(&shardID, "shardID", "s", 0, "id of the shard to which this node belongs, for example, 0")
	pflag.IntVarP(&nodeID, "nodeID", "n", 0, "id of this node, for example, 0")
	pflag.IntVarP(&modID, "modID", "m", 3, "choice Committee Method,for example, 0, [CLPA_Broker,CLPA,Broker,Relay,Broker_b2e] ")
	pflag.BoolVarP(&isClient, "client", "c", false, "whether this node is a client")
	pflag.BoolVarP(&isGen, "gen", "g", false, "generation bat")
	pflag.Parse()

	if isGen {
		fmt.Printf("[BlockEmulator] Generating launch scripts: shardNum=%d, nodesInShard=%d, mode=%d (%s)\n",
			shardNum, nodeNum, modID, modeName(modID))
		build.GenerateBatFile(nodeNum, shardNum, modID)
		build.GenerateShellFile(nodeNum, shardNum, modID)
		fmt.Println("[BlockEmulator] Done. Batch files (.bat for Windows, .sh for Linux/macOS) have been written to the current directory.")
		return
	}

	if isClient {
		fmt.Printf("[BlockEmulator] Starting Supervisor: shardNum=%d, nodesInShard=%d, mode=%d (%s)\n",
			shardNum, nodeNum, modID, modeName(modID))
		build.BuildSupervisor(uint64(nodeNum), uint64(shardNum), uint64(modID))
	} else {
		fmt.Printf("[BlockEmulator] Starting node #%d in shard #%d (shardNum=%d, nodesInShard=%d, mode=%d (%s))\n",
			nodeID, shardID, shardNum, nodeNum, modID, modeName(modID))
		build.BuildNewPbftNode(uint64(nodeID), uint64(nodeNum), uint64(shardID), uint64(shardNum), uint64(modID))
	}
}
