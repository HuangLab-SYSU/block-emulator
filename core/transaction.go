// Definition of transaction

package core

import (
	"blockEmulator/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"time"
)

type TXmig1 struct {
	Address     string //`json:"address"`
	FromshardID uint64 //`json:"fromshardID"`
	ToshardID   uint64 //`json:"toshardID"`
	// Request_Time int64
	// CommitTime   int64
	// ID           int
}

type TXmig2 struct {
	Txmig1  *TXmig1
	MPmig1  bool
	State   *AccountState
	MPstate bool
	// H       int
	// Address string   `json:"address"`
	// Value   *big.Int `json:"value"`
}

type TXann struct {
	Txmig2    *TXmig2
	MPmig2    bool
	State     *AccountState
	MPstate   bool
	H         int
	Address   string `json:"address"`
	ToshardID int    `json:"toshardID"`
}

type TXns struct {
	Txann   *TXann
	MPann   bool
	State   *AccountState
	MPstate bool
	H       int
	Address string   `json:"address"`
	Change  *big.Int `json:"value"`
}

type Transaction struct {
	Sender    utils.Address
	Recipient utils.Address
	Nonce     uint64
	Signature []byte // not implemented now.
	Value     *big.Int
	TxHash    []byte

	Time time.Time // TimeStamp the tx proposed.

	// used in transaction relaying
	Relayed bool
	// used in broker, if the tx is not a broker1 or broker2 tx, these values should be empty.
	HasBroker      bool
	SenderIsBroker bool
	OriginalSender utils.Address
	FinalRecipient utils.Address
	RawTxHash      []byte
}

func (tx *Transaction) PrintTx() string {
	vals := []interface{}{
		tx.Sender[:],
		tx.Recipient[:],
		tx.Value,
		string(tx.TxHash[:]),
	}
	res := fmt.Sprintf("%v\n", vals)
	return res
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

// Decode transaction
func DecodeTx(to_decode []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}

// new a transaction
func NewTransaction(sender, recipient string, value *big.Int, nonce uint64, proposeTime time.Time) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Value:     value,
		Nonce:     nonce,
		Time:      proposeTime,
	}

	hash := sha256.Sum256(tx.Encode())
	tx.TxHash = hash[:]
	tx.Relayed = false
	tx.FinalRecipient = ""
	tx.OriginalSender = ""
	tx.RawTxHash = nil
	tx.HasBroker = false
	tx.SenderIsBroker = false
	return tx
}
