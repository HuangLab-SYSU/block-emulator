package query

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/storage"
	"github.com/boltdb/bolt"
)

func initStorage(ShardID, NodeID uint64) *storage.Storage {
	pcc := &params.ChainConfig{
		ChainID: ShardID,
		NodeID:  NodeID,
		ShardID: ShardID,
	}
	return storage.NewStorage(pcc)
}

func QueryBlocks(ShardID, NodeID uint64) []*core.Block {
	db := initStorage(ShardID, NodeID).DataBase
	defer db.Close()
	blocks := make([]*core.Block, 0)
	err1 := db.View(func(tx *bolt.Tx) error {
		bbucket := tx.Bucket([]byte("block"))
		if err := bbucket.ForEach(func(k, v []byte) error {
			res := core.DecodeB(v)
			blocks = append(blocks, res)
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
	if err1 != nil {
		err1.Error()
	}
	return blocks
}

func QueryBlock(ShardID, NodeID, Number uint64) *core.Block {
	db := initStorage(ShardID, NodeID).DataBase
	defer db.Close()
	block := new(core.Block)
	err1 := db.View(func(tx *bolt.Tx) error {
		bbucket := tx.Bucket([]byte("block"))
		if err := bbucket.ForEach(func(k, v []byte) error {
			res := core.DecodeB(v)
			if res.Header.Number == Number {
				block = res
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
	if err1 != nil {
		err1.Error()
	}
	return block
}

func QueryNewestBlock(ShardID, NodeID uint64) *core.Block {
	storage := initStorage(ShardID, NodeID)
	defer storage.DataBase.Close()
	hash, _ := storage.GetNewestBlockHash()
	block, _ := storage.GetBlock(hash)
	return block
}

func QueryBlockTxs(ShardID, NodeID, Number uint64) []*core.Transaction {
	db := initStorage(ShardID, NodeID).DataBase
	defer db.Close()
	block := new(core.Block)
	err1 := db.View(func(tx *bolt.Tx) error {
		bbucket := tx.Bucket([]byte("block"))
		if err := bbucket.ForEach(func(k, v []byte) error {
			res := core.DecodeB(v)
			if res.Header.Number == Number {
				block = res
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
	if err1 != nil {
		err1.Error()
	}
	return block.Body
}
