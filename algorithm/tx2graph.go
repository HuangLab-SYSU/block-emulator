package algorithm

import (
	"blockEmulator/core"
	"encoding/hex"
)

func Pagerank_Tx2graph_And_Addrs(txs []*core.Transaction) (map[string]map[string]int, []string) {
	graph := make(map[string]map[string]int)
	Addr_Exist := make(map[string]bool)
	addrs := make([]string, 0)

	for _, tx := range txs {
		from, to := hex.EncodeToString(tx.Sender), hex.EncodeToString(tx.Recipient)
		if _, ok := graph[from]; !ok {
			graph[from] = make(map[string]int)
		}

		if _, ok := graph[from][to]; !ok {
			graph[from][to] = 1
		} else {
			graph[from][to]++
		}

		if _, ok := graph[to]; !ok {
			graph[to] = make(map[string]int)
		}

		if _, ok := graph[to][from]; !ok {
			graph[to][from] = 1
		} else {
			graph[to][from]++
		}

		if !Addr_Exist[from] {
			Addr_Exist[from] = true
			addrs = append(addrs, from)
		}

		if !Addr_Exist[to] {
			Addr_Exist[to] = true
			addrs = append(addrs, to)
		}
	}

	return graph, addrs
}
