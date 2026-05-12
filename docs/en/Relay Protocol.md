
# 一、Background

- The Relay mechanism is derived from Monoxide, which uses messaging to process cross-shard transactions by inter-shard messaging.
  
- A cross-shard transaction needs to be processed separately by the original shard and the destination shard, with the source shard updating the status of the sending transaction party and the destination shard updating the status of the transaction receiver

# 二、Implementation
- The client first sends the transaction to the shard where the sender of the transaction is located.
  
- The consensus node packages and broadcasts the transaction, updates the status tree before the transaction is on the chain, and only updates the account status (sender) of the current shard in the cross-shard transaction if the transaction is a cross-shard transaction, and sends the current cross-shard transaction in the form of a message to the shard that is not in the current shard account (receiver) in the transaction.
  
- After the shard master node of the recipient detects the cross-shard transaction, it will add the cross-shard transaction to the transaction pool and wait for the package to be put on the chain.
  
- After the transaction is packaged and put on the chain, complete the cross-shard transaction.


# 三、Code


## 1. The node transaction pool maintains **RelayPool** structure and Stage the cross-shard transaction

a. AddRelayTx(tx *Transaction, shardID uint64): Responsible for adding transactions to the **Relay** pool of the corresponding shard.

b. ClearRelayPool(): Responsible for clearing **Relay** pool.


## 2. Before going to the chain
When the blockchain is reached by consensus, before going to the chain, detect whether there is a cross-shard transaction in the blockchain (determine whether the recipient of the transaction belongs to the shard), if it is a cross-shard transaction, only update the status of the accounts contained in the shard in the transaction, and add the current cross-shard transaction to RelayPool.

## 3. Traversal processing
After the node completes the traversal processing of all pre-committed transactions, it builds a message based on the shard ID of the cross-shard transaction and the cross-shard transaction of the corresponding shard in RelayPool. Relay the message and send the message to the cross-shard transaction processing node (master node in PBFT) of the corresponding shard according to the shard ID.

```
1 // send relay txs
2 for sid := uint64(0); sid < rphm.pbftNode.3 pbftChainConfig.ShardNums; sid++ {
3        if sid == rphm.pbftNode.ShardID {
4                continue
5        }
6       relay := message.Relay{
7                Txs:           rphm. pbftNode.CurChain.Txpool.RelayPool[sid],
8               SenderShardID: rphm.  pbftNode.ShardID,
9               SenderSeq:     rphm. pbftNode.sequenceID,
10        }
11       rByte, err := json.Marshal(relay)
12       if err != nil {
13                log.Panic()
14        }
15        msg_send := message.MergeMessage (message.CRelay, rByte)
16        go networks.TcpDial(msg_send, rphm.pbftNode.ip_nodeTable[sid][0])
17        rphm.pbftNode.pl.Plog.Printf("S%dN%d : sended relay txs to %d\n", rphm.pbftNode.ShardID, rphm.pbftNode.NodeID, sid)
18 }
19 rphm.pbftNode.CurChain.Txpool.ClearRelayPool()
```
## 4. Detection and Parse message

Once Cross-shard transaction processing node detected **message.Relay**, it will Parse the message and add the corresponding cross-shard transaction to the current transaction pool.

```
1 // receive relay transaction, which is for cross shard txs
2 func (rrom *RawRelayOutsideModule) handleRelay(content []byte) {
3    relay := new(message.Relay)
4    err := json.Unmarshal(content, relay)
5    if err != nil {
6       log.Panic(err)
7    }
8    rrom.pbftNode.pl.Plog.Printf("S%dN%d : has received relay txs from shard %d, the senderSeq is %d\n", rrom.pbftNode.ShardID, rrom.pbftNode.NodeID, relay.SenderShardID, relay.SenderSeq)
9    rrom.pbftNode.CurChain.Txpool.AddTxs2Pool(relay.Txs)
10    rrom.pbftNode.seqMapLock.Lock()
11   rrom.pbftNode.seqIDMap[relay.SenderShardID] = relay.SenderSeq
12    rrom.pbftNode.seqMapLock.Unlock()
13    rrom.pbftNode.pl.Plog.Printf("S%dN%d : has handled relay txs msg\n", rrom.pbftNode.ShardID, rrom.pbftNode.NodeID)
14 }
```
## 5. same to step 2

The node completes the normal packaging of the transaction and the consensus is uploaded to the chain, because the judgment standard of cross-shard transaction is whether the recipient belongs to the shard, so the node receiving the cross-shard transaction will not add the received cross-shard transaction to the RelayPool again, but only update the transaction receiver status to complete the cross-shard transaction processing.

