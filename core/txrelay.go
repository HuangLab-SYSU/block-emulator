package core

import (
	"blockEmulator/account"
	"bytes"
	"encoding/gob"
	"log"
)

type TXrelay struct {
	Txcs      *Transaction
	MPcs      *ProofDB
	State     *account.AccountState
	MPstate   *ProofDB
	H         int
}

// Encode transaction for storing
func (tx *TXrelay) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeTXrelay(to_decode []byte) *TXrelay {
	var tx TXrelay

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}
