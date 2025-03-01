package pbft_all

import (
	"blockEmulator/chain"
	"blockEmulator/consensus_shard/pbft_all/dataSupport"
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/networks"
	"encoding/json"
	"log"
	"time"
)

type SHARD_CLUSTER struct {
	cdm      *dataSupport.Data_supportCLPA
	pbftNode *PbftConsensusNode
}

// receive relay transaction, which is for cross shard txs
func (crom *SHARD_CLUSTER) handleRelay(content []byte) {
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

func (crom *SHARD_CLUSTER) handleRelayWithProof(content []byte) {
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

func (crom *SHARD_CLUSTER) handleInjectTx(content []byte) {
	it := new(message.InjectTxs)
	err := json.Unmarshal(content, it)
	if err != nil {
		log.Panic(err)
	}
	crom.pbftNode.CurChain.Txpool.AddTxs2Pool(it.Txs)
	crom.pbftNode.pl.Plog.Printf("S%dN%d : has handled injected txs msg, txs: %d \n", crom.pbftNode.ShardID, crom.pbftNode.NodeID, len(it.Txs))
}

// func sendMsg(val int) {
// 	sii := message.TXAUX_1_MSG{
// 		Msg: core.TXmig1{
// 			Address: "",
// 		},
// 	}
// 	sByte, err := json.Marshal(sii)
// 	if err != nil {
// 		log.Panic()
// 	}
// 	msg_send := message.MergeMessage(message.TXaux_1, sByte)
// 	go networks.TcpDial(msg_send, crom.pbftNode.ip_nodeTable[sid][0])
// }

func (crom *SHARD_CLUSTER) handlePartitionMsg(content []byte) {
	pm := new(message.PartitionModifiedMap)
	err := json.Unmarshal(content, pm)
	if err != nil {
		log.Panic()
	}

	// PartitionModified变量 key: 账户的编号, val: 分片编号
	for key, val := range pm.PartitionModified {
		if crom.pbftNode.CurChain.Get_PartitionMap(key) == crom.pbftNode.ShardID {
			// 发送TXaux2（发送至编号为val的分片）
			txau1 := core.TXmig1{
				Address:     key,
				FromshardID: crom.pbftNode.ShardID,
				ToshardID:   val,
			}
			sii := message.TXAUX_2_MSG{
				Msg: core.TXmig2{
					Txmig1: txau1,
					MPmig1: true,
					State: core.CLUSTER_ACCOUNT_STATE{
						Key:      key,
						Location: val,
					},
					MPstate:         true,
					TimeoutDuration: 10 * time.Second,
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

			// 更新账户在源分片中的分片编号为目标分片 (未实现)
		}
	}

	crom.cdm.ModifiedMap = append(crom.cdm.ModifiedMap, pm.PartitionModified)
	crom.pbftNode.pl.Plog.Printf("S%dN%d : has received partition message\n", crom.pbftNode.ShardID, crom.pbftNode.NodeID)
	// crom.cdm.PartitionOn = true
}

func (crom *SHARD_CLUSTER) handleTXaux_2(content []byte) {
	data := new(message.TXAUX_2_MSG)
	err := json.Unmarshal(content, data)
	if err != nil {
		log.Panic()
	}
	if time.Since(data.Msg.StartTime) >= data.Msg.TimeoutDuration {
		return
	}
	if !data.Msg.MPmig1 || !data.Msg.MPstate {
		return
	}
	// accout_key := data.Msg.Txmig1.Address
	// dest_shard := data.Msg.Txmig1.ToshardID
	// 更新账户的分片状态 （未完成）
	// 发送TXann到源分片
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
	msg_send := message.MergeMessage(message.TXaux_2, sByte)
	// 发送到源分片，即发过来的分片编号
	go networks.TcpDial(msg_send, crom.pbftNode.ip_nodeTable[data.Sender][0])
}

func (crom *SHARD_CLUSTER) handleTXann(content []byte) {
	data := new(message.TXANN_MSG)
	err := json.Unmarshal(content, data)
	if err != nil {
		log.Panic()
	}
}

func (crom *SHARD_CLUSTER) handleTXns(content []byte) {

}

func (crom *SHARD_CLUSTER) HandleMessageOutsidePBFT(msgType message.MessageType, content []byte) bool {
	switch msgType {
	case message.CPartitionMsg:
		crom.handlePartitionMsg(content)
	case message.CRelay:
		crom.handleRelay(content)
	case message.CRelayWithProof:
		crom.handleRelayWithProof(content)
	case message.CInject:
		crom.handleInjectTx(content)

	case message.TXaux_2:
		crom.handleTXaux_2(content)
	case message.TXann:
		crom.handleTXann(content)
	case message.TXns:
		crom.handleTXns(content)

	default:
	}
	return true
}
