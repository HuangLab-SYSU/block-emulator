# Version Updates

## 2023/05/21
Debugs:
1. **Problem**: The function *blockChain.AddAccount* has no effect to the storage.  
- **Reason**: The implement of this func cannot change the merkle root...
- **Solution**: Use "Virtual transactions" to replace the *AddAccount* operation. For a just added account, we consider it as a virtual transaction (whose *Sender* or *Recepient* is "00000000000000"), so that we can use *GetUpdateStatusTrie* to do this operation. 
- **Future**: Solve this problem without "Virtual transactions", because this implementation cost more. 

## 2023/05/19
Debugs: 
1. **Problem**: The init_balances of some accounts are wrong.  
- **Reason**: The init_balances are set by shallow copy ("="), but *Balance* is a pointer (*big.Int)
- **Solution**: Replace the shallow copy with deep copy (with an "adding-zero" form)
