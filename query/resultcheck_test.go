package query

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFinalResult(t *testing.T) {
	// check the final result after running BlockEmulator

	// get the result from Dataset
	accountBalance := loadFinalResultFromDataset()
	acCorrect := make(map[string]bool)

	// convert keys to list
	accounts := make([]string, len(accountBalance))
	i := 0
	for key := range accountBalance {
		accounts[i] = key
		i++
	}

	// check the result from BlockEmulator
	for sid := 0; sid < params.ShardNum; sid++ {
		mptfp := "../" + params.DatabaseWrite_path + "mptDB/ldb/s" + strconv.FormatUint(uint64(sid), 10) + "/n0"
		chaindbfp := "../" + params.DatabaseWrite_path + fmt.Sprintf("chainDB/S%d_N%d", sid, 0)
		aslist := QueryAccountStateList(chaindbfp, mptfp, uint64(sid), 0, accounts)
		for idx, as := range aslist {
			if as != nil && as.Balance.Cmp(accountBalance[accounts[idx]]) == 0 {
				acCorrect[accounts[idx]] = true
			}
		}
	}
	fmt.Println("Results from BlockEmulator: # of correct accounts", len(acCorrect))
	if len(acCorrect) == len(accountBalance) {
		fmt.Println("test pass")
	} else if len(accountBalance)-len(acCorrect) < params.BrokerNum {
		fmt.Printf("%d err accounts, they maybe brokers", len(accountBalance)-len(acCorrect))
	} else {
		log.Panic("Err, too many wrong accounts", len(accountBalance)-len(acCorrect))
	}
}

func data2tx(data []string, nonce uint64) (*core.Transaction, bool) {
	if data[6] == "0" && data[7] == "0" && len(data[3]) > 16 && len(data[4]) > 16 && data[3] != data[4] {
		val, ok := new(big.Int).SetString(data[8], 10)
		if !ok {
			log.Panic("new int failed\n")
		}
		tx := core.NewTransaction(data[3][2:], data[4][2:], val, nonce, time.Now())
		return tx, true
	}
	return &core.Transaction{}, false
}

func loadFinalResultFromDataset() map[string]*big.Int {
	accountBalance := make(map[string]*big.Int)
	txfile, err := os.Open("../" + params.DatasetFile)
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
	fmt.Println("Results from dataset file: # of accounts", len(accountBalance))

	return accountBalance
}
