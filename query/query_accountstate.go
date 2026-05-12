package query

import (
	"blockEmulator/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
	"log"
	"strconv"
)

func QueryAccountState(ShardID, NodeID uint64, address string) *core.AccountState {
	fp := "./record/ldb/s" + strconv.FormatUint(ShardID, 10) + "/n" + strconv.FormatUint(NodeID, 10)
	db, _ := rawdb.NewLevelDBDatabase(fp, 0, 1, "accountState", false)
	triedb := trie.NewDatabaseWithConfig(db, &trie.Config{
		Cache:     0,
		Preimages: true,
	})
	defer db.Close()
	storage := initStorage(ShardID, NodeID)
	curHash, err := storage.GetNewestBlockHash()
	curb, err := storage.GetBlock(curHash)
	st, err := trie.New(trie.TrieID(common.BytesToHash(curb.Header.StateRoot)), triedb)
	if err != nil {
		log.Panic()
	}
	asenc, _ := st.Get([]byte(address))
	var state_a *core.AccountState
	state_a = core.DecodeAS(asenc)
	return state_a
}
