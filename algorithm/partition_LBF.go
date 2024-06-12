package algorithm

import (
	"blockEmulator/account"
	"fmt"
	"math/rand"
	"time"
)

// LBF算法状态，state of constraint label propagation algorithm
type LBFState struct {
	NetGraph          Graph          // 需运行LBF算法的图
	PartitionMap      map[Vertex]int // 记录分片信息的 map，某个节点属于哪个分片
	Alpha             float64        // alpha
	ShardNum          int            // 分片数目
	AvgWeight         float64        // 平均分片负载
}

// 加入节点
func (ls *LBFState) AddVertex(v Vertex) {
	ls.NetGraph.AddVertex(v)
}

// 加入边，uni=1代表单向边，否则双向边
func (ls *LBFState) AddEdge(u, v Vertex, uni int) {
	// 如果没有点，则增加点
	if _, ok := ls.NetGraph.VertexSet[u]; !ok {
		ls.AddVertex(u)
	}
	if _, ok := ls.NetGraph.VertexSet[v]; !ok {
		ls.AddVertex(v)
	}
	ls.NetGraph.AddEdge(u, v, uni)
}

// 输出LBF
func (ls *LBFState) PrintLBF() {
	ls.NetGraph.PrintGraph()
	println()
}

// 根据当前交易关系图，计算平均分片负载（所有节点权重加起来 / ShardNum），同时也计算每个节点负载
func (ls *LBFState) ComputeAvgWeight() {
	sum := 0
	for _, neighbors := range ls.NetGraph.EdgeSet {
		sum += len(neighbors)
	}
	ls.AvgWeight = float64(sum) / float64(ls.ShardNum)
}

// 设置参数
func (ls *LBFState) Init_LBFState(alpha float64, sn int) {
	ls.Alpha = alpha
	ls.ShardNum = sn
	ls.PartitionMap = make(map[Vertex]int)
}

// LBF 划分算法
func (ls *LBFState) LBF_Partition() ([]string, map[string]int) {
	//时间戳作为种子，所以每次随机数都是随机的
	rand.Seed(time.Now().UnixNano())

	ls.ComputeAvgWeight()
	fmt.Println(ls.AvgWeight)
	res := make(map[string]int)
	addrs := []string{}
	ls.PartitionMap = make(map[Vertex]int)
	for i := 0; i < ls.ShardNum-1; i++ {
		cur_sum := 0.0
		for cur_sum < ls.AvgWeight && len(ls.NetGraph.Vertexs) > 0 {
			//随机生成正整数
			index := rand.Intn(len(ls.NetGraph.Vertexs))
			vertex := ls.NetGraph.Vertexs[index]
			ls.PartitionMap[vertex] = i
			if i != account.Addr2Shard(vertex.Addr) {
				addrs = append(addrs, vertex.Addr)
				res[vertex.Addr] = i
			}
			cur_sum += float64(len(ls.NetGraph.EdgeSet[vertex]))
			ls.NetGraph.Vertexs = append(ls.NetGraph.Vertexs[:index], ls.NetGraph.Vertexs[index+1:]...)
		}
	}
	// 剩下的节点给到最后一个分片
	for _,vertex := range ls.NetGraph.Vertexs {
		ls.PartitionMap[vertex] = ls.ShardNum-1
		if ls.ShardNum - 1 !=  account.Addr2Shard(vertex.Addr) {
			addrs = append(addrs, vertex.Addr)
			res[vertex.Addr] = ls.ShardNum-1 
		}
	}
	return addrs, res
}
