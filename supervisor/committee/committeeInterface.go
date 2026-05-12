package committee

import "blockEmulator/message"

type CommitteeModule interface {
	HandleBlockInfo(*message.BlockInfoMsg)
	MsgSendingControl()
	HandleOtherMessage([]byte)
	Result_save()
}
