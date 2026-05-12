# CLPA算法详解

目前在 BlockEmulator 中配置了账户划分的功能，默认可选的账户划分配置算法叫 CLPA（Constrained Label Propagation Algorithm ），此方法来自于 SRDS2022 收录的论文《[Achieving Scalability and Load Balance across Blockchain Shards for State Sharding](https://ieeexplore.ieee.org/document/9996899)》。由于区块链的账户划分方式基本是固定的，这导致有些分片会出现交易过载（热分片），而有些分片内存在的交易数量寥寥的情况。在这种时候，一个合适的账户重新划分算法就能够动态的调整账户所在的分片，以达到降低跨分片交易数量的效果。

## CLPA 算法 - 数据结构

-   CLPA 算法有关的部分均在 **package partition** 中实现。

-   CLPA 算法本质上是一个图划分算法，因此实现 CLPA 算法之前首先需要实现一个 Graph 的数据类，：

```Go
// 图中的顶点
type Vertex struct {
    Addr string // 节点具体属性
    // 其他属性待补充
}
// 图
type Graph struct {
    VertexSet map[Vertex]bool     // 顶点集合
    EdgeSet   map[Vertex][]Vertex // 记录顶点与顶点间是否存在边，邻接表形式
}
```

同时，定义如下的 CLPAState 类 (以论文原文为具体标准)：

```Go
type CLPAState struct {
    NetGraph          Graph          // 需运行CLPA算法的图
    PartitionMap      map[Vertex]int // 记录分片信息的 map，某个账户属于哪个分片
    Edges2Shard       []int          // Shard 相邻接的边数，对应论文中的 total weight of edges associated with label k
    VertexsNumInShard []int          // Shard 内节点的数目
    WeightPenalty     float64        // 权重惩罚，对应论文中的 beta
    MinEdges2Shard    int            // 最少的 Shard 邻接边数，最小的 total weight of edges associated with label k
    MaxIterations     int            // 最大迭代次数，constraint，对应论文中的\tau
    CrossShardEdgeNum int            // 跨分片边的总数
    ShardNum          int            // 分片数目
}
```

需要说明的是，将 CLPA 算法应用到分片区块链场景下时，顶点（Vertex）指的是账户（account），边（edge）指的是交易（transaction），重新调整图的过程其实就是重新调整账户在不同分片的过程。

## CLPA 算法函数

上述的 CLPAState 类中定义了一些方法，主要对外调用的方法有：

```Go
// 加入顶点，需要将它默认归到一个分片中
func (cs *CLPAState) AddVertex(v Vertex)
// 加入边，需要将它的端点（如果不存在）默认归到一个分片中
func (cs *CLPAState) AddEdge(u, v Vertex)
// 设置参数
func (cs *CLPAState) Init_CLPAState(wp float64, mIter, sn int)
// CLPA 划分算法
func (cs *CLPAState) CLPA_Partition() (map[string]uint64, int)
```

## CLPA 算法建立流程

通过上述的 CLPA 算法函数的描述，用户可以轻松地调用 CLPA 算法。

1. 首先，用户先要创建CLPAState类，并且设置参数（weightPenalty=0.5，MaxIterations=100，与原论文保持一致）。
   1. ```Go
      cs = new(partition.CLPAState)
      cs.Init_CLPAState(weightPenalty_input, MaxIterations_input, shardNumber_input)
      ```
2. 接下来，用户将交易作为边，置入 CLPAState 中：
   1. ```Go
      var u, v Vertex
      cs.AddEdge(u, v) // 如果 u 或 v 不存在于当前 CLPAState 中，则顶点也会被置入CLPAState中，并以静态分片的方式置入对应分片
      ```
3. 插入若干条边之后，用户可以运行 CLPA 算法，以得到更新后的账户划分表和该情况下跨分片交易的数目：
   1. ```Go
      // dirtyMap: 算法返回的是更新后的 dirtyMap，一个账户存在于此 Map 中，仅当它在本次更新中改变了所属分片
      // crossTxNum: 该划分情况下的跨分片交易数目
      dirtyMap, crossTxNum := cs.CLPA_Partition()
      ```
4. 由于 dirtyMap 仅仅记录了某一次运行 CLPA 算法更新后的结果。如果需要多次使用 CLPA 算法，算法结果的接收方应当维护一个记录着所有运行结果的 PartitionMap，每次接收到 dirtyMap，PartitionMap 都会得到更新：
   1. ```Go
      var PartitionMap map[string]uint64 // 保存了自 CLPA 算法运行以来的所有账户变动（每个账户的所属分片被维护为最新版本）
      for key, val := range dirtyMap {
          PartitionMap[key] = val
      }
      ```

## CLPA Blockchain 测试样例运行

### **1. Supervisor设置**

按照如下方法创建 Supervisor，即可以让它具备 CLPA committee 的功能。

```Go
// 构建 Supervisor 对系统进行监听。
// 参数3 CLPA 指定了 Supervisor 具备 CLPA committee 的功能
var spv Supervisor
spv.NewSupervisor(supervisor_ip, chainConfig, "CLPA", testModules...)
```

## **2. Worker Nodes 设置**

按照如下方法创建 Supervisor，即可以让它们能够接收、处理 CLPA 相关的消息。

```Go
// 构建 Supervisor 对系统进行监听。
// 参数3 CLPA 指定了 Worker 能够接收、处理 CLPA 相关的消息
var worker PbftConsensusNode
worker.NewPbftNode(sid, nid, pcc, "CLPA").
```

### **3. 批处理文件运行**

通过 .bat 文件运行多个 Worker Nodes 和 Supervisor 即可。
