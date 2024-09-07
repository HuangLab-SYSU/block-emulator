package query

import (
	"blockEmulator/chain"
	"blockEmulator/core"
	"blockEmulator/params"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func TestQuery(t *testing.T) {
	// pre-build a blockchain
	buildBlockChain()
	fmt.Println("Now a new blockchain is generated.")

	// query block data from the database files (boltDB)
	blocks := QueryBlocks(0, 0)
	fmt.Println("The number of blocks in this shard:", len(blocks))

	block_a := QueryBlock(0, 0, 2)
	block_a.PrintBlock()

	block_b := QueryNewestBlock(0, 0)
	block_b.PrintBlock()

	// query tx data from the database files
	txs := QueryBlockTxs(0, 0, 3)
	fmt.Println(len(txs))

	// query account state from level db
	mptfp := params.DatabaseWrite_path + "mptDB/ldb/s0/n0"
	chaindbfp := params.DatabaseWrite_path + fmt.Sprintf("chainDB/S%d_N%d", 0, 0)
	accountState := QueryAccountState(chaindbfp, mptfp, 0, 0, "00000000001")
	fmt.Println("The account balance of 00000000001:", accountState.Balance)

	clearBlockchainData()
}

func buildBlockChain() {
	accounts := []string{"00000000001", "00000000002", "00000000003", "00000000004", "00000000005", "00000000006"}
	as := make([]*core.AccountState, 0)
	for idx := range accounts {
		as = append(as, &core.AccountState{
			Balance: big.NewInt(int64(idx*1000000) + 1000000),
		})
	}
	fp := params.DatabaseWrite_path + "mptDB/ldb/s0/N0"
	fmt.Println(fp)
	db, err := rawdb.NewLevelDBDatabase(fp, 0, 1, "accountState", false)
	if err != nil {
		log.Panic(err)
	}
	params.ShardNum = 1
	pcc := &params.ChainConfig{
		ChainID:        0,
		NodeID:         0,
		ShardID:        0,
		Nodes_perShard: 1,
		ShardNums:      1,
		BlockSize:      uint64(params.MaxBlockSize_global),
		BlockInterval:  uint64(params.Block_Interval),
		InjectSpeed:    uint64(params.InjectSpeed),
	}
	CurChain, _ := chain.NewBlockChain(pcc, db)
	CurChain.PrintBlockChain()
	CurChain.AddAccounts(accounts, as, 0)
	CurChain.Txpool.AddTx2Pool(core.NewTransaction("00000000001", "00000000002", big.NewInt(100000), 1, time.Now()))

	for i := 0; i < 4; i++ {
		b := CurChain.GenerateBlock(int32(i))
		CurChain.AddBlock(b)
	}

	astates := CurChain.FetchAccounts(accounts)
	for idx, state := range astates {
		fmt.Println(accounts[idx], state.Balance)
	}
	CurChain.CloseBlockChain()
}

func clearBlockchainData() {
	// clear test data file
	err := os.RemoveAll(params.ExpDataRootDir)
	if err != nil {
		fmt.Printf("Failed to delete directory: %v\n", err)
		return
	}
}
