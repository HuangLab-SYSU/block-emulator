package test

import (
	"encoding/csv"
	"fmt"
)

var (
	path   = "0to999999_BlockTransaction.csv"
	writer *csv.Writer
)

func Test_random() {
	a := float64(7)
	b := 2
	fmt.Printf("%v\n", a/float64(b))

	// csvFile, err := os.Create("len4_test1.csv")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// writer = csv.NewWriter(csvFile)
	// writer.Write([]string{"sender", "receiver", "value", "fee"})
	// writer.Flush()

	// // txs := make([]*core.Transaction, 0)
	// file, err := os.Open(path)
	// if err != nil {
	// 	log.Panic()
	// }
	// defer file.Close()
	// r := csv.NewReader(file)
	// _, err = r.Read()
	// if err != nil {
	// 	log.Panic()
	// }
	// cnt := 0
	// for {
	// 	row, err := r.Read()
	// 	// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
	// 	if err != nil && err != io.EOF {
	// 		log.Panic()
	// 	}
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if row[7] == "1" || row[8] == "0" {
	// 		continue
	// 	}
	// 	// gasPrice := new(big.Int)
	// 	// gasPrice.SetString(row[10], 10)
	// 	// gasUsed := new(big.Int)
	// 	// gasUsed.SetString(row[11], 10)

	// 	// gasFee := new(big.Int)
	// 	// gasFee.Mul(gasPrice, gasUsed)
	// 	gasPrice, err := strconv.ParseFloat(row[10], 64)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	gasUsed, err := strconv.ParseFloat(row[11], 64)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	gasFee := gasPrice * gasUsed

	// 	writer.Write([]string{row[3], row[4], row[8], fmt.Sprintf("%v", gasFee)})

	// 	cnt += 1
	// 	if cnt > 100 {
	// 		break
	// 	}
	// }
	// writer.Flush()

}
