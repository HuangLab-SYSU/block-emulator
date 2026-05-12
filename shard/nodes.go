// definition of node and shard

package shard

import (
	"fmt"
)

type Node struct {
	NodeID  uint64
	ShardID uint64
	IPaddr  string
}

func (n *Node) PrintNode() {
	v := []interface{}{
		n.NodeID,
		n.ShardID,
		n.IPaddr,
	}
	fmt.Printf("%v\n", v)
}
