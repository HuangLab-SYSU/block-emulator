package pbft

import (
	"blockEmulator/account"
	"blockEmulator/algorithm"
	"blockEmulator/chain"
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/utils"
	"bytes"
	"crypto/rand"
	"encoding/csv"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	// blocklog, txlog, ctxlog *csv.Writer
	blocklog, blockexetimelog, txlog, epochChangelog, txmig2log, txmig1log *csv.Writer
	pbftbefore, InitTime, lastpropose                                      int64
)

// //本地消息池（模拟持久化层），只有确认提交成功后才会存入此池
// var localMessagePool = []Message{}

type node struct {
	//节点ID
	nodeID string
	//节点监听地址
	addr     string
	CurChain *chain.BlockChain
}

type Pbft struct {
	//节点信息
	Node node
	//每笔请求自增序号
	sequenceID int
	//锁
	lock sync.Mutex
	//确保系统停止出块的锁
	epochLock sync.Mutex
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
	//区块高度到区块摘要的映射,目前只给主节点使用
	height2Digest map[int]string
	nodeTable     map[string]string
	Stop          chan int

	//每笔请求自增序号
	msequenceID int
	//临时消息池，消息摘要对应消息本体
	mmessagePool map[string]*Request

	//每笔请求自增序号
	ssequenceID int
	//临时消息池，消息摘要对应消息本体
	smessagePool map[string]*Request

	// 轮次
	epoch int
	// epoch锁
	pbftlock sync.Mutex

	// 延迟发送tryout，因此要有个池子存着
	TryOutPool [][]*core.TXmig2
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
	p.epoch = 1

	if config.Cross_Chain && !config.Fail {
		p.TryOutPool = make([][]*core.TXmig2, config.Shard_num)
	}

	if config.Stop_When_Migrating {
		p.msequenceID = 0
		p.mmessagePool = make(map[string]*Request)

		p.ssequenceID = 0
		p.smessagePool = make(map[string]*Request)

		inacc_pool = make([]*address_and_balance, 0)
		intx_pool = make([]*core.Transaction, 0)

		sendoutfinish = make(chan int)
	}

	return p
}

func NewLog(shardID string) {
	csvFile, err := os.Create("./log/" + shardID + "_block.csv")
	if err != nil {
		log.Panic(err)
	}
	// defer csvFile.Close()
	blocklog = csv.NewWriter(csvFile)
	blocklog.Write([]string{"timestamp", "blockHeight", "tx_total", "tx_normal", "tx_relay_sender", "tx_relay_receiver", "tx_committed", "pbfttime", "mig1", "mig2", "ann", "ns", "locked_txs"})
	blocklog.Flush()

	csvFile, err = os.Create("./log/" + shardID + "_blockexetime.csv")
	if err != nil {
		log.Panic(err)
	}
	blockexetimelog = csv.NewWriter(csvFile)
	blockexetimelog.Write([]string{"blockHeight", "tx_total", "mig_total", "mig_time", "ann_total", "ann_pool_time", "ann_time"})
	blockexetimelog.Flush()

	if params.Config.Stop_When_Migrating {
		epochChangeFile, err := os.Create("./log/" + shardID + "_epochChange.csv")
		if err != nil {
			log.Panic(err)
		}
		// defer csvFile.Close()
		epochChangelog = csv.NewWriter(epochChangeFile)
		epochChangelog.Write([]string{"epoch", "start", "end", "time"})
		epochChangelog.Flush()
	}

	csvFile, err = os.Create("./log/" + shardID + "_transaction.csv")
	if err != nil {
		log.Panic(err)
	}
	txlog = csv.NewWriter(csvFile)
	txlog.Write([]string{"sender", "recipient", "txid", "blockHeight", "request_time", "2nd_request", "commit-request", "lock", "unlock", "lock2", "unlock2", "isSuccess", "SenLock", "RecLock", "HalfLock", "SenSupposeHeight", "RecSupposeHeight", "RelayLock"})
	txlog.Flush()

	csvFile, err = os.Create("./log/" + shardID + "_mig1.csv")
	if err != nil {
		log.Panic(err)
	}
	txmig1log = csv.NewWriter(csvFile)
	txmig1log.Write([]string{"acc", "source", "target", "blockHeight", "CommitTime"})
	txmig1log.Flush()

	csvFile, err = os.Create("./log/" + shardID + "_mig2.csv")
	if err != nil {
		log.Panic(err)
	}
	txmig2log = csv.NewWriter(csvFile)
	txmig2log.Write([]string{"account", "time"})
	txmig2log.Flush()
}

func (p *Pbft) handleRequest(data []byte) {
	//切割消息，根据消息命令调用不同的功能
	cmd, content := splitMessage(data)
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
	case cClient:
		go p.handleTxFromClient(content)
	case cRelay:
		go p.handleRelay(content)
	case cNewMap:
		go p.handleNewMap(content)
	case cTXmig1:
		go p.handleMig2(content)
	case cAnnounce:
		go p.handleAnnounce(content)
	case cMigrateWanted:
		go p.handleMigrateWanted()
	case cCaP:
		p.handleChangesAndPendings(content)
	case cEpochCh:
		go p.handleEpochChange()
	case cBalanceAndPending:
		p.handleBalancesAndPendings(content)
	case cSure:
		p.handleSure(content)
	case cStop:
		p.Stop <- 1
	case cLLT:
		p.handleLLT(content)
	}
}

// 只有主节点可以调用
// 生成一个区块并发起共识
func (p *Pbft) Propose() {
	config := params.Config
	for {
		if config.Stop_When_Migrating {
			p.epochLock.Lock() //确保系统处于非停止状态
		}
		// time.Sleep(time.Duration(config.Block_interval) * time.Millisecond)
		for {
			p.pbftlock.Lock()
			if int(time.Now().Unix()-InitTime) >= (p.epoch-1)*config.Block_interval && int(time.Now().Unix()-lastpropose) >= config.Block_interval {
				p.epoch++
				p.pbftlock.Unlock()
				break
			}
			p.pbftlock.Unlock()
		}
		p.sequenceLock.Lock() //通过锁强制要求上一个区块commit完成新的区块才能被提出
		// if p.epoch != 2 {
		// 	time.Sleep(time.Duration(utils.RandInt0To3(InitTime+int64(p.epoch))) * time.Second)  //随机等待0~3秒
		// }

		lastpropose = time.Now().Unix()
		// fmt.Printf("要提出第%v个块，时间为%v\n", p.epoch, lastpropose)
		// p.propose()
		p.propose1()
	}
}

func (p *Pbft) propose() {
	r := &Request{}
	r.Timestamp = time.Now().Unix()
	r.Message.ID = getRandom()

	block := p.Node.CurChain.GenerateBlock(p.sequenceID)
	encoded_block := block.Encode()

	pbftbefore = time.Now().Unix()

	r.Message.Content = encoded_block

	// //添加信息序号
	// p.sequenceIDAdd()
	//获取消息摘要
	digest := getDigest(r)
	fmt.Println("已将request存入临时消息池")
	//存入临时消息池
	p.messagePool[digest] = r

	//拼接成PrePrepare，准备发往follower节点
	pp := PrePrepare{r, digest, p.sequenceID, "Block"}
	p.height2Digest[p.sequenceID] = digest
	b, err := json.Marshal(pp)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Println("正在向其他节点进行进行PrePrepare广播 ...")
	//进行PrePrepare广播
	p.broadcast(cPrePrepare, b)
	// fmt.Println("PrePrepare广播完成")

}

