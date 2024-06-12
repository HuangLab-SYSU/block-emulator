package core

import (
	"blockEmulator/account"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
)

type TXmig2 struct {
	Txmig1  *TXmig1
	MPmig1  *ProofDB
	State   *account.AccountState
	MPstate *ProofDB
	H       int
	Address string   `json:"address"`
	Value   *big.Int `json:"value"`
}

// Encode transaction for storing
func (tx *TXmig2) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeTXmig2(to_decode []byte) *TXmig2 {
	var tx TXmig2

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}

func (tx *TXmig2) Hash() []byte {
	hash := sha256.Sum256(tx.Encode())
	return hash[:]
}
