package message

import (
	"blockEmulator/core"
)

const (
	TXaux_1      MessageType = "CLuTXaux1"
	TXaux_2      MessageType = "CLuTXaux2"
	TXann        MessageType = "CLuTXann"
	TXns         MessageType = "CLuTXns"
	ScourceQuery MessageType = "CLuSourQ"
	DestReply    MessageType = "CLuDestR"
)

type TXAUX_1_MSG struct {
	Msg    core.TXmig1
	Sender uint64
}

type TXAUX_2_MSG struct {
	Msg    core.TXmig2
	Sender uint64
}

type TXANN_MSG struct {
	Msg    core.TXann
	Sender uint64
}

type TXNS_MSG struct {
	Msg    core.TXns
	Sender uint64
}

type CLU_SOURCE_QUERY struct {
	AccountKey string
	Sender     uint64
}

type CLU_DEST_REPLY struct {
	AccountKey      string
	AccountLocation uint64
	Sender          uint64
}
