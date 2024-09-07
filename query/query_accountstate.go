package query

import (
	"blockEmulator/core"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/trie"
)

func QueryAccountState(chainDBfp, mptfp string, ShardID, NodeID uint64, address string) *core.AccountState {
	_, err := os.Stat(mptfp)
	if os.IsNotExist(err) {
		log.Panic("No filepath")
	}
	// recover the leveldb
	db, _ := rawdb.Open(rawdb.OpenOptions{Type: "leveldb", Directory: mptfp, Namespace: "accountState"})
	triedb := trie.NewDatabaseWithConfig(db, &trie.Config{
		Cache:     0,
		Preimages: true,
	})
	defer db.Close()
	// fetch the newest blockHeader
	storage := initStorage(chainDBfp, ShardID, NodeID)
	defer storage.DataBase.Close()
	curHash, _ := storage.GetNewestBlockHash()
	curb, _ := storage.GetBlock(curHash)
	// recover the MPT
	st, err := trie.New(trie.TrieID(common.BytesToHash(curb.Header.StateRoot)), triedb)
	if err != nil {
		log.Panic()
	}
	// fetch the account
	asenc, _ := st.Get([]byte(address))
	if asenc == nil {
		return nil
	}
	state_a := core.DecodeAS(asenc)
	return state_a
}

func QueryAccountStateList(chainDBfp, mptfp string, ShardID, NodeID uint64, addresses []string) []*core.AccountState {
	_, err := os.Stat(mptfp)
	if os.IsNotExist(err) {
		log.Panic("No filepath")
	}

	ret := make([]*core.AccountState, len(addresses))
	// recover the leveldb
	db, _ := rawdb.Open(rawdb.OpenOptions{Type: "leveldb", Directory: mptfp, Namespace: "accountState"})
	triedb := trie.NewDatabaseWithConfig(db, &trie.Config{
		Cache:     0,
		Preimages: true,
	})
	defer db.Close()
	// fetch the newest blockHeader
	storage := initStorage(chainDBfp, ShardID, NodeID)
	defer storage.DataBase.Close()
	curHash, _ := storage.GetNewestBlockHash()
	curb, _ := storage.GetBlock(curHash)
	// recover the MPT
	st, err := trie.New(trie.TrieID(common.BytesToHash(curb.Header.StateRoot)), triedb)
	if err != nil {
		log.Panic()
	}
	// fetch the accounts
	for i, address := range addresses {
		asenc, _ := st.Get([]byte(address))
		if asenc == nil {
			ret[i] = nil
		} else {
			ret[i] = core.DecodeAS(asenc)
		}
	}

	return ret
}
