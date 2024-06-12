package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func Addr2Shard(senderAddr string) int {
	// 只取地址后五位已绝对够用
	senderAddr = senderAddr[len(senderAddr)-5:]
	num, err := strconv.ParseInt(senderAddr, 16, 32)
	// num, err := strconv.ParseInt(senderAddr, 10, 32)
	if err != nil {
		log.Panic()
	}
	return int(num) % 2
}

func Acc2shard(path string) {
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
		sid2 := Addr2Shard(row[1][2:])
		if sid2 == 0 {
			cnt++
		}
	}
	fmt.Println(cnt)
}