// 处理预准备消息
func (p *Pbft) handlePrePrepare(content []byte) {
	// fmt.Println("本节点已接收到主节点发来的PrePrepare ...")
	//	//使用json解析出PrePrepare结构体
	pp := new(PrePrepare)
	err := json.Unmarshal(content, pp)
	if err != nil {
		log.Panic(err)
	}

	pbftType := pp.Type

	if digest := getDigest(pp.RequestMessage); digest != pp.Digest {
		fmt.Println("信息摘要对不上，拒绝进行prepare广播")
	} else if (pbftType == "Block" && p.sequenceID != pp.SequenceID) || (params.Config.Stop_When_Migrating && pbftType == "EpochChange" && p.msequenceID != pp.SequenceID) || (params.Config.Stop_When_Migrating && pbftType == "AccState" && p.ssequenceID != pp.SequenceID) {
		fmt.Println("消息序号对不上，拒绝进行prepare广播")
	} else if pbftType == "Block" && !p.Node.CurChain.IsBlockValid(core.DecodeBlock(pp.RequestMessage.Message.Content)) {
		// todo
		fmt.Println("区块不合法，拒绝进行prepare广播")
	} else {
		// //序号赋值
		// p.sequenceID = pp.SequenceID
		//将信息存入临时消息池
		// fmt.Println("已将消息存入临时节点池")
		if pbftType == "Block" {
			p.messagePool[pp.Digest] = pp.RequestMessage
		} else if pbftType == "EpochChange" {
			p.mmessagePool[pp.Digest] = pp.RequestMessage
		} else {
			p.smessagePool[pp.Digest] = pp.RequestMessage
		}

		//拼接成Prepare
		pre := Prepare{pp.Digest, pp.SequenceID, p.Node.nodeID, pbftType}
		bPre, err := json.Marshal(pre)
		if err != nil {
			log.Panic(err)
		}
		//进行准备阶段的广播
		// fmt.Println("正在进行Prepare广播 ...")
		p.broadcast(cPrepare, bPre)
		fmt.Println("Prepare广播完成")
	}
}

// 处理准备消息
func (p *Pbft) handlePrepare(content []byte) {
	//使用json解析出Prepare结构体
	pre := new(Prepare)
	err := json.Unmarshal(content, pre)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Printf("本节点已接收到%s节点发来的Prepare ... \n", pre.NodeID)
	pbftType := pre.Type
	ok := false
	if pbftType == "Block" {
		_, ok = p.messagePool[pre.Digest]
	} else if pbftType == "EpochChange" {
		_, ok = p.mmessagePool[pre.Digest]
	} else {
		_, ok = p.smessagePool[pre.Digest]
	}

	if !ok {
		fmt.Println("当前临时消息池无此摘要，拒绝执行commit广播")
	} else if (pbftType == "Block" && p.sequenceID != pre.SequenceID) || (params.Config.Stop_When_Migrating && pbftType == "EpochChange" && p.msequenceID != pre.SequenceID) || (params.Config.Stop_When_Migrating && pbftType == "AccState" && p.ssequenceID != pre.SequenceID) {
		fmt.Printf("%v消息序号对不上，拒绝执行commit广播\n", pbftType)
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
			specifiedCount = 2 * malicious_num
		}
		//如果节点至少收到了2f个prepare的消息（包括自己）,并且没有进行过commit广播，则进行commit广播
		p.lock.Lock()
		//获取消息源节点的公钥，用于数字签名验证
		if count >= specifiedCount && !p.isCommitBordcast[pre.Digest] {
			// fmt.Println("本节点已收到至少2f个节点(包括本地节点)发来的Prepare信息 ...")

			c := Commit{pre.Digest, pre.SequenceID, p.Node.nodeID, pbftType}
			bc, err := json.Marshal(c)
			if err != nil {
				log.Panic(err)
			}
			//进行提交信息的广播
			// fmt.Println("正在进行commit广播")
			p.broadcast(cCommit, bc)
			p.isCommitBordcast[pre.Digest] = true
			fmt.Println("commit广播完成")
		}
		p.lock.Unlock()
	}
}

