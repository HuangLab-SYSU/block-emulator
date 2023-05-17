
**BrokerChain** cross-sharding mechanism, from the paper "BrokerChain: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding" included in INFOCOM2022. The paper proposes the BrokerChain cross-shard transaction mechanism for the problem of a large number of cross-distribution transactions in the current sharded blockchain system. That is, the state of the account selected as the broker is split so that it exists in each shard, and when there is a cross-shard transaction between the shards, the broker can be used to handle it. As illustrated in the following figure:

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/broker.png" width=700>

<center> A toy example of Broker</center>

Broker account is an account using state segmentation technology, so that each shard in the existence of a broker. When a user submits a cross-shard transaction, such as account A in one shard that initiates a transaction to account B in another shard, then this cross-shard transaction can be split into two intra-shard transactions, one between account A and Broker account C in the same shard, and another between account B and Broker account C in the same shard. When the transaction is completed, implement cross-shard transaction processing to reduce the number of cross-shard transactions. 

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/transactions.png" width=700>

 The overall process for a cross-shard transaction in BrokerChain involves five steps:

1. The sender account A sends the raw transaction message $\theta_{raw}$ to the Broker account $C$.
   
2. Broker account $C$ receives the raw message $\theta_{raw}$ and sends $\theta_{1}$ ​to Shard1 and Shard2, respectively, where the sender account and the receiver account reside.
   
3. Shard1 node receives  $\theta_{1}$, validates the transaction's legality, and constructs the transaction between accounts A and C. If the transaction is successfully committed to the blockchain, Shard1 sends a Confirm  $\theta_{1}$  message to the Broker account $C$.
   
4. Broker account $C$ receives the Confirm $\theta_{1}$  message, verifies its validity, and sends $\theta_{2}$  to Shard2.
   
5. Shard2 node receives $\theta_{2}$, validates the transaction's legality, constructs the transaction between accounts $C$ and $B$, and sends a Confirm
$\theta_{2}$ message to the Broker account $C$ to complete the cross-shard transaction.

# Broker design

## Architecture Design

- **Sender**: The sender is assumed by the supervisor.
  - Responsible for detecting whether the received transaction is a cross-shard transaction, and if it is a cross-shard transaction, generate $θ_{raw}$ and send it to the broker.
  
  - Responsible for processing $θ_{1}$from the broker, the main steps are validation messages, generating tx1 on-chip transactions, and adding transaction pools.
  
  - Responsible for processing $θ_{2}$ from the broker, the main steps are verifying the message, generating tx2 on-chip transactions, and adding a transaction pool.
  
  - When the new block is packaged on the chain, it is responsible for filtering BrokerTx, and generating $Confirm θ_{1}$ and $Confirm θ_{2}$to send to the Broker client respectively - currently Sender plays the role of miner nodes in the blockchain (such as the master node in PBFT, responsible for receiving transactions from the client and sending messages to the slave node and the broker client).

- **Broker**:  The Broker is assumed by the supervisor and it is in charged of  the request of cross-shard transaction.
  
  - **BrokerAccount**: In the sharded blockchain network, in order to facilitate the division of blockchain accounts, according to the number of shards, the broker will have a corresponding number of accounts, each account corresponds to a shard, and the brokerAccount does not participate in state migration
    - If A → B, A in 1 shard, B in 2 shards, after the broker detects the cross-shard transaction, after a series of information exchanges, the cross-shard transaction becomes A → BrokerAccount and BrokerAccount → B .

  

## Implementation and Design

### Detailed Explaination


The implementation of **BrokerChain** can be splited into two parts: The implementation of the broker itself and the implementation of the broker-related message structure.

