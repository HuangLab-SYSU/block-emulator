package shard

import (
	"blockEmulator/account"
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/pbft"
	"blockEmulator/utils"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"time"
)

type Node struct {
	P *pbft.Pbft //代表分片内的那部分
}

// var (
// 	// txs            []*core.Transaction
// 	queue_newtxlog *csv.Writer
// )

func NewNode() *Node {
	node := new(Node)
	node.P = pbft.NewPBFT()

	go node.P.TcpListen() //启动节点
	block := node.P.Node.CurChain.CurrentBlock
	fmt.Printf("current block: \n")
	block.PrintBlock()

	config := params.Config
	account.Account2Shard = make(map[string]int)
	account.AccountInOwnShard = make(map[string]bool)
	account.BalanceBeforeOut = make(map[string]*big.Int)
	account.Outing_Acc_Before_Announce = make(map[string]bool)
	account.Outing_Acc_After_Announce = make(map[string]bool)
	account.Lock_Acc = make(map[string]bool)

	//将初始账户的映射弄好
	for _, address := range params.Init_addrs {
		account.Account2Shard[address] = utils.Addr2Shard(address)
		if utils.Addr2Shard(address) == params.ShardTable[config.ShardID] {
			account.AccountInOwnShard[address] = true
		}
	}
	fmt.Println("本分片账户长度")
	fmt.Println(len(account.AccountInOwnShard))
	params.Init_addrs = []string{}

	if config.Cross_Chain {
		if config.ShardID == "S0" {
			config.Block_interval = 5
		} else {
			config.Block_interval = 15
		}
	}

	fmt.Println()
	fmt.Println("启动完成！")
	// fmt.Println(account.Account2Shard)
	// fmt.Println(account.AccountInOwnShard)
	fmt.Println()

	if config.NodeID == "N0" {

		pbft.NewLog(config.ShardID)
		//不是由客户发送交易
		if !config.ClientSendTX {
			core.Txs = LoadTxsWithShard(config.Path)
		}

		// time.Sleep(6000 * time.Millisecond)
		if !config.Bu_Tong_Bi_Li_2 {
			// 正常情况
			if !config.Bu_Tong_Bi_Li && !config.Bu_Tong_Shi_Jian && !config.Fail && !config.Cross_Chain {
				// go node.P.Node.CurChain.Tx_pool.InjectTxs2Shard(params.ShardTable[config.ShardID])
				node.P.Node.CurChain.Tx_pool.NewInjectTxs2Shard(params.ShardTable[config.ShardID])
				fmt.Println("注入完成！")
			} else {
				if config.ShardID == "S0" {
					node.P.Node.CurChain.TXmig1_pool.Queue = append(node.P.Node.CurChain.TXmig1_pool.Queue, &core.TXmig1{Address: "489338d5e8d42e8c923d1f47361d979503d4ad68", FromshardID: 0, ToshardID: 1})
				}
				if config.Bu_Tong_Bi_Li {
					// go node.P.Node.CurChain.Tx_pool.InjectTxs2Shard(params.ShardTable[config.ShardID])
					node.P.Node.CurChain.Tx_pool.NewInjectTxs2Shard(params.ShardTable[config.ShardID])
				} else {
					node.P.Node.CurChain.Tx_pool.NewInjectTxs2Shard(params.ShardTable[config.ShardID])
				}
				fmt.Println("注入完成！")
				fmt.Printf("队列长度：%v\n", len(node.P.Node.CurChain.Tx_pool.Queue))
			}
		} else {
			if config.ShardID == "S0" {
				// addr := []string{"489338d5e8d42e8c923d1f47361d979503d4ad60","489338d5e8d42e8c923d1f47361d979503d4ad62","489338d5e8d42e8c923d1f47361d979503d4ad64","489338d5e8d42e8c923d1f47361d979503d4ad66","489338d5e8d42e8c923d1f47361d979503d4ad68","489338d5e8d42e8c923d1f47361d979503d4ad6a","489338d5e8d42e8c923d1f47361d979503d4ad6c","489338d5e8d42e8c923d1f47361d979503d4ad6e","489338d5e8d42e8c923d1f47361d979503d4ad70","489338d5e8d42e8c923d1f47361d979503d4ad72"}
				// addr := []string{"489338d5e8d42e8c923d1f47361d979503d4ad68"}
				// for _,v := range addr {
				// 	node.P.Node.CurChain.TXmig1_pool.Queue = append(node.P.Node.CurChain.TXmig1_pool.Queue, &core.TXmig1{Address: v, FromshardID: 0, ToshardID: 1})
				// }
				fmt.Println("不用自己来啦")
			}
			fmt.Println("注入完成！")
		}

		if config.Pressure && config.ShardID == "S0" {
			core.OutAccs = LoadOutAccs("./accountCount.csv")
			node.P.Node.CurChain.TXmig1_pool.NewInjectOutAccs2Shard()
			fmt.Println("注入完成！")
		}

		//每个分片同时开始注入交易和出块，因此等到离初始时间5s才开始
		for time.Now().Unix()-pbft.InitTime != 2 {

		}
		// fmt.Println("要准备出块咯")
		if config.Bu_Tong_Bi_Li_2 {
			if !config.ClientSendTX {
				// go node.P.Node.CurChain.Tx_pool.InjectTxs2Shard(params.ShardTable[config.ShardID])
				go node.P.Node.CurChain.Tx_pool.NewInjectTxs2Shard(params.ShardTable[config.ShardID])
			} else {
				fmt.Println("由客户端发送交易!")
			}
		}
		for time.Now().Unix()-pbft.InitTime != 4 {

		}
		fmt.Println("要准备出块咯")
		go node.P.Propose()

	}

	return node
}