// 处理提交确认消息
func (p *Pbft) handleCommit(content []byte) {
	//使用json解析出Commit结构体
	c := new(Commit)
	err := json.Unmarshal(content, c)
	if err != nil {
		log.Panic(err)
	}
	pbftType := c.Type
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

	if count > malicious_num*2 && !p.isReply[c.Digest] {
		// fmt.Println("本节点已收到至少2f + 1 个节点(包括本地节点)发来的Commit信息 ...")
		//将消息信息，提交到本地消息池中！
		ok := false
		if pbftType == "Block" {
			_, ok = p.messagePool[c.Digest]
		} else if pbftType == "EpochChange" {
			_, ok = p.mmessagePool[c.Digest]
		} else {
			_, ok = p.smessagePool[c.Digest]
		}

		if pbftType == "Block" && !ok {
			// 1. 如果本地消息池里没有这个消息，说明节点落后于其他节点，向主节点请求缺失的区块
			p.isReply[c.Digest] = true
			p.requestLock.Lock()
			p.requestBlocks(p.sequenceID, c.SequenceID)
		} else if !ok {
			// todo
		} else {
			// 2.
			if pbftType == "Block" {
				r := p.messagePool[c.Digest]
				encoded_block := r.Message.Content
				block := core.DecodeBlock(encoded_block)
				pbftend := time.Now().UnixMilli()
				// 本地化存储。修改内存、存储至硬盘
				outbalance := p.Node.CurChain.AddBlock(block)
				fmt.Printf("编号为 %d 的区块已加入本地区块链！", p.sequenceID)
				curBlock := p.Node.CurChain.CurrentBlock
				fmt.Printf("curBlock: \n")
				curBlock.PrintBlock()

				if p.Node.nodeID == "N0" {
					if !params.Config.Stop_When_Migrating && !params.Config.Lock_Acc_When_Migrating {
						account.BalanceBeforeOutLock.Lock()
						for k, v := range outbalance {
							account.BalanceBeforeOut[k] = new(big.Int)
							account.BalanceBeforeOut[k].Set(v)
						}
						account.BalanceBeforeOutLock.Unlock()
					}

					tx_total := len(block.Transactions)
					relayCount := 0
					//已上链交易集
					commit_txs := []*core.Transaction{}
					//要发送relay交易集合
					relaytxs := make([]*core.Transaction, 0)
					for _, v := range block.Transactions {
						_, toid := account.Addr2Shard(hex.EncodeToString(v.Sender)), account.Addr2Shard(hex.EncodeToString(v.Recipient))
						//若交易接收者属于本分片才加入已上链交易集
						if params.Config.Lock_Acc_When_Migrating {
							account.Lock_Acc_Lock.Lock()
							if params.ShardTable[params.Config.ShardID] == toid {
								if account.Lock_Acc[hex.EncodeToString(v.Recipient)] {
									v.IsRelay = true
									// if v.LockTime > 0 {
									// 	v.LockTime2 = pbftend
									// } else {
									// 	v.LockTime = pbftend
									// }
									// v.RecLock = true
									// if params.Config.RelayLock {
									// 	v.Relay_Lock = true
									// }
									// p.Node.CurChain.Tx_pool.Locking_TX_Pools[hex.EncodeToString(v.Recipient)] = append(p.Node.CurChain.Tx_pool.Locking_TX_Pools[hex.EncodeToString(v.Recipient)], v)
								} else {
									commit_txs = append(commit_txs, v)
									v.CommitTime = pbftend
									var s string
									s = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v", hex.EncodeToString(v.Sender), hex.EncodeToString(v.Recipient), v.Id, block.Header.Number, v.RequestTime-InitTime*1000, v.Second_RequestTime-InitTime*1000, v.CommitTime-v.RequestTime, v.LockTime-InitTime*1000, v.UnlockTime-InitTime*1000, v.LockTime2-InitTime*1000, v.UnlockTime2-InitTime*1000, v.Success, v.SenLock, v.RecLock, v.HalfLock, v.Sen_Suppose_on_chain, v.Rec_Suppose_on_chain, v.Relay_Lock)
									txlog.Write(strings.Split(s, " "))
									// latency := pbftend - v.RequestTime
									// if fromid != toid {
									// 	ctxs := fmt.Sprintf("%v %v", v.Id, latency)
									// 	ctxlog.Write(strings.Split(ctxs, " "))
									// }
								}
							}
							account.Lock_Acc_Lock.Unlock()
						} else {
							if params.ShardTable[params.Config.ShardID] == toid {
								commit_txs = append(commit_txs, v)
								v.CommitTime = pbftend
								var s string
								s = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v %v", hex.EncodeToString(v.Sender), hex.EncodeToString(v.Recipient), v.Id, block.Header.Number, v.RequestTime-InitTime*1000, v.Second_RequestTime-InitTime*1000, v.CommitTime-v.RequestTime, v.LockTime-InitTime*1000, v.UnlockTime-InitTime*1000, v.LockTime2-InitTime*1000, v.UnlockTime2-InitTime*1000, v.Success, v.SenLock, v.RecLock, v.HalfLock, v.Sen_Suppose_on_chain, v.Rec_Suppose_on_chain, v.Relay_Lock)
								txlog.Write(strings.Split(s, " "))
								// latency := pbftend - v.RequestTime
								// if fromid != toid {
								// 	ctxs := fmt.Sprintf("%v %v", v.Id, latency)
								// 	ctxlog.Write(strings.Split(ctxs, " "))
								// }
							}
						}

						// var s string
						// s = fmt.Sprintf("%v %v %v %v %v %v", hex.EncodeToString(v.Sender), hex.EncodeToString(v.Recipient), v.Id, block.Header.Number, v.RequestTime-InitTime, pbftend-v.RequestTime)
						// txlog.Write(strings.Split(s, " "))
						if v.IsRelay {
							relayCount++
						} else if toid != params.ShardTable[params.Config.ShardID] {
							relayCount++
							v.IsRelay = true
							relaytxs = append(relaytxs, v)
						}
					}

					//记录接收账户的时间
					for _, mig := range block.TXmig2s {
						acc := mig.Address
						s := fmt.Sprintf("%v %v", acc, pbftend-InitTime*1000)
						txmig2log.Write(strings.Split(s, " "))
					}

					//记录TXmig1的时间
					if !params.Config.Stop_When_Migrating {
						for _, v := range block.TXmig1s {
							// ([]string{"acc", "source", "target", "blockHeight", "CommitTime"})
							s := fmt.Sprintf("%v %v %v %v %v", v.Address, params.Config.ShardID, v.ToshardID, block.Header.Number, pbftend-InitTime*1000)
							txmig1log.Write(strings.Split(s, " "))
						}
					}

					mig_begin := time.Now().UnixMicro()
					// 若要锁账户，就把账户锁住
					if params.Config.Lock_Acc_When_Migrating {
						account.Lock_Acc_Lock.Lock()
						for _, v := range block.TXmig1s {
							account.Lock_Acc[v.Address] = true
							p.Node.CurChain.Tx_pool.Locking_TX_Pools[v.Address] = make([]*core.Transaction, 0)
						}
						account.Lock_Acc_Lock.Unlock()
					}

					// 不停不锁，要把账户半锁
					if !params.Config.Stop_When_Migrating && !params.Config.Lock_Acc_When_Migrating {
						account.Outing_Acc_Before_Announce_Lock.Lock()
						for _, v := range block.TXmig1s {
							account.Outing_Acc_Before_Announce[v.Address] = true
							p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools[v.Address] = make([]*core.Transaction, 0)
						}
						account.Outing_Acc_Before_Announce_Lock.Unlock()
					}

					// // 超时
					// if params.Config.Fail && params.Config.Fail_Time == p.sequenceID && params.Config.ShardID == "S0" {
					// 	if params.Config.Lock_Acc_When_Migrating { //锁
					// 		account.Lock_Acc_Lock.Lock()
					// 		account.Lock_Acc["489338d5e8d42e8c923d1f47361d979503d4ad68"] = false
					// 		for _, v := range p.Node.CurChain.Tx_pool.Locking_TX_Pools["489338d5e8d42e8c923d1f47361d979503d4ad68"] {
					// 			v.UnlockTime = pbftend
					// 			v.Success = false
					// 		}
					// 		p.Node.CurChain.Tx_pool.AddTxs(p.Node.CurChain.Tx_pool.Locking_TX_Pools["489338d5e8d42e8c923d1f47361d979503d4ad68"])
					// 		delete(p.Node.CurChain.Tx_pool.Locking_TX_Pools, "489338d5e8d42e8c923d1f47361d979503d4ad68")
					// 		account.Lock_Acc_Lock.Unlock()
					// 	} else { // 不停不锁
					// 		account.Outing_Acc_Before_Announce_Lock.Lock()
					// 		account.Outing_Acc_Before_Announce["489338d5e8d42e8c923d1f47361d979503d4ad68"] = false
					// 		for _, v := range p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools["489338d5e8d42e8c923d1f47361d979503d4ad68"] {
					// 			v.UnlockTime = pbftend
					// 			v.Success = false
					// 		}
					// 		p.Node.CurChain.Tx_pool.AddTxs(p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools["489338d5e8d42e8c923d1f47361d979503d4ad68"])
					// 		delete(p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools, "489338d5e8d42e8c923d1f47361d979503d4ad68")
					// 		account.Outing_Acc_Before_Announce_Lock.Unlock()
					// 	}
					// }

					// build trie from the triedb (in disk)
					st1, err := trie.New(trie.TrieID(common.BytesToHash(block.Header.StateRoot)), p.Node.CurChain.Triedb)
					if err != nil {
						log.Panic(err)
					}
					st2, err := trie.New(trie.TrieID(common.BytesToHash(block.Header.StateRoot)), p.Node.CurChain.Triedb)
					if err != nil {
						log.Panic(err)
					}
					st3, err := trie.New(trie.TrieID(common.BytesToHash(block.Header.StateRoot)), p.Node.CurChain.Triedb)
					if err != nil {
						log.Panic(err)
					}
					st4, err := trie.New(trie.TrieID(common.BytesToHash(block.Header.StateRoot)), p.Node.CurChain.Triedb)
					if err != nil {
						log.Panic(err)
					}

					// use a memory trie database to do this, instead of disk database
					triedb := trie.NewDatabase(rawdb.NewMemoryDatabase())
					migTree := trie.NewEmpty(triedb)
					for _, tx := range block.TXmig1s {
						migTree.Update(tx.Hash(), tx.Encode())
					}
					for _, tx := range block.TXmig2s {
						migTree.Update(tx.Hash(), tx.Encode())
					}
					for _, tx := range block.Anns {
						migTree.Update(tx.Hash(), tx.Encode())
					}
					for _, tx := range block.NSs {
						migTree.Update(tx.Hash(), tx.Encode())
					}

					//锁交易(马上锁，交易池中都要锁)
					if !params.Config.Not_Lock_immediately {
						p.Node.CurChain.Tx_pool.LockTX()
					}

					//处理要被彻底迁移出去的账户
					jiaoyichitime, totaltime := p.handleAnns(block.Anns, st1, migTree)

					// 发送relay交易，不再等待
					if len(relaytxs) != 0 {
						go p.TryRelay(relaytxs, block.Transactions, st2, p.sequenceID)
					}

					if !params.Config.Stop_When_Migrating && !params.Config.Fail {
						// 发送迁出账户给对应分片
						if params.Config.Cross_Chain {
							p.TryTXmig1(block.TXmig1s, outbalance, st3, migTree)
						} else if len(block.TXmig1s) != 0 {
							go p.TryTXmig1(block.TXmig1s, outbalance, st3, migTree)
						}
						// 通知各分片，账户已在本分片
						if len(block.TXmig2s) != 0 {
							go p.TryAnnounce(block.TXmig2s, st4, migTree)
						}
					}

					migend := time.Now().UnixMicro()
					s := fmt.Sprintf("%v %v %v %v %v %v %v", block.Header.Number, tx_total, len(block.TXmig1s)+len(block.TXmig2s)+len(block.Anns)+len(block.NSs), migend-mig_begin, len(block.Anns), jiaoyichitime, totaltime)
					blockexetimelog.Write(strings.Split(s, " "))
					blockexetimelog.Flush()

					locked_cnt := 0
					if params.Config.Lock_Acc_When_Migrating {
						account.Lock_Acc_Lock.Lock()
						for _, locked := range p.Node.CurChain.Tx_pool.Locking_TX_Pools {
							locked_cnt += len(locked)
						}
						account.Lock_Acc_Lock.Unlock()
					} else if !params.Config.Stop_When_Migrating {
						account.Outing_Acc_Before_Announce_Lock.Lock()
						for _, locked := range p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools {
							locked_cnt += len(locked)
						}
						account.Outing_Acc_Before_Announce_Lock.Unlock()
					}

					txlog.Flush()
					txmig1log.Flush()
					txmig2log.Flush()
					// ctxlog.Flush()
					s = fmt.Sprintf("%v %v %v %v %v %v %v %v %v %v %v %v %v", pbftend-InitTime*1000, block.Header.Number, tx_total, tx_total-relayCount, tx_total-len(commit_txs), relayCount-tx_total+len(commit_txs), len(commit_txs), pbftend-pbftbefore, len(block.TXmig1s), len(block.TXmig2s), len(block.Anns), len(block.NSs), locked_cnt)
					blocklog.Write(strings.Split(s, " "))
					blocklog.Flush()

					commitTX_numofNS := Txs_and_Num_of_New_State{
						Txs:              commit_txs,
						BlockSize:        len(block.TXmig1s) + len(block.TXmig2s) + len(block.Anns) + len(block.NSs) + len(block.Transactions),
						Num_of_New_State: len(block.NSs),
					}
					//主节点向客户端发送已确认上链的交易集
					c, err := json.Marshal(commitTX_numofNS)
					if err != nil {
						log.Panic(err)
					}
					m := jointMessage(cReply, c)
					utils.TcpDial(m, params.ClientAddr)

					if p.sequenceID == 1 && (params.Config.Bu_Tong_Bi_Li || params.Config.Bu_Tong_Shi_Jian || (params.Config.Fail && !params.Config.Bu_Tong_Bi_Li_2) || params.Config.Cross_Chain) {
						t := time.Now().UnixMilli()
						for _, tx := range p.Node.CurChain.Tx_pool.Queue {
							tx.RequestTime = t
						}
					}
				}
				p.isReply[c.Digest] = true
				p.sequenceID += 1

				if p.Node.nodeID == "N0" {
					p.sequenceLock.Unlock()
					if params.Config.Stop_When_Migrating { //若迁移时要暂停
						p.epochLock.Unlock()
					}
				}

			} else if pbftType == "EpochChange" {
				r := p.mmessagePool[c.Digest]
				encoded_new := r.Message.Content
				var new map[string]int
				decoder := gob.NewDecoder(bytes.NewReader(encoded_new))
				err := decoder.Decode(&new)
				if err != nil {
					log.Panic(err)
				}

				// 将account2shard换成新的并得到out账户及映射
				out := make(map[string]int)
				account.Account2ShardLock.Lock()
				for addr, shard := range new {
					account.Account2Shard[addr] = shard
					if account.AccountInOwnShard[addr] && shard != params.ShardTable[params.Config.ShardID] {
						out[addr] = shard
						delete(account.AccountInOwnShard, addr)
					}
				}
				// account.Account2Shard = new
				account.Account2ShardLock.Unlock()

				// // 得到out账户及映射
				// out := make(map[string]int)
				// account.Account2ShardLock.Lock()
				// for k := range account.AccountInOwnShard {
				// 	if new[k] != params.ShardTable[params.Config.ShardID] {
				// 		out[k] = new[k]
				// 		delete(account.AccountInOwnShard, k)
				// 	}
				// }
				// account.Account2ShardLock.Unlock()

				if p.Node.nodeID == "N0" {
					//整理要发出去的账户的状态和交易，并发到对应分片主节点
					go p.SendOut(out)
				}

				// //将出去的账户从状态树中删除, 并本地化存储
				// for k, _ := range out {
				// 	fmt.Println(k)
				// 	hex_address, _ := hex.DecodeString(k)
				// 	p.Node.CurChain.StatusTrie.Delete(hex_address)
				// 	fmt.Printf("\n删除账户 %v 成功\n", k)
				// 	p.Node.CurChain.StatusTrie.PrintState()

				// }
				// p.Node.CurChain.Storage.UpdateStateTree(p.Node.CurChain.StatusTrie)
				p.msequenceID += 1

				p.isReply[c.Digest] = true

				if p.Node.nodeID == "N0" {
					flag := false
					//计算进来的个数
					inacc_count_lock.Lock()
					for k, v := range new {
						account.Account2ShardLock.Lock()
						if v == params.ShardTable[params.Config.ShardID] && !account.AccountInOwnShard[k] {
							inacc_count++
							flag = true
						}
						account.Account2ShardLock.Unlock()
					}

					if flag {
						fmt.Printf("\n有要进来的，此时inacc_count数量为: %v\n\n", inacc_count)
					}

					p.sequenceLock.Unlock()

					if flag && inacc_count == 0 { //有要进来的，且进完了
						inacc_count_lock.Unlock()
						p.spropose(inacc_pool)
					} else {
						inacc_count_lock.Unlock()
						if !flag { //没有要进来的
							go p.SendSure()
						}

					}

				}

			} else {
				r := p.smessagePool[c.Digest]
				encoded_stateNtx := r.Message.Content
				//解码成用户地址和状态
				in_StateAndTx := []*address_and_balance{}
				decoder := gob.NewDecoder(bytes.NewReader(encoded_stateNtx))
				err := decoder.Decode(&in_StateAndTx)
				if err != nil {
					log.Panic(err)
				}
				// build trie from the triedb (in disk)
				st, err := trie.New(trie.TrieID(common.BytesToHash(p.Node.CurChain.CurrentBlock.Header.StateRoot)), p.Node.CurChain.Triedb)
				if err != nil {
					log.Panic(err)
				}
				for _, v := range in_StateAndTx {
					//修改状态树
					hex_address, _ := hex.DecodeString(v.Address)

					account_state := &account.AccountState{
						Balance: new(big.Int).Set(v.Balance),
						Migrate: -1,
					}
					st.Update(hex_address, account_state.Encode())

					//将新账户放入account.AccountInOwnShard中，并设为true
					account.Account2ShardLock.Lock()
					account.AccountInOwnShard[v.Address] = true
					account.Account2ShardLock.Unlock()
				}
				// commit the memory trie to the database in the disk
				rt, ns := st.Commit(false)
				err = p.Node.CurChain.Triedb.Update(trie.NewWithNodeSet(ns))
				if err != nil {
					log.Panic()
				}
				err = p.Node.CurChain.Triedb.Commit(rt, false)
				if err != nil {
					log.Panic(err)
				}

				if p.Node.nodeID == "N0" {
					//将交易放入交易池
					p.Node.CurChain.Tx_pool.AddTxs(intx_pool)
					//发送完成指令给各个主节点
					fmt.Println("准备sendsure！")
					p.SendSure()

					// 将两个pool置空
					inacc_pool = make([]*address_and_balance, 0)
					intx_pool = make([]*core.Transaction, 0)
				}

				p.ssequenceID += 1

				p.isReply[c.Digest] = true

				if p.Node.nodeID == "N0" {
					p.sequenceLock.Unlock()
				}

			}

		}

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
		p.Node.CurChain.AddBlock(sb.Blocks[id-sb.StartID])
		fmt.Printf("编号为 %d 的区块已加入本地区块链！", id)
		curBlock := p.Node.CurChain.CurrentBlock
		fmt.Printf("curBlock: \n")
		curBlock.PrintBlock()
	}
	p.sequenceID = sb.EndID + 1
	p.requestLock.Unlock()
}

