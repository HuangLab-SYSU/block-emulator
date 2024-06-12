package core

import (
	"blockEmulator/params"
	"sync"
)

type TXann_pool struct {
	// 迁入队列
	Queue []*TXann
	//锁
	lock sync.Mutex
}

func NewTXannPool() *TXann_pool {
	pool := TXann_pool{
		Queue: make([]*TXann, 0),
	}
	return &pool
}

func (pool *TXann_pool) AddTXann(ann *TXann) {
	pool.lock.Lock()
	pool.Queue = append(pool.Queue, ann)
	pool.lock.Unlock()
}

func (pool *TXann_pool) AddTXanns(anns []*TXann) {
	for _, ann := range anns {
		pool.AddTXann(ann)
	}
}

func (pool *TXann_pool) FetchTXanns2Pack(quota int) (anns []*TXann, left int) {

	max := quota
	anns_cnt := max
	pool.lock.Lock()
	if len(pool.Queue) < max {
		anns_cnt = len(pool.Queue)
	}
	anns = make([]*TXann, 0)
	anns = append(anns, pool.Queue[:anns_cnt]...)
	pool.Queue = pool.Queue[anns_cnt:]
	pool.lock.Unlock()
	left = max-anns_cnt
	return
}

func (pool *TXann_pool) FetchTXanns2Pack2() (anns []*TXann) {
	config := params.Config
	ann_cnt := config.MaxAnnSize
	pool.lock.Lock()
	if len(pool.Queue) < config.MaxAnnSize {
		ann_cnt = len(pool.Queue)
	}
	anns = make([]*TXann, 0)
	anns = append(anns, pool.Queue[:ann_cnt]...)
	pool.Queue = pool.Queue[ann_cnt:]
	pool.lock.Unlock()
	return
}
