package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func Newcsv(path string) [][]string {
	txs := make([][]string, 0)
	file, err := os.Open(path)
	if err != nil {
		log.Panic()
	}
	defer file.Close()
	r := csv.NewReader(file)
	// for i:=0; i<1000000;i++{
	// 	_, err = r.Read()
	// 	if err != nil {
	// 		log.Panic()
	// 	}
	// }

	_, err = r.Read()
	if err != nil {
		log.Panic()
	}
	

	for i := 0; i < 10; i++ {
		row, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}
		txs = append(txs, row)

	}
	return txs
}

func WriteToCsv(txs [][]string) {
	//数据写入到csv文件

	//首行
	titles := "sender,receiver,value,fee\n"

	var stringBuilder strings.Builder
	stringBuilder.WriteString(titles)

	for j:=0; j<17; j++ {
		var i int
		for i = 0; i < 10; i++ {
			var dataString string
			if i == 0 {
				dataString = fmt.Sprintf("%v,%v,%v,%v\n", txs[i][0], txs[i][1], 1, 0)
			}else {
				dataString = fmt.Sprintf("%v,%v,%v,%v\n", txs[i][0], txs[i][1], 0, 0)
			}
			stringBuilder.WriteString(dataString)
		}
	}
	
	filename := "./len4_test9.csv"
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModeAppend|os.ModePerm)
	dataString := stringBuilder.String()
	file.WriteString(dataString)
	file.Close()
}