// 向除自己外的其他节点进行广播
func (p *Pbft) broadcast(cmd command, content []byte) {
	for i := range p.nodeTable {
		if i == p.Node.nodeID {
			continue
		}
		message := jointMessage(cmd, content)
		go utils.TcpDial(message, p.nodeTable[i])
	}
}

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

// 返回一个十位数的随机数，作为msgid
func getRandom() int {
	x := big.NewInt(30000000000)
	for {
		result, err := rand.Int(rand.Reader, x)
		if err != nil {
			log.Panic(err)
		}
		if result.Int64() > 1000000000 {
			return int(result.Int64())
		}
	}
}

// relay
func (p *Pbft) TryRelay(relaytxs, alltxs []*core.Transaction, st *trie.Trie, num int) {
	fmt.Printf("\n第%v个块开始准备relay\n\n", num)
	// use a memory trie database to do this, instead of disk database
	triedb := trie.NewDatabase(rawdb.NewMemoryDatabase())
	transactionTree := trie.NewEmpty(triedb)
	for _, tx := range alltxs {
		transactionTree.Update(tx.TxHash, tx.Encode())
	}

	config := params.Config
	//}
	relaypool := make([][]*core.TXrelay, config.Shard_num)
	for _, v := range relaytxs {
		shardID := account.Addr2Shard(hex.EncodeToString(v.Recipient))
		proofDB1 := &core.ProofDB{}
		err1 := transactionTree.Prove(v.TxHash, 0, proofDB1)
		if err1 != nil {
			log.Panic(err1)
		}

		proofDB2 := &core.ProofDB{}
		err2 := st.Prove(v.Recipient, 0, proofDB2)
		if err2 != nil {
			log.Panic(err2)
		}
		relaypool[shardID] = append(relaypool[shardID], &core.TXrelay{Txcs: v, MPcs: proofDB1, State: account.DecodeAccountState(st.Get(v.Recipient)), MPstate: proofDB2})
	}

	for k, v := range params.NodeTable {
		if k == config.ShardID {
			continue
		}
		txs := relaypool[params.ShardTable[k]]
		if len(txs) != 0 {
			target_leader := v["N0"]
			r := Relay{
				Txs:     txs,
				ShardID: config.ShardID,
			}
			bc, err := json.Marshal(r)
			if err != nil {
				log.Panic(err)
			}

			// fmt.Printf("正在向分片%v的主节点发送relay交易\n", k)
			message := jointMessage(cRelay, bc)
			go utils.TcpDial(message, target_leader)
		}
	}
	fmt.Printf("\n第%v个块结束relay\n\n", num)
}

