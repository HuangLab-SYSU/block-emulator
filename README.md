# block-emulator
This is the official handbook of **BlockEmulator**(version 1.0).

# Introduction 

## 1. The Background of BlockEmulator


Initiated by **[HuangLab](http://xintelligence.pro/)**  (a research group in the School of Software Engineering, Sun Yat-sen University, China), ***BlockEmulator** is a blockchain testbed that enables researchers to verify their proposed new protocols and mechanisms. It supports multiple consensus protocols and particularly the cross-shard mechanism.


The main purpose of this testbed is to help users (researchers, students, etc.) quickly verify their own blockchain consensus protocols and blockchain-sharing mechanisms. **BlockEmulator** is designed as an experimental platform that adopts lightweight system architecture.

It simplifies the implementation of industrial-class blockchains since **BlockEmulator** only implements the core functions of a blockchain, including the transaction pool, block packaging, consensus protocols, and on-chain transaction storage. It also supports common consensus protocols, such as Practical Byzantine Fault Tolerance (PBFT) and Proof of Work (PoW).

In particular, BlockEmulator offers the system-level design and implementation for blockchain sharding mechanisms. For example, the cross-shard transaction mechanisms implemented by BlockEmulator include the following two representative solutions, i.e., i) **Relay transaction mechanism** proposed by **Monoxide** (NSDI'2019), and the **BrokerChain** protocol proposed by **BrokerChain** (INFOCOM'2022).

BlockEmulator is oriented toward blockchain researchers because it provides a blockchain experimental platform for quickly implementing their own algorithms, protocols, and mechanisms. It also offers very helpful functions for researchers to help them collect experimental data, facilitating their plotting experimental figures.

This document aims to provide an easy-to-follow user manual of BlockEmulator.

## 2. Related Work

The following papers from HuangLab's publications have adopted **BlockEmulator** as an experimental tool.

1. **BrokerChain**: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding (published at INFOCOM 2022) [PDF](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding)
   
2. Achieving Scalability and Load Balance across Blockchain Shards for State Sharding (published at SRDS 2022) [PDF](https://ieeexplore.ieee.org/document/9996899)
   
3. **tMPT**: Reconfiguration across Blockchain Shards via Trimmed Merkle Patricia Trie (published at IWQoS 2023) [PDF](https://www.researchgate.net/publication/370633426_tMPT_Reconfiguration_across_Blockchain_Shards_via_Trimmed_Merkle_Patricia_Trie)
   
4. **MVCom**: Scheduling Most Valuable Committees for the Large-Scale Sharded Blockchain (published at ICDCS 2021) [PDF](https://ieeexplore.ieee.org/document/9546408)


## 3. Highlights of BlockEmulator:

1. **Lightweight**: **BlockEmulator** is a lightweight testbed platform for blockchain experiments.

2. **Fast Setup**: **BlockEmulator** enables users to set up their environments quickly and supports remote deployment on Cloud.

3. **Customization**: **BlockEmulator** is implemented by using Goland language, which supports users' customization and modification.

4. **Apt to Experiments**: **BlockEmulator** supports the playback of historical transactions of mainstream blockchains (such as Ethereum), as well as automatic output experimental log files. From those log files, researchers can interpret plenty of metrics such as system throughput, confirmation latency of transactions, the queueing of the transaction pool, etc. These features are very helpful for researchers and students to facilitate their experimental data collection and plotting of experimental charts.

## 4. Getting Started to Use BlockEmulator

Quickly get started with BlockEmulator through the following document. Please refer to::
[BlockEmulator help document](https://j4s9dl19cd.feishu.cn/docx/AEXIdlSJAob2Y4xHp2PcquE0nvd)



# Handbook of BlockEmulator

## 1. Environments Configuration

### 1.1 The downloading of codes & datasets.

BlockEmulator dataset: The dataset is from [Ethereum-On-chain data](https://xblock.pro/#/dataset/14), which is published by [xblock.pro](https://xblock.pro/#/).

BlockEmulator code: [github project](https://github.com/ChenQinde/block-emulator/tree/bc)

### 1.2 Installing and Configuring BlockEmulator

1. The version of Goland should be at least 1.18.

2. In mainland China, you can use the following commands to configure domestic mirroring:
   
 ```
1 go env -w GO111MODULE=on
2 go env -w GOPROXY=https://goproxy.cn,direct
 ```

## 2. The Architecture  of **BlockEmulator**


### 2.1 System Architecture


### 2.2  Parameter Configuration of **BlockEmulator**


### 2.3 The Design of Data Structures


### 2.4 The Workflow of **BlockEmulator**


### 2.5 Built-in Protocols/Algorithms
A. Relay transaction mechanism [NSDI'2019 Monoxide]

[Monoxide: Scale out blockchains with asynchronous consensus zones](https://www.usenix.org/system/files/nsdi19-wang-jiaping.pdf)

B. BrokerChain protocol [INFOCOM2022 BrokerChain] 

[BrokerChain: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding)

C. Account-partition mechanism, a.k.a CLPA algorithm [SRDS'2O22 Achieving]

 [Achieving Scalability and Load Balance across Blockchain Shards for State Sharding](https://www.researchgate.net/publication/364255361_Achieving_Scalability_and_Load_Balance_across_Blockchain_Shards_for_State_Sharding?_sg%5B0%5D=jct9Kh4wJ0gp8kw3vbKR-frMZZQqhPPgEG7PBm7uWFTog-2CO72kgHEI_1urmA1uyLTPpo4ShKXSPbbuyC4n8gci6HhxurEyIfsmstpC.K_dZJY3-a0pIJ9sG-0nbOTzKExJ4fKWWlKwABVI1l2AIUVtDsKV85U7wSYVc_8y2h667U-dif-amFlR2YfDATw)

### 2.7 The Measuremental Metrics of **BlockEmulator**


### 2.8 The Log Files of **BlockEmulator**

## Reference
  - [NSDI'2019 Monoxide] J. Wang and H. Wang, “Monoxide: Scale out blockchains with asynchronous consensus zones,” in 16th USENIX Symposium on Networked Systems Design and Implementation (NSDI 19), 2019, pp. 95–112.
  
  - [INFOCOM'2022 BrokerChain] H. Huang, X. Peng, J. Zhan, S. Zhang, Y. Lin, Z. Zheng, and S. Guo, “Brokerchain: A cross-shard blockchain protocol for account/balance-based state sharding,” in Proc. of IEEE Conference on Computer Communications (INFOCOM’22), 2022, pp. 1–10. 
  
  - [SRDS'2022 Achieving] C. Li, H. Huang, Y. Zhao, X. Peng, R. Yang, Z. Zheng, and S. Guo, “Achieving scalability and load balance across blockchain shards for state sharding,” in Proc. of 2022 41st International Symposium on Reliable Distributed Systems (SRDS’22), 2022, pp. 284–294.


# Contributors Page



   


