package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
)

// StateAccount is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type AccountState struct {
	Nonce    uint64
	Balance  *big.Int
	Root     []byte // merkle root of the storage trie
	CodeHash []byte
}

func (s *AccountState) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(s)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeAccountState(to_decode []byte) *AccountState {
	var state AccountState

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&state)
	if err != nil {
		log.Panic(err)
	}

	return &state
}

func (s *AccountState) Hash() []byte {
	hash := sha256.Sum256(s.Encode())
	return hash[:]
}

// 从账户上扣钱
func (s *AccountState) Deduct(value *big.Int) {
	// todo 判断判断
	s.Balance = s.Balance.Sub(s.Balance, value)
}

// 往账户上打钱
func (s *AccountState) Deposit(value *big.Int) {
	// todo 判断判断
	s.Balance = s.Balance.Add(s.Balance, value)
}