// relay
func (p *Pbft) TryRelay2(relaytxs []*core.TXrelay) {
	config := params.Config
	relaypool := make([][]*core.TXrelay, config.Shard_num)
	for _, v := range relaytxs {
		shardID := account.Addr2Shard(hex.EncodeToString(v.Txcs.Recipient))
		relaypool[shardID] = append(relaypool[shardID], v)
	}

	for k, v := range params.NodeTable {
		if k == config.ShardID {
			continue
		}
		txs := relaypool[params.ShardTable[k]]
		if len(txs) != 0 {
			target_leader := v["N0"]
			r := Relay{
				Txs:     txs,
				ShardID: config.ShardID,
			}
			bc, err := json.Marshal(r)
			if err != nil {
				log.Panic(err)
			}

			// fmt.Printf("正在向分片%v的主节点发送relay交易\n", k)
			message := jointMessage(cRelay, bc)
			go utils.TcpDial(message, target_leader)
		}
	}
}

func (p *Pbft) TrySendTX(txs []*core.Transaction) {
	config := params.Config

	sendpool := make([][]*core.Transaction, config.Shard_num)
	for _, v := range txs {
		shardID := account.Account2Shard[hex.EncodeToString(v.Sender)]
		sendpool[shardID] = append(sendpool[shardID], v)
	}

	for k, v := range params.NodeTable {
		if k == config.ShardID {
			continue
		}
		txs := sendpool[params.ShardTable[k]]
		if len(txs) != 0 {
			target_leader := v["N0"]
			r := TxFromClient{
				Txs: txs,
			}
			bc, err := json.Marshal(r)
			if err != nil {
				log.Panic(err)
			}

			// fmt.Printf("正在向分片%v的主节点发送relay交易\n", k)
			message := jointMessage(cClient, bc)
			go utils.TcpDial(message, target_leader)
		}
	}
}

func (p *Pbft) handleRelay(content []byte) {
	relay := new(Relay)
	err := json.Unmarshal(content, relay)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Printf("\n本节点已接收到分片%v发来的relay交易，数量为：%v \n\n", relay.ShardID, len(relay.Txs))
	//如果接收者不属于本分片了，则发给目标分片
	self_shardID := params.ShardTable[params.Config.ShardID]
	relaytx2 := make([]*core.TXrelay, 0)
	p.Node.CurChain.Tx_pool.Relaypoollock.Lock()
	txcss := make([]*core.Transaction, 0)
	for _, tx := range relay.Txs {
		txcs := tx.Txcs
		target_shardID := account.Addr2Shard(hex.EncodeToString(txcs.Recipient))
		account.Lock_Acc_Lock.Lock()
		if !params.Config.Not_Lock_immediately && account.Lock_Acc[hex.EncodeToString(txcs.Recipient)] {
			if txcs.LockTime > 0 {
				txcs.LockTime2 = time.Now().UnixMilli()
			} else {
				txcs.LockTime = time.Now().UnixMilli()
			}
			txcs.RecLock = true
			txcs.Second_RequestTime = time.Now().UnixMilli()
			p.Node.CurChain.Tx_pool.Locking_TX_Pools[hex.EncodeToString(txcs.Recipient)] = append(p.Node.CurChain.Tx_pool.Locking_TX_Pools[hex.EncodeToString(txcs.Recipient)], txcs)
			account.Lock_Acc_Lock.Unlock()
			continue
		}
		account.Lock_Acc_Lock.Unlock()
		if target_shardID == self_shardID {
			txcss = append(txcss, txcs)
			txcs.Second_RequestTime = time.Now().UnixMilli()
		} else {
			relaytx2 = append(relaytx2, tx)
		}
	}
	p.Node.CurChain.Tx_pool.AddTxs(txcss)
	p.Node.CurChain.Tx_pool.Relaypoollock.Unlock()
	if len(relaytx2) != 0 {
		p.TryRelay2(relaytx2)
	}
	// fmt.Printf("\n本节点结束接收relay交易，数量为：%v \n\n", len(relay.Txs))

}

