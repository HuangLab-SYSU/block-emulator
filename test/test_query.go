package test

import (
	"blockEmulator/query"
	"fmt"
)

func TestQueryBlocks() {
	blocks := query.QueryBlocks(2, 0)
	fmt.Println(len(blocks))
}

func TestQueryBlock() {
	block := query.QueryBlock(0, 1, 22)
	fmt.Println(len(block.Body))
}

func TestQueryNewestBlock() {
	blocks := query.QueryNewestBlock(0, 0)
	fmt.Println(blocks.Header.Number)
}

func TestQueryBlockTxs() {
	txs := query.QueryBlockTxs(0, 0, 3)
	fmt.Println(len(txs))
}

func TestQueryAccountState() {
	accountState := query.QueryAccountState(0, 0, "32be343b94f860124dc4fee278fdcbd38c102d88")
	fmt.Println(accountState.Balance)
}
