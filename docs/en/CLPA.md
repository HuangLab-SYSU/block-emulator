The detailed explanation of CLPA algorithm

- It is the sixth section of the second chapter of the **BlockEmulator** English introduction document
# CLPA algorithm---data structure

The default optional account division configuration algorithm is called CLPA (Constrained Label Propagation Algorithm), which is derived from the paper *Achieving Scalability and Load Balance across Blockchain Shards for State Sharding* included in SRDS2022.


The account division method of the blockchain is basically fixed incurs the phenomenon in which some shards is transaction overloaded (hot sharding) while some shards have feww transactions.
 At this time, a suitable account redivision algorithm can dynamically adjust the shards in which the account is located to achieve the effect of reducing the number of cross-shard transactions.

 ## 1.1 The related parts of CLPA is implemented in **package partition**

 ## 1.2   The CLPA algorithm is essentially a graph division algorithm, so before implementing the CLPA algorithm, you first need to implement a graph data class:

 ```
 1 // The vertex of graph
2 type Vertex struct {
3    Addr string // The attribute of node 
4    // other attributes
5}
6 // graph
7 type Graph struct {
8    VertexSet map[Vertex]bool     // the set of vertex
9    EdgeSet   map[Vertex][]Vertex // to record whether there exists edge between vertexs with the form of adjacency table
10 }
```

 ## 1.3 Meanwhile, to define the follows class of CLPAState:

 ```
1 type CLPAState struct {
2    NetGraph          Graph          // The graph that needs to execute CLPA algorithm 
3    PartitionMap      map[Vertex]int // The map records the information of shards (the account belongs to which shards)
4   Edges2Shard       []int          // the number of adjacency edges of Shard
5    VertexsNumInShard []int          // the number of ndoes of Shard 
6    WeightPenalty     float64        // the weight of penalty
7    MinEdges2Shard    int            // minimum  number of adjacency edges of Shard
8    MaxIterations     int            // maximum iterations
9    CrossShardEdgeNum int            // The total number of cross-shard edges
10    ShardNum          int            // The number of shard
11 }
 ```
 It should be noted that when applying the CLPA algorithm to the sharded blockchain scenario, **Vertex** refers to the account, **edge** refers to the transaction, and the process of readjusting the graph is actually the process of readjusting the account in different shards.

 # CLPA algorithm---function

The class **CLPAState** defines some methods, some of which are used to external calls:

```
1 // To join an vertex, it needs to be grouped into a shard by default
2 func (cs *CLPAState) AddVertex(v Vertex)
3 // To join an edge, its endpoints (if they do not exist) need to be grouped into a shard by default
4 func (cs *CLPAState) AddEdge(u, v Vertex)
5 // set up parameters
6 func (cs *CLPAState) Init_CLPAState(wp float64, mIter, sn int)
7 // CLPA partition algorithm 
8 func (cs *CLPAState) CLPA_Partition() (map[string]uint64, int)
```

 # CLPA algorithm---the procession of establishment
 According to the above description of CLPA algorithm function, users can call CLPA algorithm easily.

1. Firstly, the user must first create a CLPAState class and set the parameters (weightPenalty=0.5, MaxIterations=100, consistent with the original paper).
```
1 cs = new(partition.CLPAState)
2 cs.Init_CLPAState(weightPenalty_input, MaxIterations_input,shardNumber_input)
```
2. Next, the user places the transaction as an edge in CLPAState:
```
1 var u, v Vertex
2 cs.AddEdge(u, v) // If $u$ or $v$ does not exist in the current CLPAState, the vertex is also placed in CLPAState and the corresponding shard is placed as a static shard
```

3. After inserting several edges, the user can run the CLPA algorithm to get the updated account division table and the number of cross-shard transactions in this case:

```
1 // dirtyMap: The return of this algorithm is the updated **dirtyMap**, in which exists an account.
2 // crossTxNum: The number of corss-shard transaction under this scenario
3 dirtyMap, crossTxNum := cs.CLPA_Partition()
```

4. Since dirtyMap only records the results of a CLPA algorithm update run once. If the CLPA algorithm needs to be used multiple times, the receiver of the algorithm result should maintain a PartitionMap that records all the results of the run, and the PartitionMap will be updated each time a dirtyMap is received:

```
1 var PartitionMap map[string]uint64 // save the changes of all accounts since the CLPA algorithm run.
2 for key, val := range dirtyMap {
3   PartitionMap[key] = val
4 }
```

 # CLPA Blockchain---demo


 ## The settings of **Supervisor**
Create a supervisor as follows, you can make it have the functionality of the CLPA committee:
```
1 // Create **Supervisor** to monitor the system
2 // Parameter 3, CLPA, specifies that the supervisor has the functionality of the CLPA committee
3 var spv Supervisor
4 spv.NewSupervisor(supervisor_ip, chainConfig, "CLPA", testModules...)
```
 ## The settings of **Worker Nodes**
Create a supervisor as follows, you can make them receive and deal with the message of CLPA.
```
// Create **Supervisor** to monitor the system
//Parameter 3, CLPA, specifies that the Worker has the functionality of the CLPA receiving and procession
var worker PbftConsensusNode
worker.NewPbftNode(sid, nid, pcc, "CLPA").
```
 ## The execution of bat files
To run the multiple **Worker Nodes** and **Supervisor** by **.bat** files.