func (p *Pbft) TryTXmig1(txmig1s []*core.TXmig1, outbalance map[string]*big.Int, st, migTree *trie.Trie) {
	config := params.Config

	txmig2pool := make([][]*core.TXmig2, config.Shard_num)
	for _, v := range txmig1s {
		shardID := v.ToshardID
		hex_address, _ := hex.DecodeString(v.Address)
		proofDB1 := &core.ProofDB{}
		err1 := migTree.Prove(v.Hash(), 0, proofDB1)
		if err1 != nil {
			log.Panic(err1)
		}

		proofDB2 := &core.ProofDB{}
		err2 := st.Prove(hex_address, 0, proofDB2)
		if err2 != nil {
			log.Panic(err2)
		}

		encoded := st.Get(hex_address)
		if encoded == nil {
			log.Panic()
		}

		if config.Cross_Chain && !config.Fail && len(txmig1s) != 0 {
			p.TryOutPool[shardID] = append(p.TryOutPool[shardID], &core.TXmig2{Txmig1: v, MPmig1: proofDB1, State: account.DecodeAccountState(encoded), MPstate: proofDB2, Address: v.Address, Value: new(big.Int).Set(outbalance[v.Address])})
		} else {
			txmig2pool[shardID] = append(txmig2pool[shardID], &core.TXmig2{Txmig1: v, MPmig1: proofDB1, State: account.DecodeAccountState(encoded), MPstate: proofDB2, Address: v.Address, Value: new(big.Int).Set(outbalance[v.Address])})
		}
	}

	if config.Cross_Chain && !config.Fail {
		if p.sequenceID >= 3 {
			txmig2pool = p.TryOutPool
			p.TryOutPool = make([][]*core.TXmig2, config.Shard_num)
		} else {
			return
		}
	}

	if len(txmig2pool) != 0 {
		for k, v := range params.NodeTable {
			if k == config.ShardID {
				continue
			}
			txmig2s := txmig2pool[params.ShardTable[k]]
			if len(txmig2s) != 0 {
				target_leader := v["N0"]
				o := Mig2{
					TXmig2s: txmig2s,
					ShardID: config.ShardID,
				}
				bc, err := json.Marshal(o)
				if err != nil {
					log.Panic(err)
				}

				// fmt.Printf("正在向分片%v的主节点发送迁出\n", k)
				message := jointMessage(cTXmig1, bc)
				go utils.TcpDial(message, target_leader)
			}
		}
	}
}

func (p *Pbft) handleMig2(content []byte) {
	mig2 := new(Mig2)
	err := json.Unmarshal(content, mig2)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println()
	fmt.Printf("本节点已接收到分片%v发来的迁移请求 \n", mig2.ShardID)
	fmt.Println(mig2.TXmig2s[0].Value)
	fmt.Println()

	p.Node.CurChain.TXmig2_pool.AddTXmig2s(mig2.TXmig2s)

	// if !params.Config.Cross_Chain {
	// 	p.Node.CurChain.InAccout1_pool.AddIn1s(in1.Outs)
	// }else {
	// 	t := time.Now().Unix()
	// 	go p.wait_N_addIn1s(t, in1.Outs)
	// }

}

func (p *Pbft) wait_N_addIn1s(t int64, migs []*core.TXmig2) {

	if params.Config.ShardID == "S0" {
		for time.Now().Unix()-t != 6*15+1 {

		}
	} else {
		for time.Now().Unix()-t != 2*5+1 {

		}
	}

	p.Node.CurChain.TXmig2_pool.AddTXmig2s(migs)

}

func (p *Pbft) TryAnnounce(txmig2s []*core.TXmig2, st, migTree *trie.Trie) {
	config := params.Config

	// addrs := make([]string, 0)
	txanns := make([]*core.TXann, 0)
	for _, v := range txmig2s {
		hex_address, _ := hex.DecodeString(v.Address)

		proofDB1 := &core.ProofDB{}
		err1 := migTree.Prove(v.Hash(), 0, proofDB1)
		if err1 != nil {
			log.Panic(err1)
		}

		proofDB2 := &core.ProofDB{}
		err2 := st.Prove(hex_address, 0, proofDB2)
		if err2 != nil {
			log.Panic(err2)
		}

		encoded := st.Get(hex_address)
		if encoded == nil {
			log.Panic()
		}
		txanns = append(txanns, &core.TXann{Txmig2: v, MPmig2: proofDB1, State: account.DecodeAccountState(encoded), MPstate: proofDB2, Address: v.Address, ToshardID: params.ShardTable[config.ShardID]})
	}

	if config.ClientSendTX {
		a := Announce{
			TXanns:  txanns,
			ShardID: config.ShardID,
		}
		bc, err := json.Marshal(a)
		if err != nil {
			log.Panic(err)
		}

		// fmt.Printf("正在向分片%v的主节点发送通知\n", k)
		message := jointMessage(cAnnounce, bc)
		go utils.TcpDial(message, params.ClientAddr)

	}

	for k, v := range params.NodeTable {
		if k == config.ShardID {
			continue
		}
		target_leader := v["N0"]
		a := Announce{
			// Addrs:   addrs,
			TXanns:  txanns,
			ShardID: config.ShardID,
		}
		bc, err := json.Marshal(a)
		if err != nil {
			log.Panic(err)
		}

		// fmt.Printf("正在向分片%v的主节点发送通知\n", k)
		message := jointMessage(cAnnounce, bc)
		go utils.TcpDial(message, target_leader)
	}

}

// func (p *Pbft) handleAnnounce(content []byte) {
// 	announce := new(Announce)
// 	err := json.Unmarshal(content, announce)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	fmt.Printf("本节点已接收到分片%v发来通知 \n", announce.ShardID)

// 	p.Node.CurChain.Tx_pool.Relaypoollock.Lock()
// 	p.sequenceLock.Lock()
// 	fmt.Println(1)
// 	p.Node.CurChain.Tx_pool.Lock.Lock()
// 	fmt.Println(2)
// 	account.Account2ShardLock.Lock()
// 	fmt.Println(3)

// 	nowqueue := []*core.Transaction{}
// 	core.TxPoolDeepCopy(&nowqueue, p.Node.CurChain.Tx_pool.Queue)

// 	cAps := make(map[string]*ChangeAndPending)
// 	for _, v := range announce.Addrs {
// 		// 将该账户映射到目标分片
// 		account.Account2Shard[v] = params.ShardTable[announce.ShardID]

// 		// 若该账户目前在本分片，需要
// 		//   1.从本分片账户列表删除，
// 		////   2.收集关于该账户的交易，
// 		//   3.计算账户余额的变化，
// 		//   4.将该账户放入内存中的 删除账户队列中， 这个队列将在下一个区块里面
// 		//   5.将账户的余额变化与pending交易放入map[string]struct中，待会一起发送给目标分片
// 		if account.AccountInOwnShard[v] {

// 			// 将该账户设为要出，并分配内存存储要进交易池的交易
// 			account.Outing_Acc_After_Announce_Lock.Lock()
// 			account.Outing_Acc_After_Announce[v] = true
// 			p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools[v] = make([]*core.Transaction, 0)
// 			account.Outing_Acc_After_Announce_Lock.Unlock()

