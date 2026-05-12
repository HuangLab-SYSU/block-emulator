package message

import "blockEmulator/core"

// if transaction relaying is used, this message is used for sending sequence id, too
type Relay struct {
	Txs           []*core.Transaction
	SenderShardID uint64
	SenderSeq     uint64
}
