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
	Sender int
}

type TXAUX_2_MSG struct {
	Msg core.TXmig2
}

type TXANN struct {
	Msg core.TXann
}

type TXNS struct {
	Msg core.TXns
}
