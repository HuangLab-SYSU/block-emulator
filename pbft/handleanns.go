package pbft

import (
	"blockEmulator/account"
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/trie"
)

func (p *Pbft) handleAnns(anns []*core.TXann, st *trie.Trie, migTree *trie.Trie) (int64, int64) {
	if len(anns) == 0 {
		return 0,0
	}
	fmt.Printf("handleAnns时是第%v个块\n", p.epoch-1)
	// unchanged := 0
	p.Node.CurChain.Tx_pool.Relaypoollock.Lock()
	p.Node.CurChain.Tx_pool.Lock.Lock()
	account.Account2ShardLock.Lock()
	begintime := time.Now().UnixMicro()

	// nowqueue := []*core.Transaction{}
	// core.TxPoolDeepCopy(&nowqueue, p.Node.CurChain.Tx_pool.Queue)
	cAps := make(map[string]map[string]*ChangeAndPending) // map[ToShard][Addr]*ChangeAndPending
	// nss := make(map[string][]*core.TXns)                  // map[ToShard][Addr]*ChangeAndPending
	isdelete := make(map[string]string)                   // map[Addr]ToShard

	//   1.从本分片账户列表删除，
	////   2.收集关于该账户的交易，
	//   3.计算账户余额的变化，
	////   4.将该账户放入内存中的 删除账户队列中， 这个队列将在下一个区块里面
	//   5.将账户的余额变化与pending交易放入map[string]struct中，待会一起发送给目标分片
	for _, ann := range anns {
		// 原本就不在这个分片，改table就好
		if !account.AccountInOwnShard[ann.Address] {
			account.Account2Shard[ann.Address] = ann.ToshardID
			continue
		}
		// hex_address, _ := hex.DecodeString(ann.Address)

		// encoded := st.Get(hex_address)
		// if encoded == nil {
		// 	log.Panic()
		// }
		// state := account.DecodeAccountState(encoded)

		// //   3.计算账户余额的变化
		// change := big.NewInt(0)
		// if !params.Config.Lock_Acc_When_Migrating {
		// 	account.BalanceBeforeOutLock.Lock()
		// 	change = new(big.Int).Sub(state.Balance, account.BalanceBeforeOut[ann.Address])
		// 	//  可以从内存删除该账户在迁移1时的余额
		// 	delete(account.BalanceBeforeOut, ann.Address)
		// 	account.BalanceBeforeOutLock.Unlock()
		// }


		// if change.BitLen() == 0 {
		// 	unchanged++
		// }else {
		// 	mpbegin := time.Now().UnixMicro()

		// 	proofDB1 := &core.ProofDB{}
		// 	err1 := migTree.Prove(ann.Hash(), 0, proofDB1)
		// 	if err1 != nil {
		// 		log.Panic(err1)
		// 	}

		// 	proofDB2 := &core.ProofDB{}
		// 	err2 := st.Prove(hex_address, 0, proofDB2)
		// 	if err2 != nil {
		// 		log.Panic(err2)
		// 	}

		// 	mpend := time.Now().UnixMicro()


			
		// 	ns := &core.TXns{Txann: ann, MPann: proofDB1, State: state, MPstate: proofDB2, Address: ann.Address}
			
		// 	ns.Change = new(big.Int).Set(change)

		// 	// 5. 将账户的余额变化放入列表，待会发送
		// 	if _, ok := nss[params.ShardTableInt2Str[ann.ToshardID]]; !ok {
		// 		nss[params.ShardTableInt2Str[ann.ToshardID]] = make([]*core.TXns, 0)
		// 	}
		// 	nss[params.ShardTableInt2Str[ann.ToshardID]] = append(nss[params.ShardTableInt2Str[ann.ToshardID]], ns)
		// 	mptime += mpend - mpbegin
		// }

		

		isdelete[ann.Address] = params.ShardTableInt2Str[ann.ToshardID]
		cAp := &ChangeAndPending{}
		cAp.PendingTxs = make([]*core.Transaction, 0)

		//   1.从本分片账户列表删除，
		account.Account2Shard[ann.Address] = ann.ToshardID
		delete(account.AccountInOwnShard, ann.Address)

		if params.Config.Lock_Acc_When_Migrating { //如果是锁机制则将其解锁
			account.Lock_Acc_Lock.Lock()
			delete(account.Lock_Acc, ann.Address)
			account.Lock_Acc_Lock.Unlock()
		} else if !params.Config.Stop_When_Migrating { //如果是不锁机制则将其半锁解除
			account.Outing_Acc_Before_Announce_Lock.Lock()
			delete(account.Outing_Acc_Before_Announce, ann.Address)
			account.Outing_Acc_Before_Announce_Lock.Unlock()
		}

		// 2.如果是lock机制，将关于账户的被锁交易中，sender为本分片的交易放回交易池中
		if params.Config.Lock_Acc_When_Migrating {
			j := 0
			account.Lock_Acc_Lock.Lock()
			for _, tx := range p.Node.CurChain.Tx_pool.Locking_TX_Pools[ann.Address] {
				if tx.UnlockTime > 0 {
					tx.UnlockTime2 = time.Now().UnixMilli()
				} else {
					tx.UnlockTime = time.Now().UnixMilli()
				}
				// 如果不是(relay或relayLock)且v为接收者，则放回交易池
				if hex.EncodeToString(tx.Recipient) == ann.Address && !tx.IsRelay && !tx.Relay_Lock{
					if _,ok := isdelete[hex.EncodeToString(tx.Sender)]; ok || account.Lock_Acc[hex.EncodeToString(tx.Sender)] {
						tx.SenLock = true
						p.Node.CurChain.Tx_pool.Locking_TX_Pools[hex.EncodeToString(tx.Sender)] = append(p.Node.CurChain.Tx_pool.Locking_TX_Pools[hex.EncodeToString(tx.Sender)], tx)
						continue
					}
					// if tx.UnlockTime > 0 {
					// 	tx.UnlockTime2 = time.Now().UnixMilli()
					// } else {
					// 	tx.UnlockTime = time.Now().UnixMilli()
					// }
					tx.Success = true
					p.Node.CurChain.Tx_pool.Queue = append(p.Node.CurChain.Tx_pool.Queue, tx)
					continue
				}
				// 否则，留在被锁交易中
				p.Node.CurChain.Tx_pool.Locking_TX_Pools[ann.Address][j] = tx
				j++
			}
			p.Node.CurChain.Tx_pool.Locking_TX_Pools[ann.Address] = p.Node.CurChain.Tx_pool.Locking_TX_Pools[ann.Address][:j]
			account.Lock_Acc_Lock.Unlock()
		}


		// //   4.将该账户放入内存中的 删除账户队列中， 这个队列将在下一个区块里面
		// p.Node.CurChain.Delete_pool.AddDelete(del.Address, params.ShardTable[announce.ShardID])

		//   5.将pending交易放入列表中，待会一起发送给目标分片
		if _, ok := cAps[params.ShardTableInt2Str[ann.ToshardID]]; !ok {
			cAps[params.ShardTableInt2Str[ann.ToshardID]] = make(map[string]*ChangeAndPending)
		}
		cAps[params.ShardTableInt2Str[ann.ToshardID]][ann.Address] = cAp
		
		// nss[params.ShardTableInt2Str[ann.ToshardID]] = append(nss[params.ShardTableInt2Str[ann.ToshardID]], ns)
	}

	// account.Account2ShardLock.Unlock()
	// p.Node.CurChain.Tx_pool.Lock.Unlock()
	// p.Node.CurChain.Tx_pool.Relaypoollock.Unlock()


	//   2.将交易池中关于账户的交易收集起来
	jiaoyichitime := time.Now().UnixMicro()
	if params.Config.Lock_Acc_When_Migrating {
		for _, caps := range cAps {
			for addr, cAp := range caps {
				account.Lock_Acc_Lock.Lock()
				cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Locking_TX_Pools[addr]...)
				delete(p.Node.CurChain.Tx_pool.Locking_TX_Pools, addr)
				account.Lock_Acc_Lock.Unlock()
			}
		}
	} else if !params.Config.Stop_When_Migrating {
		for _, caps := range cAps {
			for addr, cAp := range caps {
				account.Outing_Acc_Before_Announce_Lock.Lock()
				cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools[addr]...)
				delete(p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools, addr)
				account.Outing_Acc_Before_Announce_Lock.Unlock()
			}
		}
	}

	// for _, tx := range nowqueue {
	// 	// 如果sender是要出去的账户，且不是relaytx, 则收集
	// 	if shard, ok := isdelete[hex.EncodeToString(tx.Sender)]; ok && !tx.IsRelay && !tx.Relay_Lock {
	// 		cAps[shard][hex.EncodeToString(tx.Sender)].PendingTxs = append(cAps[shard][hex.EncodeToString(tx.Sender)].PendingTxs, tx)
	// 		continue
	// 	}

	// 	//  如果recipient是要出去的账户，且已经是relaytx，则收集
	// 	if shard, ok := isdelete[hex.EncodeToString(tx.Recipient)]; ok && (tx.IsRelay || tx.Relay_Lock) {
	// 		// 在半锁机制中，即使这些交易没有被锁，也通过这样做个标记，表明是被迁移过去的
	// 		tx.RecLock = true
	// 		cAps[shard][hex.EncodeToString(tx.Recipient)].PendingTxs = append(cAps[shard][hex.EncodeToString(tx.Recipient)].PendingTxs, tx)
	// 		continue
	// 	}
	// 	// 否则，啥也不做
	// }
	
	j := 0
	for _, tx := range p.Node.CurChain.Tx_pool.Queue {
		// 如果sender是要出去的账户，且不是relaytx, 则收集
		if shard, ok := isdelete[hex.EncodeToString(tx.Sender)]; ok && !tx.IsRelay && !tx.Relay_Lock {
			cAps[shard][hex.EncodeToString(tx.Sender)].PendingTxs = append(cAps[shard][hex.EncodeToString(tx.Sender)].PendingTxs, tx)
			continue
		}

		//  如果recipient是要出去的账户，且已经是relaytx，则收集
		if shard, ok := isdelete[hex.EncodeToString(tx.Recipient)]; ok && (tx.IsRelay || tx.Relay_Lock) {
			// 在半锁机制中，即使这些交易没有被锁，也通过这样做个标记，表明是被迁移过去的
			tx.RecLock = true
			cAps[shard][hex.EncodeToString(tx.Recipient)].PendingTxs = append(cAps[shard][hex.EncodeToString(tx.Recipient)].PendingTxs, tx)
			continue
		}
		// 否则，留在交易池中
		p.Node.CurChain.Tx_pool.Queue[j] = tx
		j++
	}
	p.Node.CurChain.Tx_pool.Queue = p.Node.CurChain.Tx_pool.Queue[:j]

	endtime := time.Now().UnixMicro()
	fmt.Printf("\nhandleannounce交易池操作时差为：%v\n", endtime-jiaoyichitime)
	// fmt.Printf("handleannounce证明操作时差为：%v\n", mptime)
	fmt.Printf("handleannounce两锁时差为：%v\n\n", endtime-begintime)

	account.Account2ShardLock.Unlock()
	p.Node.CurChain.Tx_pool.Lock.Unlock()
	p.Node.CurChain.Tx_pool.Relaypoollock.Unlock()

	for _, caps := range cAps {
		for addr, cAp := range caps {
			account.Outing_Acc_After_Announce_Lock.Lock()
			cAp.PendingTxs = append(cAp.PendingTxs, p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools[addr]...)
			delete(account.Outing_Acc_After_Announce, addr)
			delete(p.Node.CurChain.Tx_pool.Outing_After_Announce_TX_Pools, addr)
			account.Outing_Acc_After_Announce_Lock.Unlock()
		}
	}

	// for _, ann := range anns {
	// 	// 原本就不在这个分片
	// 	if _,ok := isdelete[ann.Address];!ok {
	// 		continue
	// 	}

	// 	hex_address, _ := hex.DecodeString(ann.Address)

	// 	encoded := st.Get(hex_address)
	// 	if encoded == nil {
	// 		log.Panic()
	// 	}
	// 	state := account.DecodeAccountState(encoded)

	// 	//   3.计算账户余额的变化
	// 	change := big.NewInt(0)
	// 	if !params.Config.Lock_Acc_When_Migrating {
	// 		account.BalanceBeforeOutLock.Lock()
	// 		change = new(big.Int).Sub(state.Balance, account.BalanceBeforeOut[ann.Address])
	// 		//  可以从内存删除该账户在迁移1时的余额
	// 		delete(account.BalanceBeforeOut, ann.Address)
	// 		account.BalanceBeforeOutLock.Unlock()
	// 	}


	// 	if change.BitLen() == 0 {
	// 		unchanged++
	// 	}else {
	// 		// mpbegin := time.Now().UnixMicro()
	// 		proofDB1 := &core.ProofDB{}
	// 		err1 := migTree.Prove(ann.Hash(), 0, proofDB1)
	// 		if err1 != nil {
	// 			log.Panic(err1)
	// 		}

	// 		proofDB2 := &core.ProofDB{}
	// 		err2 := st.Prove(hex_address, 0, proofDB2)
	// 		if err2 != nil {
	// 			log.Panic(err2)
	// 		}

	// 		// mpend := time.Now().UnixMicro()
	// 		ns := &core.TXns{Txann: ann, MPann: proofDB1, State: state, MPstate: proofDB2, Address: ann.Address}
	// 		ns.Change = new(big.Int).Set(change)

	// 		// 5. 将账户的余额变化放入列表，待会发送
	// 		if _, ok := nss[params.ShardTableInt2Str[ann.ToshardID]]; !ok {
	// 			nss[params.ShardTableInt2Str[ann.ToshardID]] = make([]*core.TXns, 0)
	// 		}
	// 		nss[params.ShardTableInt2Str[ann.ToshardID]] = append(nss[params.ShardTableInt2Str[ann.ToshardID]], ns)
	// 		// mptime += mpend - mpbegin
	// 	}
	// }

	// for shard, caps := range cAps {
	// 	//将所有涉及账户的余额变化和pending交易发送给目标分片
	// 	p.TrySendChangesAndPendings(nss[shard], caps, shard)
	// }

	// // 将 交易池中的交易发给 client
	// c, err := json.Marshal(unchanged)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// m := jointMessage(cUnchangedState, c)
	// utils.TcpDial(m, params.ClientAddr)
	// fmt.Println("本节点已将 不变的账户数量 发送到client\n")
	go p.mpAndSend(isdelete, anns, cAps, st, migTree)

	return endtime-jiaoyichitime, endtime-begintime
}

