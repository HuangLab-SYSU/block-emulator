# Relay 跨分片交易技术文档

## 一、背景

- Relay 机制来源于 Monoxide ，其采用消息传递的方式，由分片间消息传递进行跨分片交易的处理
- 一笔跨分片交易需要被源分片和目标分片分别进行处理，① 源分片对发送交易方进行状态更新，② 目的分片对交易接收方状态进行更新

## 二、实现

1. 客户端将交易发送到交易发送者所在分片
2. 节点生成区块，如果交易为跨分片交易，则仅更新跨分片交易中发送者的账户状态，并将当前跨分片交易以消息的形式发送到交易接受者所在分片
3. 接受者所在分片主节点检测到跨分片交易之后，将跨分片交易加入到交易池，等待打包上链
4. 交易打包上链之后，完成跨分片交易



### 三、代码

1. 节点交易池维护`RelayPool`结构，负责暂存跨分片交易，并提供`RelayPool`的操作函数
   1. `AddRelayTx(tx *Transaction, shardID uint64)`：负责将交易添加到对应分片的Relay池中
   2. `ClearRelayPool()`：负责清空Relay池
2. 当区块链经过共识，上链之前，检测区块链中是否存在跨分片交易（判断交易接收者是否属于该分片），如果是跨分片交易，仅仅更新交易中分片内所包含的账户的状态，并将当前跨分片交易添加到`RelayPool`。
3. 节点完成所有预提交交易的遍历处理后，根据跨分片交易所属分片ID和`RelayPool`中对应分片的跨分片交易，构建 `message.Relay` 消息，并根据分片ID将消息发送到对应分片的跨分片交易处理节点（PBFT中为主节点）。

```Go
// send relay txs
for sid := uint64(0); sid < rphm.pbftNode.pbftChainConfig.ShardNums; sid++ {
        if sid == rphm.pbftNode.ShardID {
                continue
        }
        relay := message.Relay{
                Txs:           rphm.pbftNode.CurChain.Txpool.RelayPool[sid],
                SenderShardID: rphm.pbftNode.ShardID,
                SenderSeq:     rphm.pbftNode.sequenceID,
        }
        rByte, err := json.Marshal(relay)
        if err != nil {
                log.Panic()
        }
        msg_send := message.MergeMessage(message.CRelay, rByte)
        go networks.TcpDial(msg_send, rphm.pbftNode.ip_nodeTable[sid][0])
        rphm.pbftNode.pl.Plog.Printf("S%dN%d : sended relay txs to %d\n", rphm.pbftNode.ShardID, rphm.pbftNode.NodeID, sid)
}
rphm.pbftNode.CurChain.Txpool.ClearRelayPool()
```

4. 对应分片的跨分片交易处理节点（PBFT中为主节点）检测到 `message.Relay` 消息后，解析消息得到交易数据，然后将解析得到的交易添加到当前节点交易池中。

```Go
// receive relay transaction, which is for cross shard txs
func (rrom *RawRelayOutsideModule) handleRelay(content []byte) {
    relay := new(message.Relay)
    err := json.Unmarshal(content, relay)
    if err != nil {
       log.Panic(err)
    }
    rrom.pbftNode.pl.Plog.Printf("S%dN%d : has received relay txs from shard %d, the senderSeq is %d\n", rrom.pbftNode.ShardID, rrom.pbftNode.NodeID, relay.SenderShardID, relay.SenderSeq)
    rrom.pbftNode.CurChain.Txpool.AddTxs2Pool(relay.Txs)
    rrom.pbftNode.seqMapLock.Lock()
    rrom.pbftNode.seqIDMap[relay.SenderShardID] = relay.SenderSeq
    rrom.pbftNode.seqMapLock.Unlock()
    rrom.pbftNode.pl.Plog.Printf("S%dN%d : has handled relay txs msg\n", rrom.pbftNode.ShardID, rrom.pbftNode.NodeID)
}
```

5. 同2，节点完成交易的正常打包、共识上链，因为跨分片交易判断标准为接收者是否属于该分片，所以接收跨分片交易的节点不会将接收到的跨分片交易再次加入到`RelayPool`中，而是仅更新交易接收者状态，完成跨分片交易处理。