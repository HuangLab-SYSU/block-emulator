# Version Updates

## 2023/05/19
Debugs: 
1. **Problem**: The init_balances of some accounts are wrong.  
- **Reason**: The init_balances are set by shallow copy ("="), but *Balance* is a pointer (*big.Int)
- **Solution**: Replace the shallow copy with deep copy (with an "adding-zero" form)
