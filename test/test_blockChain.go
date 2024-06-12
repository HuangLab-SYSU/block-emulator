package test

import (
	"blockEmulator/chain"
	"blockEmulator/params"
	"fmt"
)

func Test_blockChain() {
	config := &params.ChainConfig{ChainID: 77}
	bc, _ := chain.NewBlockChain(config)

	curBlock := bc.CurrentBlock
	fmt.Printf("curBlock: \n")
	curBlock.PrintBlock()

	// stateTree := bc.StatusTrie
	// fmt.Printf("stateTree: \n")
	// stateTree.PrintState()

	// for {
	// 	time.Sleep(2 * time.Second)
	// 	chain.GenerateTxs(bc)
	// 	newBlock := bc.GenerateBlock()
	// 	bc.AddBlock(newBlock)

	// 	curBlock := bc.CurrentBlock
	// 	fmt.Printf("curBlock: \n")
	// 	curBlock.PrintBlock()

	// 	stateTree := bc.StatusTrie
	// 	fmt.Printf("stateTree: \n")
	// 	stateTree.PrintState()
	// }

	// chain.GenerateTxs(bc)
	// newBlock := bc.GenerateBlock()
	// bc.AddBlock(newBlock)

	// curBlock = bc.CurrentBlock
	// fmt.Printf("curBlock: \n")
	// curBlock.PrintBlock()

	// stateTree = bc.StatusTrie
	// fmt.Printf("stateTree: \n")
	// stateTree.PrintState()

}
