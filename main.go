package main

import (
	"blockEmulator/test"
)

// "fmt"

// GO111MODULE=on go run main.go

func main() {
	// test.Test_account()
	// test.Test_blockChain()
	// test.Test_pbft()
	// test.Test_node()
	// test.Test_random()
	// test.Test_pool()
	// fmt.Println(len("a1e4380a3b1f749673e270229993ee55f35663b4"))
	// path := "len4_test8.csv"
	// test.Acc2shard(path)
	// test.LoadTx(path)
	// txs := test.Newcsv(path)
	// test.WriteToCsv(txs)
	// test.F()
	// time.Sleep(10000*time.Millisecond)
	// test.TxDelayCsv()
	// test.Goru()

	test.Test_shard()

	// test.Test_DB("S0", "N0")
	// fmt.Println()
	// test.Test_DB("S0", "N1")
	// fmt.Println()
	// test.Test_DB("S1", "N0")
	// fmt.Println()
	// test.Test_DB("S1", "N1")
}
