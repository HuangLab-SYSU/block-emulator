package pbft_all

import "blockEmulator/message"

// Define operations in a PBFT.
// This may be varied by different consensus protocols.

type ExtraOpInConsensus interface {
	// mining / message generation
	HandleinPropose() (bool, *message.Request)
	// checking
	HandleinPrePrepare(*message.PrePrepare) bool
	// nothing necessary
	HandleinPrepare(*message.Prepare) bool
	// confirming
	HandleinCommit(*message.Commit) bool
	// do for need
	HandleReqestforOldSeq(*message.RequestOldMessage) bool
	// do for need
	HandleforSequentialRequest(*message.SendOldMessage) bool
}

// Define operations among some PBFTs.
// This may be varied by different consensus protocols.
type OpInterShards interface {
	// operation inter-shards
	HandleMessageOutsidePBFT(message.MessageType, []byte) bool
}
