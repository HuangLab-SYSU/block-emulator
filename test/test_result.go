package test

import (
	"blockEmulator/chain"
	"blockEmulator/core"
	"blockEmulator/params"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func data2tx(data []string, nonce uint64) (*core.Transaction, bool) {
	if data[6] == "0" && data[7] == "0" && len(data[3]) > 16 && len(data[4]) > 16 && data[3] != data[4] {
		val, ok := new(big.Int).SetString(data[8], 10)
		if !ok {
			log.Panic("new int failed\n")
		}
		gasPrice, ok := new(big.Int).SetString(data[10], 10)
		gasUsed, ok := new(big.Int).SetString(data[11], 10)
		fee := gasPrice.Mul(gasPrice, gasUsed)
		tx := core.NewTransaction(data[3][2:], data[4][2:], val, nonce, fee)
		return tx, true
	}
	return &core.Transaction{}, false
}

func Ttestresult(ShardNums int) {
	accountBalance := make(map[string]*big.Int)
	acCorrect := make(map[string]bool)

	txfile, err := os.Open(params.FileInput)
	if err != nil {
		log.Panic(err)
	}
	defer txfile.Close()
	reader := csv.NewReader(txfile)
	nowDataNum := 0
	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		if nowDataNum == params.TotalDataSize {
			break
		}
		if tx, ok := data2tx(data, uint64(nowDataNum)); ok {
			nowDataNum++
			if _, ok := accountBalance[tx.Sender]; !ok {
				accountBalance[tx.Sender] = new(big.Int)
				accountBalance[tx.Sender].Add(accountBalance[tx.Sender], params.Init_Balance)
			}

			if _, ok := accountBalance[tx.Recipient]; !ok {
				accountBalance[tx.Recipient] = new(big.Int)
				accountBalance[tx.Recipient].Add(accountBalance[tx.Recipient], params.Init_Balance)
			}
			if accountBalance[tx.Sender].Cmp(tx.Value) != -1 {
				accountBalance[tx.Sender].Sub(accountBalance[tx.Sender], tx.Value)
				accountBalance[tx.Recipient].Add(accountBalance[tx.Recipient], tx.Value)
			}
		}
	}
	fmt.Println(len(accountBalance))
	for sid := 0; sid < ShardNums; sid++ {
		fp := "./record/ldb/s" + strconv.FormatUint(uint64(sid), 10) + "/n0"
		db, err := rawdb.NewLevelDBDatabase(fp, 0, 1, "accountState", false)
		if err != nil {
			log.Panic(err)
		}
		pcc := &params.ChainConfig{
			ChainID:        uint64(sid),
			NodeID:         0,
			ShardID:        uint64(sid),
			Nodes_perShard: uint64(params.NodesInShard),
			ShardNums:      uint64(ShardNums),
			BlockSize:      uint64(params.MaxBlockSize_global),
			BlockInterval:  uint64(params.Block_Interval),
			InjectSpeed:    uint64(params.InjectSpeed),
		}
		CurChain, _ := chain.NewBlockChain(pcc, db)
		for key, val := range accountBalance {
			v := CurChain.FetchAccounts([]string{key})
			if val.Cmp(v[0].Balance) == 0 {
				acCorrect[key] = true
			}
		}
		CurChain.CloseBlockChain()
	}
	fmt.Println(len(acCorrect))
	if len(acCorrect) == len(accountBalance) {
		fmt.Println("test pass")
	} else {
		fmt.Println(len(accountBalance)-len(acCorrect), "accounts errs, they may be brokers~;")
		fmt.Println("if the number of err accounts is too large, the mechanism has bugs")
	}
}