func LoadTxsWithShard(path string) []*core.Transaction {
	txs := make([]*core.Transaction, 0)
	txid := 0

	// for i := 0; i < 2; i++ {

	file, err := os.Open(path)
	if err != nil {
		log.Panic()
	}
	defer file.Close()

	r := csv.NewReader(file)
	_, err = r.Read()
	if err != nil {
		log.Panic()
	}
	for i := 0; i < 1000000; i++ {
	// for i := 0; i < 500000; i++ {
		// for i:=0;i<56010;i++{
		row, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}

		// if params.Config.Algorithm && i<50000 {
		// 	continue
		// }

		// 所有交易读入内存（不再只是读入本分片交易
		sender, _ := hex.DecodeString(row[1][2:])
		recipient, _ := hex.DecodeString(row[2][2:])
		value := new(big.Int)
		value, ok := value.SetString(row[3], 64)
		if !ok {
			log.Panic()
		}

		txs = append(txs, &core.Transaction{
			Sender:             sender,
			Recipient:          recipient,
			Value:              value,
			Id:                 txid,
			RequestTime:        -1,
			Second_RequestTime: -1,
			CommitTime:         -1,
			LockTime:           -1,
			UnlockTime:         -1,
		})

		txid += 1
	}

	// }

	fmt.Printf("%d\n", len(txs))
	txs[0].PrintTx()
	return txs
}

func LoadOutAccs(path string) []*core.TXmig1 {
	outs := make([]*core.TXmig1, 0)
	outid := 0

	file, err := os.Open(path)
	if err != nil {
		log.Panic()
	}
	defer file.Close()

	r := csv.NewReader(file)
	_, err = r.Read()
	if err != nil {
		log.Panic()
	}
	for {
		row, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}

		if row[2] == "1" {
			continue
		}

		// 所有交易读入内存（不再只是读入本分片交易
		addr := row[1][2:]
		outs = append(outs, &core.TXmig1{
			Address:      addr,
			ToshardID:    1,
			Request_Time: -1,
			CommitTime:   -1,
			ID:           outid,
		})

		outid += 1
		if outid == 5000 {
			break
		}
	}

	// }

	fmt.Printf("%d\n", len(outs))
	return outs
}

// func LoadTxsWithShard2(path string) []*core.Transaction {
// 	txs := make([]*core.Transaction, 0)
// 	txid := 0

// 	// for j:=0;j<10;j++{

// 		for i:=0;i<len(params.Init_addrs);i++  {

// 			// 所有交易读入内存（不再只是读入本分片交易
// 			sender, _ := hex.DecodeString(params.Init_addrs[i])
// 			recipient, _ := hex.DecodeString(params.Init_addrs[0])
// 			if i!=len(params.Init_addrs)-1{
// 				recipient, _ = hex.DecodeString(params.Init_addrs[i+1])
// 			}

// 			value := float64(1)

// 			txs = append(txs, &core.Transaction{
// 				Sender:    sender,
// 				Recipient: recipient,
// 				Value:     value,
// 				Id:        txid,
// 			})

// 			txid += 1
// 		}

// 	// }

// 	fmt.Printf("%d\n", len(txs))
// 	txs[0].PrintTx()
// 	return txs
// }

// func InjectTxs2Shard(pool *core.Tx_pool, sid int) {
// 	cnt := 0
// 	inject_speed := params.Config.Inject_speed
// 	for {
// 		time.Sleep(1000 * time.Millisecond)
// 		upperBound := utils.Min(cnt+inject_speed, len(txs))

// 		queue_len_before := len(pool.Queue)
// 		for i := cnt; i < upperBound; i++ {
// 			if account.Addr2Shard(hex.EncodeToString(txs[i].Sender)) == sid { //是本分片的才进入交易池
// 				txs[i].RequestTime = time.Now().Unix()
// 				pool.AddTx(txs[i])
// 			}
// 		}
// 		cnt = upperBound
// 		if cnt == len(txs) {
// 			break
// 		}
// 	}
// }
