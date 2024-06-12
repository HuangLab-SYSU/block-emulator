package pbft

import (
	"blockEmulator/core"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
)

//<REQUEST,o,t,c>
type Request struct {
	Message
	Timestamp int64
}

type Message struct {
	Content []byte
	ID      int
}

//<<PRE-PREPARE,v,n,d>,m>
type PrePrepare struct {
	RequestMessage *Request
	Digest         string
	SequenceID     int
	Type           string // BLOCK，EpochChange，AccState
	// Sign           []byte
}

//<PREPARE,v,n,d,i>
type Prepare struct {
	Digest     string
	SequenceID int
	NodeID     string
	Type       string // BLOCK，EpochChange，AccState
	// Sign       []byte
}

//<COMMIT,v,n,D(m),i>
type Commit struct {
	Digest     string
	SequenceID int
	NodeID     string
	Type       string // Block，EpochChange，AccState
	// Sign       []byte
}

type RequestBlocks struct {
	StartID  int
	EndID    int
	ServerID string
	NodeID   string
}

type SendBlocks struct {
	StartID int
	EndID   int
	Blocks  []*core.Block
	NodeID  string
}

//<REPLY,v,t,c,i,r>
type Reply struct {
	MessageID int
	NodeID    string
	Result    bool
}

type Relay struct {
	Txs     []*core.TXrelay
	ShardID string
}

type TxFromClient struct {
	Txs     []*core.Transaction
}

type Mig2 struct {
	TXmig2s    []*core.TXmig2
	ShardID string
}

type Announce struct {
	TXanns   []*core.TXann
	ShardID string
}

type ChangeAndPending struct {
	Change     float64
	PendingTxs []*core.Transaction
}

type ChangesAndPendings struct {
	TXnss   []*core.TXns    
	List    map[string]*ChangeAndPending
	ShardID string
}

type BalanceAndPending struct {
	Balance    *big.Int
	PendingTxs []*core.Transaction
}

type BalancesAndPendings struct {
	List    map[string]*BalanceAndPending
	ShardID string
}

type NaM struct {
	New_Addrs  []string
	New_Addr2shard map[string]int
}

type Txs_and_Num_of_New_State struct {
	Txs []*core.Transaction
	BlockSize int
	Num_of_New_State int
}

// low latency
type LL_Block struct {
	ShardID int
	Block   *core.Block
	//后半数量
	Cnum int
	//片内数量
	Nnum int
	//后半总延迟
	Clatency int
	//片内总延迟
	Nlatency int
	//片内平均手续费
	Nfee float64
	//已进入系统片内数量
	Sys_Nnum int
	//已进入系统跨分片数量
	Sys_Cnum int
}

const prefixCMDLength = 12

type command string

const (
	// cRequest    command = "request"
	cPrePrepare   command = "preprepare"
	cPrepare      command = "prepare"
	cCommit       command = "commit"
	cRequestBlock command = "requestBlock"
	cSendBlock    command = "sendBlock"
	cReply        command = "reply"

	cClient   command = "client"
	cRelay    command = "relay"
	cTXmig1     command = "txmig1"
	cAnnounce command = "announce"
	cCaP      command = "csAps"

	cStop    command = "stop"
	cLLBlock command = "LLblock"
	cLLT     command = "LLT"

	cNewMap command = "newMap"
	cEpochCh command = "epochchange"
	cMigrateWanted command = "migwanted"
	cPendingTXs command = "pendingTXs"
	cNumOfMigratedTXsAddrs     command = "migTXsAddr#"
	cUnchangedState     command = "#ofUnchange"

	cBalanceAndPending command = "balance&TX"
	cSure	command = "sure"
)

//默认前十二位为命令名称
func jointMessage(cmd command, content []byte) []byte {
	b := make([]byte, prefixCMDLength)
	for i, v := range []byte(cmd) {
		b[i] = v
	}
	joint := make([]byte, 0)
	joint = append(b, content...)
	return joint
}

//默认前十二位为命令名称
func splitMessage(message []byte) (cmd string, content []byte) {
	cmdBytes := message[:prefixCMDLength]
	newCMDBytes := make([]byte, 0)
	for _, v := range cmdBytes {
		if v != byte(0) {
			newCMDBytes = append(newCMDBytes, v)
		}
	}
	cmd = string(newCMDBytes)
	content = message[prefixCMDLength:]
	return
}

//对消息详情进行摘要
func getDigest(request *Request) string {
	b, err := json.Marshal(request)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	//进行十六进制字符串编码
	return hex.EncodeToString(hash[:])
}
