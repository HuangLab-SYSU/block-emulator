package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
)

var (
	path = "0to999999_BlockTransaction.csv"
)

func Test_random() {
	// txs := make([]*core.Transaction, 0)
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
	cnt := 0
	for {
		row, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}
		gasPrice := new(big.Int)
		gasPrice.SetString(row[10], 10)
		gasUsed := new(big.Int)
		gasUsed.SetString(row[11], 10)

		z := new(big.Int)
		z.Mul(gasPrice, gasUsed)

		fmt.Printf("%v %v %v\n", gasPrice, gasUsed, z)
		cnt += 1
		if cnt > 10 {
			break
		}
	}

}