// 			cAp := &ChangeAndPending{}
// 			cAp.PendingTxs = make([]*core.Transaction, 0)

// 			//   1.从本分片账户列表删除，
// 			delete(account.AccountInOwnShard, v)

// 			if params.Config.Lock_Acc_When_Migrating {
// 				account.Lock_Acc_Lock.Lock()
// 				delete(account.Lock_Acc, v)
// 				account.Lock_Acc_Lock.Unlock()
// 			}else if !params.Config.Stop_When_Migrating {
// 				account.Outing_Acc_Before_Announce_Lock.Lock()
// 				delete(account.Outing_Acc_Before_Announce, v)
// 				account.Outing_Acc_Before_Announce_Lock.Unlock()
// 			}

// 			// 2.如果是lock机制，将关于账户的被锁交易中，sender为本分片的交易放回交易池中
// 			if params.Config.Lock_Acc_When_Migrating {
// 				j := 0
// 				account.Lock_Acc_Lock.Lock()
// 				for _, tx := range p.Node.CurChain.Tx_pool.Locking_TX_Pools[v] {
// 					// 如果不是relay且v为接收者，则放回交易池
// 					if hex.EncodeToString(tx.Recipient) == v && !tx.IsRelay {
// 						tx.UnlockTime = time.Now().UnixMilli()
// 						tx.Success    = true
// 						p.Node.CurChain.Tx_pool.Queue = append(p.Node.CurChain.Tx_pool.Queue, tx)
// 						continue
// 					}
// 					// 否则，留在被锁交易中
// 					p.Node.CurChain.Tx_pool.Locking_TX_Pools[v][j] = tx
// 					j++
// 				}
// 				p.Node.CurChain.Tx_pool.Locking_TX_Pools[v] = p.Node.CurChain.Tx_pool.Locking_TX_Pools[v][:j]
// 				account.Lock_Acc_Lock.Unlock()
// 			}

// 			//   3.计算账户余额的变化，
// 			if !params.Config.Lock_Acc_When_Migrating {
// 				hex_address, _ := hex.DecodeString(v)
// 				decoded, success := p.Node.CurChain.StatusTrie.Get(hex_address)
// 				if !success {
// 					log.Panic()
// 				}
// 				accoount_state := account.DecodeAccountState(decoded)
// 				account.BalanceBeforeOutLock.Lock()
// 				cAp.Change = accoount_state.Balance - account.BalanceBeforeOut[v]
// 				//  可以从内存删除该账户在迁移1时的余额
// 				delete(account.BalanceBeforeOut, v)
// 				account.BalanceBeforeOutLock.Unlock()
// 			} else {
// 				cAp.Change = 0
// 			}

// 			//   4.将该账户放入内存中的 删除账户队列中， 这个队列将在下一个区块里面
// 			p.Node.CurChain.Delete_pool.AddDelete(v, params.ShardTable[announce.ShardID])

// 			//   5.将账户的余额变化与pending交易放入列表中，待会一起发送给目标分片
// 			cAps[v] = cAp
// 		}

// 	}

// 	account.Account2ShardLock.Unlock()
// 	p.Node.CurChain.Tx_pool.Lock.Unlock()
// 	p.sequenceLock.Unlock()
// 	p.Node.CurChain.Tx_pool.Relaypoollock.Unlock()

// 	// //   2.将交易池中关于账户的交易收集起来
// 	// for _,cAp := range cAps {
// 	// 	if params.Config.Lock_Acc_When_Migrating {
// 	// 		account.Lock_Acc_Lock.Lock()
// 	// 		cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Locking_TX_Pools[cAp.Address]...)
// 	// 		delete(p.Node.CurChain.Tx_pool.Locking_TX_Pools, cAp.Address)
// 	// 		account.Lock_Acc_Lock.Unlock()
// 	// 	}

// 	// 	for _, tx := range nowqueue {
// 	// 		// 如果sender是要出去的账户，且不是relaytx, 则收集
// 	// 		if hex.EncodeToString(tx.Sender) == cAp.Address && !tx.IsRelay {
// 	// 			cAp.PendingTxs = append(cAp.PendingTxs, tx)
// 	// 			continue
// 	// 		}
// 	// 		//  如果recipient是要出去的账户，且已经是relaytx，则收集
// 	// 		if hex.EncodeToString(tx.Recipient) == cAp.Address && tx.IsRelay {
// 	// 			cAp.PendingTxs = append(cAp.PendingTxs, tx)
// 	// 			continue
// 	// 		}
// 	// 		// 否则，要么与该账户无关，要么该账户是recipient，但是sender部分还没处理，因此啥也不做
// 	// 	}

// 	// 	account.Outing_Acc_After_Announce_Lock.Lock()
// 	// 	cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools[cAp.Address]...)
// 	// 	delete(account.Outing_Acc_After_Announce, cAp.Address)
// 	// 	delete(p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools, cAp.Address)
// 	// 	account.Outing_Acc_After_Announce_Lock.Unlock()
// 	// }

// 	//   2.将交易池中关于账户的交易收集起来
// 	if params.Config.Lock_Acc_When_Migrating {
// 		for addr, cAp := range cAps {
// 			account.Lock_Acc_Lock.Lock()
// 			cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Locking_TX_Pools[addr]...)
// 			delete(p.Node.CurChain.Tx_pool.Locking_TX_Pools, addr)
// 			account.Lock_Acc_Lock.Unlock()
// 		}
// 	}else if !params.Config.Stop_When_Migrating {
// 		for addr, cAp := range cAps {
// 			account.Outing_Acc_Before_Announce_Lock.Lock()
// 			cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools[addr]...)
// 			delete(p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools, addr)
// 			account.Outing_Acc_Before_Announce_Lock.Unlock()
// 		}
// 	}

// 	for _, tx := range nowqueue {
// 		// 如果sender是要出去的账户，且不是relaytx, 则收集
// 		if _, ok := cAps[hex.EncodeToString(tx.Sender)]; ok && !tx.IsRelay {
// 			cAps[hex.EncodeToString(tx.Sender)].PendingTxs = append(cAps[hex.EncodeToString(tx.Sender)].PendingTxs, tx)
// 			continue
// 		}

// 		//  如果recipient是要出去的账户，且已经是relaytx，则收集
// 		if _, ok := cAps[hex.EncodeToString(tx.Recipient)]; ok && tx.IsRelay {
// 			cAps[hex.EncodeToString(tx.Recipient)].PendingTxs = append(cAps[hex.EncodeToString(tx.Recipient)].PendingTxs, tx)
// 			continue
// 		}
// 		// 否则，啥也不做
// 	}

// 	for addr, cAp := range cAps {
// 		account.Outing_Acc_After_Announce_Lock.Lock()
// 		cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools[addr]...)
// 		delete(account.Outing_Acc_After_Announce, addr)
// 		delete(p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools, addr)
// 		account.Outing_Acc_After_Announce_Lock.Unlock()
// 	}

// 	//将所有涉及账户的余额变化和pending交易发送给目标分片
// 	p.TrySendChangesAndPendings(cAps, announce.ShardID)
// }

