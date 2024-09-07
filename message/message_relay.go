package message

import (
	"blockEmulator/chain"
	"blockEmulator/core"
)

// if transaction relaying is used, this message is used for sending sequence id, too
type Relay struct {
	Txs           []*core.Transaction
	SenderShardID uint64
	SenderSeq     uint64
}

// This struct is similar to Relay. Nodes receiving this message will validate this proof first.
type RelayWithProof struct {
	Txs           []*core.Transaction
	TxProofs      []chain.TxProofResult
	SenderShardID uint64
	SenderSeq     uint64
}
