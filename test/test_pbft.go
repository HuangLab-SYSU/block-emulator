//
//
// 将没有分片的情况视为分片数量为1的情况，此文件作废，参照test_shard.go
//
//

package test

// import (
// 	"blockEmulator/core"
// 	"blockEmulator/params"
// 	"blockEmulator/pbft"
// 	"encoding/csv"
// 	"encoding/hex"
// 	"fmt"
// 	"io"
// 	"log"
// 	"math/big"
// 	"os"
// 	"time"
// )

// var (
// 	txs []*core.Transaction
// )

// func Test_pbft() {
// 	// 命令：go run main.go 1(shard_num) 1(malicious_num) N0 len3.csv
// 	// 或 go run main.go 1(shard_num) 1(malicious_num) N1(/N2/N3)
// 	if len(os.Args) < 4 {
// 		log.Panic("输入的参数有误！")
// 	}
// 	if os.Args[3] == "N0" {
// 		if len(os.Args) != 5 {
// 			log.Panic("输入的参数有误！")
// 		} else {
// 			txs = LoadTxs(os.Args[4])
// 		}
// 	}

// 	nodeID := os.Args[1]
// 	shardID := "S0"
// 	if addr, ok := params.NodeTable[shardID][nodeID]; ok {
// 		p := pbft.NewPBFT(shardID, nodeID, addr)
// 		go p.TcpListen() //启动节点
// 		block := p.Node.CurChain.CurrentBlock
// 		fmt.Printf("current block: \n")
// 		block.PrintBlock()

// 		// 定时打包区块
// 		if nodeID == "N0" {
// 			go InjectTxs(p.Node.CurChain.Tx_pool)
// 			go p.Propose()
// 		}
// 	} else {
// 		log.Fatal("无此节点编号！")
// 	}
// 	select {}
// }

// func LoadTxs(path string) []*core.Transaction {
// 	txs := make([]*core.Transaction, 0)
// 	file, err := os.Open(path)
// 	if err != nil {
// 		log.Panic()
// 	}
// 	defer file.Close()
// 	r := csv.NewReader(file)
// 	_, err = r.Read()
// 	if err != nil {
// 		log.Panic()
// 	}
// 	for {
// 		row, err := r.Read()
// 		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
// 		if err != nil && err != io.EOF {
// 			log.Panic()
// 		}
// 		if err == io.EOF {
// 			break
// 		}
// 		sender, _ := hex.DecodeString(row[0][2:])
// 		recipient, _ := hex.DecodeString(row[1][2:])
// 		value := new(big.Int)
// 		value.SetString(row[2], 10)
// 		txs = append(txs, &core.Transaction{
// 			Sender:    sender,
// 			Recipient: recipient,
// 			Value:     value,
// 		})
// 	}
// 	fmt.Printf("%d\n", len(txs))
// 	txs[0].PrintTx()
// 	return txs
// }

// func InjectTxs(pool *core.Tx_pool) {
// 	// todo 过滤掉一些交易，只保留普通转账交易
// 	cnt := 0
// 	inject_speed := 10
// 	for {
// 		time.Sleep(5000 * time.Millisecond)
// 		upperBound := min(cnt+inject_speed, len(txs))
// 		pool.AddTxs(txs[cnt:upperBound])
// 		// fmt.Printf("注入交易%v\n", txs[cnt:upperBound])
// 		cnt = upperBound
// 		if cnt == len(txs) {
// 			break
// 		}
// 	}
// }
