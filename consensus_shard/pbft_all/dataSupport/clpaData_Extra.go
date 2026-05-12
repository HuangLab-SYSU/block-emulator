package dataSupport

import (
	"blockEmulator/core"
	"blockEmulator/message"
	"sync"
)

type Data_supportCLPA struct {
	ModifiedMap             []map[string]uint64                   // record the modified map from the decider(s)
	AccountTransferRound    uint64                                // denote how many times accountTransfer do
	ReceivedNewAccountState map[string]*core.AccountState         // the new accountState From other Shards
	ReceivedNewTx           []*core.Transaction                   // new transactions from other shards' pool
	AccountStateTx          map[uint64]*message.AccountStateAndTx // the map of accountState and transactions, pool
	PartitionOn             bool                                  // judge nextEpoch is partition or not

	PartitionReady map[uint64]bool // judge whether all shards has done all txs
	P_ReadyLock    sync.Mutex      // lock for ready

	ReadySeq     map[uint64]uint64 // record the seqid when the shard is ready
	ReadySeqLock sync.Mutex        // lock for seqMap

	CollectOver bool       // judge whether all txs is collected or not
	CollectLock sync.Mutex // lock for collect
}

func NewCLPADataSupport() *Data_supportCLPA {
	return &Data_supportCLPA{
		ModifiedMap:             make([]map[string]uint64, 0),
		AccountTransferRound:    0,
		ReceivedNewAccountState: make(map[string]*core.AccountState),
		ReceivedNewTx:           make([]*core.Transaction, 0),
		AccountStateTx:          make(map[uint64]*message.AccountStateAndTx),
		PartitionOn:             false,
		PartitionReady:          make(map[uint64]bool),
		CollectOver:             false,
		ReadySeq:                make(map[uint64]uint64),
	}
}
