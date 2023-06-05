# Version Updates

## 2023/05/25
1. **New Features**: We have added query function. 
- After the experiment is completed, users can query the account balance now. 
- After the experiment is completed, users can query the information of blocks and transactions on the blockchain now. 

## 2023/05/24
1. Future work of [**2023/05/21 Debugs - 1**](#20230521)
- 2023/05/21: **Future**: We will solve this problem without "Virtual transactions", because this implementation cost more.
- Now this function *blockChain.AddAccount* can operate correctly without invoking *GetUpdateStatusTrie* (the resource cost is reduced). 

## 2023/05/21
Debugs:
1. **Problem**: The function *blockChain.AddAccount* has no effect to the storage.  
- **Reason**: The implementation of this func cannot change the merkle root...
- **Solution**: We use "Virtual transactions" to replace the *AddAccount* operation. For a just added account, we consider it as a virtual transaction (whose *Sender* or *Recipient* is "00000000000000"), so that we can use *GetUpdateStatusTrie* to do this operation. 
- **Future**: We will solve this problem without "Virtual transactions", because this implementation cost more.

2. **Problem**: The transactions migrate to the incorrect shard (in CLPA + Broker mechanism)
- **Reason**: The judge function is incorrect, leading to the wrong behaviors. 
- **Solution**: We have added a new attribute *SenderIsBroker* in broker tx, to identify whether the sender is a broker account (if *HasBroker && !SenderIsBroker* is true, then the *recipient* is a broker account), and modified the *sendAccounts_and_Txs* function in *accountTransfermod_Broker.go* file with this attribute. 

## 2023/05/19
Debugs: 
1. **Problem**: The init_balances of some accounts were wrong.  
- **Reason**: The init_balances were set by shallow copy ("="), but *Balance* is a pointer (*big.Int)
- **Solution**: We have replaced the shallow copy with deep copy (with an "adding-zero" form)

