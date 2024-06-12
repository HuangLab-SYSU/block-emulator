package core

import (
	"bytes"
	"encoding/gob"
	"log"
)

type HashEnc struct {
	Hash []byte
	Enc  []byte
}

type ProofDB []HashEnc

func (pdb *ProofDB) Put(key []byte, value []byte) error {
	*pdb = append(*pdb, HashEnc{key, value})
	return nil
}

// Encode transaction for storing
func (pdb *ProofDB) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(pdb)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func (pdb ProofDB) Delete(key []byte) error {
	return nil
}