1.  The implementation of the broker itself

    The implementation of the broker itself, including the information that the broker should have, and some of its related variables, are used to store the information generated during processing. 

    ```
    1 type Broker struct {
    2   BrokerRawMegs  map[string]*message.BrokerRawMeg
    3    ChainConfig    *params.ChainConfig
    4    BrokerAddress  []string
    5    RawTx2BrokerTx map[string][]string
    6 }
    ```

    | Variance|Type|Explaination|
    |----|-------|--------|
    |BrokerRawMegs|map[string]*BrokerRawMeg |The abstract of RawMeg and the map of *BrokerRwaMeg
    |ChainConfig|*params.ChainConfig|The configuration information of system|
    |BrokerAddress|[]string|Broker address|
    |rawTx2BrokerTx|map[string][]string|Mapping of the original cross-shard transaction and the two intra-slice transactions produced by the broker.|

    Additional mehtods:

    ```
    1 // generate a new broker
    2 func NewBroker(pcc *params.Cha in config)
    3 
    4 // get the digest of rawMeg
    5 func getBrokerRawMagDigest(r *message.BrokerRawMeg) []byte 
    6
    7 // Get parition (if not exist, return default)
    8 func fetchModifiedMap(key string) uint64 
    9
    10 // Handle the raw messsage 
    11 func handleBrokerRawMag(brokerRawMags []*message.BrokerRawMeg) 
    12
    13 // Handle the tx1 
    14 func handleTx1ConfirmMag(mag1confirms []*message.Mag1Confirm) 
    15
    16 // Handle the tx2
    17 func handleTx2ConfirmMag(mag2confirms []*message.Mag2Confirm)
    18 //init broker address
    19 func initBrokerAddr(num int) []string  
    ```
2. Message Structure
   
    $\theta_raw$: dentoes the original message that the sender sends to Broker.
    ```
    1 type BrokerRawMeg struct {
    2 Tx        *core.Transaction
    3 Broker    utils.Address
    4 Hlock     uint64 //ignore
    5 Snonce    uint64 //ignore
    6 Bnonce    uint64 //ignore
    7 Signature []byte // not implemented now.
    8 }
    ```

    $\theta_1$:The broker sends to the shard miner node to which the transaction initiator account belongs.
    ```
    1 type BrokerType1Meg struct {
    2 RawMeg   *BrokerRawMeg
    3 Hcurrent uint64        //ignore
    4 Signature []byte // not implemented now.
    5 Broker   utils.Address 
    6 }
    ```
    $Confirm \theta_1$: The shard miner node to which the transaction initiator account belongs will send Tx1 to the broker after it is uploaded to the chain.

    $\theta_{2}$: Broker original transaction receiver shard miner node sending

    ```
    1 type BrokerType2Meg struct {
    2 RawMeg *BrokerRawMeg
    3 Signature []byte // not implemented now.
    4 Broker utils.Address 
    5  }
    ```

    $Confirm \theta_{2}$: The original transaction receiver shard miner node, after uploading Tx2 to the chain, sends it to the broker

    ```
    type mag2Confirm struct {
    Tx2Hash []byte
    RawMeg  *BrokerRawMeg
    }
    ```
    **Methods**

    |handleBrokerRawMag(brokerRawMags []*message.BrokerRawMeg) |Process the raw message from the miner node, generate and send $\theta_1$|
    |-----|---------|
    |handleTx1ConfirmMag(mag1confirms []*message.Mag1Confirm) |Process $Confirm \theta_{1}$ from the miner node, send $\theta_{2}$|
    |handleTx2ConfirmMag(mag2confirms []*message.Mag2Confirm) |Process $Confirm \theta_{1}$ from the miner node, record the results|
    |fetchModifiedMap(key string) uint64 | Used to judge which shard the transaction account belongs to|

    **Structure**

    The Broker part is mainly a staging pool of $Confirm θ_{1}$ and $Confirm θ_{2}$ messages, because the two on-chip transactions generated by the broker are generated at the time of transaction injection, while the Confirm message needs to be processed after the transaction is on the chain

    ```
    1 brokerConfirm1Pool = make(map[string]*mag1Confirm)
    2 brokerConfirm2Pool = make(map[string]*mag2Confirm)
    ```

    |dealTxByBroker(txs []*core.Transaction) (itxs []*core.Transaction)|Broker transactions are used for cross-shard transaction processing, the input is a generation injection transaction, the output is an on-chip transaction, and a RawMeg is generated and sent to the broker|
    |----|--------|
    |handleBrokerType1Mes(brokerType1Megs []*message.BrokerType1Meg)  | Process $θ_{1}$ from the broker, parse the message and generate the on-chip transaction Tx1, and join the transaction pool |
    |handleBrokerType2Mes(brokerType2Megs []*message.BrokerType2Meg)  |Process $θ_{2}$ from the broker, parse the message and generate the on-chip transaction Tx1, and join the transaction pool |
    |createConfirm(txs []*core.Transaction)|After On-chain, the tranasction records respectively generate $Confirmθ_{1}$ and $Confirmθ_{2}$ that shall be processed by Broker | 




