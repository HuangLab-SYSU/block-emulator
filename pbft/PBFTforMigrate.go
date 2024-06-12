package pbft

import (
	"blockEmulator/account"
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/utils"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

type address_and_balance struct {
	Address string
	Balance *big.Int
}

var (
	inacc_count          int
	inacc_count_lock     sync.Mutex
	inacc_pool           []*address_and_balance
	intx_pool            []*core.Transaction
	sure_count           int // 完成的分片的数量
	sure_count_lock      sync.Mutex
	epochchangeStartTime int64
	sendoutfinish        chan int
)

// func (p *Pbft) mpropose(new map[string]int) {
// 	for {
// 		p.pbftlock.Lock()
// 		if int(time.Now().Unix() - InitTime) > p.epoch*params.Config.Block_interval   && int(time.Now().Unix() - lastpropose) >= params.Config.Block_interval + utils.RandInt0To3(){
// 			p.epoch++
// 			p.pbftlock.Unlock()
// 			break
// 		}
// 		p.pbftlock.Unlock()
// 	}
// 	p.sequenceLock.Lock() //通过锁强制要求上一个PBFTcommit完成新propose才能被提出
// 	lastpropose = time.Now().Unix()
// 	fmt.Println("开始epochChange")

// 	r := &Request{}
// 	r.Timestamp = time.Now().Unix()
// 	r.Message.ID = getRandom()

// 	// 对新的映射进行编码
// 	var buff bytes.Buffer
// 	enc := gob.NewEncoder(&buff)
// 	err := enc.Encode(new)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	r.Message.Content = buff.Bytes()

// 	pbftbefore = time.Now().Unix()

// 	//获取消息摘要
// 	digest := getDigest(r)
// 	fmt.Println("已将request存入临时消息池")
// 	//存入临时消息池
// 	p.mmessagePool[digest] = r

// 	//拼接成PrePrepare，准备发往follower节点
// 	pp := PrePrepare{r, digest, p.msequenceID, "EpochChange"}
// 	b, err := json.Marshal(pp)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	// fmt.Println("正在向其他节点进行进行PrePrepare广播 ...")
// 	//进行PrePrepare广播
// 	p.broadcast(cPrePrepare, b)
// 	// fmt.Println("PrePrepare广播完成")
// }

// func (p *Pbft) spropose(inaccpool []*address_and_balance) {
// 	for {
// 		p.pbftlock.Lock()
// 		if int(time.Now().Unix() - InitTime) > p.epoch*params.Config.Block_interval   && int(time.Now().Unix() - lastpropose) >= params.Config.Block_interval + utils.RandInt0To3(){
// 			p.epoch++
// 			p.pbftlock.Unlock()
// 			break
// 		}
// 		p.pbftlock.Unlock()
// 	}
// 	fmt.Println("\n提出spropose1\n")
// 	p.sequenceLock.Lock() //通过锁强制要求上一个PBFTcommit完成新propose才能被提出
// 	lastpropose = time.Now().Unix()
// 	fmt.Println("\n提出spropose2\n")

// 	r := &Request{}
// 	r.Timestamp = time.Now().Unix()
// 	r.Message.ID = getRandom()

// 	// 对新的映射进行编码
// 	var buff bytes.Buffer
// 	enc := gob.NewEncoder(&buff)
// 	err := enc.Encode(inaccpool)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	r.Message.Content = buff.Bytes()

// 	pbftbefore = time.Now().Unix()

// 	//获取消息摘要
// 	digest := getDigest(r)
// 	fmt.Println("已将request存入临时消息池")
// 	//存入临时消息池
// 	p.smessagePool[digest] = r

// 	//拼接成PrePrepare，准备发往follower节点
// 	pp := PrePrepare{r, digest, p.ssequenceID, "AccState"}
// 	b, err := json.Marshal(pp)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	// fmt.Println("正在向其他节点进行进行PrePrepare广播 ...")
// 	//进行PrePrepare广播
// 	p.broadcast(cPrePrepare, b)
// 	// fmt.Println("PrePrepare广播完成")
// }

// 一个分片一个节点
func (p *Pbft) mpropose(new map[string]int) {
	for {
		p.pbftlock.Lock()
		if int(time.Now().Unix()-InitTime) > p.epoch*params.Config.Block_interval && int(time.Now().Unix()-lastpropose) >= params.Config.Block_interval+utils.RandInt0To3(InitTime+int64(p.epoch)) {
			p.epoch++
			p.pbftlock.Unlock()
			break
		}
		p.pbftlock.Unlock()
	}
	p.sequenceLock.Lock() //通过锁强制要求上一个PBFTcommit完成新propose才能被提出
	lastpropose = time.Now().Unix()
	fmt.Println("开始epochChange")

	// 对新的映射进行编码
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(new)
	if err != nil {
		log.Panic(err)
	}
	content := buff.Bytes()
	pbftType := "EpochChange"

	p.commit1(content, pbftType)
}

func (p *Pbft) spropose(inaccpool []*address_and_balance) {
	for {
		p.pbftlock.Lock()
		if int(time.Now().Unix()-InitTime) > p.epoch*params.Config.Block_interval && int(time.Now().Unix()-lastpropose) >= params.Config.Block_interval+utils.RandInt0To3(InitTime+int64(p.epoch)) {
			p.epoch++
			p.pbftlock.Unlock()
			break
		}
		p.pbftlock.Unlock()
	}
	fmt.Println("\n提出spropose1\n")
	p.sequenceLock.Lock() //通过锁强制要求上一个PBFTcommit完成新propose才能被提出
	lastpropose = time.Now().Unix()
	fmt.Println("\n提出spropose2\n")

	// 对新的映射进行编码
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(inaccpool)
	if err != nil {
		log.Panic(err)
	}
	content := buff.Bytes()
	pbftType := "AccState"

	p.commit1(content, pbftType)
}

func (p *Pbft) SendOut(out map[string]int) {
	config := params.Config
	ccc := time.Now().UnixMicro()

	outpool := make([]*BalancesAndPendings, config.Shard_num)
	// 遍历交易，寻找到要出去的交易
	j := 0
	p.Node.CurChain.Tx_pool.Lock.Lock()
	// build trie from the triedb (in disk)
	st, err := trie.New(trie.TrieID(common.BytesToHash(p.Node.CurChain.CurrentBlock.Header.StateRoot)), p.Node.CurChain.Triedb)
	if err != nil {
		log.Panic(err)
	}
	for _, tx := range p.Node.CurChain.Tx_pool.Queue {

		// 如果sender是要出去的账户，且还不是后半，则收集
		if shard, ok := out[hex.EncodeToString(tx.Sender)]; ok && !tx.IsRelay {
			if outpool[shard] == nil {
				outpool[shard] = &BalancesAndPendings{}
				outpool[shard].List = make(map[string]*BalanceAndPending)
			}
			if _, ok2 := outpool[shard].List[hex.EncodeToString(tx.Sender)]; !ok2 {
				encoded_state := st.Get(tx.Sender)
				if encoded_state == nil {
					log.Panic()
				}
				state := account.DecodeAccountState(encoded_state)
				balance := new(big.Int).Set(state.Balance)
				outpool[shard].List[hex.EncodeToString(tx.Sender)] = &BalanceAndPending{Balance: balance, PendingTxs: []*core.Transaction{}}
			}
			outpool[shard].List[hex.EncodeToString(tx.Sender)].PendingTxs = append(outpool[shard].List[hex.EncodeToString(tx.Sender)].PendingTxs, tx)
			continue
		}

		//  如果recipient是要出去的账户，且已经是后半，则收集
		if shard, ok := out[hex.EncodeToString(tx.Recipient)]; ok && tx.IsRelay {
			if outpool[shard] == nil {
				outpool[shard] = &BalancesAndPendings{}
				outpool[shard].List = make(map[string]*BalanceAndPending)
			}
			if _, ok2 := outpool[shard].List[hex.EncodeToString(tx.Recipient)]; !ok2 {
				encoded_state := st.Get(tx.Recipient)
				if encoded_state == nil {
					log.Panic()
				}
				state := account.DecodeAccountState(encoded_state)
				balance := state.Balance
				outpool[shard].List[hex.EncodeToString(tx.Recipient)] = &BalanceAndPending{Balance: balance, PendingTxs: []*core.Transaction{}}
			}
			outpool[shard].List[hex.EncodeToString(tx.Recipient)].PendingTxs = append(outpool[shard].List[hex.EncodeToString(tx.Recipient)].PendingTxs, tx)
			continue
		}

		// 否则，这笔交易不用出去
		p.Node.CurChain.Tx_pool.Queue[j] = tx
		j++
	}
	p.Node.CurChain.Tx_pool.Queue = p.Node.CurChain.Tx_pool.Queue[:j]
	p.Node.CurChain.Tx_pool.Lock.Unlock()

	// 可能有的要出去的账户没有交易需要出去，所以上面没有遍历到，因此再遍历出去的账户，看看有没有没加到的
	for addr, shard := range out {
		if outpool[shard] == nil {
			outpool[shard] = &BalancesAndPendings{}
			outpool[shard].List = make(map[string]*BalanceAndPending)
		}
		if _, ok := outpool[shard].List[addr]; !ok {
			outpool[shard].List[addr] = &BalanceAndPending{Balance: big.NewInt(0), PendingTxs: []*core.Transaction{}}
		}
	}
	fmt.Printf("\n%vms\n\n", time.Now().UnixMicro()-ccc)

	// for k, v := range out {
	// 	if outpool[v] == nil {
	// 		outpool[v] = &BalancesAndPendings{}
	// 	}
	// 	hex_out, _ := hex.DecodeString(k)
	// 	encoded_state, ok := p.Node.CurChain.StatusTrie.Get(hex_out)
	// 	if !ok {
	// 		log.Panic()
	// 	}
	// 	state := account.DecodeAccountState(encoded_state)
	// 	balance := state.Balance

	// 	pendingTXs := make([]*core.Transaction, 0)
	// 	//   2.将交易池中关于账户的交易收集起来
	// 	j := 0
	// 	p.Node.CurChain.Tx_pool.Lock.Lock()
	// 	hhh := time.Now().UnixMicro()
	// 	for _, tx := range p.Node.CurChain.Tx_pool.Queue {
	// 		// 如果sender是要出去的账户，且还不是后半, 则收集
	// 		if hex.EncodeToString(tx.Sender) == k && !tx.IsRelay{
	// 			pendingTXs = append(pendingTXs, tx)
	// 			continue
	// 		}
	// 		//  如果recipient是要出去的账户，且已经是后半，则收集
	// 		if hex.EncodeToString(tx.Recipient) == k && tx.IsRelay {
	// 			pendingTXs = append(pendingTXs, tx)
	// 			continue
	// 		}
	// 		// 否则，要么与该账户无关，要么该账户是recipient，但是sender部分还没处理
	// 		p.Node.CurChain.Tx_pool.Queue[j] = tx
	// 		j++
	// 	}
	// 	p.Node.CurChain.Tx_pool.Queue = p.Node.CurChain.Tx_pool.Queue[:j]
	// 	ccc += time.Now().UnixMicro()-hhh
	// 	p.Node.CurChain.Tx_pool.Lock.Unlock()

	// 	outpool[v].List = append(outpool[v].List, &BalanceAndPending{Address: k, Balance: balance, PendingTxs: pendingTXs})
	// }
	// fmt.Printf("\n%vms\n\n",ccc)

	for k, v := range params.NodeTable {
		if k == config.ShardID {
			continue
		}
		outs := outpool[params.ShardTable[k]]
		if outs != nil {
			target_leader := v["N0"]
			outs.ShardID = config.ShardID
			bc, err := json.Marshal(outs)
			if err != nil {
				log.Panic(err)
			}

			fmt.Printf("正在向分片%v的主节点发送迁出\n", k)
			message := jointMessage(cBalanceAndPending, bc)
			go utils.TcpDial(message, target_leader)
		}
	}
	sendoutfinish <- 1
}

func (p *Pbft) handleBalancesAndPendings(content []byte) {
	baps := new(BalancesAndPendings)
	err := json.Unmarshal(content, baps)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("\n本节点已接收到分片%v发来的账户状态和交易 \n\n", baps.ShardID)

	inacc_count_lock.Lock()
	//将pending放入intx池，将acc放入intacc池
	for addr, bap := range baps.List {
		inacc_pool = append(inacc_pool, &address_and_balance{Address: addr, Balance: bap.Balance})
		intx_pool = append(intx_pool, bap.PendingTxs...)

		//计数-1
		inacc_count--
	}

	fmt.Printf("\n收到状态后，inacc_count数量为：%v \n\n", inacc_count)

	//若计数为0，则发起片内共识
	if inacc_count == 0 {
		fmt.Printf("本节点已接收到所有账户和交易\n")
		p.spropose(inacc_pool)
	}
	inacc_count_lock.Unlock()

}

func (p *Pbft) SendSure() {
	<-sendoutfinish
	shardid := params.ShardTable[params.Config.ShardID]
	int, err := json.Marshal(shardid)
	message := jointMessage(cSure, int)
	if err != nil {
		log.Panic(err)
	}
	for k, v := range params.NodeTable {
		if k == params.Config.ShardID {
			continue
		}
		target_leader := v["N0"]
		fmt.Printf("\n正在向分片%v的主节点发送sure消息\n\n", k)
		go utils.TcpDial(message, target_leader)
	}

	sure_count_lock.Lock()
	sure_count++
	if sure_count == params.Config.Shard_num {
		sure_count = 0
		p.epochLock.Unlock()
		// fmt.Printf("\nepochchange时长为：%v\n\n", time.Now().Unix()-epochchangeStartTime)
		s := fmt.Sprintf("%v %v %v %v", p.msequenceID, epochchangeStartTime, time.Now().Unix(), time.Now().Unix()-epochchangeStartTime)
		epochChangelog.Write(strings.Split(s, " "))
		epochChangelog.Flush()
	}
	sure_count_lock.Unlock()
}

func (p *Pbft) handleSure(content []byte) {
	var shardid int
	err := json.Unmarshal(content, &shardid)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("\n收到 分片%v 发来的sure消息\n\n", shardid)

	sure_count_lock.Lock()
	sure_count++
	if sure_count == params.Config.Shard_num {
		sure_count = 0
		p.epochLock.Unlock()
		// fmt.Printf("\nepochchange时长为：%v\n\n", time.Now().Unix()-epochchangeStartTime)
		s := fmt.Sprintf("%v %v %v %v", p.msequenceID, epochchangeStartTime, time.Now().Unix(), time.Now().Unix()-epochchangeStartTime)
		epochChangelog.Write(strings.Split(s, " "))
		epochChangelog.Flush()
	}
	sure_count_lock.Unlock()
}
