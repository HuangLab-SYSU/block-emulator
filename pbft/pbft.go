package pbft

import (
	"blockEmulator/chain"
	"blockEmulator/core"
	"blockEmulator/new_trie"
	"blockEmulator/params"
	"blockEmulator/utils"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	blocklog, txlog *csv.Writer
	queuelog        *csv.Writer
	timelog         *csv.Writer
	reconfigdatalog *csv.Writer
)

// //本地消息池（模拟持久化层），只有确认提交成功后才会存入此池
// var localMessagePool = []Message{}

type node struct {
	//节点ID
	nodeID string
	//节点监听地址
	addr     string
	CurChain *chain.BlockChain
	// 中心分片主节点 用于测量重组时间
	time1         time.Time //开始时间
	time2         time.Time //结束时间
	time_send     time.Time
	time_trans    time.Time
	time_receieve time.Time
}

type Pbft struct {
	//节点信息
	Node node
	//每笔请求自增序号
	sequenceID int
	//共识+重组为一个周期
	epochID int
	//锁
	lock sync.Mutex
	//确保消息不重复发送的锁
	sequenceLock sync.Mutex
	//确保落后节点对同一区块最多只向主节点请求一次
	requestLock sync.Mutex
	//临时消息池，消息摘要对应消息本体
	messagePool map[string]*Request
	//存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
	prePareConfirmCount map[string]map[string]bool
	//存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	commitConfirmCount map[string]map[string]bool
	//该笔消息是否已进行Commit广播
	isCommitBordcast map[string]bool
	//该笔消息是否已对客户端进行Reply
	isReply map[string]bool
	//区块高度到区块摘要的映射
	height2Digest map[int]string
	nodeTable     map[string]string
	Stop          chan int

	height2BlockHash map[int][]byte

	// a.确保共识结束后再开始重组周期 & 重组完成后再返回重组结束给新分片主节点
	reconfigLock sync.Mutex
	// 标记是否为共识周期
	// b.重组周期开始时关闭Propose协程(即共识周期)
	consensusFlag bool

	// 确保先发送重组信息再接收重组信息的锁
	reconfigBlockLock sync.Mutex
	// 中心分片主节点用于存放收到的reconfigStart数量，若满，则重组阶段开始
	reconfigStartCount map[string]bool
	// 存放收到的reconfigDone数量
	reconfigCount map[string]bool
	// 存放收到的nextEpochReply数量
	nextEpochCount int
	// // 若节点接收到重组完成reply，flag=1
	// flag bool
	// 重组开始
	reconfigStart chan int
	// 主节点收到其他节点的重组完成消息后再进行重组
	// 中心分片收到其他分片重组完成信息后再结束重组过程
	reconfigDone chan int
	// 重组完成后再修改本地nodetable
	reconfigOver chan int
	// 若节点固定时间内未接收到重组完成reply，重发
	time_to_reissue int64

	// 确保主节点处理完relay消息后再进行重组
	relay chan int
	// 节点用于存放收到的relay数量，若满，则重组阶段开始
	relayCount                int
	num_new_account           int
	addrMap                   map[string]int
	proposeBlock              chan int
	num_account_send          int
	num_account_all           int
	num_accounts_visit        int
	num_active_accounts_visit int
}

func NewPBFT() *Pbft {
	config := params.Config
	p := new(Pbft)
	p.Node.nodeID = config.NodeID
	p.Node.addr = params.NodeTable[config.ShardID][config.NodeID]

	p.Node.CurChain, _ = chain.NewBlockChain(config)
	p.sequenceID = p.Node.CurChain.CurrentBlock.Header.Number + 1
	p.messagePool = make(map[string]*Request)
	p.prePareConfirmCount = make(map[string]map[string]bool)
	p.commitConfirmCount = make(map[string]map[string]bool)
	p.isCommitBordcast = make(map[string]bool)
	p.isReply = make(map[string]bool)
	p.height2Digest = make(map[int]string)
	p.nodeTable = params.NodeTable[config.ShardID]
	p.Stop = make(chan int, 0)
	p.height2BlockHash = make(map[int][]byte)

	p.epochID = 0
	p.reconfigStartCount = make(map[string]bool)
	p.reconfigCount = make(map[string]bool)
	p.nextEpochCount = 0
	// p.reconfigReply = false
	// p.flag = false
	p.reconfigStart = make(chan int, 0)
	p.reconfigDone = make(chan int, 0)
	p.reconfigOver = make(chan int, 0)
	p.consensusFlag = true
	p.time_to_reissue = 60

	p.relay = make(chan int, 0)
	p.relayCount = 0
	p.num_new_account = 0
	p.addrMap = make(map[string]int)
	p.proposeBlock = make(chan int, 0)
	p.num_account_send = 0
	p.num_account_all = 0
	p.num_accounts_visit = 0
	p.num_active_accounts_visit = 0
	return p
}

func NewLog(shardID string) {
	csvFile, err := os.Create("./log/" + shardID + "_block.csv")
	if err != nil {
		log.Panic(err)
	}

	blocklog = csv.NewWriter(csvFile)
	blocklog.Write([]string{"timestamp", "blockHeight", "tx_total", "tx_normal", "tx_relay_first_half", "tx_relay_second_half", "num_new_account", "num_send_account", "num_all_account", "num_accounts_visit", "num_active_accounts_visit"})
	blocklog.Flush()

	csvFile, err = os.Create("./log/" + shardID + "_transaction.csv")
	if err != nil {
		log.Panic(err)
	}
	txlog = csv.NewWriter(csvFile)
	txlog.Write([]string{"txid", "blockHeight", "request_time", "commit_time", "delay"})
	txlog.Flush()

	csvFile, err = os.Create("./log/" + shardID + "_queue_length.csv")
	if err != nil {
		log.Panic(err)
	}
	queuelog = csv.NewWriter(csvFile)
	queuelog.Write([]string{"timestamp", "queue_length"})
	queuelog.Flush()

	csvFile, err = os.Create("./log/" + shardID + "_reconfig_data.csv")
	if err != nil {
		log.Panic(err)
	}
	reconfigdatalog = csv.NewWriter(csvFile)
	reconfigdatalog.Write([]string{"block_size", "state_size"})
	reconfigdatalog.Flush()

	csvFile.Close()
}
func NewCenterLog(shardID string) {
	csvFile, err := os.Create("./log/" + "reconfig_time_per_epoch.csv")
	if err != nil {
		log.Panic(err)
	}
	timelog = csv.NewWriter(csvFile)
	timelog.Write([]string{"Shard", "reconfig_time_per_epoch", "time_send", "time_trans", "time_receieve", "time_consensus"})
	timelog.Flush()

	csvFile, err = os.Create("./log/" + shardID + "_reconfig_data.csv")
	if err != nil {
		log.Panic(err)
	}
	reconfigdatalog = csv.NewWriter(csvFile)
	reconfigdatalog.Write([]string{"reconfigData_size"})
	reconfigdatalog.Flush()

	csvFile.Close()
}

