# 数据结构设计

## 1.基础数据结构模块。

Core 模块，顾名思义，就是系统中比较核心又很基础的部分，其包含常见的各类基本数据的定义以及其对应的函数和使用方法。在此处，我们设计了四个常用的数据结构及其常用方法，分别为账户、区块、交易和交易池。

### A. 账户模块

账户在 BlockEmulator 中，分为两个部分：

第一个部分是账户本身，包含数据账户的地址以及公钥。

```Go
type Account struct {
   AcAddress utils.Address
   PublicKey []byte
}
```

第二个部分是账户状态，包含数账户的地址、Nonce 值、余额、存储根以及代码哈希，其中后两个值是单独提供给智能合约使用的。

```Go
// AccoutState record the details of an account, it will be saved in status trie
type AccountState struct {
   AcAddress   utils.Address // this part is not useful, abort
   Nonce       uint64
   Balance     *big.Int
   StorageRoot []byte // only for smart contract account
   CodeHash    []byte // only for smart contract account
}
```

同时账户模块也提供对应的函数，提供了账户该有的基础功能，包含账户余额的增加、账户状态的编解码等方法。

```Go
// Reduce the balance of an account
func (as *AccountState) Deduct(val *big.Int) bool

// Increase the balance of an account
func (s *AccountState) Deposit(value *big.Int)

// Encode AccountState in order to store in the MPT
func (as *AccountState) Encode() []byte

// Decode AccountState
func DecodeAS(b []byte) *AccountState

// Hash AccountState for computing the MPT Root
func (as *AccountState) Hash() []byte
```

### B. 区块模块

区块在 BlockEmulator 中，分为两个部分：

第一个部分是区块头，包含如下的信息，父区块的哈希、状态根、交易根、交易的数量、时间以及区块的生成者（矿工）。

```Go
// The definition of blockheader
type BlockHeader struct {
        ParentBlockHash []byte
        StateRoot       []byte
        TxRoot          []byte
        Number          uint64
        Time            time.Time
        Miner           uint64
}
```

同时，也包含了不少相关的功能函数:

```Go
// Encode blockHeader for storing further
func (bh *BlockHeader) Encode() []byte

// Decode blockHeader
func DecodeBH(b []byte) *BlockHeader 

// Hash the blockHeader
func (bh *BlockHeader) Hash() []byte

// Print the information of the blockheader
func (bh *BlockHeader) PrintBlockHeader() 
```

第二个部分是区块本身，包含区块头、区块体（被打包的交易合集）以及哈希值。

```Go
// The definition of block
type Block struct {
        Header *BlockHeader
        Body   []*Transaction
        Hash   []byte
}
```

同时，对于区块本身，也提供了不少与之相关的功能函数：

```Go
// generate a new block
func NewBlock(bh *BlockHeader, bb []*Transaction) *Block 

// print information of the block
func (b *Block) PrintBlock() string 

// Encode Block for storing
func (b *Block) Encode() []byte 

// Decode Block
func DecodeB(b []byte) *Block 
```

### C. 交易模块

交易在 BlockEmulator 中，包含了 7 个基本的字段，包含交易的发起者、交易的接受者、交易的 Nonce 值、签名、交易的值、交易的哈希以及交易加入到交易池的时间（用于统计交易的接受时间）。

同时因为在 BlockEmulator 中预设了不少跨分片相关的设计，需要在交易中增加特别的功能，如 Relayed 字段表示是否被 relay 池所接收，这也是 relay 跨分片交易机制所需要的功能，如字段 OriginalSender 、 字段FinalRecipient 以及字段 RawTxHash 都是用于 BrokerChain 跨分片交易机制所需要的特别设计。

```Go
// Definition of transaction
type Transaction struct {
   Sender    utils.Address
   Recipient utils.Address
   Nonce     uint64
   Signature []byte 
   Value     *big.Int
   TxHash    []byte
   
   // the time adding in pool
   Time time.Time 

   // used in transaction relaying
   Relayed bool
   
   // used in broker, if the tx is not a broker1 or broker2 tx, these values should be empty.
   OriginalSender utils.Address
   FinalRecipient utils.Address
   RawTxHash      []byte
}
```

对于交易本身，也设置了不少相关的功能函数：

```Go
// Print the information of the transaction
func (tx *Transaction) PrintTx() string 

// Encode transaction for storing
func (tx *Transaction) Encode() []byte 

// Decode transaction
func DecodeTx(to_decode []byte) *Transaction 

// new a transaction
func NewTransaction(sender, recipient string, value *big.Int, nonce uint64) *Transaction
```

### D. 交易池模块

交易池在 BlockEmulator 中，包含三个字段，交易队列、Relay 交易池（为 Relay 机制所设计）以及锁（用以保证多线程操作保持一致）。

```Go
type TxPool struct {
   TxQueue   []*Transaction            // transaction Queue
   RelayPool map[uint64][]*Transaction //designed for sharded blockchain, from Monoxide
   lock      sync.Mutex
}
```

同时，在 BlockEmulator 中也提供了不少交易池的相关功能函数。

```Go
// Generate 
func NewTxPool() *TxPool

// Add a transaction to the pool (consider the queue only)
func (txpool *TxPool) AddTx2Pool(tx *Transaction) 

// Add a list of transactions to the pool
func (txpool *TxPool) AddTxs2Pool(txs []*Transaction)

// add transactions into the pool head
func (txpool *TxPool) AddTxs2Pool_Head(tx []*Transaction)

// Pack transactions for a proposal
func (txpool *TxPool) PackTxs(max_txs uint64) []*Transaction 

// Relay transactions
func (txpool *TxPool) AddRelayTx(tx *Transaction, shardID uint64) 

// txpool get locked
func (txpool *TxPool) GetLocked() 

// txpool get unlocked
func (txpool *TxPool) GetUnlocked()

// get the length of tx queue
func (txpool *TxPool) GetTxQueueLen() int 

// get the length of ClearRelayPool
func (txpool *TxPool) ClearRelayPool() 

// abort ! Pack relay transactions from relay pool
func (txpool *TxPool) PackRelayTxs(shardID, minRelaySize, maxRelaySize uint64) ([]*Transaction, bool)

// abort ! Transfer transactions when re-sharding
func (txpool *TxPool) TransferTxs(addr utils.Address) []*Transaction 
```