package partition

import (
	"blockEmulator/params"
	"encoding/csv"
	"io"
	"log"
	"os"
	"testing"
)

func TestClpa(t *testing.T) {
	k := new(CLPAState)
	k.Init_CLPAState(0.5, 100, 4)

	txfile, err := os.Open("../" + params.DatasetFile)
	if err != nil {
		log.Panic(err)
	}

	defer txfile.Close()
	reader := csv.NewReader(txfile)
	datanum := 0

	// read transactions
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
			s := Vertex{
				Addr: data[3][2:],
			}
			r := Vertex{
				Addr: data[4][2:],
			}
			k.AddEdge(s, r)
			datanum++
		}
	}

	k.CLPA_Partition()
}