func OpenLog(shardID string) {
	csvFile, err := os.OpenFile("./log/"+shardID+"_block.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	blocklog = csv.NewWriter(csvFile)

	csvFile, err = os.OpenFile("./log/"+shardID+"_transaction.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	txlog = csv.NewWriter(csvFile)

	csvFile, err = os.OpenFile("./log/"+shardID+"_queue_length.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	queuelog = csv.NewWriter(csvFile)

	// csvFile, err = os.OpenFile("./log/"+"reconfig_time_per_epoch.csv", os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// timelog = csv.NewWriter(csvFile)

	// csvFile.Close()
}

func (p *Pbft) handleRequest(data []byte) {
	//切割消息，根据消息命令调用不同的功能
	cmd, content := splitMessage(data)
	fmt.Println(command(cmd))
	switch command(cmd) {
	case cPrePrepare:
		p.handlePrePrepare(content)
	case cPrepare:
		p.handlePrepare(content)
	case cCommit:
		p.handleCommit(content)
	case cRequestBlock:
		p.handleRequestBlock(content)
	case cSendBlock:
		p.handleSendBlock(content)
	case cRelay:
		go p.handleRelay(content)
	case cStop:
		p.Stop <- 1

	case cReconfigStart:
		go p.handleReconfigStart(content)
	case cReconfig:
		go p.handleReconfig_Stale(content)
		// go p.handleReconfig_New(content)
	// case cReconfigBlock:
	// 	go p.handleReconfigBlock(content)
	case cReconfigTries:
		go p.handleReconfigTries(content)
	// 中心分片重组传输数据
	case cReconfigCenter:
		go p.handleReconfigCenter(content)
	case cReconfigDone:
		go p.handleReconfigDone(content)
		// case cReconfigDoneReply:
		// 	go p.handleReconfigDoneReply(content)
	case cNextEpochStart:
		go p.NextEpochReply()
	case cNextEpochStartReply:
		go p.handleNextEpochReply(content)
	case cSendAccount:
		go p.handleSendAccount(content)
	case cHandleSendAccount:
		go p.handlehandleSendAccount(content)

	}

}

// 只有主节点可以调用
// 生成一个区块并发起共识
func (p *Pbft) Propose() {
	// p.reconfigFlag <- 1
	fmt.Printf("===当前周期为：%d ===\n", p.epochID)

	p.nodeTable = params.NodeTable[params.Config.ShardID]
	config := params.Config
	for i := 0; i < num_block_per_epoch; i++ {
		fmt.Println("Proposing a block …………")
		p.sequenceLock.Lock() //通过锁强制要求上一个区块commit完成新的区块才能被提出
		time.Sleep(time.Duration(config.Block_interval) * time.Millisecond)

		p.proposeBlock = make(chan int, 0)
		p.propose()

	}

}

func (p *Pbft) propose() {
	r := &Request{}
	r.Timestamp = time.Now().Unix()
	r.Message.ID = getRandom()

	block, num_accounts_visit, num_active_accounts_visit, accounts_send, accounts_all := p.Node.CurChain.GenerateBlock() // 打印len(pool.Queue)
	p.num_account_send += len(accounts_send)
	p.num_account_all += len(accounts_all)
	p.num_accounts_visit = num_accounts_visit
	p.num_active_accounts_visit = num_active_accounts_visit

	fmt.Println(accounts_send)
	fmt.Println(p.epochID)
	if p.epochID > 0 {
		p.SendAccounts(accounts_send, p.Node.addr)
		<-p.proposeBlock
	}

	encoded_block := block.Encode()

	r.Message.Content = encoded_block

	// //添加信息序号
	// p.sequenceIDAdd()
	//获取消息摘要
	digest := getDigest(r)
	fmt.Println("已将request存入临时消息池")
	//存入临时消息池
	p.messagePool[digest] = r

	//拼接成PrePrepare，准备发往follower节点
	pp := PrePrepare{r, digest, p.sequenceID}
	p.height2Digest[p.sequenceID] = digest
	b, err := json.Marshal(pp)
	if err != nil {
		// fmt.Println("Propose出错")
		log.Panic(err)
	}
	fmt.Println("正在向其他节点进行进行PrePrepare广播 ...")
	//进行PrePrepare广播
	p.broadcast(cPrePrepare, b)
	fmt.Println("PrePrepare广播完成")
}

// 处理预准备消息
func (p *Pbft) handlePrePrepare(content []byte) {
	fmt.Println("本节点已接收到主节点发来的PrePrepare ...")
	//	//使用json解析出PrePrepare结构体
	pp := new(PrePrepare)
	err := json.Unmarshal(content, pp)
	if err != nil {
		// fmt.Println("handlePrePrepare出错")
		log.Panic(err)
	}

	if digest := getDigest(pp.RequestMessage); digest != pp.Digest {
		fmt.Println("信息摘要对不上，拒绝进行prepare广播")
	} else if p.sequenceID != pp.SequenceID {
		fmt.Println("消息序号对不上，拒绝进行prepare广播")
	} else if !p.Node.CurChain.IsBlockValid(core.DecodeBlock(pp.RequestMessage.Message.Content)) {
		// todo
		fmt.Println("区块不合法，拒绝进行prepare广播")
	} else {
		// //序号赋值
		// p.sequenceID = pp.SequenceID
		p.height2Digest[p.sequenceID] = pp.Digest

		//将信息存入临时消息池
		fmt.Println("已将消息存入临时节点池")
		p.messagePool[pp.Digest] = pp.RequestMessage
		//拼接成Prepare
		pre := Prepare{pp.Digest, pp.SequenceID, p.Node.nodeID}
		bPre, err := json.Marshal(pre)
		if err != nil {
			// fmt.Println("handlePrePrepare出错")
			log.Panic(err)
		}
		//进行准备阶段的广播
		// fmt.Println("正在进行Prepare广播 ...")
		p.broadcast(cPrepare, bPre)
		// fmt.Println("Prepare广播完成")
	}
}

// 处理准备消息
func (p *Pbft) handlePrepare(content []byte) {
	//使用json解析出Prepare结构体
	pre := new(Prepare)
	err := json.Unmarshal(content, pre)
	if err != nil {
		// fmt.Println("handlePrepare出错")
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s节点发来的Prepare ... \n", pre.NodeID)

	if _, ok := p.messagePool[pre.Digest]; !ok {
		fmt.Println("当前临时消息池无此摘要，拒绝执行commit广播")
	} else if p.sequenceID != pre.SequenceID {
		fmt.Println("消息序号对不上，拒绝执行commit广播")
	} else {
		p.setPrePareConfirmMap(pre.Digest, pre.NodeID, true)
		count := 0
		for range p.prePareConfirmCount[pre.Digest] {
			count++
		}
		//因为主节点不会发送Prepare，所以不包含自己
		specifiedCount := 0
		if p.Node.nodeID == "N0" {
			specifiedCount = 2 * malicious_num
		} else {
			specifiedCount = 2*malicious_num - 1
		}
		//如果节点至少收到了2f个prepare的消息（包括自己）,并且没有进行过commit广播，则进行commit广播
		p.lock.Lock()
		//获取消息源节点的公钥，用于数字签名验证
		if count >= specifiedCount && !p.isCommitBordcast[pre.Digest] {
			fmt.Println("本节点已收到至少 2f 个节点(包括本地节点)发来的Prepare信息 ...")

			c := Commit{pre.Digest, pre.SequenceID, p.Node.nodeID}
			bc, err := json.Marshal(c)
			if err != nil {
				// fmt.Println("handlePrepare出错")
				log.Panic(err)
			}
			//进行提交信息的广播
			fmt.Println("正在进行commit广播")
			p.broadcast(cCommit, bc)
			p.isCommitBordcast[pre.Digest] = true
			// fmt.Println("commit广播完成")
			p.lock.Unlock()
			p.handleCommit(bc)

		} else {
			p.lock.Unlock()
		}
	}
}

// 处理提交确认消息
func (p *Pbft) handleCommit(content []byte) {
	//使用json解析出Commit结构体
	c := new(Commit)
	err := json.Unmarshal(content, c)
	if err != nil {
		// fmt.Println("handleCommit出错")
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s节点发来的Commit ... \n", c.NodeID)

	// if _, ok := p.prePareConfirmCount[c.Digest]; !ok {
	// 	fmt.Println("当前prepare池无此摘要，拒绝将信息持久化到本地消息池")
	// } else if p.sequenceID != c.SequenceID {
	// if p.sequenceID != c.SequenceID {
	// 	fmt.Println("消息序号对不上，拒绝将信息持久化到本地消息池")
	// } else {
	p.setCommitConfirmMap(c.Digest, c.NodeID, true)
	count := 0
	for range p.commitConfirmCount[c.Digest] {
		count++
	}
	//如果节点至少收到了2f+1个commit消息（包括自己）,并且节点没有回复过,并且已进行过commit广播，则提交信息至本地消息池，并reply成功标志至客户端！
	p.lock.Lock()
	require_cnt := malicious_num * 2
	if p.sequenceID != c.SequenceID {
		require_cnt += 1
	}

	if count > malicious_num*2 && !p.isReply[c.Digest] {
		fmt.Println("本节点已收到至少2f + 1 个节点(包括本地节点)发来的Commit信息 ...")
		//将消息信息，提交到本地消息池中！
		if _, ok := p.messagePool[c.Digest]; !ok {
			// 1. 如果本地消息池里没有这个消息，说明节点落后于其他节点，向主节点请求缺失的区块
			p.isReply[c.Digest] = true
			p.requestLock.Lock()
			p.requestBlocks(p.sequenceID, c.SequenceID)
		} else {
			// 2.本地消息池里有这个消息，将交易区块上链
			r := p.messagePool[c.Digest]
			encoded_block := r.Message.Content
			block := core.DecodeBlock(encoded_block)
			p.num_new_account += p.Node.CurChain.AddBlock(block, p.consensusFlag, p.epochID)
			fmt.Printf("编号为 %d 的区块已加入本地区块链！\n", p.sequenceID)
			curBlock := p.Node.CurChain.CurrentBlock
			fmt.Printf("curBlock: \n")
			curBlock.PrintBlock()
			fmt.Printf("Block打印完毕 \n")

			p.height2BlockHash[p.sequenceID] = block.GetHash()

			if p.epochID != 0 {
				for _, v := range block.Transactions {
					// 判断list中是否有该地址，若无则向中心分片请求
					// 对Sender
					if _, ok := p.addrMap[hex.EncodeToString(v.Sender)]; !ok {

					}
					p.addrMap[hex.EncodeToString(v.Sender)] = p.epochID

					//若交易接收者属于本分片才加入已上链交易集
					if params.ShardTable[params.Config.ShardID] == utils.Addr2Shard(hex.EncodeToString(v.Recipient)) {
						if _, ok := p.addrMap[hex.EncodeToString(v.Recipient)]; !ok {

						}
						p.addrMap[hex.EncodeToString(v.Recipient)] = p.epochID
					}

				}
			}

			if p.Node.nodeID == "N0" {

				OpenLog(params.Config.ShardID)
				tx_total := len(block.Transactions)
				now := time.Now().Unix()
				relayCount := 0
				//已上链交易集
				commit_ids := []int{}
				for _, v := range block.Transactions {
					//若交易接收者属于本分片才加入已上链交易集
					if params.ShardTable[params.Config.ShardID] == utils.Addr2Shard(hex.EncodeToString(v.Recipient)) {
						commit_ids = append(commit_ids, v.Id)
						s := fmt.Sprintf("%v %v %v %v %v", v.Id, block.Header.Number, v.RequestTime, now, now-v.RequestTime)
						txlog.Write(strings.Split(s, " "))
					}
					//若交易接收者不属于本分片，转发至其他分片
					if utils.Addr2Shard(hex.EncodeToString(v.Sender)) != utils.Addr2Shard(hex.EncodeToString(v.Recipient)) {
						relayCount++
					}
				}
				txlog.Flush()
				now = time.Now().UnixMilli()
				s := fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v", now, block.Header.Number, tx_total, tx_total-relayCount, tx_total-len(commit_ids), relayCount-tx_total+len(commit_ids), p.num_new_account, p.num_account_send, p.num_account_all, p.num_accounts_visit, p.num_active_accounts_visit)
				blocklog.Write(strings.Split(s, " "))
				blocklog.Flush()

				//主节点向客户端发送已确认上链的交易集
				c, err := json.Marshal(commit_ids)
				if err != nil {
					// fmt.Println("handleCommit出错")
					log.Panic(err)
				}
				m := jointMessage(cReply, c)
				utils.TcpDial(m, params.ClientAddr)

				queue_len := len(p.Node.CurChain.Tx_pool.Queue)
				err = queuelog.Write([]string{fmt.Sprintf("%v", now), fmt.Sprintf("%v", queue_len)})
				if err != nil {
					// fmt.Println("handleCommit出错")
					log.Panic(err)
				}
				queuelog.Flush()
				p.ReconfigStart()

			}
			p.isReply[c.Digest] = true

			p.sequenceID += 1
			if p.Node.nodeID == "N0" {
				// 发送relay tx
				p.TryRelay()
				p.sequenceLock.Unlock()
			}
		}
		fmt.Println("---------当前出块结束----------")
	}
	p.lock.Unlock()

	// }
}

func (p *Pbft) requestBlocks(startID, endID int) {
	r := RequestBlocks{
		StartID:  startID,
		EndID:    endID,
		ServerID: "N0",
		NodeID:   p.Node.nodeID,
	}
	bc, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("正在请求区块高度%d到%d的区块\n", startID, endID)
	message := jointMessage(cRequestBlock, bc)
	go utils.TcpDial(message, p.nodeTable[r.ServerID])
}

// 目前只有主节点会接收和处理这个请求，假设主节点拥有完整的全部区块
func (p *Pbft) handleRequestBlock(content []byte) {
	rb := new(RequestBlocks)
	err := json.Unmarshal(content, rb)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s节点发来的 requestBlock ... \n", rb.NodeID)
	blocks := make([]*core.Block, 0)
	for id := rb.StartID; id <= rb.EndID; id++ {
		if _, ok := p.height2Digest[id]; !ok {
			fmt.Printf("主节点没有找到高度%d对应的区块摘要！\n", id)
		}
		if r, ok := p.messagePool[p.height2Digest[id]]; !ok {
			fmt.Printf("主节点没有找到高度%d对应的区块！\n", id)
			log.Panic()
		} else {
			encoded_block := r.Message.Content
			block := core.DecodeBlock(encoded_block)
			blocks = append(blocks, block)
		}
	}
	p.SendBlocks(rb, blocks)

}

func (p *Pbft) SendBlocks(rb *RequestBlocks, blocks []*core.Block) {
	s := SendBlocks{
		StartID: rb.StartID,
		EndID:   rb.EndID,
		Blocks:  blocks,
		NodeID:  rb.ServerID,
	}
	bc, err := json.Marshal(s)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("正在向节点%s发送区块高度%d到%d的区块\n", rb.NodeID, s.StartID, s.EndID)
	message := jointMessage(cSendBlock, bc)
	go utils.TcpDial(message, p.nodeTable[rb.NodeID])
}

func (p *Pbft) handleSendBlock(content []byte) {
	sb := new(SendBlocks)
	err := json.Unmarshal(content, sb)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到%s发来的%d到%d的区块 \n", sb.NodeID, sb.StartID, sb.EndID)
	for id := sb.StartID; id <= sb.EndID; id++ {
		p.Node.CurChain.AddBlock(sb.Blocks[id-sb.StartID], p.consensusFlag, p.epochID)
		p.height2BlockHash[id] = sb.Blocks[id-sb.StartID].GetHash()
		fmt.Printf("编号为 %d 的区块已加入本地区块链！\n", id)
		curBlock := p.Node.CurChain.CurrentBlock
		fmt.Printf("curBlock: \n")
		curBlock.PrintBlock()
	}
	p.sequenceID = sb.EndID + 1
	p.requestLock.Unlock()
}

// 向除自己外的其他节点进行广播(本分片)
func (p *Pbft) broadcast(cmd command, content []byte) {
	message := jointMessage(cmd, content)
	for i := range p.nodeTable {
		if i == p.Node.nodeID {
			continue
		}
		go utils.TcpDial(message, p.nodeTable[i])
	}
}

// 向所有主节点进行广播
func (p *Pbft) broadcastToMain(cmd command, content []byte) {

	message := jointMessage(cmd, content)
	for i, node := range params.NodeTable {
		if i == params.Config.ShardID {
			continue
		}
		fmt.Printf("==========正在向节点%s发送消息======\n", node["N0"])
		go utils.TcpDial(message, node["N0"])
	}
}

// 向所有节点进行广播
func (p *Pbft) broadcastToAll(cmd command, content []byte) {
	message := jointMessage(cmd, content)
	for _, node := range params.NodeTable {
		// if shardID != "SC" {
		// }
		for nodeID, _ := range node {
			go utils.TcpDial(message, node[nodeID])
		}
	}
}

// 主节点向中心分片发送共识完成信息
func (p *Pbft) ReconfigStart() {

	fmt.Println("=============向中心分片发送共识完成信息==========")

	shardID := p.Node.CurChain.ChainConfig.ShardID
	sendip := p.Node.addr
	height := fmt.Sprint(p.sequenceID - 1)
	// time := time.Now()

	bc, err := json.Marshal(shardID + "_" + sendip + "_" + height)
	if err != nil {
		log.Panic(err)
	}
	message := jointMessage(cReconfigStart, bc)
	go utils.TcpDial(message, params.NodeTable["SC"]["N0"])

}

// 中心分片主节点收到消息，开始重组
func (p *Pbft) handleReconfigStart(content []byte) {
	var shardID string
	err := json.Unmarshal(content, &shardID)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到一次%s分片发来的ReconfigStart ... \n", string(shardID))
	p.setReconfigStartCountMap(string(shardID), true)
	count := 0
	for range p.reconfigStartCount {
		count++
	}
	fmt.Printf("----Count: %d------\n", count)

	num_Wshard := len(params.NodeTable) - 1 // 各分片从节点数量
	if count == num_Wshard*num_block_per_epoch {
		// if count == num_Wshard {
		fmt.Println("本节点已收到所有分片发来的ReconfigStart信息 ...")
		time.Sleep(time.Duration(200) * time.Millisecond)
		go p.CenterNode()
		p.reconfigStartCount = make(map[string]bool)
	}
}

// 中心分片主节点进行重组方案广播
func (p *Pbft) Reconfig(new_nodes map[string]map[string]string) {

	fmt.Printf("===当前周期为：%d ===\n", p.epochID)

	r := &Reconfig{}
	r.ID = getRandom()
	r.Content = new_nodes

	b, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("正在向其他节点广播重组方案 ...")
	p.broadcastToMain(cReconfig, b) // 发给其他分片主节点
	p.broadcast(cReconfig, b)       // 发给本分片其他节点
	fmt.Println("重组方案广播完成")

	p.epochID++

	// p.sequenceLock.Lock() //通过锁强制要求上一个区块commit完成新的区块才能被提出
	// <-p.reconfigDone //	用于中心分片计算重组时间
	// p.handleReconfig_New(b)
}

// content 为 新的重组方案
// 主节点开始重组条件
// - 结束handleCommit，区块上链
// - 结束handleRelay, relay交易加入交易池
func (p *Pbft) handleReconfig_Stale(content []byte) {
	if p.consensusFlag { // 若重组周期收到reconfig消息则忽略
		// 初始化channel
		// p.reconfigStart = make(chan int, 0)
		p.reconfigDone = make(chan int, 0)
		p.reconfigOver = make(chan int, 0)
		// p.reconfigFlag = make(chan int, 0)
		// p.consensusFlag = false

		fmt.Println("本节点已接收到主分片发来的Reconfig ...")
		p.sequenceLock.Lock() // 保证共识阶段结束再重组
		p.sequenceLock.Unlock()
		p.reconfigBlockLock.Lock()

		// time.Sleep(time.Duration(params.Config.Relay_time) * time.Millisecond)

		// 进入重组阶段，停止共识
		p.consensusFlag = false
		p.Node.time1 = time.Now() // 记录重组时间
		fmt.Println("----------time1----------")

		// 使用json解析出Reconfig结构体
		pp := new(Reconfig)
		err := json.Unmarshal(content, pp)
		if err != nil {
			log.Panic(err)
		}

		// p.reconfigFlag = make(chan int, 0)
		fmt.Println("============进入重组阶段============")

		// 仅工作分片进行重组
		if p.Node.CurChain.ChainConfig.ShardID != "SC" {

			p.Node.CurChain.Storage.PrintBucket()

			if p.Node.CurChain.ChainConfig.NodeID == "N0" {
				fmt.Println("等待共识完成……………………")
				<-p.relay
				p.reconfigLock.Lock()           //通过锁强制要求共识周期结束后才能重组
				p.broadcast(cReconfig, content) // 发给本分片其他节点
				fmt.Println("------- 重组消息发送完成，开始重组！-------")

				chain := p.Node.CurChain
				WriteReconfigBlockSizeLog(chain.ChainConfig.ShardID, chain.Storage.BlockBucketSize, chain.Storage.StateTreeBucketSize)

			} else {
				p.reconfigLock.Lock() //通过锁强制要求共识周期结束后才能重组
				fmt.Println("------ 收到主节点发来消息，开始重组！------")
			}

			// 1、交易区块发送和接收
			// var send_addr, receive_addr string
			// 找到自己被哪个节点替换
			var send_addr string
			for shardID, node := range pp.Content {
				for nodeID, addr := range node {
					// fmt.Printf("分片%v 节点编号%v 节点地址%v \n", shardID, nodeID, addr)
					if shardID == p.Node.CurChain.ChainConfig.ShardID && nodeID == p.Node.CurChain.ChainConfig.NodeID {
						send_addr = addr
					}
					// if addr ==  {
					// 	receive_addr = addr
					// }
				}
			}

			// 发送历史交易信息
			// // b、发送状态树
			// stateTree := &new_trie.N_Trie{}
			// new_trie.DeepCopy(stateTree, p.Node.CurChain.StatusTrie)
			stateTree := p.Node.CurChain.StatusTrie.ForReconfig(p.epochID)
			fmt.Printf("============重组消息（状态树）生成完毕============\n")

			p.ReconfigTries(send_addr, stateTree, p.Node.CurChain.Tx_pool)
			p.reconfigBlockLock.Unlock()

			// 1 数据处理与发送时间
			p.Node.time_send = time.Now()

			p.reconfigStart <- 1 // 先发送完消息(执行ReconfigBlock)再处理接收的消息(执行handleReconfigBlock)
			fmt.Printf("============等待更新============\n")

			// 2、接收信息并更新状态树
			// p.handleReconfigBlock()

			// 3、重组信息共识

			// 更新完本地区块信息和状态树，最后更新节点和分片ID相关信息
			<-p.reconfigOver

			stale_shard := p.Node.CurChain.ChainConfig.ShardID
			stale_node := p.Node.CurChain.ChainConfig.NodeID

			// 更改NodeTable
			params.NodeTable = pp.Content

			for shardID, node := range params.NodeTable {
				for nodeID, addr := range node {
					// fmt.Printf("分片%v 节点编号%v 节点地址%v \n", shardID, nodeID, addr)
					if addr == p.Node.addr {
						params.Config.ShardID = shardID
						params.Config.NodeID = nodeID
					}
				}
			}

			p.Node.nodeID = params.Config.NodeID
			p.nodeTable = params.NodeTable[params.Config.ShardID]
			fmt.Println(p.Node.addr)
			fmt.Println(params.Config.ShardID)
			fmt.Println(p.Node.nodeID)

			// 重组阶段结束
			p.reconfigLock.Unlock() //通过锁强制要求重组周期结束后才能共识
			p.consensusFlag = true

			if stale_node == "N0" {
				// 工作分片主节点记录各次重组耗时
				p.Node.time2 = time.Now()
				WriteReconfigTimeLog(stale_shard, p.Node.time2.Sub(p.Node.time1), p.Node.time_send.Sub(p.Node.time1), p.Node.time_trans.Sub(p.Node.time1), p.Node.time_receieve.Sub(p.Node.time1), p.Node.time2.Sub(p.Node.time1))
				fmt.Println("----------time2----------")

				// handleReconfigBlock(content []byte)
				// p.reconfigCount = make(map[string]bool)
			}

			// handleReconfigBlock(content []byte)
			// p.reconfigDone <- 1
			// p.reconfigCount = make(map[string]bool)

			if params.Config.NodeID == "N0" {
				p.NextEopch()
				p.relay = make(chan int, 0)
			}

		} else {
			// 中心分片仅更改NodeTable
			params.NodeTable = pp.Content
			p.epochID++
			// config := params.Config
			// time.Sleep(time.Duration(config.Reconfig_time) * time.Millisecond)
		}

		fmt.Println("============重组阶段结束============")

	}

}

// 主节点向其他节点广播
func (p *Pbft) NextEopch() {
	p.epochID++

	fmt.Println("正在向其他节点广播下一轮共识过程开始 ...")

	p.broadcast(cNextEpochStart, nil)

	fmt.Println("下一轮共识过程开始广播完成")
}

// 其他节点向主节点返回重组结束消息
func (p *Pbft) NextEpochReply() {
	p.epochID++
	fmt.Println("接收到主节点发来的消息……………………")
	p.reconfigLock.Lock()
	fmt.Println("====重组结束，进入next epoch共识===")
	nodeID := p.Node.CurChain.ChainConfig.NodeID
	bc, err := json.Marshal(nodeID)
	if err != nil {
		log.Panic(err)
	}
	message := jointMessage(cNextEpochStartReply, bc)
	go utils.TcpDial(message, p.nodeTable["N0"])

	p.reconfigLock.Unlock()

}

// 主节点收到消息，开始共识
func (p *Pbft) handleNextEpochReply(content []byte) {
	var nodeID string
	err := json.Unmarshal(content, &nodeID)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("本节点已接收到%s节点发来的NextEpochReply ... \n", string(nodeID))
	p.nextEpochCount++

	// fmt.Printf("---------Count: %d--------------\n", count)

	num_node := len(p.nodeTable) - 1 // 各分片从节点数量
	if p.nextEpochCount == num_node {
		fmt.Println("本节点已收到所有节点发来的NextEpochReply信息 ...")
		time.Sleep(time.Duration(200) * time.Millisecond)
		if p.Node.CurChain.ChainConfig.ShardID != "SC" {
			go p.Propose()
		}
		p.nextEpochCount = 0
	}
}

// // 发送重组信息，即全部历史交易信息
// func (p *Pbft) ReconfigBlocks(addr string, blocks []*core.Block, tx_pool *core.Tx_pool) {

// 	s := SendBlocks{
// 		StartID: 1,
// 		EndID:   p.sequenceID - 1,
// 		Blocks:  blocks,
// 		NodeID:  p.Node.CurChain.ChainConfig.NodeID,
// 	}
// 	r := ReconfigBlockMessage{
// 		SendBlocks: s,
// 		Tx_pool:    tx_pool,
// 	}

// 	bc, err := json.Marshal(r)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	fmt.Printf("正在向节点 %s 发送区块高度%d到%d的区块\n", addr, r.StartID, r.EndID)
// 	message := jointMessage(cReconfigBlock, bc)
// 	go utils.TcpDial(message, addr)
// 	fmt.Printf("已发送\n")

// }

// // 收到历史交易信息
// func (p *Pbft) handleReconfigBlock(content []byte) {

// 	rm := new(ReconfigBlockMessage)
// 	err := json.Unmarshal(content, rm)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	blocks := rm.SendBlocks
// 	tx_pool := rm.Tx_pool

// 	bl, err := json.Marshal(blocks)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	// 先发送完消息再处理接收的消息
// 	<-p.reconfigStart

// 	// fmt.Printf("======收到消息，等待lock=======\n")
// 	p.reconfigBlockLock.Lock()
// 	fmt.Printf("----收到消息，开始更新本地区块和状态信息----\n")

// 	// fmt.Printf("before before: len(pool.Queue): %d\n", len(p.Node.CurChain.Tx_pool.Queue))
// 	// 清空本地交易区块及状态树存储
// 	p.Node.CurChain.Storage.ClearStorage()

// 	fmt.Printf("before: len(pool.Queue): %d\n", len(p.Node.CurChain.Tx_pool.Queue))
// 	p.Node.CurChain, _ = chain.NewBlockChain(params.Config)
// 	// p.Node.CurChain.NewBlockChainAfterReconfig(params.Config)

// 	// <-p.reconfigOver
// 	p.sequenceID = p.Node.CurChain.CurrentBlock.Header.Number + 1

// 	fmt.Printf("------清除本地数据后，sequenceID: %d------\n", p.sequenceID)

// 	// 添加新交易区块
// 	p.AddBlock(bl)
// 	p.reconfigBlockLock.Unlock()
// 	fmt.Printf("------添加新区块后，sequenceID: %d------\n", p.sequenceID)

// 	// 更新交易池
// 	p.Node.CurChain.Tx_pool = tx_pool
// 	fmt.Printf("after: len(pool.Queue): %d\n", len(p.Node.CurChain.Tx_pool.Queue))

// 	// 主节点先阻塞，接收完其他节点重组完成的消息再返回消息给中心分片
// 	if p.Node.CurChain.ChainConfig.NodeID == "N0" {
// 		<-p.reconfigDone
// 		fmt.Println("本节点已收到所有分片发来的重组完成信息 ...")

// 		p.reconfigOver <- 1

// 	} else {
// 		p.ReconfigDone()
// 	}
// 	// 普通节点重组完成后，返回消息给主节点；
// 	// 主节点更新后，本分片重组完毕，返回消息给中心分片

// }

// 发送重组信息，即全部历史交易信息
func (p *Pbft) ReconfigTries(addr string, trie *new_trie.N_Trie, tx_pool *core.Tx_pool) {
	// time.Sleep()

	r := ReconfigTrieMessage{
		Trie:    trie,
		Tx_pool: tx_pool,
	}
	bc, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("正在向节点 %s 发送状态树\n", addr)
	message := jointMessage(cReconfigTries, bc)
	go utils.TcpDial(message, addr)
	fmt.Printf("已发送\n")
}

// 收到历史交易信息
func (p *Pbft) handleReconfigTries(content []byte) {

	timeSend := p.Node.CurChain.Storage.StateTreeBucketSize / (1024 * bandwidth)
	time.Sleep(time.Duration(timeSend) * time.Second)
	// 2 数据传输时间
	p.Node.time_trans = time.Now()

	rm := new(ReconfigTrieMessage)
	json.Unmarshal(content, rm)
	// err := json.Unmarshal(content, rm)
	// if err != nil {
	// 	log.Panic(err)
	// }
	trie := rm.Trie
	tx_pool := rm.Tx_pool

	// 先发送完消息再处理接收的消息
	<-p.reconfigStart

	p.reconfigBlockLock.Lock()
	fmt.Printf("----收到消息，开始更新本地区块和状态信息----\n")

	// 打印本地区块信息
	p.Node.CurChain.Storage.PrintBucket()
	// 清空本地交易区块及状态树存储
	p.Node.CurChain.Storage.ClearStorage()

	// chain_temp := p.Node.CurChain
	// WriteReconfigBlockSizeLog(chain_temp.ChainConfig.ShardID, chain_temp.Storage.BlockBucketSize, chain_temp.Storage.StateTreeBucketSize)

	fmt.Printf("before: len(pool.Queue): %d\n", len(p.Node.CurChain.Tx_pool.Queue))
	p.Node.CurChain, _ = chain.NewBlockChain(params.Config)
	p.sequenceID = p.Node.CurChain.CurrentBlock.Header.Number + 1
	fmt.Printf("------清除本地数据后，sequenceID: %d------\n", p.sequenceID)

	// 更改状态树
	p.Node.CurChain.StatusTrie = trie
	p.reconfigBlockLock.Unlock()
	fmt.Printf("------添加状态树后，sequenceID: %d------\n", p.sequenceID)

	// 更新交易池
	p.Node.CurChain.Tx_pool = tx_pool
	fmt.Printf("after: len(pool.Queue): %d\n", len(p.Node.CurChain.Tx_pool.Queue))

	// 3 处理新数据时间
	p.Node.time_receieve = time.Now()

	// 主节点先阻塞，接收完其他节点重组完成的消息再返回消息给中心分片
	if p.Node.CurChain.ChainConfig.NodeID == "N0" {

		<-p.reconfigDone
		fmt.Println("本节点已收到所有分片发来的重组完成信息 ...")

		p.reconfigOver <- 1

	} else {
		p.ReconfigDone()
	}
	// 普通节点重组完成后，返回消息给主节点；
	// 主节点更新后，本分片重组完毕，返回消息给中心分片
}

// 中心分片发送重组消息
func (p *Pbft) ReconfigCenter(send_addr string) {
	r := ReconfigCenterMessage{
		NodeTable: p.nodeTable,
	}
	bc, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("正在向节点 %s 发送NodeTable\n", send_addr)
	message := jointMessage(cReconfigCenter, bc)
	go utils.TcpDial(message, send_addr)
	fmt.Printf("已发送\n")

}

// 收到中心分片消息
func (p *Pbft) handleReconfigCenter(content []byte) {

	rm := new(ReconfigCenterMessage)
	err := json.Unmarshal(content, rm)
	if err != nil {
		log.Panic(err)
	}

	// 先发送完消息再处理接收的消息
	<-p.reconfigStart

	p.reconfigBlockLock.Lock()
	// 更改相关信息
	p.reconfigBlockLock.Unlock()

	// 主节点先阻塞，接收完其他节点重组完成的消息再返回消息给中心分片
	if p.Node.CurChain.ChainConfig.NodeID == "N0" {
		<-p.reconfigDone
		fmt.Println("本节点已收到所有分片发来的重组完成信息 ...")

		p.reconfigOver <- 1
		fmt.Println("=======test=====")

	} else {
		p.ReconfigDone()
	}

}

// 向上层节点报告重组完成信息
func (p *Pbft) ReconfigDone() {
	// 发送消息：ReconfigDone+我的ip

	// 普通节点返回消息给主节点；
	if p.Node.CurChain.ChainConfig.NodeID != "N0" {
		nodeID := p.Node.CurChain.ChainConfig.NodeID
		bc, err := json.Marshal(nodeID)
		if err != nil {
			log.Panic(err)
		}
		message := jointMessage(cReconfigDone, bc)
		go utils.TcpDial(message, p.nodeTable["N0"])
		fmt.Println("------ReconfigDone消息已发送------")

		p.reconfigOver <- 1

		// time1 := time.Now().Unix()
		// for !p.flag {
		// 	time2 := time.Now().Unix()
		// 	if time2-time1 >= p.time_to_reissue {
		// 		go utils.TcpDial(message, p.nodeTable["N0"])
		// 		time1 = time2
		// 	}
		// }
	}

}

// 各分片主节点接收各节点传来的消息
func (p *Pbft) handleReconfigDone(content []byte) {

	//使用json解析出nodeID
	var nodeID string
	err := json.Unmarshal(content, &nodeID)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("本节点已接收到%s节点发来的ReconfigDone ... \n", string(nodeID))
	p.setReconfigCountMap(string(nodeID), true)
	count := 0
	for range p.reconfigCount {
		count++
	}

	// num_Wshard := len(params.NodeTable) - 1     // 工作分片数量
	num_node := len(params.NodeTable["S0"]) - 1 // 各分片从节点数量

	if count == num_node {
		fmt.Println("本节点已收到所有节点发来的重组完成信息 ...")
		// handleReconfigBlock(content []byte)
		p.reconfigDone <- 1
		p.reconfigCount = make(map[string]bool)

	}

}

// // 主节点返回reply消息给发来reconfigdone的节点
// func (p *Pbft) ReconfigDoneReply(content []byte) {

// 	bc, err := json.Marshal(nil)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	message := jointMessage(cReconfigDoneReply, bc)
// 	go utils.TcpDial(message, string(content))

// }

// func (p *Pbft) handleReconfigDoneReply(content []byte) {

// 	p.flag = true
// }

// 为多重映射开辟赋值
func (p *Pbft) setPrePareConfirmMap(val, val2 string, b bool) {
	if _, ok := p.prePareConfirmCount[val]; !ok {
		p.prePareConfirmCount[val] = make(map[string]bool)
	}
	p.prePareConfirmCount[val][val2] = b
}

// 为多重映射开辟赋值
func (p *Pbft) setCommitConfirmMap(val, val2 string, b bool) {
	if _, ok := p.commitConfirmCount[val]; !ok {
		p.commitConfirmCount[val] = make(map[string]bool)
	}
	p.commitConfirmCount[val][val2] = b
}

// 为多重映射开辟赋值
func (p *Pbft) setReconfigCountMap(val string, b bool) {
	if _, ok := p.reconfigCount[val]; !ok {
		p.reconfigCount[val] = false
	}
	p.reconfigCount[val] = b
}

// // 为多重映射开辟赋值
// func (p *Pbft) setNextEpochCountMap(val string, b bool) {
// 	if _, ok := p.nextEpochCount[val]; !ok {
// 		p.nextEpochCount[val] = false
// 	}
// 	p.nextEpochCount[val] = b
// }

// 为多重映射开辟赋值
func (p *Pbft) setReconfigStartCountMap(val string, b bool) {
	if _, ok := p.reconfigStartCount[val]; !ok {
		p.reconfigStartCount[val] = false
	}
	p.reconfigStartCount[val] = b
}

// // 为多重映射开辟赋值
// func (p *Pbft) setRelayCountMap(val string, b bool) {
// 	if _, ok := p.relayCount[val]; !ok {
// 		p.relayCount[val] = false
// 	}
// 	p.relayCount[val] = b
// }

// 返回一个十位数的随机数，作为msgid
func getRandom() int {
	// x := big.NewInt(10000000000)
	// for {
	// 	result, err := rand.Int(rand.Reader, x)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	if result.Int64() > 1000000000 {
	// 		return int(result.Int64())
	// 	}
	// }

	string1 := fmt.Sprintf("%010v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000))
	string2 := fmt.Sprintf("%010v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000))
	result, _ := strconv.Atoi(string1 + string2)
	return result
}

// relay
func (p *Pbft) TryRelay() {
	fmt.Println("===========TryRelay=============")
	config := params.Config
	// for p.consensusFlag {
	// 仅在非重组周期进行重组

	// <-p.reconfigFlag
	// time.Sleep(time.Duration(config.Relay_interval) * time.Millisecond)
	for k, v := range params.NodeTable {
		if k == config.ShardID || k == "SC" { // 不转发给本分片和中心分片
			continue
		}
		if txs, isEnough := p.Node.CurChain.Tx_pool.FetchRelayTxs(k); isEnough {
			target_leader := v["N0"]
			r := Relay{
				Txs:     txs,
				ShardID: k,
			}
			bc, err := json.Marshal(r)
			if err != nil {
				log.Panic(err)
			}

			fmt.Printf("正在向分片%v的主节点发送relay交易\n", k)
			// fmt.Printf("发送交易%v\n", txs)
			message := jointMessage(cRelay, bc)
			go utils.TcpDial(message, target_leader)

		} else {
			target_leader := v["N0"]

			bc, err := json.Marshal(nil)
			if err != nil {
				log.Panic(err)
			}
			message := jointMessage(cRelay, bc)
			go utils.TcpDial(message, target_leader)

		}
	}

	// }
}

func (p *Pbft) handleRelay(content []byte) {
	// 重组阶段不处理relay消息
	p.reconfigLock.Lock()
	relay := new(Relay)
	err := json.Unmarshal(content, relay)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到分片%v发来的relay交易 \n", relay.ShardID)
	if relay != nil {
		p.Node.CurChain.Tx_pool.AddTxsToTop(relay.Txs)
	}

	p.relayCount++
	fmt.Println(p.relayCount)
	if p.relayCount == num_block_per_epoch*(params.Config.Shard_num-1) {
		p.relay <- 1
		p.relayCount = 0
	}

	p.reconfigLock.Unlock()
}

// func (p *Pbft) AddBlock(content []byte) {
// 	// handleSendBlock函数无锁版
// 	sb := new(SendBlocks)
// 	err := json.Unmarshal(content, sb)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	fmt.Printf("本节点已接收到%s发来的%d到%d的区块 \n", sb.NodeID, sb.StartID, sb.EndID)
// 	for id := sb.StartID; id <= sb.EndID; id++ {
// 		p.Node.CurChain.AddBlock(sb.Blocks[id-sb.StartID], p.consensusFlag, p.epochID)
// 		p.height2BlockHash[id] = sb.Blocks[id-sb.StartID].GetHash()
// 		fmt.Printf("编号为 %d 的区块已加入本地区块链！\n", id)
// 		// curBlock := p.Node.CurChain.CurrentBlock
// 		// fmt.Printf("curBlock: \n")
// 		// curBlock.PrintBlock()
// 	}
// 	p.sequenceID = sb.EndID + 1
// }

func WriteReconfigTimeLog(shardID string, time time.Duration, time1 time.Duration, time2 time.Duration, time3 time.Duration, time4 time.Duration) {
	csvFile, err := os.OpenFile("./log/"+"reconfig_time_per_epoch.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	timelog = csv.NewWriter(csvFile)
	time_temp := []string{shardID, fmt.Sprint(time.Seconds()), fmt.Sprint(time1.Seconds()), fmt.Sprint(time2.Seconds()), fmt.Sprint(time3.Seconds()), fmt.Sprint(time4.Seconds())}
	// time_temp := []string{fmt.Sprint(time)}
	// len := len(time_temp[0])
	// time_temp[0] = time_temp[0][:len-1]
	timelog.Write([]string(time_temp))
	timelog.Flush()
	// csvFile.Close()
}

func WriteReconfigBlockSizeLog(shardID string, blockSize int, stateSize int) {
	csvFile, err := os.OpenFile("./log/"+shardID+"_reconfig_data.csv", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	reconfigdatalog = csv.NewWriter(csvFile)
	reconfigDataSize_temp := []string{fmt.Sprint(blockSize), fmt.Sprint(stateSize)}
	reconfigdatalog.Write([]string(reconfigDataSize_temp))
	reconfigdatalog.Flush()
	// csvFile.Close()
}

func (p *Pbft) CenterNode() {

	new_nodes := StaleGenReconfigProposal()
	// new_nodes := NewGenReconfigProposal()

	for shardID, node := range new_nodes {
		for nodeID, addr := range node {
			fmt.Printf("分片%v 节点编号%v 节点地址%v \n", shardID, nodeID, addr)
		}
	}

	// time.Sleep(time.Duration(params.Config.Reconfig_interval) * time.Millisecond)

	// <-node.P.reconfigFlag
	// 2、广播重组方案
	p.Reconfig(new_nodes)

	// 3、修改本地NodeTable
	params.NodeTable = new_nodes
}

// 旧重组方案生成（含中心分片重组）
func StaleGenReconfigProposal() map[string]map[string]string {
	// 可改为可验证随机数
	timestamp := time.Now().Unix()

	var nodes []string
	for shardID, node := range params.NodeTable {
		if shardID != "SC" {
			for _, addr := range node {
				nodes = append(nodes, addr)
			}

		}
	}
	// fmt.Println(nodes)

	// 随机打乱数组
	rand.Seed(timestamp)
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	// fmt.Println("新工作分片节点")
	// fmt.Println(nodes)

	new_nodes := make(map[string]map[string]string)
	temp := 0
	for shardID, node := range params.NodeTable {
		new_nodes[shardID] = make(map[string]string)
		for nodeID, _ := range node {
			if shardID != "SC" {
				new_nodes[shardID][nodeID] = nodes[temp]
				temp++
			} else {
				new_nodes[shardID][nodeID] = node[nodeID]
			}
		}
	}

	return new_nodes

}

// 新重组方案生成（含中心分片重组）
func NewGenReconfigProposal() map[string]map[string]string {
	// time.Sleep(1000 * time.Millisecond)
	// 可改为可验证随机数
	timestamp := time.Now().Unix()
	// timestamp2 := getRandom()
	num_Wshard := len(params.NodeTable) - 1 // 工作分片数量
	num_node := len(params.NodeTable["SC"]) // 各分片节点数量

	n_choice := int(timestamp) % num_node // 从第几个开始选(0,1,2,3)
	n_choice2 := num_node / num_Wshard    // 每个分片选几个
	n_choice3 := (n_choice + n_choice2 - 1) % num_node
	// extra := num_node % num_Wshard // 可能多余的节点，放到中心分片
	// fmt.Println(n_choice)
	// fmt.Println(n_choice3)

	var nodes [][]string
	for shardID, node := range params.NodeTable {
		if shardID != "SC" {
			var row []string
			for _, addr := range node {
				row = append(row, addr)
			}
			nodes = append(nodes, row)
		}
	}
	var row []string
	for _, addr := range params.NodeTable["SC"] {
		row = append(row, addr)
	}
	nodes = append(nodes, row)
	// fmt.Println(nodes)

	var nodes_mix []string
	var centShard []string
	if n_choice < n_choice3 {
		for i := 0; i < num_Wshard; i++ {
			for j := 0; j < len(nodes[i]); j++ {
				if j < n_choice || j > n_choice3 {
					nodes_mix = append(nodes_mix, nodes[i][j])
				} else {
					centShard = append(centShard, nodes[i][j])
				}
			}
		}
	} else {
		for i := 0; i < num_Wshard; i++ {
			for j := 0; j < len(nodes[i]); j++ {
				if j >= n_choice || j <= n_choice3 {
					centShard = append(centShard, nodes[i][j])
				} else {
					nodes_mix = append(nodes_mix, nodes[i][j])
				}
			}
		}
	}
	for j := 0; j < len(nodes[num_Wshard]); j++ {
		nodes_mix = append(nodes_mix, nodes[num_Wshard][j])
	}
	// fmt.Println(nodes_mix)

	// 随机打乱数组
	rand.Seed(timestamp)
	rand.Shuffle(len(nodes_mix), func(i, j int) {
		nodes_mix[i], nodes_mix[j] = nodes_mix[j], nodes_mix[i]
	})
	// fmt.Println("新工作分片节点")
	// fmt.Println(nodes_mix)
	// fmt.Println("新中心分片节点")
	// fmt.Println(centShard)

	new_nodes := make(map[string]map[string]string)
	temp1, temp2 := 0, 0
	for shardID, node := range params.NodeTable {
		new_nodes[shardID] = make(map[string]string)
		for nodeID, _ := range node {
			if shardID != "SC" {
				new_nodes[shardID][nodeID] = nodes_mix[temp1]
				temp1++
			} else {
				new_nodes[shardID][nodeID] = centShard[temp2]
				temp2++
			}
		}
	}
	// for shardID, node := range params.NodeTable {
	// 	for nodeID, addr := range node {
	// 		fmt.Printf("分片%v 节点编号%v 节点地址%v \n", shardID, nodeID, addr)
	// 	}
	// }
	return new_nodes
	// return params.NodeTable

}

func (p *Pbft) SendAccounts(accounts []string, nodeID string) {
	s := SendAccounts{
		Accounts: accounts,
		NodeID:   nodeID,
	}
	bc, err := json.Marshal(s)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("正在向主分片请求账户信息\n")
	message := jointMessage(cSendAccount, bc)
	go utils.TcpDial(message, params.NodeTable["SC"]["N0"])
}

func (p *Pbft) handleSendAccount(content []byte) {
	sa := new(SendAccounts)
	err := json.Unmarshal(content, sa)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(sa.NodeID)
	accounts := sa.Accounts

	stateTree := p.Node.CurChain.StatusTrie
	var reply []string
	for _, account := range accounts {
		if _, ok := stateTree.Get([]byte(account)); !ok { // 若原状态树中不存在发送账户
			reply = append(reply, "1")
		} else {
			reply = append(reply, "0")
		}
	}

	bc, err := json.Marshal(reply)
	if err != nil {
		log.Panic(err)
	}
	message := jointMessage(cHandleSendAccount, bc)
	go utils.TcpDial(message, sa.NodeID)
}

func (p *Pbft) handlehandleSendAccount(content []byte) {
	var reply []string
	err := json.Unmarshal(content, &reply)
	if err != nil {
		log.Panic(err)
	}

	p.proposeBlock <- 1
}