func (p *Pbft) handleAnnounce(content []byte) {
	announce := new(Announce)
	err := json.Unmarshal(content, announce)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到分片%v发来通知 \n", announce.ShardID)

	p.sequenceLock.Lock()
	p.Node.CurChain.Tx_pool.Relaypoollock.Lock()
	fmt.Println(1)
	p.Node.CurChain.Tx_pool.Lock.Lock()
	fmt.Println(2)
	account.Account2ShardLock.Lock()
	fmt.Println(3)

	// nowqueue := []*core.Transaction{}
	// core.TxPoolDeepCopy(&nowqueue, p.Node.CurChain.Tx_pool.Queue)

	// cAps := make(map[string]*ChangeAndPending)
	for _, v := range announce.TXanns {
		// 若账户本就不属于本分片，直接将该ann加入TXann队列
		if !account.AccountInOwnShard[v.Address] {
			// account.Account2Shard[v] = params.ShardTable[announce.ShardID]
			p.Node.CurChain.TXann_pool.AddTXann(v)
			continue
		}

		// 若该账户目前在本分片，需要
		////   1.从本分片账户列表删除，
		////   2.收集关于该账户的交易，
		////   3.计算账户余额的变化，
		//   4.将该账户放入内存中的 删除账户队列中， 这个队列将在下一个区块里面
		////   5.将账户的余额变化与pending交易放入map[string]struct中，待会一起发送给目标分片
		if account.AccountInOwnShard[v.Address] {

			// 将该账户设为要出，并分配内存存储要进交易池的交易
			account.Outing_Acc_After_Announce_Lock.Lock()
			account.Outing_Acc_After_Announce[v.Address] = true
			p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools[v.Address] = make([]*core.Transaction, 0)
			account.Outing_Acc_After_Announce_Lock.Unlock()

			//   4.将该账户放入内存中的 删除账户队列中， 这个队列将在下一个区块里面
			p.Node.CurChain.TXann_pool.AddTXann(v)
		}

	}

	account.Account2ShardLock.Unlock()
	p.Node.CurChain.Tx_pool.Lock.Unlock()
	p.Node.CurChain.Tx_pool.Relaypoollock.Unlock()
	p.sequenceLock.Unlock()
}

func (p *Pbft) TrySendChangesAndPendings(nss []*core.TXns, caps map[string]*ChangeAndPending, target_shardID string) {
	config := params.Config

	target_leader := params.NodeTable[target_shardID]["N0"]
	cp := ChangesAndPendings{
		TXnss:   nss,
		List:    caps,
		ShardID: config.ShardID,
	}

	bc, err := json.Marshal(cp)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("时间：%v  正在向分片%v的主节点发送新\n", time.Now().UnixMilli()-InitTime*1000, target_shardID)
	message := jointMessage(cCaP, bc)
	go utils.TcpDial(message, target_leader)

}

func (p *Pbft) handleChangesAndPendings(content []byte) {
	fmt.Printf("时间：%v  本节点已接收到新\n", time.Now().UnixMilli())
	caps := new(ChangesAndPendings)
	err := json.Unmarshal(content, caps)
	if err != nil {
		log.Panic(err)
	}
	// fmt.Printf("本节点已接收到分片%v发来的余额变化和交易 \n", caps.ShardID)

	//将pending放入交易池，将changes放入change池
	for _, cap := range caps.List {

		for _, v := range cap.PendingTxs {
			// to := hex.EncodeToString(v.Recipient)
			// //全锁
			// if params.Config.Lock_Acc_When_Migrating {
			// 	account.Lock_Acc_Lock.Lock()
			// 	if account.Lock_Acc[to] {
			// 		v.RecLock = true
			// 		p.Node.CurChain.Tx_pool.Locking_TX_Pools[to] = append(p.Node.CurChain.Tx_pool.Locking_TX_Pools[to], v)
			// 		account.Lock_Acc_Lock.Unlock()
			// 		continue
			// 	}
			// 	account.Lock_Acc_Lock.Unlock()
			// }
			// //半锁什么都不用做

			// if v.UnlockTime > 0 {
			// 	v.UnlockTime2 = time.Now().UnixMilli()
			// } else {
			// 	v.UnlockTime = time.Now().UnixMilli()
			// }
			v.Success = true
			p.Node.CurChain.Tx_pool.Lock.Lock()
			p.Node.CurChain.Tx_pool.AddTx(v)
			p.Node.CurChain.Tx_pool.Lock.Unlock()
		}
		// if cap.Change != 0 {
		// }
		// fmt.Println()
		// fmt.Println(cap.Change)
		// fmt.Println()

		// 收到5k个啦！！
	}
	if !p.Node.CurChain.ChainConfig.Lock_Acc_When_Migrating {
		p.Node.CurChain.TXns_pool.AddTXnss(caps.TXnss)
	}

	// 将 交易池中的交易发给 client
	c, err := json.Marshal(len(caps.List))
	if err != nil {
		log.Panic(err)
	}
	m := jointMessage(cNumOfMigratedTXsAddrs, c)
	utils.TcpDial(m, params.ClientAddr)
	fmt.Println("本节点已将 迁移过来的pending交易的 账户数量发送到client\n")
}

func (p *Pbft) handleMigrateWanted() {
	fmt.Println("本节点已接收到client发来准备迁移通知")

	p.sequenceLock.Lock()
	p.Node.CurChain.Tx_pool.Relaypoollock.Lock()
	fmt.Println(1)
	p.Node.CurChain.Tx_pool.Lock.Lock()
	fmt.Println(2)
	account.Account2ShardLock.Lock()
	fmt.Println(3)

	// 将 交易池中的交易发给 client
	c, err := json.Marshal(&p.Node.CurChain.Tx_pool.Queue)
	if err != nil {
		log.Panic(err)
	}
	m := jointMessage(cPendingTXs, c)
	utils.TcpDial(m, params.ClientAddr)
	fmt.Println("本节点已将pending发送到client\n")

	account.Account2ShardLock.Unlock()
	p.Node.CurChain.Tx_pool.Lock.Unlock()
	p.Node.CurChain.Tx_pool.Relaypoollock.Unlock()
	p.sequenceLock.Unlock()
}

func (p *Pbft) handleEpochChange() {
	p.epochLock.Lock()
	epochchangeStartTime = time.Now().Unix()
	// 用算法决定新的映射表
	account.Account2ShardLock.Lock()
	new := algorithm.Algorithm2(account.Account2Shard, params.ShardTable[params.Config.ShardID])
	account.Account2ShardLock.Unlock()
	p.mpropose(new)
}

func (p *Pbft) handleNewMap(content []byte) {
	nam := NaM{}
	err := json.Unmarshal(content, &nam)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到客户端发送的新映射\n")
	new_addr2shard := nam.New_Addr2shard
	new_addrs := nam.New_Addrs

	if params.Config.Stop_When_Migrating {
		p.epochLock.Lock()
		epochchangeStartTime = time.Now().Unix()
		// // 用算法决定新的映射表
		// account.Account2ShardLock.Lock()
		// new := algorithm.Algorithm2(account.Account2Shard, params.ShardTable[params.Config.ShardID])
		// account.Account2ShardLock.Unlock()

		p.mpropose(new_addr2shard)
	} else {
		// 找到自己要出去的
		account.Account2ShardLock.Lock()
		for _, addr := range new_addrs {
			if account.AccountInOwnShard[addr] && new_addr2shard[addr] != params.ShardTable[params.Config.ShardID] {
				p.Node.CurChain.TXmig1_pool.AddTXmig1(&core.TXmig1{Address: addr, FromshardID: params.ShardTable[params.Config.ShardID], ToshardID: new_addr2shard[addr], Request_Time: time.Now().UnixMilli()})
			}
		}
		account.Account2ShardLock.Unlock()
	}

}

// 接收到初始时间
func (p *Pbft) handleLLT(content []byte) {
	var inittime int64
	err := json.Unmarshal(content, &inittime)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("本节点已接收到客户端发送的初始时间：%v\n", inittime)
	InitTime = inittime
	core.InitTime = inittime
}
