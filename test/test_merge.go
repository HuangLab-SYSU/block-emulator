package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	path1 = []string{"./log/S0_transaction.csv", "./log/S1_transaction.csv", "./log/S2_transaction.csv", "./log/S3_transaction.csv"}
)

type TxDelay struct {
	//交易ID
	txid 					int
	//交易延迟
	delay 					int
	//跨分片交易第一次上链时间
	first_on_chain_time 	int64
	//跨分片交易第二次进入队列时间
	second_request_time 	int64
}

// func readdelay() ([]TxDelay, []TxDelay) {
// 	cnt := 0
// 	firsttx := make([]TxDelay, 0)
// 	secondtx := make([]TxDelay, 0)
// 	for i := 0; i < len(path1); i++ {
// 		file, err := os.Open(path1[i])
// 		if err != nil {
// 			log.Panic()
// 		}
// 		fmt.Println(i)
// 		r := csv.NewReader(file)
// 		_, err = r.Read()
// 		if err != nil {
// 			log.Panic()
// 		}
// 		for {
// 			row, err := r.Read()
// 			// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
// 			if err != nil && err != io.EOF {
// 				log.Panic(err)
// 			}
// 			if err == io.EOF {
// 				break
// 			}

// 			senderid, _ := strconv.Atoi(row[7])
// 			receiverid,_ := strconv.Atoi(row[8])
// 			//片内交易不管
// 			if senderid == receiverid {
// 				cnt++
// 				continue
// 			}

// 			id, _ := strconv.Atoi(row[0])
// 			delay, _ := strconv.Atoi(row[4])
// 			tx := TxDelay{txid:id, delay:delay}
			
// 			if senderid == i {
// 				firsttx = append(firsttx, tx)
// 			} else {
// 				secondtx = append(secondtx, tx)
// 			}
// 		}
// 		file.Close()
// 	}
// 	fmt.Printf("cnt: %v\n",cnt)
// 	return firsttx, secondtx
// }

func order(txs []TxDelay) {
	//升序排序
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].txid < txs[j].txid
	})
}

func TxDelayCsv() {
	ntx, temp_first, second := readdelay2()
	order(temp_first)
	order(second)
	order(ntx)
	
	first := make([]TxDelay, 0)
	j:=0
	for i := 0; i<len(second); i++ {
		
		for temp_first[j].txid!=second[i].txid {
			j++
		}
		first = append(first, temp_first[j])
	}

	if len(first) != len(second) {
		fmt.Println("不对！")
		fmt.Printf("len1: %v, len2: %v\n", len(first), len(second))
		return
	}


	//首行
	// var titles string
	titles := "txid,sender_delay,delay,time_diff,1stOnChainTo2ndEntry\n"

	var stringBuilder strings.Builder
	stringBuilder.WriteString(titles)

	for i := 0; i < len(first); i++ {
		dataString := fmt.Sprintf("%v,%v,%v,%v,%v\n", first[i].txid, first[i].delay, second[i].delay, second[i].delay-first[i].delay, second[i].second_request_time - first[i].first_on_chain_time)
		stringBuilder.WriteString(dataString)
	}
	filename := "./log/跨分片交易延迟及时间差.csv"
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModeAppend|os.ModePerm)
	dataString := stringBuilder.String()
	file.WriteString(dataString)
	file.Close()

	//首行
	// var titles string
	titles = "txid,delay\n"

	var stringBuilder2 strings.Builder
	stringBuilder2.WriteString(titles)

	for i := 0; i < len(ntx); i++ {
		dataString := fmt.Sprintf("%v,%v\n", ntx[i].txid, ntx[i].delay)
		stringBuilder2.WriteString(dataString)
	}
	filename = "./log/片内交易延迟.csv"
	file, _ = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModeAppend|os.ModePerm)
	dataString = stringBuilder2.String()
	file.WriteString(dataString)
	file.Close()
}


func readdelay2() (ntx, ctx1, ctx2 []TxDelay) {

	for i := 0; i < len(path1); i++ {
		file, err := os.Open(path1[i])
		if err != nil {
			log.Panic()
		}
		fmt.Println(i)
		r := csv.NewReader(file)
		_, err = r.Read()
		if err != nil {
			log.Panic()
		}
		for {
			row, err := r.Read()
			// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
			if err != nil && err != io.EOF {
				log.Panic(err)
			}
			if err == io.EOF {
				break
			}

			senderid, _ := strconv.Atoi(row[8])
			receiverid,_ := strconv.Atoi(row[9])
			id, _ := strconv.Atoi(row[0])
			delay, _ := strconv.Atoi(row[5])
			firsttime, _ := strconv.Atoi(row[4])
			secondetime, _ := strconv.ParseInt(row[3], 10, 64)
			tx := TxDelay{txid:id, delay:delay, first_on_chain_time: int64(firsttime), second_request_time: secondetime}
			//片内交易不管
			if senderid == receiverid {
				ntx = append(ntx, tx)
			}else if  senderid == i {
				ctx1 = append(ctx1, tx)
			} else {
				ctx2 = append(ctx2, tx)
			}
		}
		file.Close()
	}
	fmt.Printf("ntx,ctx1,ctx2: %v %v %v\n",len(ntx), len(ctx1), len(ctx2))
	return ntx, ctx1, ctx2
}