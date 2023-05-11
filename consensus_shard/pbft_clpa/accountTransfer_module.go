// account transfer happens when the leader received the re-partition message.
// leaders send the infos about the accounts to be transferred to other leaders, and
// handle them.

package pbft_clpa

import (
	"blockEmulator/core"
	"blockEmulator/message"
	"blockEmulator/networks"
	"encoding/json"
	"log"
	"time"
)

// this message used in propose stage, so it will be invoked by InsidePBFT_Module
func (cphm *CLPAPbftInsideExtraHandleMod) sendPartitionReady() {
	cphm.cdm.pReadyLock.Lock()
	cphm.cdm.partitionReady[cphm.pbftNode.ShardID] = true
	cphm.cdm.pReadyLock.Unlock()

	pr := message.PartitionReady{
		FromShard: cphm.pbftNode.ShardID,
		NowSeqID:  cphm.pbftNode.sequenceID,
	}
	pByte, err := json.Marshal(pr)
	if err != nil {
		log.Panic()
	}
	send_msg := message.MergeMessage(message.CPartitionReady, pByte)
	for sid := 0; sid < int(cphm.pbftNode.pbftChainConfig.ShardNums); sid++ {
		if sid != int(pr.FromShard) {
			networks.TcpDial(send_msg, cphm.pbftNode.ip_nodeTable[uint64(sid)][0])
		}
	}
	cphm.pbftNode.pl.plog.Print("Ready for partition\n")
}

// get whether all shards is ready, it will be invoked by InsidePBFT_Module
func (cphm *CLPAPbftInsideExtraHandleMod) getPartitionReady() bool {
	cphm.cdm.pReadyLock.Lock()
	defer cphm.cdm.pReadyLock.Unlock()
	cphm.pbftNode.seqMapLock.Lock()
	defer cphm.pbftNode.seqMapLock.Unlock()
	cphm.cdm.readySeqLock.Lock()
	defer cphm.cdm.readySeqLock.Unlock()

	flag := true
	for sid, val := range cphm.pbftNode.seqIDMap {
		if rval, ok := cphm.cdm.readySeq[sid]; !ok || (rval-1 != val) {
			flag = false
		}
	}
	return len(cphm.cdm.partitionReady) == int(cphm.pbftNode.pbftChainConfig.ShardNums) && flag
}

// send the transactions and the accountState to other leaders
func (cphm *CLPAPbftInsideExtraHandleMod) sendAccounts_and_Txs() {
	// generate accout transfer and txs message
	accountToFetch := make([]string, 0)
	lastMapid := len(cphm.cdm.modifiedMap) - 1
	for key, val := range cphm.cdm.modifiedMap[lastMapid] {
		if val != cphm.pbftNode.ShardID && cphm.pbftNode.CurChain.Get_PartitionMap(key) == cphm.pbftNode.ShardID {
			accountToFetch = append(accountToFetch, key)
		}
	}
	asFetched := cphm.pbftNode.CurChain.FetchAccounts(accountToFetch)
	// send the accounts to other shards
	cphm.pbftNode.CurChain.Txpool.GetLocked()
	for i := uint64(0); i < cphm.pbftNode.pbftChainConfig.ShardNums; i++ {
		if i == cphm.pbftNode.ShardID {
			continue
		}
		addrSend := make([]string, 0)
		addrSet := make(map[string]bool)
		asSend := make([]*core.AccountState, 0)
		for idx, addr := range accountToFetch {
			if cphm.cdm.modifiedMap[lastMapid][addr] == i {
				addrSend = append(addrSend, addr)
				addrSet[addr] = true
				asSend = append(asSend, asFetched[idx])
			}
		}
		// fetch transactions to it, after the transactions is fetched, delete it in the pool
		txSend := make([]*core.Transaction, 0)
		head := 0
		tail := len(cphm.pbftNode.CurChain.Txpool.TxQueue)
		for head < tail {
			ptx := cphm.pbftNode.CurChain.Txpool.TxQueue[head]
			// if this is a normal transaction or ctx1 before re-sharding && the addr is correspond
			_, ok := addrSet[ptx.Sender]
			condition1 := ok && !ptx.Relayed
			// if this tx is ctx2
			_, ok = addrSet[ptx.Recipient]
			condition2 := ok && ptx.Relayed
			if condition1 || condition2 {
				txSend = append(txSend, ptx)
				tail--
				cphm.pbftNode.CurChain.Txpool.TxQueue[head] = cphm.pbftNode.CurChain.Txpool.TxQueue[tail]
			} else {
				head++
			}
			cphm.pbftNode.CurChain.Txpool.TxQueue = cphm.pbftNode.CurChain.Txpool.TxQueue[:tail]
		}

		cphm.pbftNode.pl.plog.Printf("The txSend to shard %d is generated \n", i)
		ast := message.AccountStateAndTx{
			Addrs:        addrSend,
			AccountState: asSend,
			FromShard:    cphm.pbftNode.ShardID,
			Txs:          txSend,
		}
		aByte, err := json.Marshal(ast)
		if err != nil {
			log.Panic()
		}
		send_msg := message.MergeMessage(message.AccountState_and_TX, aByte)
		networks.TcpDial(send_msg, cphm.pbftNode.ip_nodeTable[i][0])
		cphm.pbftNode.pl.plog.Printf("The message to shard %d is sent\n", i)
	}
	cphm.pbftNode.CurChain.Txpool.GetUnlocked()
}

