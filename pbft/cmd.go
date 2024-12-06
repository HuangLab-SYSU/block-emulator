package pbft

import (
	"blockEmulator/core"
	"blockEmulator/new_trie"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
)

type Reconfig struct {
	Content map[string]map[string]string
	ID      int
}

// <REQUEST,o,t,c>
type Request struct {
	Message
	Timestamp int64
	//相当于clientID
	// ClientAddr string
}

type Message struct {
	Content []byte
	ID      int
}

// <<PRE-PREPARE,v,n,d>,m>
type PrePrepare struct {
	RequestMessage *Request
	Digest         string
	SequenceID     int
	// Sign           []byte
}

// <PREPARE,v,n,d,i>
type Prepare struct {
	Digest     string
	SequenceID int
	NodeID     string
	// Sign       []byte
}

// <COMMIT,v,n,D(m),i>
type Commit struct {
	Digest     string
	SequenceID int
	NodeID     string
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

type SendAccounts struct {
	Accounts []string
	NodeID   string
}

type ReconfigBlockMessage struct {
	SendBlocks
	Tx_pool *core.Tx_pool
}

type ReconfigTrieMessage struct {
	Trie    *new_trie.N_Trie
	Tx_pool *core.Tx_pool
}

type ReconfigCenterMessage struct {
	NodeTable map[string]string
}

// <REPLY,v,t,c,i,r>
type Reply struct {
	MessageID int
	NodeID    string
	Result    bool
}

type Relay struct {
	Txs     []*core.Transaction
	ShardID string
}

const prefixCMDLength = 12

type command string

const (
	// cRequest    		command = "request"
	cPrePrepare          command = "preprepare"
	cPrepare             command = "prepare"
	cCommit              command = "commit"
	cRequestBlock        command = "requestBlock"
	cSendBlock           command = "sendBlock"
	cReply               command = "reply"
	cReconfig            command = "reconf"
	cReconfigBlock       command = "reconfblock"
	cReconfigTries       command = "reconftries"
	cReconfigCenter      command = "reconfcent"
	cReconfigDone        command = "reconfdone"
	cReconfigDoneReply   command = "reconfrpl"
	cRelay               command = "relay"
	cNextEpochStart      command = "nxtepo"
	cNextEpochStartReply command = "nxteporpl"
	cReconfigStart       command = "reconst"
	cStop                command = "stop"
	cSendAccount         command = "sendacc"
	cHandleSendAccount   command = "hdsendacc"
)

// 默认前十二位为命令名称
func jointMessage(cmd command, content []byte) []byte {
	b := make([]byte, prefixCMDLength)
	for i, v := range []byte(cmd) {
		b[i] = v
	}
	joint := make([]byte, 0)
	joint = append(b, content...)
	return joint
}

// 默认前十二位为命令名称
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

// 对消息详情进行摘要
func getDigest(request *Request) string {
	b, err := json.Marshal(request)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	//进行十六进制字符串编码
	return hex.EncodeToString(hash[:])
}
