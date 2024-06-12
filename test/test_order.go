package test

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

type transaction struct {
	Sender    []byte  `json:"sender"`
	Recipient []byte  `json:"recipient"`
	Value     float64 `json:"value"`

	GasFee      float64
	Utility     float64 // 用于在交易池中排序
	TxHash      []byte
	Id          int
	RequestTime int64
}

func LoadTxsWithShard(path string) []*transaction {
	txs := make([]*transaction, 0)
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
	txid := 0
	for {
		row, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}
		// 所有交易读入内存（不再只是读入本分片交易
		sender, _ := hex.DecodeString(row[0][2:])
		recipient, _ := hex.DecodeString(row[1][2:])
		// value := new(big.Int)
		// value.SetString(row[2], 10)
		value, _ := strconv.ParseFloat(row[2], 64)
		// fee := new(big.Int)
		// fee.SetString(row[3], 10)
		fee, _ := strconv.ParseFloat(row[3], 64)
		txs = append(txs, &transaction{
			Sender:    sender,
			Recipient: recipient,
			Value:     value,
			GasFee:    fee,
			Id:        txid,
		})
		txid += 1
	}
	return txs
}

func Reorder(queue []*transaction) {
	//降序排序
	sort.Slice(queue, func(i, j int) bool {
		return queue[i].GasFee > queue[j].GasFee
	})
	for _,tx := range queue {
		fmt.Println(tx.GasFee)
	}
}
