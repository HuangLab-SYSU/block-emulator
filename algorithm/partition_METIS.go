package algorithm

import (
	"blockEmulator/account"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// METIS算法状态，state of constraint label propagation algorithm
type METISState struct {
	NetGraph     Graph          // 需运行METIS算法的图
	PartitionMap map[Vertex]int // 记录分片信息的 map，某个节点属于哪个分片
	Alpha        float64        // alpha
	ShardNum     int            // 分片数目
	AvgWeight    float64        // 平均分片负载
}

// 加入节点
func (ls *METISState) AddVertex(v Vertex) {
	ls.NetGraph.AddVertex(v)
}

// 加入边，uni=1代表单向边，否则双向边
func (ls *METISState) AddEdge(u, v Vertex, uni int) {
	// 如果没有点，则增加点
	if _, ok := ls.NetGraph.VertexSet[u]; !ok {
		ls.AddVertex(u)
	}
	if _, ok := ls.NetGraph.VertexSet[v]; !ok {
		ls.AddVertex(v)
	}
	ls.NetGraph.AddEdge(u, v, uni)
}

// 输出METIS
func (ls *METISState) PrintMETIS() {
	ls.NetGraph.PrintGraph()
	println()
}

// 设置参数
func (ls *METISState) Init_METISState(alpha float64, sn int) {
	ls.Alpha = alpha
	ls.ShardNum = sn
	ls.PartitionMap = make(map[Vertex]int)
}

// 将图写入txt文件
func (ls *METISState) Write_to_txt() {
	// 节点数量  边数（重复交易算1边）
	// 1 code1 weight1 code2 weight2 ....
	// 1 code1 weight1 code2 weight2 ....
	// 1 code1 weight1 code2 weight2 ....
	// ....

	if len(ls.NetGraph.VertexSet) != len(ls.NetGraph.Vertexs) {
		log.Panic("图节点数不等！！")
	}
	nodeNum := len(ls.NetGraph.Vertexs)
	edgeNum := 0
	for _, edge := range ls.NetGraph.EdgeWeight {
		edgeNum += len(edge)
	}
	if edgeNum%2 != 0 {
		log.Panic("边算错啦！")
	}
	edgeNum /= 2

	txtFile, err := os.Create("./sampleGraph0.txt")
	if err != nil {
		log.Panic(err)
	}
	defer txtFile.Close()
	txtlog := bufio.NewWriter(txtFile)
	s := fmt.Sprintf("%v %v\n", nodeNum, edgeNum)
	txtlog.WriteString(s)
	txtlog.Flush()

	for _, v1 := range ls.NetGraph.EdgeWeight {
		s = "1"
		for key, v := range v1 {
			s = s + " " + strconv.Itoa(key) + " " + strconv.Itoa(v)
		}
		s += "\n"
		txtlog.WriteString(s)
		txtlog.Flush()
	}

}

// 执行Linux命令运行metis
func (ls *METISState) Metis_Shell(Input, Output, File_path string, shardNum int) {
	fmt.Println("执行MetisCpp\n")
	cmd := exec.Command(File_path, Input, Output, strconv.Itoa(shardNum))
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(output))
	fmt.Println("执行MetisCpp完毕\n")

}

// METIS 划分算法
func (ls *METISState) METIS_Partition() ([]string, map[string]int) {
	res := make(map[string]int)
	addrs := []string{}

	ls.Write_to_txt()
	InPut := "sampleGraph0.txt"
	OutPut := "MetisPartionGraph0.txt"
	ls.Metis_Shell(InPut, OutPut, "METIS/partition", ls.ShardNum)

	file, err := os.Open(OutPut)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		strArr := strings.Split(line, " ")
		addr, _ := strconv.Atoi(strArr[0])
		shard, _ := strconv.Atoi(strArr[1])
		vertex := ls.NetGraph.Vertexs[addr]
		if shard != account.Addr2Shard(vertex.Addr) {
			addrs = append(addrs, vertex.Addr)
			res[vertex.Addr] = shard
		}
	}
	return addrs, res
}
