
It is the third section of the second chapter of the **BlockEmulator** English introduction document.

# 1. Basic Data Structure Module 

Core module, as the name implies, is the more critical and fundamental part of the system, which includes the definition of various basic data types and their corresponding functions and usage methods. Here, we have designed four commonly used data structures and their common methods, namely account, block, transaction, and transaction pool.

## A. Account Module

In BlockEmulator, an account is divided into two parts:
The first part is the account itself, including the address and public key of the data account.
```
type Account struct {
   AcAddress utils.Address
   PublicKey []byte
}
```

The second part is the account state, which includes the address of the data account, the nonce value, the balance, the storage root, and the code hash, with the last two values provided separately for use by smart contracts.

```
// AccoutState record the details of an account, it will be saved in status trie
type AccountState struct {
   AcAddress   utils.Address // this part is not useful, abort
   Nonce       uint64
   Balance     *big.Int
   StorageRoot []byte // only for smart contract account
   CodeHash    []byte // only for smart contract account
}
```

The account module also provides corresponding functions, providing the basic functions that an account should have, including increasing the account balance and encoding and decoding the account state.

```
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

## B. Block Module

In BlockEmulator, a block is divided into two parts:
The first part is the block header, which includes the following information: the hash of the parent block, the state root, the transaction root, the number of transactions, the time, and the generator (miner) of the block.

```
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


At the same time, it also includes many related functional functions:
```
// Encode blockHeader for storing further
func (bh *BlockHeader) Encode() []byte

// Decode blockHeader
func DecodeBH(b []byte) *BlockHeader 

// Hash the blockHeader
func (bh *BlockHeader) Hash() []byte

// Print the information of the blockheader
func (bh *BlockHeader) PrintBlockHeader() 


The second part is the block itself, which includes the block header, block body (a collection of packaged transactions), and hash value.

// The definition of block
type Block struct {
        Header *BlockHeader
        Body   []*Transaction
        Hash   []byte
}
```

At the same time, there are many functions related to the block itself, including:

```
// generate a new block
func NewBlock(bh *BlockHeader, bb []*Transaction) *Block 
// print information of the block
func (b *Block) PrintBlock() string 
// Encode Block for storing
func (b *Block) Encode() []byte 
// Decode Block
func DecodeB(b []byte) *Block 
```


## C. Transaction module

In BlockEmulator, a transaction contains seven basic fields, including the initiator and receiver of the transaction, the nonce value of the transaction, signature, transaction value, transaction hash, and the time when the transaction is added to the transaction pool (used to calculate the acceptance time of the transaction).
At the same time, because BlockEmulator has a lot of cross-shard related designs, special functions need to be added to the transaction, such as the Relayed field indicating whether it is accepted by the relay pool, which is also required for the relay cross-shard transaction mechanism, and fields such as OriginalSender, FinalRecipient, and RawTxHash are all special designs required for the BrokerChain cross-shard transaction mechanism.

```
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

For the transaction itself, many related functions are also set up:
```
// Print the information of the transaction
func (tx *Transaction) PrintTx() string 

// Encode transaction for storing
func (tx *Transaction) Encode() []byte 

// Decode transaction
func DecodeTx(to_decode []byte) *Transaction 

// new a transaction
func NewTransaction(sender, recipient string, value *big.Int, nonce uint64) *Transaction
```

## D. Transaction pool module

The transaction pool in BlockEmulator has three fields: a transaction queue, a Relay transaction pool (designed for the Relay mechanism), and a lock (used to ensure consistency in multi-threaded operations).

```
type TxPool struct {
   TxQueue   []*Transaction            // transaction Queue
   RelayPool map[uint64][]*Transaction //designed for sharded blockchain, from Monoxide
   lock      sync.Mutex
}
```

At the same time, BlockEmulator also provides many functions related to the transaction pool.

```
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
