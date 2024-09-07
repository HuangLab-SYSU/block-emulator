package chain

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func TestBlockChain(t *testing.T) {
	accounts := []string{"000000000001", "00000000002", "00000000003", "00000000004", "00000000005", "00000000006"}
	as := make([]*core.AccountState, 0)
	for idx := range accounts {
		as = append(as, &core.AccountState{
			Balance: big.NewInt(int64(idx)),
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
	CurChain, _ := NewBlockChain(pcc, db)
	CurChain.PrintBlockChain()
	CurChain.AddAccounts(accounts, as, 0)
	CurChain.PrintBlockChain()

	astates := CurChain.FetchAccounts(accounts)
	for idx, state := range astates {
		fmt.Println(accounts[idx], state.Balance)
	}
	CurChain.CloseBlockChain()

	// clear test data file
	err = os.RemoveAll(params.ExpDataRootDir)
	if err != nil {
		fmt.Printf("Failed to delete directory: %v\n", err)
		return
	}
}
