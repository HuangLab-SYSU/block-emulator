// storage is a key-value database and its interfaces indeed
// the information of block will be saved in storage

package storage

import (
	"blockEmulator/core"
	"blockEmulator/params"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

type Storage struct {
	dbFilePath            string // path to the database
	blockBucket           string // bucket in bolt database
	blockHeaderBucket     string // bucket in bolt database
	newestBlockHashBucket string // bucket in bolt database
	DataBase              *bolt.DB
}

// new a storage, build a bolt datase
func NewStorage(cc *params.ChainConfig) *Storage {
	_, errStat := os.Stat("./record")
	if os.IsNotExist(errStat) {
		errMkdir := os.Mkdir("./record", os.ModePerm)
		if errMkdir != nil {
			log.Panic(errMkdir)
		}
	} else if errStat != nil {
		log.Panic(errStat)
	}

	s := &Storage{
		dbFilePath:            "./record/" + strconv.FormatUint(cc.ShardID, 10) + "_" + strconv.FormatUint(cc.NodeID, 10) + "_database",
		blockBucket:           "block",
		blockHeaderBucket:     "blockHeader",
		newestBlockHashBucket: "newestBlockHash",
	}

	db, err := bolt.Open(s.dbFilePath, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// create buckets
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.blockBucket))
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
	s.DataBase = db
	return s
}

// update the newest block in the database
func (s *Storage) UpdateNewestBlock(newestbhash []byte) {
	err := s.DataBase.Update(func(tx *bolt.Tx) error {
		nbhBucket := tx.Bucket([]byte(s.newestBlockHashBucket))
		// the bucket has the only key "OnlyNewestBlock"
		err := nbhBucket.Put([]byte("OnlyNewestBlock"), newestbhash)
		if err != nil {
			log.Panic()
		}
		return nil
	})
	if err != nil {
		log.Panic()
	}
	fmt.Println("The newest block is updated")
}

// add a blockheader into the database
func (s *Storage) AddBlockHeader(blockhash []byte, bh *core.BlockHeader) {
	err := s.DataBase.Update(func(tx *bolt.Tx) error {
		bhbucket := tx.Bucket([]byte(s.blockHeaderBucket))
		err := bhbucket.Put(blockhash, bh.Encode())
		if err != nil {
			log.Panic()
		}
		return nil
	})
	if err != nil {
		log.Panic()
	}
}

// add a block into the database
func (s *Storage) AddBlock(b *core.Block) {
	err := s.DataBase.Update(func(tx *bolt.Tx) error {
		bbucket := tx.Bucket([]byte(s.blockBucket))
		err := bbucket.Put(b.Hash, b.Encode())
		if err != nil {
			log.Panic()
		}
		return nil
	})
	if err != nil {
		log.Panic()
	}
	s.AddBlockHeader(b.Hash, b.Header)
	s.UpdateNewestBlock(b.Hash)
	fmt.Println("Block is added")
}

// read a blockheader from the database
func (s *Storage) GetBlockHeader(bhash []byte) (*core.BlockHeader, error) {
	var res *core.BlockHeader
	err := s.DataBase.View(func(tx *bolt.Tx) error {
		bhbucket := tx.Bucket([]byte(s.blockHeaderBucket))
		bh_encoded := bhbucket.Get(bhash)
		if bh_encoded == nil {
			return errors.New("the block is not existed")
		}
		res = core.DecodeBH(bh_encoded)
		return nil
	})
	return res, err
}

// read a block from the database
func (s *Storage) GetBlock(bhash []byte) (*core.Block, error) {
	var res *core.Block
	err := s.DataBase.View(func(tx *bolt.Tx) error {
		bbucket := tx.Bucket([]byte(s.blockBucket))
		b_encoded := bbucket.Get(bhash)
		if b_encoded == nil {
			return errors.New("the block is not existed")
		}
		res = core.DecodeB(b_encoded)
		return nil
	})
	return res, err
}

// read the Newest block hash
func (s *Storage) GetNewestBlockHash() ([]byte, error) {
	var nhb []byte
	err := s.DataBase.View(func(tx *bolt.Tx) error {
		bhbucket := tx.Bucket([]byte(s.newestBlockHashBucket))
		// the bucket has the only key "OnlyNewestBlock"
		nhb = bhbucket.Get([]byte("OnlyNewestBlock"))
		if nhb == nil {
			return errors.New("cannot find the newest block hash")
		}
		return nil
	})
	return nhb, err
}
