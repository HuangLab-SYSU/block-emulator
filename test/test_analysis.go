package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// const (
// 	Shard_no = 3
// )

func addrtoShard(Addr string, Shard_no int) int {
	// 只取地址后五位已绝对够用
	Addr = Addr[len(Addr)-5:]
	num, err := strconv.ParseInt(Addr, 16, 32)
	// num, err := strconv.ParseInt(senderAddr, 10, 32)
	if err != nil {
		log.Panic()
	}
	return int(num) % Shard_no
}

func Load(path string) {
	// for Shard_no:=2;Shard_no<33;Shard_no++ {
		//写首行
		//首行
		titles := "分片,总交易数量\n"

		var stringBuilder strings.Builder
		stringBuilder.WriteString(titles)

		// filename := "./log for len4_test5/Motivation1_3_" + strconv.Itoa(Shard_no) + "个分片.csv"
		filename := "./log for len4_test5/Motivation1_1.csv"
		file1, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModeAppend|os.ModePerm)

		for Shard_no := 2; Shard_no < 33; Shard_no++ {
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

			// txs := make([]int, Shard_no)
			tx,ctx := 0.0, 0.0
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
				// sender, _ := hex.DecodeString(row[0][2:])
				// recipient, _ := hex.DecodeString(row[1][2:])
				if addrtoShard(row[0][2:], Shard_no) != addrtoShard(row[1][2:], Shard_no) {
					ctx++
				}
				// } else {
				// 	txs[addrtoShard(row[0][2:], Shard_no)]++
				// 	txs[addrtoShard(row[1][2:], Shard_no)]++
				// }
				tx++
				

			}

			// for i:=0;i<Shard_no;i++ {
			// 	dataStrings := fmt.Sprintf("%v,%v\n", i, txs[i])
			// 	stringBuilder.WriteString(dataStrings)
			// }
			dataStrings := fmt.Sprintf("%v,%v,%v%%,%v%%\n", Shard_no, tx, 100*(tx-ctx)/tx, 100*ctx/tx)
			stringBuilder.WriteString(dataStrings)
		
		}

	// fmt.Println(tx, tx-ctx, ctx)
	dataString := stringBuilder.String()
	file1.WriteString(dataString)
	file1.Close()
	// }


	
}
