package pbft_clpa

import (
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/networks"
	"blockEmulator/params"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// handle
type PbftInsideExtraHandleMod interface {
	HandleinPropose() (bool, *message.Request)
	HandleinPrePrepare(*message.PrePrepare) bool
	HandleinPrepare(*message.Prepare) bool
	HandleinCommit(*message.Commit) bool
	HandleReqestforOldSeq(*message.RequestOldMessage) bool
	HandleforSequentialRequest(*message.SendOldMessage) bool
}

type CLPAPbftInsideExtraHandleMod struct {
	cdm      *Data_supportCLPA
	pbftNode *PbftConsensusNode
}

// propose request with different types
func (cphm *CLPAPbftInsideExtraHandleMod) HandleinPropose() (bool, *message.Request) {
	if cphm.cdm.partitionOn {
		cphm.sendPartitionReady()
		for !cphm.getPartitionReady() {
			time.Sleep(time.Second)
		}
		// send accounts and txs
		cphm.sendAccounts_and_Txs()
		// propose a partition
		for !cphm.getCollectOver() {
			time.Sleep(time.Second)
		}
		return cphm.proposePartition()
	}

	// ELSE: propose a block
	block := cphm.pbftNode.CurChain.GenerateBlock()
	r := &message.Request{
		RequestType: message.BlockRequest,
		ReqTime:     time.Now(),
	}
	r.Msg.Content = block.Encode()
	return true, r

}

// the diy operation in preprepare
func (cphm *CLPAPbftInsideExtraHandleMod) HandleinPrePrepare(ppmsg *message.PrePrepare) bool {
	// judge whether it is a partitionRequest or not
	isPartitionReq := ppmsg.RequestMsg.RequestType == message.PartitionReq

	if isPartitionReq {
		// after some checking
		cphm.pbftNode.pl.plog.Printf("S%dN%d : a partition block\n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID)
	} else {
		// the request is a block
		if cphm.pbftNode.CurChain.IsValidBlock(core.DecodeB(ppmsg.RequestMsg.Msg.Content)) != nil {
			cphm.pbftNode.pl.plog.Printf("S%dN%d : not a valid block\n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID)
			return false
		}
	}
	cphm.pbftNode.pl.plog.Printf("S%dN%d : the pre-prepare message is correct, putting it into the RequestPool. \n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID)
	cphm.pbftNode.requestPool[string(ppmsg.Digest)] = ppmsg.RequestMsg
	// merge to be a prepare message
	return true
}

// the operation in prepare, and in pbft + tx relaying, this function does not need to do any.
func (cphm *CLPAPbftInsideExtraHandleMod) HandleinPrepare(pmsg *message.Prepare) bool {
	fmt.Println("No operations are performed in Extra handle mod")
	return true
}

// the operation in commit.
func (cphm *CLPAPbftInsideExtraHandleMod) HandleinCommit(cmsg *message.Commit) bool {
	r := cphm.pbftNode.requestPool[string(cmsg.Digest)]
	// requestType ...
	if r.RequestType == message.PartitionReq {
		// if a partition Requst ...
		atm := message.DecodeAccountTransferMsg(r.Msg.Content)
		cphm.accountTransfer_do(atm)
		return true
	}
	// if a block request ...
	block := core.DecodeB(r.Msg.Content)
	cphm.pbftNode.pl.plog.Printf("S%dN%d : adding the block %d...now height = %d \n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID, block.Header.Number, cphm.pbftNode.CurChain.CurrentBlock.Header.Number)
	cphm.pbftNode.CurChain.AddBlock(block)
	cphm.pbftNode.pl.plog.Printf("S%dN%d : added the block %d... \n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID, block.Header.Number)
	cphm.pbftNode.CurChain.PrintBlockChain()

	// now try to relay txs to other shards (for main nodes)
	if cphm.pbftNode.NodeID == cphm.pbftNode.view {
		cphm.pbftNode.pl.plog.Printf("S%dN%d : main node is trying to send relay txs at height = %d \n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID, block.Header.Number)
		// generate relay pool and collect txs excuted
		txExcuted := make([]*core.Transaction, 0)
		cphm.pbftNode.CurChain.Txpool.RelayPool = make(map[uint64][]*core.Transaction)
		relay1Txs := make([]*core.Transaction, 0)
		for _, tx := range block.Body {
			rsid := cphm.pbftNode.CurChain.Get_PartitionMap(tx.Recipient)
			if rsid != cphm.pbftNode.ShardID {
				ntx := tx
				ntx.Relayed = true
				cphm.pbftNode.CurChain.Txpool.AddRelayTx(ntx, rsid)
				relay1Txs = append(relay1Txs, tx)
			} else {
				txExcuted = append(txExcuted, tx)
			}
		}
		// send relay txs
		for sid := uint64(0); sid < cphm.pbftNode.pbftChainConfig.ShardNums; sid++ {
			if sid == cphm.pbftNode.ShardID {
				continue
			}
			relay := message.Relay{
				Txs:           cphm.pbftNode.CurChain.Txpool.RelayPool[sid],
				SenderShardID: cphm.pbftNode.ShardID,
				SenderSeq:     cphm.pbftNode.sequenceID,
			}
			rByte, err := json.Marshal(relay)
			if err != nil {
				log.Panic()
			}
			msg_send := message.MergeMessage(message.CRelay, rByte)
			go networks.TcpDial(msg_send, cphm.pbftNode.ip_nodeTable[sid][0])
			cphm.pbftNode.pl.plog.Printf("S%dN%d : sended relay txs to %d\n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID, sid)
		}
		cphm.pbftNode.CurChain.Txpool.ClearRelayPool()
		// send txs excuted in this block to the listener
		// add more message to measure more metrics
		bim := message.BlockInfoMsg{
			BlockBodyLength: len(block.Body),
			ExcutedTxs:      txExcuted,
			Epoch:           int(cphm.cdm.accountTransferRound),
			Relay1Txs:       relay1Txs,
			Relay1TxNum:     uint64(len(relay1Txs)),
			SenderShardID:   cphm.pbftNode.ShardID,
			ProposeTime:     r.ReqTime,
			CommitTime:      time.Now(),
		}
		bByte, err := json.Marshal(bim)
		if err != nil {
			log.Panic()
		}
		msg_send := message.MergeMessage(message.CBlockInfo, bByte)
		go networks.TcpDial(msg_send, cphm.pbftNode.ip_nodeTable[params.DeciderShard][0])
		cphm.pbftNode.pl.plog.Printf("S%dN%d : sended excuted txs\n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID)
		cphm.pbftNode.writeCSVline([]string{strconv.Itoa(len(txExcuted)), strconv.Itoa(int(bim.Relay1TxNum))})
	}
	return true
}

func (cphm *CLPAPbftInsideExtraHandleMod) HandleReqestforOldSeq(*message.RequestOldMessage) bool {
	fmt.Println("No operations are performed in Extra handle mod")
	return true
}

// the operation for sequential requests
func (cphm *CLPAPbftInsideExtraHandleMod) HandleforSequentialRequest(som *message.SendOldMessage) bool {
	if int(som.SeqStartHeight-som.SeqEndHeight) != len(som.OldRequest) {
		cphm.pbftNode.pl.plog.Printf("S%dN%d : the SendOldMessage message is not enough\n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID)
	} else { // add the block into the node pbft blockchain
		for height := som.SeqStartHeight; height <= som.SeqEndHeight; height++ {
			r := som.OldRequest[height-som.SeqStartHeight]
			if r.RequestType == message.BlockRequest {
				b := core.DecodeB(r.Msg.Content)
				cphm.pbftNode.CurChain.AddBlock(b)
			}
		}
		cphm.pbftNode.sequenceID = som.SeqEndHeight + 1
		cphm.pbftNode.CurChain.PrintBlockChain()
	}
	return true
}
