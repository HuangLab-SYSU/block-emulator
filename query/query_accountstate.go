package query

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
)

func QueryAccountState(ShardID, NodeID uint64, address string) *core.AccountState {
	fp := params.DatabaseWrite_path + "mptDB/ldb/s" + strconv.FormatUint(ShardID, 10) + "/n" + strconv.FormatUint(NodeID, 10)
	db, _ := rawdb.NewLevelDBDatabase(fp, 0, 1, "accountState", false)
	triedb := trie.NewDatabaseWithConfig(db, &trie.Config{
		Cache:     0,
		Preimages: true,
	})
	defer db.Close()
	storage := initStorage(ShardID, NodeID)
	curHash, _ := storage.GetNewestBlockHash()
	curb, _ := storage.GetBlock(curHash)
	st, err := trie.New(trie.TrieID(common.BytesToHash(curb.Header.StateRoot)), triedb)
	if err != nil {
		log.Panic()
	}
	asenc, _ := st.Get([]byte(address))
	state_a := core.DecodeAS(asenc)
	return state_a
}
