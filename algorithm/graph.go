// 图的相关操作
package algorithm

// 图中的结点，即区块链网络中参与交易的账户
type Vertex struct {
	Addr string // 账户地址
	// 其他属性待补充
}

// 描述当前区块链交易集合的图
type Graph struct {
	VertexSet map[Vertex]int // 节点集合，其实是 set
	Vertexs []Vertex
	EdgeSet map[Vertex][]Vertex // 记录节点与节点间是否存在交易，邻接表
	EdgeWeight []map[int]int   // 记录节点与节点间交易权重，邻接表（【code】map【code】weight）
	// lock      sync.RWMutex       //锁，但是每个储存节点各自存储一份图，不需要此
}

// 创建节点
func (v *Vertex) ConstructVertex(s string) {
	v.Addr = s
}

// 增加图中的点
func (g *Graph) AddVertex(v Vertex) {
	if g.VertexSet == nil {
		g.VertexSet = make(map[Vertex]int)
	}
	if _,ok := g.VertexSet[v]; !ok {
		g.VertexSet[v] = len(g.Vertexs)
		g.Vertexs = append(g.Vertexs, v)
		g.EdgeWeight = append(g.EdgeWeight, make(map[int]int))
	}
}

// 增加图中的边, uni=1代表单向边
func (g *Graph) AddEdge(u, v Vertex, uni int) {
	// 如果没有点，则增加边，权恒定为 1
	if _, ok := g.VertexSet[u]; !ok {
		g.AddVertex(u)
	}
	if _, ok := g.VertexSet[v]; !ok {
		g.AddVertex(v)
	}
	if g.EdgeSet == nil {
		g.EdgeSet = make(map[Vertex][]Vertex)
	}
	if uni==0 {
		// 无向图，使用双向边，u也要指向v
		g.EdgeSet[u] = append(g.EdgeSet[u], v)
		g.EdgeWeight[g.VertexSet[u]][g.VertexSet[v]]++
	}
	g.EdgeSet[v] = append(g.EdgeSet[v], u)
	g.EdgeWeight[g.VertexSet[v]][g.VertexSet[u]]++
}

// 复制图
func (dst *Graph) CopyGraph(src Graph) {
	dst.VertexSet = make(map[Vertex]int)
	for key,val := range src.VertexSet {
		dst.VertexSet[key] = val
	}
	if src.EdgeSet != nil {
		dst.EdgeSet = make(map[Vertex][]Vertex)
		for v := range src.VertexSet {
			dst.EdgeSet[v] = make([]Vertex, len(src.EdgeSet[v]))
			copy(dst.EdgeSet[v], src.EdgeSet[v])
		}
	}
}

// 输出图
func (g Graph) PrintGraph() {
	for v := range g.VertexSet {
		print(v.Addr, " ")
		print("edge:")
		for _, u := range g.EdgeSet[v] {
			print(" ", u.Addr, "\t")
		}
		println()
	}
	println()
}
