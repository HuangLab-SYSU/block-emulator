package test

import (
	"blockEmulator/partition"
	"encoding/csv"
	"io"
	"log"
	"os"
)

func Test_CLPA() {
	k := new(partition.CLPAState)
	k.Init_CLPAState(0.5, 100, 4)

	txfile, err := os.Open("../2000000to2999999_BlockTransaction.csv")
	if err != nil {
		log.Panic(err)
	}

	defer txfile.Close()
	reader := csv.NewReader(txfile)
	datanum := 0
	reader.Read()
	for {
		data, err := reader.Read()
		if err == io.EOF || datanum == 200000 {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		if data[6] == "0" && data[7] == "0" && len(data[3]) > 16 && len(data[4]) > 16 && data[3] != data[4] {
			s := partition.Vertex{
				Addr: data[3][2:],
			}
			r := partition.Vertex{
				Addr: data[4][2:],
			}
			k.AddEdge(s, r)
			datanum++
		}
	}

	k.CLPA_Partition()

	print(k.CrossShardEdgeNum)
}
