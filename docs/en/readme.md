
# Handbook of BlockEmulator

## 1. Environments Configuration

### 1.1 The downloading of codes & datasets.

BlockEmulator dataset: The dataset is from [Ethereum-On-chain data](https://xblock.pro/#/dataset/14), which is published by [xblock.pro](https://xblock.pro/#/).

BlockEmulator code: [github project](https://github.com/HuangLab-SYSU/block-emulator/)

### 1.2 Installing and Configuring BlockEmulator

1. The version of Goland should be at least 1.18.

2. In mainland China, you can use the following commands to configure domestic mirroring:
   
 ```
1 go env -w GO111MODULE=on
2 go env -w GOPROXY=https://goproxy.cn,direct
 ```

## 2. The Architecture  of **BlockEmulator**


### 2.1 [System Architecture](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/System%20Architecture.md)


### 2.2  [Parameter Configuration of **BlockEmulator**](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/ParametersConfig.md)


### 2.3 [The Design of Data Structures](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/dataStructure.md)


### 2.4 [The Workflow of **BlockEmulator**](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/executionFlow.md)


### 2.5 Built-in Protocols/Algorithms
A. [Relay transaction mechanism](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/Relay%20Protocol.md) [NSDI'2019 Monoxide]

[Monoxide: Scale out blockchains with asynchronous consensus zones](https://www.usenix.org/system/files/nsdi19-wang-jiaping.pdf)

B. [BrokerChain protocol](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/BrokerChain.md) [INFOCOM2022 BrokerChain] 

[BrokerChain: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding)

C. Account-partition mechanism, a.k.a [CLPA algorithm](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/CLPA.md) [SRDS'2O22 Achieving]

 [Achieving Scalability and Load Balance across Blockchain Shards for State Sharding](https://www.researchgate.net/publication/364255361_Achieving_Scalability_and_Load_Balance_across_Blockchain_Shards_for_State_Sharding?_sg%5B0%5D=jct9Kh4wJ0gp8kw3vbKR-frMZZQqhPPgEG7PBm7uWFTog-2CO72kgHEI_1urmA1uyLTPpo4ShKXSPbbuyC4n8gci6HhxurEyIfsmstpC.K_dZJY3-a0pIJ9sG-0nbOTzKExJ4fKWWlKwABVI1l2AIUVtDsKV85U7wSYVc_8y2h667U-dif-amFlR2YfDATw)

### 2.7 [The Measuremental Metrics of **BlockEmulator**](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/indicators.md)


### 2.8 [The Log Files of **BlockEmulator**](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/logfile.md)

## Reference
  - [NSDI'2019 Monoxide] J. Wang and H. Wang, “Monoxide: Scale out blockchains with asynchronous consensus zones,” in 16th USENIX Symposium on Networked Systems Design and Implementation (NSDI 19), 2019, pp. 95–112.
  
  - [INFOCOM'2022 BrokerChain] H. Huang, X. Peng, J. Zhan, S. Zhang, Y. Lin, Z. Zheng, and S. Guo, “Brokerchain: A cross-shard blockchain protocol for account/balance-based state sharding,” in Proc. of IEEE Conference on Computer Communications (INFOCOM’22), 2022, pp. 1–10. 
  
  - [SRDS'2022 Achieving] C. Li, H. Huang, Y. Zhao, X. Peng, R. Yang, Z. Zheng, and S. Guo, “Achieving scalability and load balance across blockchain shards for state sharding,” in Proc. of 2022 41st International Symposium on Reliable Distributed Systems (SRDS’22), 2022, pp. 284–294.


