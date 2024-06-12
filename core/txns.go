package core

import (
	"blockEmulator/account"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
)

type TXns struct {
	Txann   *TXann
	MPann   *ProofDB
	State   *account.AccountState
	MPstate *ProofDB
	H       int
	Address string   `json:"address"`
	Change  *big.Int `json:"value"`
}

// Encode transaction for storing
func (tx *TXns) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeTXns(to_decode []byte) *TXns {
	var tx TXns

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}

func (tx *TXns) Hash() []byte {
	hash := sha256.Sum256(tx.Encode())
	return hash[:]
}