// fetch collect infos
func (cphm *CLPAPbftInsideExtraHandleMod) getCollectOver() bool {
	cphm.cdm.collectLock.Lock()
	defer cphm.cdm.collectLock.Unlock()
	return cphm.cdm.collectOver
}

// propose a partition message
func (cphm *CLPAPbftInsideExtraHandleMod) proposePartition() (bool, *message.Request) {
	cphm.pbftNode.pl.plog.Printf("S%dN%d : begin partition proposing\n", cphm.pbftNode.ShardID, cphm.pbftNode.NodeID)
	// add all data in pool into the set
	for _, at := range cphm.cdm.accountStateTx {
		for i, addr := range at.Addrs {
			cphm.cdm.receivedNewAccountState[addr] = at.AccountState[i]
		}
		cphm.cdm.receivedNewTx = append(cphm.cdm.receivedNewTx, at.Txs...)
	}
	// propose, send all txs to other nodes in shard
	cphm.pbftNode.CurChain.Txpool.AddTxs2Pool(cphm.cdm.receivedNewTx)

	atmaddr := make([]string, 0)
	atmAs := make([]*core.AccountState, 0)
	for key, val := range cphm.cdm.receivedNewAccountState {
		atmaddr = append(atmaddr, key)
		atmAs = append(atmAs, val)
	}
	atm := message.AccountTransferMsg{
		ModifiedMap:  cphm.cdm.modifiedMap[cphm.cdm.accountTransferRound],
		Addrs:        atmaddr,
		AccountState: atmAs,
		ATid:         uint64(len(cphm.cdm.modifiedMap)),
	}
	atmbyte := atm.Encode()
	r := &message.Request{
		RequestType: message.PartitionReq,
		Msg: message.RawMessage{
			Content: atmbyte,
		},
		ReqTime: time.Now(),
	}
	return true, r
}

// all nodes in a shard will do accout Transfer, to sync the state trie
func (cphm *CLPAPbftInsideExtraHandleMod) accountTransfer_do(atm *message.AccountTransferMsg) {
	// change the partition Map
	cnt := 0
	for key, val := range atm.ModifiedMap {
		cnt++
		cphm.pbftNode.CurChain.Update_PartitionMap(key, val)
	}
	cphm.pbftNode.pl.plog.Printf("%d key-vals are updated\n", cnt)
	// add the account into the state trie
	cphm.pbftNode.CurChain.AddAccounts(atm.Addrs, atm.AccountState)

	if uint64(len(cphm.cdm.modifiedMap)) != atm.ATid {
		cphm.cdm.modifiedMap = append(cphm.cdm.modifiedMap, atm.ModifiedMap)
	}
	cphm.cdm.accountTransferRound = atm.ATid
	cphm.cdm.accountStateTx = make(map[uint64]*message.AccountStateAndTx)
	cphm.cdm.receivedNewAccountState = make(map[string]*core.AccountState)
	cphm.cdm.receivedNewTx = make([]*core.Transaction, 0)
	cphm.cdm.partitionOn = false

	cphm.cdm.collectLock.Lock()
	cphm.cdm.collectOver = false
	cphm.cdm.collectLock.Unlock()

	cphm.cdm.pReadyLock.Lock()
	cphm.cdm.partitionReady = make(map[uint64]bool)
	cphm.cdm.pReadyLock.Unlock()

	cphm.pbftNode.CurChain.PrintBlockChain()
}
