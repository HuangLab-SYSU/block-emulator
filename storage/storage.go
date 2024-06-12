package storage

import (
	"blockEmulator/core"
	"errors"
	"log"

	"blockEmulator/params"

	"github.com/boltdb/bolt"
)

const (
	dbFile                = "blockchain_db"
	blocksBucket          = "blocks"
	blockHeaderBucket     = "blockHeaders"
	newestBlockHashBucket = "newBlockHash"
	stateTreeBucket       = "stateTree"
)

type Storage struct {
	dbFile                string
	blocksBucket          string
	blockHeaderBucket     string
	newestBlockHashBucket string
	DB                    *bolt.DB
}

func NewStorage(chainConfig *params.ChainConfig) *Storage {
	s := &Storage{
		dbFile:                chainConfig.ShardID + "_" + chainConfig.NodeID + "_" + dbFile,
		blocksBucket:          blocksBucket,
		blockHeaderBucket:     blockHeaderBucket,
		newestBlockHashBucket: newestBlockHashBucket,
	}

	db, err := bolt.Open(s.dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.blocksBucket))
		if err != nil {
			log.Panic("create blocksBucket failed")
		}

		_, err = tx.CreateBucketIfNotExists([]byte(s.blockHeaderBucket))
		if err != nil {
			log.Panic("create blockHeaderBucket failed")
		}

		_, err = tx.CreateBucketIfNotExists([]byte(s.newestBlockHashBucket))
		if err != nil {
			log.Panic("create newestBlockHashBucket failed")
		}

		return nil
	})

	s.DB = db

	return s
}

func (s *Storage) AddBlock(block *core.Block) {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		blockBucket := tx.Bucket([]byte(s.blocksBucket))
		err := blockBucket.Put(block.Hash, block.Encode())
		if err != nil {
			log.Panic()
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	s.AddBlockHeader(block.Hash, block.Header)

	s.UpdateNewestBlockHash(block.Hash)

}

func (s *Storage) GetBlock(blockHash []byte) (*core.Block, error) {
	var block *core.Block

	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash[:])

		if blockData == nil {
			return errors.New("block is not found")
		}

		block = core.DecodeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

func (s *Storage) AddBlockHeader(blockHash []byte, header *core.BlockHeader) {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		blockHeaderBucket := tx.Bucket([]byte(s.blockHeaderBucket))
		err := blockHeaderBucket.Put(blockHash, header.Encode())
		if err != nil {
			log.Panic()
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (s *Storage) GetBlockHeader(blockHash []byte) (*core.BlockHeader, error) {
	var blockHeader *core.BlockHeader

	err := s.DB.View(func(tx *bolt.Tx) error {
		bh := tx.Bucket([]byte(blockHeaderBucket))

		blockHeaderData := bh.Get(blockHash[:])

		if blockHeaderData == nil {
			return errors.New("blockHeader is not found")
		}

		blockHeader = core.DecodeBlockHeader(blockHeaderData)

		return nil
	})
	if err != nil {
		return blockHeader, err
	}

	return blockHeader, nil
}

// func (s *Storage) UpdateStateTree(statusTree *trie.Trie) {
// 	err := s.DB.Update(func(tx *bolt.Tx) error {
// 		stateTreeBucket := tx.Bucket([]byte(stateTreeBucket))
// 		err := stateTreeBucket.Put([]byte(s.stateTreeBucket), statusTree.Encode())
// 		if err != nil {
// 			log.Panic()
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

// func (s *Storage) GetStatusTree() (*trie.Trie, error) {
// 	var stateTree *trie.Trie

// 	err := s.DB.View(func(tx *bolt.Tx) error {
// 		st := tx.Bucket([]byte(stateTreeBucket))

// 		stateTreeData := st.Get([]byte(s.stateTreeBucket))
// 		if stateTreeData == nil {
// 			return errors.New("stateTree is not found")
// 		}
// 		stateTree = trie.DecodeStateTree(stateTreeData)

// 		return nil
// 	})
// 	if err != nil {
// 		return stateTree, err
// 	}

// 	return stateTree, nil
// }

func (s *Storage) UpdateNewestBlockHash(blockHash []byte) {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		newestBlockHashBucket := tx.Bucket([]byte(newestBlockHashBucket))
		// the bucket has the only key "newBlockHash"
		err := newestBlockHashBucket.Put([]byte(s.newestBlockHashBucket), blockHash)
		if err != nil {
			log.Panic()
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (s *Storage) GetNewestBlockHash() ([]byte, error) {
	var newestBlockHash []byte

	err := s.DB.View(func(tx *bolt.Tx) error {
		nh := tx.Bucket([]byte(newestBlockHashBucket))
		// the bucket has the only key "newBlockHash"
		newestBlockHash = nh.Get([]byte(newestBlockHashBucket))
		if newestBlockHash == nil {
			return errors.New("newestBlockHash is not found")
		}

		return nil
	})

	return newestBlockHash, err
}
