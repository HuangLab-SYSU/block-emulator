package pbft_all

import "blockEmulator/message"

type PbftInsideExtraHandleMod interface {
	HandleinPropose() (bool, *message.Request)
	HandleinPrePrepare(*message.PrePrepare) bool
	HandleinPrepare(*message.Prepare) bool
	HandleinCommit(*message.Commit) bool
	HandleReqestforOldSeq(*message.RequestOldMessage) bool
	HandleforSequentialRequest(*message.SendOldMessage) bool
}

type PbftOutsideHandleMod interface {
	HandleMessageOutsidePBFT(message.MessageType, []byte) bool
}
