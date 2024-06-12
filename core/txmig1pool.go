package core

import (
	"blockEmulator/params"
	"sync"
	"time"
)

type TXmig1_pool struct {
	// 迁出队列
	Queue []*TXmig1
	//锁
	lock sync.Mutex
}

func NewTXmig1Pool() *TXmig1_pool {
	pool := TXmig1_pool{
		Queue: make([]*TXmig1, 0),
	}
	return &pool
}

func (pool *TXmig1_pool) AddTXmig1(txmig1 *TXmig1) {
	pool.lock.Lock()
	pool.Queue = append(pool.Queue, txmig1)
	pool.lock.Unlock()
}

func (pool *TXmig1_pool) AddTXmig1s(txmig1s []*TXmig1) {
	for _, txmig1 := range txmig1s {
		pool.AddTXmig1(txmig1)
	}
}

func (pool *TXmig1_pool) FetchTXmig1s2Pack() (outs []*TXmig1, left int) {
	config := params.Config
	if !config.Bu_Tong_Bi_Li_2{
		pool.lock.Lock()
		outs = make([]*TXmig1, 0)
		outs = append(outs, pool.Queue...)
		pool.Queue = []*TXmig1{}
		pool.lock.Unlock()
	}else {
		max := config.MaxMigSize
		out_cnt := max
		pool.lock.Lock()
		if len(pool.Queue) < max {
			out_cnt = len(pool.Queue)
		}
		outs = make([]*TXmig1, 0)
		outs = append(outs, pool.Queue[:out_cnt]...)
		pool.Queue = pool.Queue[out_cnt:]
		pool.lock.Unlock()
		left = max-out_cnt
	}
	
	return
}

func (pool *TXmig1_pool) FetchTXmig1s2Pack2() (txmig1s []*TXmig1) {
	config := params.Config
	out_cnt := config.MaxMig1Size
	pool.lock.Lock()
	if len(pool.Queue) < config.MaxMig1Size {
		out_cnt = len(pool.Queue)
	}
	txmig1s = make([]*TXmig1, 0)
	txmig1s = append(txmig1s, pool.Queue[:out_cnt]...)
	pool.Queue = pool.Queue[out_cnt:]
	pool.lock.Unlock()
	return
}

func (pool *TXmig1_pool) NewInjectOutAccs2Shard() {
	for i := 0; i < len(OutAccs); i++ {
		OutAccs[i].Request_Time = time.Now().UnixMilli()
		pool.AddTXmig1(OutAccs[i])
	}
	OutAccs = []*TXmig1{}
}
