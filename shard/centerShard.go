package shard

import (
	"blockEmulator/params"
	"blockEmulator/pbft"
)

//	type ReconfigProposal struct {
//		params.NodeTable
//	}

func NewCenterNode() *Node {
	node := new(Node)
	node.P = pbft.NewPBFT()
	config := params.Config

	go node.P.TcpListen() //启动节点

	if config.NodeID == "N0" {
		pbft.NewCenterLog(config.ShardID)
		// time.Sleep(time.Duration(config.Reconfig_interval) * time.Millisecond)
		// node.P.CenterNode()
	}

	return node
}
