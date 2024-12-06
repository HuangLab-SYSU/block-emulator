package core

import (
	"blockEmulator/utils"
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"encoding/gob"
	"encoding/hex"
)

type Transaction struct {
	Sender    []byte   `json:"sender"`
	Recipient []byte   `json:"recipient"`
	Value     *big.Int `json:"value"`
	TxHash    []byte
	Id        int
	RequestTime int64
}

func (tx *Transaction) PrintTx() {
	vals := []interface{}{
		hex.EncodeToString(tx.Sender),
		utils.Addr2Shard(hex.EncodeToString(tx.Sender)),
		hex.EncodeToString(tx.Recipient),
		utils.Addr2Shard(hex.EncodeToString(tx.Recipient)),
		tx.Value,
		hex.EncodeToString(tx.TxHash),
	}
	fmt.Printf("%v\n", vals)
}

func NewTransaction(sender, to []byte, value *big.Int) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Recipient: to,
		Value:     value,
	}

	tx.TxHash = tx.Hash()

	return tx
}

func (tx *Transaction) Hash() []byte {
	hash := sha256.Sum256(tx.Encode())
	return hash[:]
}

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
