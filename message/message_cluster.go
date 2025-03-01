package message

import (
	"blockEmulator/core"
)

const (
	TXaux_1 MessageType = "TXaux1"
	TXaux_2 MessageType = "TXaux2"
	TXann   MessageType = "TXann"
	TXns    MessageType = "TXns"
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
