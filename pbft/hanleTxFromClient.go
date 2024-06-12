package pbft

import (
	"blockEmulator/account"
	"blockEmulator/core"
	"blockEmulator/params"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"
)

func (p *Pbft) handleTxFromClient(content []byte) {
	txsFromClient := new(TxFromClient)
	err := json.Unmarshal(content, txsFromClient)
	if err != nil {
		log.Panic(err)
	}
	tx2 := make([]*core.Transaction, 0)
	p.Node.CurChain.Tx_pool.Lock.Lock()
	account.Account2ShardLock.Lock()
	j := 0
	self_shardID := params.ShardTable[params.Config.ShardID]
	if !params.Config.Not_Lock_immediately {
		for _, tx := range txsFromClient.Txs {
			senderSID := account.Account2Shard[hex.EncodeToString(tx.Sender)]
			if senderSID == self_shardID {
				from, to := hex.EncodeToString(tx.Sender), hex.EncodeToString(tx.Recipient)
				//全锁
				if !params.Config.Stop_When_Migrating && params.Config.Lock_Acc_When_Migrating {
					account.Lock_Acc_Lock.Lock()
					if account.Lock_Acc[from] {
						tx.LockTime = time.Now().UnixMilli()
						tx.SenLock = true
						p.Node.CurChain.Tx_pool.Locking_TX_Pools[from] = append(p.Node.CurChain.Tx_pool.Locking_TX_Pools[from], tx)
						account.Lock_Acc_Lock.Unlock()
						continue
					}
					if account.Lock_Acc[to] {
						tx.LockTime = time.Now().UnixMilli()
						tx.RecLock = true
						p.Node.CurChain.Tx_pool.Locking_TX_Pools[to] = append(p.Node.CurChain.Tx_pool.Locking_TX_Pools[to], tx)
						account.Lock_Acc_Lock.Unlock()
						continue
					}
					account.Lock_Acc_Lock.Unlock()
				}
				//半锁
				if !params.Config.Stop_When_Migrating && !params.Config.Lock_Acc_When_Migrating {
					account.Outing_Acc_Before_Announce_Lock.Lock()
					if account.Outing_Acc_Before_Announce[from] {
						tx.LockTime = time.Now().UnixMilli()
						p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools[from] = append(p.Node.CurChain.Tx_pool.Outing_Before_Announce_TX_Pools[from], tx)
						account.Outing_Acc_Before_Announce_Lock.Unlock()
						continue
					}
					account.Outing_Acc_Before_Announce_Lock.Unlock()
				}
				txsFromClient.Txs[j] = tx
				j++
			} else {
				tx2 = append(tx2, tx)
			}
		}
	}

	if params.Config.Not_Lock_immediately {
		for _, tx := range txsFromClient.Txs {
			senderSID := account.Account2Shard[hex.EncodeToString(tx.Sender)]
			if senderSID == self_shardID {
				txsFromClient.Txs[j] = tx
				j++
			} else {
				tx2 = append(tx2, tx)
			}
		}
	}
	
	txsFromClient.Txs = txsFromClient.Txs[:j]
	p.Node.CurChain.Tx_pool.Queue = append(p.Node.CurChain.Tx_pool.Queue, txsFromClient.Txs...)

	if len(tx2) != 0 {
		p.TrySendTX(tx2)
	}
	account.Account2ShardLock.Unlock()
	p.Node.CurChain.Tx_pool.Lock.Unlock()
}
