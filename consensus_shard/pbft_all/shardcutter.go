package pbft_all

import (
	"blockEmulator/chain"
	"blockEmulator/consensus_shard/pbft_all/dataSupport"
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/networks"
	"blockEmulator/params"
	"encoding/json"
	"log"
	"sync"
	"time"
)

type source_query struct {
	mu              sync.Mutex
	receivedData    bool
	AccountKey      string
	AccountLocation uint64
}

type SHARD_CUTTER struct {
	cdm      *dataSupport.Data_supportCLPA
	pbftNode *PbftConsensusNode
	sq       source_query
}

// receive relay transaction, which is for cross shard txs
func (crom *SHARD_CUTTER) handleRelay(content []byte) {
	relay := new(message.Relay)
	err := json.Unmarshal(content, relay)
	if err != nil {
		log.Panic(err)
	}
	crom.pbftNode.pl.Plog.Printf("S%dN%d : has received relay txs from shard %d, the senderSeq is %d\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID, relay.SenderShardID, relay.SenderSeq)
	crom.pbftNode.CurChain.Txpool.AddTxs2Pool(relay.Txs)
	crom.pbftNode.seqMapLock.Lock()
	crom.pbftNode.seqIDMap[relay.SenderShardID] = relay.SenderSeq
	crom.pbftNode.seqMapLock.Unlock()
	crom.pbftNode.pl.Plog.Printf("S%dN%d : has handled relay txs msg\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
}

func (crom *SHARD_CUTTER) handleRelayWithProof(content []byte) {
	rwp := new(message.RelayWithProof)
	err := json.Unmarshal(content, rwp)
	if err != nil {
		log.Panic(err)
	}
	crom.pbftNode.pl.Plog.Printf("S%dN%d : has received relay txs & proofs from shard %d, the senderSeq is %d\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID, rwp.SenderShardID, rwp.SenderSeq)
	// validate the proofs of txs
	isAllCorrect := true
	for i, tx := range rwp.Txs {
		if ok, _ := chain.TxProofVerify(tx.TxHash, &rwp.TxProofs[i]); !ok {
			isAllCorrect = false
			break
		}
	}
	if isAllCorrect {
		crom.pbftNode.CurChain.Txpool.AddTxs2Pool(rwp.Txs)
	} else {
		crom.pbftNode.pl.Plog.Println("Err: wrong proof!")
	}

	crom.pbftNode.seqMapLock.Lock()
	crom.pbftNode.seqIDMap[rwp.SenderShardID] = rwp.SenderSeq
	crom.pbftNode.seqMapLock.Unlock()
	crom.pbftNode.pl.Plog.Printf("S%dN%d : has handled relay txs msg\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
}

func (crom *SHARD_CUTTER) handleInjectTx(content []byte) {
	it := new(message.InjectTxs)
	err := json.Unmarshal(content, it)
	if err != nil {
		log.Panic(err)
	}
	crom.pbftNode.CurChain.Txpool.AddTxs2Pool(it.Txs)
	crom.pbftNode.pl.Plog.Printf("S%dN%d : has handled injected txs msg, txs: %d \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID, len(it.Txs))
}

func (crom *SHARD_CUTTER) handlePartitionReady(content []byte) {
	pr := new(message.PartitionReady)
	err := json.Unmarshal(content, pr)
	if err != nil {
		log.Panic()
	}
	crom.cdm.P_ReadyLock.Lock()
	crom.cdm.PartitionReady[pr.FromShard] = true
	crom.cdm.P_ReadyLock.Unlock()

	crom.pbftNode.seqMapLock.Lock()
	crom.cdm.ReadySeq[pr.FromShard] = pr.NowSeqID
	crom.pbftNode.seqMapLock.Unlock()
}

// when the message from other shard arriving, it should be added into the message pool
func (crom *SHARD_CUTTER) handleAccountStateAndTxMsg(content []byte) {
	at := new(message.AccountStateAndTx)
	err := json.Unmarshal(content, at)
	if err != nil {
		log.Panic()
	}
	crom.cdm.AccountStateTx[at.FromShard] = at

	if len(crom.cdm.AccountStateTx) == int(crom.pbftNode.pbftChainConfig.ShardNums)-1 {
		crom.cdm.CollectLock.Lock()
		crom.cdm.CollectOver = true
		crom.cdm.CollectLock.Unlock()
	}
}

// stage1：源分片接收来自监管节点的消息(TXaux1)，更新账户状态，然后将消息TXaux2发送给目标分片
func (crom *SHARD_CUTTER) handlePartitionMsg(content []byte) {
	pm := new(message.PartitionModifiedMap)
	err := json.Unmarshal(content, pm)
	if err != nil {
		log.Panic()
	}

	// PartitionModified变量 key: 账户的编码, val: 目标分片编号
	for key, val := range pm.PartitionModified {
		// 如果这个账户源分片为当前分片，则处理
		if crom.pbftNode.CurChain.Get_PartitionMap(key) == crom.pbftNode.ShardID {
			crom.pbftNode.pl.Plog.Printf("S%dN%d : source shard received TXaux1 from M-shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
			// 生成TXaux1
			txau1 := core.TXmig1{
				Address:     key,
				FromshardID: crom.pbftNode.ShardID,
				ToshardID:   val,
			}
			// 发送TXaux2到目标分片
			sii := message.TXAUX_2_MSG{
				Msg: core.TXmig2{
					Txmig1: txau1,
					MPmig1: true,
					State: core.CLUSTER_ACCOUNT_STATE{
						Key:      key,
						Location: val,
						SourceID: crom.pbftNode.ShardID,
					},
					MPstate:         true,
					TimeoutDuration: 1000 * time.Second,
					StartTime:       time.Now(),
				},
				Sender: crom.pbftNode.ShardID,
			}
			sByte, err := json.Marshal(sii)
			if err != nil {
				log.Panic()
			}
			msg_send := message.MergeMessage(message.TXaux_2, sByte)
			go networks.TcpDial(msg_send, crom.pbftNode.ip_nodeTable[val][0])
			crom.pbftNode.pl.Plog.Printf("S%dN%d : account %s's location is updated to: %d \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID, key, val)
			crom.pbftNode.pl.Plog.Printf("S%dN%d : source shard sended Txaux2 to dest shard\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
			crom.pbftNode.pl.Plog.Printf("S%dN%d : source shard handle stage1 done \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
		}
	}
	crom.cdm.ModifiedMap = append(crom.cdm.ModifiedMap, pm.PartitionModified)
	crom.cdm.PartitionOn = true
}

// stage2：目标分片接收到TXaux2将进行验证，验证成功则更新账户状态，将TXann发送给所有分片
func (crom *SHARD_CUTTER) handleTXaux_2(content []byte) {
	crom.pbftNode.pl.Plog.Printf("S%dN%d : dest shard received TXaux2 from source shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
	data := new(message.TXAUX_2_MSG)
	err := json.Unmarshal(content, data)
	if err != nil {
		log.Panic()
	}
	// 超时，目标分片接收到TXaux2超时，直接视为失败
	if time.Since(data.Msg.StartTime) >= data.Msg.TimeoutDuration {
		crom.pbftNode.pl.Plog.Printf("S%dN%d : account transfer time out\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
		return
	}
	if !data.Msg.MPmig1 || !data.Msg.MPstate {
		return
	}
	crom.pbftNode.pl.Plog.Printf("S%dN%d : dest shard validate TXaux1: correct \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
	crom.pbftNode.pl.Plog.Printf("S%dN%d : account %s's location is updated to: %d \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID, data.Msg.Txmig1.Address, data.Msg.Txmig1.ToshardID)

	sii := message.TXANN_MSG{
		Msg: core.TXann{
			Txmig2:  data.Msg,
			MPmig2:  true,
			State:   data.Msg.State,
			MPstate: true,
		},
		Sender: crom.pbftNode.ShardID,
	}
	sByte, err := json.Marshal(sii)
	if err != nil {
		log.Panic()
	}
	msg_send := message.MergeMessage(message.TXann, sByte)

	crom.pbftNode.pl.Plog.Printf("S%dN%d : dest shard send TXann \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)

	// 广播 TXann 给除了主节点的所有节点
	for i := uint64(0); i < uint64(params.ShardNum); i++ {
		for j := uint64(1); j < uint64(params.NodesInShard); j++ {
			go networks.TcpDial(msg_send, crom.pbftNode.ip_nodeTable[uint64(i)][uint64(j)])
		}
	}
	crom.pbftNode.pl.Plog.Printf("S%dN%d : dest shard handle stage2 done \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
}

// 目标分片处理账户转移失败类型2中源分片发过来的请求
func (crom *SHARD_CUTTER) handleSourceQuery(content []byte) {
	data := new(message.CLU_SOURCE_QUERY)
	err := json.Unmarshal(content, data)
	if err != nil {
		log.Panic()
	}
	sii := message.CLU_DEST_REPLY{
		AccountKey:      data.AccountKey,
		AccountLocation: crom.pbftNode.ShardID, // crom.pbftNode.CurChain.Get_PartitionMap(data.AccountKey),
		Sender:          crom.pbftNode.ShardID,
	}
	sByte, err := json.Marshal(sii)
	if err != nil {
		log.Panic()
	}
	msg_send := message.MergeMessage(message.TXann, sByte)
	// 发送到源分片，即发过来的分片编号
	go networks.TcpDial(msg_send, crom.pbftNode.ip_nodeTable[data.Sender][0])
	crom.pbftNode.pl.Plog.Printf("S%dN%d : dest shard has handled Query from Source shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
}

// 源分片处理账户转移失败类型2中目标分片返回的请求
func (crom *SHARD_CUTTER) handleDestReply(content []byte) {
	data := new(message.CLU_DEST_REPLY)
	err := json.Unmarshal(content, data)
	if err != nil {
		log.Panic()
	}
	crom.sq.mu.Lock()
	crom.sq.receivedData = true
	crom.sq.AccountKey = data.AccountKey
	crom.sq.AccountLocation = data.AccountLocation
	crom.sq.mu.Unlock()
	crom.pbftNode.pl.Plog.Printf("S%dN%d : source has received reply from Dest shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
}

// stage3：源/其他分片接收到消息TXann，并更新账户信息。随后将TXns发送给目标分片
func (crom *SHARD_CUTTER) handleTXann(content []byte) {

	data := new(message.TXANN_MSG)
	err := json.Unmarshal(content, data)
	if err != nil {
		log.Panic()
	}

	// 如果不是源分片的主节点，则更新账户
	if !(crom.pbftNode.ShardID == data.Msg.Txmig2.Txmig1.FromshardID && crom.pbftNode.NodeID == 0) {
		crom.pbftNode.pl.Plog.Printf(
			"S%dN%d : update account %s's location to: %d \n",
			crom.pbftNode.ShardID,
			crom.pbftNode.NodeID,
			data.Msg.State.Key,
			data.Msg.State.Location,
		)
		return
	}
	// 如果此时已经超时
	if time.Since(data.Msg.Txmig2.StartTime) > data.Msg.Txmig2.TimeoutDuration {
		crom.pbftNode.pl.Plog.Printf("S%dN%d : account transfer time out, query for dest shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)

		// 在账户转移失败的情况讨论中，源分片需要询问目标分片目前的账户状态，如果不是失败那么再进行转移
		send_msg_data := message.CLU_SOURCE_QUERY{
			AccountKey: data.Msg.State.Key,
			Sender:     crom.pbftNode.ShardID,
		}
		send_bytes, err := json.Marshal(send_msg_data)
		if err != nil {
			log.Panic()
		}
		send_msg_struct := message.MergeMessage(message.ScourceQuery, send_bytes)
		// 发送到目标分片，即发过来的分片编号
		go networks.TcpDial(send_msg_struct, crom.pbftNode.ip_nodeTable[data.Sender][0])

		// 设置receiveddata为false
		crom.sq.receivedData = false

		// 阻塞程序，直到查询结果成功
		for {
			crom.sq.mu.Lock()
			// 如果获得了结果
			if crom.sq.receivedData {
				crom.pbftNode.pl.Plog.Printf("S%dN%d : has received data from dest shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
				// 如果目标分片中，该账户状态为目标分片（检测成功），则继续
				if crom.sq.AccountLocation == data.Msg.State.Location || true {
					crom.pbftNode.pl.Plog.Printf("S%dN%d : Obtain verification from the source shard\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
					crom.sq.mu.Unlock()
					break
				} else { // 否则结束，认为账户转移失败
					crom.sq.mu.Unlock()
					return
				}
			}
			crom.sq.mu.Unlock()
			// 等待500毫秒
			time.Sleep(500 * time.Millisecond)
		}
	}
	sii := message.TXNS_MSG{
		Msg: core.TXns{
			Txann:   data.Msg,
			MPann:   true,
			State:   data.Msg.State,
			Address: data.Msg.State.Key,
		},
		Sender: crom.pbftNode.ShardID,
	}
	sByte, err := json.Marshal(sii)
	if err != nil {
		log.Panic()
	}
	msg_send := message.MergeMessage(message.TXns, sByte)

	// 将TXns发送给目标分片
	crom.pbftNode.pl.Plog.Printf("S%dN%d : shource shard send TXns to dest shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
	go networks.TcpDial(msg_send, crom.pbftNode.ip_nodeTable[data.Msg.State.Location][0])
	crom.pbftNode.pl.Plog.Printf("S%dN%d : shource shard handle stage 3 done \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
}

// stage 4：节点接收到TXns后更新账户状态
func (crom *SHARD_CUTTER) handleTXns(content []byte) {
	crom.pbftNode.pl.Plog.Printf("S%dN%d : dest shard received TXns from source shard \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
	data := new(message.TXNS_MSG)
	err := json.Unmarshal(content, data)
	if err != nil {
		log.Panic()
	}
	crom.pbftNode.pl.Plog.Printf("S%dN%d : dest shard handle stage 4 done \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
}

func (crom *SHARD_CUTTER) HandleMessageOutsidePBFT(msgType message.MessageType, content []byte) bool {
	switch msgType {
	// CLPA
	case message.CPartitionMsg:
		crom.handlePartitionMsg(content)
	case message.CRelay:
		crom.handleRelay(content)
	case message.CRelayWithProof:
		crom.handleRelayWithProof(content)
	case message.CInject:
		crom.handleInjectTx(content)
	case message.AccountState_and_TX:
		crom.handleAccountStateAndTxMsg(content)
	case message.CPartitionReady:
		crom.handlePartitionReady(content)

	// SHARD_CUTTER
	case message.TXaux_2:
		crom.handleTXaux_2(content)
	case message.TXann:
		crom.handleTXann(content)
	case message.TXns:
		crom.handleTXns(content)
	case message.ScourceQuery:
		crom.handleSourceQuery(content)
	case message.DestReply:
		crom.handleDestReply(content)

	default:
		log.Panic()
	}
	return true
}
