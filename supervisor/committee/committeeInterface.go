package committee

import "blockEmulator/message"

type CommitteeModule interface {
	AdjustByBlockInfos(*message.BlockInfoMsg)
	TxHandling()
	HandleOtherMessage([]byte)
}
