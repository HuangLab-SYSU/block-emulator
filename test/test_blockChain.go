package test

import (
	"blockEmulator/chain"
	"blockEmulator/core"
	"blockEmulator/params"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/rawdb"
)

func TestBlockChain(ShardNums int) {
	accounts := []string{"000000000001", "00000000002", "00000000003", "00000000004", "00000000005", "00000000006"}
	as := make([]*core.AccountState, 0)
	for idx := range accounts {
		as = append(as, &core.AccountState{
			Balance: big.NewInt(int64(idx)),
		})
	}
	for sid := 0; sid < 1; sid++ {
		fp := "./record/ldb/s0/N0"
		db, err := rawdb.NewLevelDBDatabase(fp, 0, 1, "accountState", false)
		if err != nil {
			log.Panic(err)
		}
		params.ShardNum = 1
		pcc := &params.ChainConfig{
			ChainID:        uint64(sid),
			NodeID:         0,
			ShardID:        uint64(sid),
			Nodes_perShard: uint64(1),
			ShardNums:      4,
			BlockSize:      uint64(params.MaxBlockSize_global),
			BlockInterval:  uint64(params.Block_Interval),
			InjectSpeed:    uint64(params.InjectSpeed),
		}
		CurChain, _ := chain.NewBlockChain(pcc, db)
		CurChain.PrintBlockChain()
		CurChain.AddAccounts(accounts, as)
		CurChain.PrintBlockChain()

		astates := CurChain.FetchAccounts(accounts)
		for _, state := range astates {
			fmt.Println(state.Balance)
		}
		CurChain.CloseBlockChain()
	}
}