func (p *Pbft) mpAndSend(isdelete map[string]string, anns []*core.TXann, cAps map[string]map[string]*ChangeAndPending, st, migTree *trie.Trie) {
	unchanged := 0
	nss := make(map[string][]*core.TXns)                  // map[ToShard][Addr]*ChangeAndPending
	
	for _, ann := range anns {
		// 原本就不在这个分片
		if _,ok := isdelete[ann.Address];!ok {
			continue
		}

		hex_address, _ := hex.DecodeString(ann.Address)

		encoded := st.Get(hex_address)
		if encoded == nil {
			log.Panic()
		}
		state := account.DecodeAccountState(encoded)

		//   3.计算账户余额的变化
		change := big.NewInt(0)
		if !params.Config.Lock_Acc_When_Migrating {
			account.BalanceBeforeOutLock.Lock()
			change = new(big.Int).Sub(state.Balance, account.BalanceBeforeOut[ann.Address])
			//  可以从内存删除该账户在迁移1时的余额
			delete(account.BalanceBeforeOut, ann.Address)
			account.BalanceBeforeOutLock.Unlock()
		}


		if change.BitLen() == 0 {
			unchanged++
		}else {
			// mpbegin := time.Now().UnixMicro()
			proofDB1 := &core.ProofDB{}
			err1 := migTree.Prove(ann.Hash(), 0, proofDB1)
			if err1 != nil {
				log.Panic(err1)
			}

			proofDB2 := &core.ProofDB{}
			err2 := st.Prove(hex_address, 0, proofDB2)
			if err2 != nil {
				log.Panic(err2)
			}

			// mpend := time.Now().UnixMicro()
			ns := &core.TXns{Txann: ann, MPann: proofDB1, State: state, MPstate: proofDB2, Address: ann.Address}
			ns.Change = new(big.Int).Set(change)

			// 5. 将账户的余额变化放入列表，待会发送
			if _, ok := nss[params.ShardTableInt2Str[ann.ToshardID]]; !ok {
				nss[params.ShardTableInt2Str[ann.ToshardID]] = make([]*core.TXns, 0)
			}
			nss[params.ShardTableInt2Str[ann.ToshardID]] = append(nss[params.ShardTableInt2Str[ann.ToshardID]], ns)
			// mptime += mpend - mpbegin
		}
	}

	for shard, caps := range cAps {
		//将所有涉及账户的余额变化和pending交易发送给目标分片
		p.TrySendChangesAndPendings(nss[shard], caps, shard)
	}

	// 将 交易池中的交易发给 client
	c, err := json.Marshal(unchanged)
	if err != nil {
		log.Panic(err)
	}
	m := jointMessage(cUnchangedState, c)
	utils.TcpDial(m, params.ClientAddr)
	fmt.Println("本节点已将 不变的账户数量 发送到client\n")
}
