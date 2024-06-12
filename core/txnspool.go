package core

import (
	"blockEmulator/params"
	"sync"
)

type TXns_pool struct {
	// 迁入队列
	Queue []*TXns
	//锁
	lock sync.Mutex
}

func NewTXnsPool() *TXns_pool {
	pool := TXns_pool{
		Queue: make([]*TXns, 0),
	}
	return &pool
}

func (pool *TXns_pool) AddTXns(ns *TXns) {
	pool.lock.Lock()
	pool.Queue = append(pool.Queue, ns)
	pool.lock.Unlock()
}

func (pool *TXns_pool) AddTXnss(nss []*TXns) {
	for _, ann := range nss {
		pool.AddTXns(ann)
	}
}

func (pool *TXns_pool) FetchTXnss2Pack(quota int) (nss []*TXns, left int) {

	max := quota
	nss_cnt := max
	pool.lock.Lock()
	if len(pool.Queue) < max {
		nss_cnt = len(pool.Queue)
	}
	nss = make([]*TXns, 0)
	nss = append(nss, pool.Queue[:nss_cnt]...)
	pool.Queue = pool.Queue[nss_cnt:]
	left = max - nss_cnt
	pool.lock.Unlock()
	return
}

func (pool *TXns_pool) FetchTXnss2Pack2() (nss []*TXns) {
	config := params.Config
	ann_cnt := config.MaxCapSize
	pool.lock.Lock()
	if len(pool.Queue) < config.MaxCapSize {
		ann_cnt = len(pool.Queue)
	}
	nss = make([]*TXns, 0)
	nss = append(nss, pool.Queue[:ann_cnt]...)
	pool.Queue = pool.Queue[ann_cnt:]
	pool.lock.Unlock()
	return
}
