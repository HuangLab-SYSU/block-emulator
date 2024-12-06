package core

import (
	"blockEmulator/params"
	"blockEmulator/utils"
	"fmt"
	"sync"
)

type Tx_pool struct {
	// 交易队列
	Queue []*Transaction
	// relay到不同分片的交易队列
	Relay_Pools map[string][]*Transaction
	//锁
	lock sync.Mutex
}

func NewTxPool() *Tx_pool {
	return &Tx_pool{
		Queue:       make([]*Transaction, 0),
		Relay_Pools: make(map[string][]*Transaction),
	}
}

func (pool *Tx_pool) AddTx(tx *Transaction) {
	pool.lock.Lock()
	pool.Queue = append(pool.Queue, tx)
	pool.lock.Unlock()
}

func (pool *Tx_pool) AddTxs(txs []*Transaction) {
	pool.lock.Lock()
	// fmt.Printf("收到交易%v\n", txs)
	pool.Queue = append(pool.Queue, txs...)
	fmt.Printf("len(pool.Queue): %d\n", len(pool.Queue))

	pool.lock.Unlock()
}

func (pool *Tx_pool) AddTxsToTop(txs []*Transaction) {
	pool.lock.Lock()
	// fmt.Printf("收到交易%v\n", txs)
	pool.Queue = append(txs, pool.Queue...)
	fmt.Printf("len(pool.Queue): %d\n", len(pool.Queue))

	pool.lock.Unlock()
}

func (pool *Tx_pool) FetchTxs2Pack() (txs []*Transaction) {
	config := params.Config
	tx_cnt := config.MaxBlockSize
	pool.lock.Lock()
	if len(pool.Queue) < config.MaxBlockSize {
		tx_cnt = len(pool.Queue)
	}
	txs = pool.Queue[:tx_cnt]
	pool.Queue = pool.Queue[tx_cnt:]
	pool.lock.Unlock()
	return
}

// relay
func (pool *Tx_pool) AddRelayTx(tx *Transaction, shardID string) {
	pool.lock.Lock()
	queue, ok := pool.Relay_Pools[shardID]
	if !ok {
		pool.Relay_Pools[shardID] = make([]*Transaction, 0)
	}
	queue = append(queue, tx)
	pool.Relay_Pools[shardID] = queue
	pool.lock.Unlock()
}

func (pool *Tx_pool) FetchRelayTxs(shardID string) (txs []*Transaction, isEnough bool) {
	config := params.Config
	pool.lock.Lock()
	if queue, ok := pool.Relay_Pools[shardID]; !ok {
		pool.lock.Unlock()
		return nil, false
	} else if len(queue) < config.MinRelayBlockSize {
		pool.lock.Unlock()
		return nil, false
	} else {
		tx_cnt := utils.Min(len(queue), config.MaxRelayBlockSize)
		txs := queue[:tx_cnt]
		pool.Relay_Pools[shardID] = queue[tx_cnt:]
		pool.lock.Unlock()
		return txs, true
	}
}
