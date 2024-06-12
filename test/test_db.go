package test

import (
	"blockEmulator/trie"
	"errors"
	"log"

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
	stateTreeBucket       string
	DB                    *bolt.DB
}

func NewStorage(shardID, nodeID string) *Storage {
	s := &Storage{
		dbFile:                shardID + "_" + nodeID + "_" + dbFile,
		blocksBucket:          blocksBucket,
		blockHeaderBucket:     blockHeaderBucket,
		newestBlockHashBucket: newestBlockHashBucket,
		stateTreeBucket:       stateTreeBucket,
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

		_, err = tx.CreateBucketIfNotExists([]byte(s.stateTreeBucket))
		if err != nil {
			log.Panic("create stateTreeBucket failed")
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

func (s *Storage) GetStatusTree() (*trie.Trie, error) {
	var stateTree *trie.Trie

	err := s.DB.View(func(tx *bolt.Tx) error {
		st := tx.Bucket([]byte(stateTreeBucket))

		stateTreeData := st.Get([]byte(s.stateTreeBucket))
		if stateTreeData == nil {
			return errors.New("stateTree is not found")
		}
		stateTree = trie.DecodeStateTree(stateTreeData)

		return nil
	})
	if err != nil {
		return stateTree, err
	}

	return stateTree, nil
}

func Test_DB(shardID, nodeID string) {
	s := NewStorage(shardID, nodeID)
	t,err := s.GetStatusTree()
	if err != nil {
		log.Panic(err)
	}
	t.PrintState()
}