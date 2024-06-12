package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type TXmig1 struct {
	Address      string `json:"address"`
	FromshardID  int    `json:"fromshardID"`
	ToshardID    int    `json:"toshardID"`
	Request_Time int64
	CommitTime   int64
	ID           int
}

// Encode transaction for storing
func (tx *TXmig1) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeTXmig1(to_decode []byte) *TXmig1 {
	var tx TXmig1

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}

func (tx *TXmig1) Hash() []byte {
	hash := sha256.Sum256(tx.Encode())
	return hash[:]
}
