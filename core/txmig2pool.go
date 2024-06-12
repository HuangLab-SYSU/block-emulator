package core

import (
	"blockEmulator/params"
	"sync"
)

type TXmig2_pool struct {
	// 迁入队列
	Queue []*TXmig2
	//锁
	lock sync.Mutex
}

func NewTXmig2Pool() *TXmig2_pool {
	pool := TXmig2_pool{
		Queue: make([]*TXmig2, 0),
	}
	return &pool
}

func (pool *TXmig2_pool) AddTXmig2(mig2 *TXmig2) {
	pool.lock.Lock()
	pool.Queue = append(pool.Queue, mig2)
	pool.lock.Unlock()
}

func (pool *TXmig2_pool) AddTXmig2s(mig2s []*TXmig2) {
	for _, mig2 := range mig2s {
		pool.AddTXmig2(mig2)
	}
}

func (pool *TXmig2_pool) FetchTXmig2s2Pack(quota int) (mig2s []*TXmig2, left int) {

	max := quota
	mig2s_cnt := max
	pool.lock.Lock()
	if len(pool.Queue) < max {
		mig2s_cnt = len(pool.Queue)
	}
	mig2s = make([]*TXmig2, 0)
	mig2s = append(mig2s, pool.Queue[:mig2s_cnt]...)
	pool.Queue = pool.Queue[mig2s_cnt:]
	pool.lock.Unlock()
	left = max-mig2s_cnt
	return
}

func (pool *TXmig2_pool) FetchTXmig2s2Pack2() (mig2s []*TXmig2) {
	config := params.Config
	mig2_cnt := config.MaxMig2Size
	pool.lock.Lock()
	if len(pool.Queue) < config.MaxMig2Size {
		mig2_cnt = len(pool.Queue)
	}
	mig2s = make([]*TXmig2, 0)
	mig2s = append(mig2s, pool.Queue[:mig2_cnt]...)
	pool.Queue = pool.Queue[mig2_cnt:]
	pool.lock.Unlock()
	return
}
