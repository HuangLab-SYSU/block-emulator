package pbft_clpa

import (
	"blockEmulator/core"
	"blockEmulator/message"
	"sync"
)

type Data_supportCLPA struct {
	modifiedMap             []map[string]uint64                   // record the modified map from the decider(s)
	accountTransferRound    uint64                                // denote how many times accountTransfer do
	receivedNewAccountState map[string]*core.AccountState         // the new accountState From other Shards
	receivedNewTx           []*core.Transaction                   // new transactions from other shards' pool
	accountStateTx          map[uint64]*message.AccountStateAndTx // the map of accountState and transactions, pool
	partitionOn             bool                                  // judge nextEpoch is partition or not

	partitionReady map[uint64]bool // judge whether all shards has done all txs
	pReadyLock     sync.Mutex      // lock for ready

	readySeq     map[uint64]uint64 // record the seqid when the shard is ready
	readySeqLock sync.Mutex        // lock for seqMap

	collectOver bool       // judge whether all txs is collected or not
	collectLock sync.Mutex // lock for collect
}

func NewCLPADataSupport() *Data_supportCLPA {
	return &Data_supportCLPA{
		modifiedMap:             make([]map[string]uint64, 0),
		accountTransferRound:    0,
		receivedNewAccountState: make(map[string]*core.AccountState),
		receivedNewTx:           make([]*core.Transaction, 0),
		accountStateTx:          make(map[uint64]*message.AccountStateAndTx),
		partitionOn:             false,
		partitionReady:          make(map[uint64]bool),
		collectOver:             false,
		readySeq:                make(map[uint64]uint64),
	}
}
