package message

import (
	"blockEmulator/core"
)

const (
	TXaux_1      MessageType = "CUTXaux1"
	TXaux_2      MessageType = "CUTXaux2"
	TXann        MessageType = "CUTXann"
	TXns         MessageType = "CUTXns"
	ScourceQuery MessageType = "CUSourQ"
	DestReply    MessageType = "CUDestR"
	ReplayAttack MessageType = "CUAttack"
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

type CU_SOURCE_QUERY struct {
	State  core.CUTTER_ACCOUNT_STATE
	Sender uint64
}

type CU_DEST_REPLY struct {
	State  core.CUTTER_ACCOUNT_STATE
	Sender uint64
}

type CU_REPLAY_ATTACK struct {
	AccountID string
	Sender    uint64
	Location  uint64
}
