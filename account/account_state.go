package account

import (
	"blockEmulator/params"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"log"
	"math/big"
	"strconv"
	"sync"
)

// StateAccount is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type AccountState struct {
	// Nonce    uint64
	Balance *big.Int
	// Root     []byte // merkle root of the storage trie
	// CodeHash []byte
	Migrate int
	Location int
}

var Account2ShardLock sync.Mutex

//  账户到分片的映射
var Account2Shard map[string]int

//  本分片的账户
var AccountInOwnShard map[string]bool

var BalanceBeforeOutLock sync.Mutex

//  正在迁出的账户在Out1时的余额
var BalanceBeforeOut map[string]*big.Int

var ComingAccountLock sync.Mutex


var Outing_Acc_Before_Announce_Lock sync.Mutex

//  正在迁出且已经还没到Announce的账户，交易池中这类账户发起的交易不会被打包，而是到专门的内存Outing_TX中
var Outing_Acc_Before_Announce map[string]bool


var Outing_Acc_After_Announce_Lock sync.Mutex

//  正在迁出且已经收到Announce的账户，这类账户的交易不会进入交易池，而是到专门的内存Outing_TX中
var Outing_Acc_After_Announce map[string]bool

var Lock_Acc_Lock sync.Mutex

//  迁出账户要对账户锁
var Lock_Acc map[string]bool

//  根据账户地址的出所在分片。若是旧账户
func Addr2Shard(senderAddr string) int {
	Account2ShardLock.Lock()
	if shardID, ok := Account2Shard[senderAddr]; ok {
		Account2ShardLock.Unlock()
		return shardID
	}
	// 只取地址后五位已绝对够用
	senderAddrlast := senderAddr[len(senderAddr)-5:]
	num, err := strconv.ParseInt(senderAddrlast, 16, 32)
	// num, err := strconv.ParseInt(senderAddrlast, 10, 32)
	if err != nil {
		log.Panic()
	}

	shardID := int(num) % params.Config.Shard_num
	Account2Shard[senderAddr] = shardID
	if shardID == params.ShardTable[params.Config.ShardID] {
		AccountInOwnShard[senderAddr] = true
	}
	Account2ShardLock.Unlock()
	return shardID
}

func (s *AccountState) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(s)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DecodeAccountState(to_decode []byte) *AccountState {
	var state AccountState

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&state)
	if err != nil {
		log.Panic(err)
	}

	return &state
}

func (s *AccountState) Hash() []byte {
	hash := sha256.Sum256(s.Encode())
	return hash[:]
}

// // 从账户上扣钱
// func (s *AccountState) Deduct(value float64) {
// 	// todo 判断判断
// 	s.Balance -= value
// }

// // 往账户上打钱
// func (s *AccountState) Deposit(value float64) {
// 	// todo 判断判断
// 	s.Balance += value
// }

func Addr2ShardDeepCopy(dst, src map[string]int) error {
	if tmp, err := json.Marshal(&src); err != nil {
		return err
	} else {
		err = json.Unmarshal(tmp, &dst)
		return err
	}
}
