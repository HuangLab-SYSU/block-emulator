package core

import (
	"blockEmulator/account"
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"encoding/gob"
	"encoding/hex"
)

type Transaction struct {
	Sender             []byte `json:"sender"`
	Recipient          []byte `json:"recipient"`
	TxHash             []byte
	Id                 int
	Success            bool
	IsRelay            bool
	SenLock            bool
	RecLock            bool
	Value              *big.Int `json:"value"`
	RequestTime        int64
	Second_RequestTime int64
	CommitTime         int64
	LockTime           int64
	UnlockTime         int64
	LockTime2          int64
	UnlockTime2        int64
	HalfLock           bool
	Rec_Suppose_on_chain   int
	Sen_Suppose_on_chain   int
	Relay_Lock         bool
}

func (tx *Transaction) PrintTx() {
	vals := []interface{}{
		hex.EncodeToString(tx.Sender),
		account.Addr2Shard(hex.EncodeToString(tx.Sender)),
		hex.EncodeToString(tx.Recipient),
		account.Addr2Shard(hex.EncodeToString(tx.Recipient)),
		tx.Value,
		// hex.EncodeToString(tx.TxHash),
	}
	fmt.Printf("%v\n", vals)
}

func (tx *Transaction) Hash() []byte {
	hash := sha256.Sum256(tx.Encode())
	return hash[:]
}

// Encode transaction for storing
func (tx *Transaction) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeTx(to_decode []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}
