package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func LoadTx(path string) {
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

	for {
		row, err := r.Read()
		// fmt.Printf("%v %v %v\n", row[0][2:], row[1][2:], row[2])
		if err != nil && err != io.EOF {
			log.Panic()
		}
		if err == io.EOF {
			break
		}

		value, _ := strconv.ParseFloat(row[2], 64)
		// fee := new(big.Int)
		// fee.SetString(row[3], 10)
		fee, _ := strconv.ParseFloat(row[3], 64)
		fmt.Printf("value: %v, fee: %v\n", value, fee)

	}

}
