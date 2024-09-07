package query

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/storage"
	"fmt"

	"github.com/boltdb/bolt"
)

func initStorage(dbfp string, ShardID, NodeID uint64) *storage.Storage {
	pcc := &params.ChainConfig{
		ChainID: ShardID,
		NodeID:  NodeID,
		ShardID: ShardID,
	}
	return storage.NewStorage(dbfp, pcc)
}

func QueryBlocks(ShardID, NodeID uint64) []*core.Block {
	dbfp := params.DatabaseWrite_path + fmt.Sprintf("chainDB/S%d_N%d", ShardID, NodeID)
	db := initStorage(dbfp, ShardID, NodeID).DataBase
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
		fmt.Println(err1.Error())
	}
	return blocks
}

func QueryBlock(ShardID, NodeID, Number uint64) *core.Block {
	dbfp := params.DatabaseWrite_path + fmt.Sprintf("chainDB/S%d_N%d", ShardID, NodeID)
	db := initStorage(dbfp, ShardID, NodeID).DataBase
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
		fmt.Println(err1.Error())
	}
	return block
}

func QueryNewestBlock(ShardID, NodeID uint64) *core.Block {
	dbfp := params.DatabaseWrite_path + fmt.Sprintf("chainDB/S%d_N%d", ShardID, NodeID)
	storage := initStorage(dbfp, ShardID, NodeID)
	defer storage.DataBase.Close()
	hash, _ := storage.GetNewestBlockHash()
	block, _ := storage.GetBlock(hash)
	return block
}

func QueryBlockTxs(ShardID, NodeID, Number uint64) []*core.Transaction {
	dbfp := params.DatabaseWrite_path + fmt.Sprintf("chainDB/S%d_N%d", ShardID, NodeID)
	db := initStorage(dbfp, ShardID, NodeID).DataBase
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
		fmt.Println(err1.Error())
	}
	return block.Body
}
