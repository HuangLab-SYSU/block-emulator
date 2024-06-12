package core

import (
	"blockEmulator/account"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type TXann struct {
	Txmig2    *TXmig2
	MPmig2    *ProofDB
	State     *account.AccountState
	MPstate   *ProofDB
	H         int
	Address   string `json:"address"`
	ToshardID int    `json:"toshardID"`
}

// Encode transaction for storing
func (tx *TXann) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeTXann(to_decode []byte) *TXann {
	var tx TXann

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}

func (tx *TXann) Hash() []byte {
	hash := sha256.Sum256(tx.Encode())
	return hash[:]
}
