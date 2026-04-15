# BlockEmulator Handbook
This is the official handbook of **BlockEmulator** (version 1.0).
We also provide a **[Chinese-version handbook](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/ch/readme_ChnVersion.md)**.

On Dec. 31, 2024, we uploaded a detailed 139-page **Chinese-version User Manual**, named "_2024Dec31-(139页)使用指南-黄华威.pdf_". Please feel free to download it from the main folder.

# Issue Log
View the issue logs here:
**[Issue logs](./docs/issueLogRoot.md)**

# Version Updates
View the updates of **BlockEmulator** here: 
**[ Version Updates ](./docs/versionUpdate.md)**


# FAQ (Frequently Asked Questions)
See **[FAQ](./docs/FAQ.md)**. 

# ========================

# The Advanced Version of BlockEmulator (BlockEmulator-X)
See **[ BlockEmulator-X ](https://github.com/HuangLab-SYSU/block-emulator-x)**.


# ========================


# Introduction 

## 1. Background of BlockEmulator


Initiated by **[HuangLab](http://xintelligence.pro/)**  (a research group in the School of Software Engineering, Sun Yat-sen University, China), ***BlockEmulator** is a blockchain testbed that enables researchers to verify their proposed new protocols and mechanisms. It supports multiple consensus protocols, particularly the blockchain sharding mechanism.


The primary purpose of this testbed is to help users (researchers, students, etc.) quickly verify their own blockchain consensus protocols and blockchain-sharding protocols. 

**BlockEmulator** is designed as an experimental platform that adopts a lightweight system architecture. It simplifies the implementation of industrial-grade blockchains because **BlockEmulator** implements only the core functions of a blockchain, including the transaction pool, block packaging, consensus protocols, and on-chain transaction storage. It also supports standard consensus protocols, such as Practical Byzantine Fault Tolerance (PBFT).

In particular, BlockEmulator offers the system-level design and implementation for blockchain-sharding mechanisms. For example, the cross-shard transaction mechanisms implemented by BlockEmulator include the following two representative solutions: i) **Relay transaction mechanism** proposed by **Monoxide** (NSDI'2019), and ii) the **BrokerChain** protocol (INFOCOM'2022) [PDF](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding).

BlockEmulator is oriented toward blockchain researchers. It provides a blockchain experimental platform for quickly implementing their own algorithms, protocols, and mechanisms. It also offers invaluable functions for collecting experimental data, facilitating the plotting of experimental figures.


## 2. Official Technical Paper & Citation

To provide an official handbook for BlockEmulator, we have written a technical paper titled "BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols" [arXiv page](https://arxiv.org/abs/2311.03612). Please cite our **TSC**-version paper if you use BlockEmulator as an experiment tool in your own paper, using the following **bib data**:

```
@article{huang2025blockemulator,
   title={BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols},
   author={Huang, Huawei and Ye, Guang and Yang, Qinglin and Chen, Qinde and Yin, Zhaokang and Luo, Xiaofei and Lin, Jianru and Zheng, Jian and Li, Taotao and  Zheng, Zibin},
   journal = {IEEE Transactions on Services Computing (TSC)},
   volume={18},
   number={2},
   pages = {690--703},
   year = {2025},
   }
```


## 3. Related Work

The following papers from HuangLab's publications have adopted **BlockEmulator** as an experimental tool.

1. **BrokerChain**: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding **(INFOCOM 2022)** 【[PDF](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding)】
2. **BrokerChain-Ext**: BrokerChain: A Blockchain Sharding Protocol by Exploiting Broker Accounts **(ToN 2025)** 【[PDF](https://www.researchgate.net/publication/390218703_BrokerChain_A_Blockchain_Sharding_Protocol_by_Exploiting_Broker_Accounts)】
3. **ShardCutter**: ShardCutter: A Blockchain Sharding Protocol Achieving Transaction Workload Balance Across State Shards **(ToN 2026)** 【[PDF](https://www.researchgate.net/publication/400699615_ShardCutter_A_Blockchain_Sharding_Protocol_achieving_Transaction_Workload_Balance_across_State_Shards)】
4. **Broker2Earn**: Towards Maximizing Broker Revenue and System Liquidity for Sharded Blockchains **(INFOCOM 2024)** 【[PDF](https://www.researchgate.net/publication/379213048_Broker2Earn_Towards_Maximizing_Broker_Revenue_and_System_Liquidity_for_Sharded_Blockchains)】
5. **LiquidityPool**: LiquidityPool: Game-Theoretic Analysis of Stakeholder Revenue in Ranking-Dependent DeFi **(WWW 2026)** 【[PDF](https://www.researchgate.net/publication/400068018_LiquidityPool_Game-Theoretic_Analysis_of_Stakeholder_Revenue_in_Ranking-Dependent_DeFi)】
6. Account Migration across Blockchain Shards using Fine-tuned Lock Mechanism **(INFOCOM 2024)** 【[PDF](https://www.researchgate.net/publication/379210418_Account_Migration_across_Blockchain_Shards_using_Fine-tuned_Lock_Mechanism)】
7. **Justitia**: An Incentive Mechanism towards the Fairness of Cross-shard Transactions **(INFOCOM 2025)** 【[PDF](http://xintelligence.pro/archives/1371)】
8. **MVCom**: Scheduling Most Valuable Committees for the Sharded Blockchain **(ToN 2023)** 【[PDF](https://www.researchgate.net/publication/370671128_Scheduling_Most_Valuable_Committees_for_the_Sharded_Blockchain)】
9. **CLPA**: Achieving Scalability and Load Balance across Blockchain Shards for State Sharding (published at SRDS 2022) [PDF](https://ieeexplore.ieee.org/document/9996899)
10. **tMPT**: Reconfiguration across Blockchain Shards via Trimmed Merkle Patricia Trie (published at IWQoS 2023) [PDF](https://www.researchgate.net/publication/370633426_tMPT_Reconfiguration_across_Blockchain_Shards_via_Trimmed_Merkle_Patricia_Trie)


## 4. Highlights of BlockEmulator

1. **Lightweight**. BlockEmulator is a lightweight testbed platform for blockchain experiments.
2. **Fast Configuration**. BlockEmulator enables users to set up their environments quickly and supports remote deployment in the Cloud.
3. **Customization**. BlockEmulator is implemented in GoLand, a language that supports user customization and modification.
4. **Easy to Conduct Experiments**. BlockEmulator supports replaying historical transactions from mainstream blockchains (such as Ethereum). It can automatically yield experimental log files. Using those log files, researchers can interpret metrics such as system throughput, transaction confirmation latency, and queueing in the transaction pool. This function is handy for researchers and students to facilitate their experimental data collection and plotting of experimental charts.

## 5. Getting Started to Use BlockEmulator

Could you quickly get started with BlockEmulator through the following document? Please refer to:
[BlockEmulator help document](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/Help%20Documents.md). More details of the user guidebook can be found at [Handbook of BlockEmulator (Eng.)](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/readme.md) or [Handbook of BlockEmulator (Chn.)](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/ch/readme_ChnVersion.md).

## Reference
- [INFOCOM'2022 BrokerChain] H. Huang, X. Peng, J. Zhan, S. Zhang, Y. Lin, Z. Zheng, and S. Guo, “Brokerchain: A cross-shard blockchain protocol for account/balance-based state sharding,” in Proc. of IEEE Conference on Computer Communications (INFOCOM’22), 2022, pp. 1–10. 
- [NSDI'2019 Monoxide] J. Wang and H. Wang, “Monoxide: Scale out blockchains with asynchronous consensus zones,” in 16th USENIX Symposium on Networked Systems Design and Implementation (NSDI 19), 2019, pp. 95–112.
- [SRDS'2022 Achieving] C. Li, H. Huang, Y. Zhao, X. Peng, R. Yang, Z. Zheng, and S. Guo, “Achieving scalability and load balance across blockchain shards for state sharding,” in Proc. of 2022 41st International Symposium on Reliable Distributed Systems (SRDS’22), 2022, pp. 284–294.


# ========================

# Contributors

## Huawei Huang

Huawei Huang is a professor (full) at Sun Yat-sen University. He is a senior member of IEEE. He received his Ph.D. in Computer Science and Engineering from the University of Aizu, Japan. His research interests include blockchain and distributed computing/protocols. He has published more than 150 papers in top journals and conferences, including IEEE/ACM ToN, JSAC, TPDS, TDSC, TMC, INFOCOM, WWW, ICDCS, SRDS, and IWQoS. He is a PI or Co-PI on several blockchain-related research projects funded by the National Key Research & Development Program of China, the National Natural Science Foundation of China (NSFC), and other funding agencies. He has published 3 blockchain monographs with Springer: "Blockchain Sharding: Theory and Practice", "From Blockchain to Web3 & Metaverse", and "Blockchain Scalability".

## Guang Ye

Guang Ye is currently a Master's student of the School of Software and Engineering at Sun Yat-Sen University. His research interests mainly include Blockchain.  Since joining HuangLab as an undergraduate intern in August 2021, he has participated in the development of BlockEmulator. He redesigned and implemented an advanced version of BlockEmulator -- BlockEmulator-X (coming soon).

## Zhaokang Yin

Zhaokang Yin is currently a Master's student of the School of Software and Engineering at Sun Yat-Sen University. His research interests mainly focus on blockchain. In October 2022, he joined HuangLab to help develop BlockEmulator.

## Jianbo Xiong

Jianbo Xiong is currently a Master's student of the School of Software and Engineering at Sun Yat-Sen University. His research interests mainly focus on blockchain. In September 2025, he joined HuangLab to help develop BlockEmulator and maintain BrokerChain Testnet.

## Jianru Lin

Jianru Lin is a research scientist and engineer in HuangLab, with extensive experience in designing and implementing decentralized systems, smart contract languages, and virtual machines. He is the translator of the Chinese edition of Highly Scalable Systems. He is also the technical mentor of HuangLab, a blockchain laboratory at Sun Yat-sen University. Visit the URL{ https://github.com/Jianru-Lin/ } to learn about his open-source projects.


## Xiaofei Luo

Xiaofei Luo is currently a postdoctoral researcher at Sun Yat-sen University. He received his Ph.D. degree in Computer Science and Engineering from the University of Aizu in March 2023. His current research interests include blockchain, payment channel networks, and reinforcement learning. His research has been published in IEEE JSAC and other well-known international journals and conferences. He participated in the development of BlockEmulator.


##  Qinglin Yang

Qinglin Yang, IEEE member, received his Ph.D. in Computer Science and Engineering from the University of Aizu in March 2021. His research interests include intelligent edge cloud computing, federated learning privacy protection, and Web3. In recent years, he has published nearly 20 papers in international academic journals/conferences. He is also a member of the editorial board of the IEEE Open Journal of the Computer Society (OJ-CS). He participated in the development of BlockEmulator. He is the co-author of the blockchain book titled "From Blockchain to Web3 & Metaverse", published by Springer, 2023.

## Qinde Chen

Qinde Chen is a Ph.D. student at the School of Software Engineering, Sun Yat-sen University. His research interests mainly include blockchain. In August 2022, he joined HuangLab and participated in the development of BlockEmulator.

## Xiaowen Peng 

Xiaowen Peng received his M.Sc. degree from the School of Computer Science and Engineering at Sun Yat-Sen University in 2023, China. His research interests mainly include blockchain. In August 2021, he participated in the development of BlockEmulator.

## Yue Lin

Yue Lin received his M.Sc. degree from the School of Computer Science and Engineering at Sun Yat-Sen University in 2024, China. His research interests mainly include Blockchain. Since joining Huanglab in May 2021, he has been involved in two papers published at well-known conferences and journals. He is an early developer of the BlockEmulator.

## Miaoyong Xu

Miaoyong Xu received his M.Sc. degree from the School of Computer Science and Engineering at Sun Yat-Sen University in 2025, China. His research interests mainly include blockchain. He joined HuangLab in October 2021 and participated in the development of BlockEmulator.

## Junhao Wu

Junhao Wu is a Ph.D student at the School of Software Engineering, Sun Yat-sen University. In September 2022, he joined HuangLab as an intern and participated in the development of BlockEmulator.

## Baozhou Xie

Baozhou Xie is currently a Master's student of the School of Software and Engineering at Sun Yat-Sen University. His research interests mainly focus on blockchain. In September 2024, he joined HuangLab to help develop BlockEmulator. He also led the development of the first version of BrokerChain Testnet.

## Feihong Hu

Feihong Hu is currently a Master's student of the School of Software and Engineering at Sun Yat-Sen University. His research interests mainly focus on blockchain. In September 2024, he joined HuangLab to help develop BlockEmulator.


## Yang Zhou

Yang Zhou is currently a Master's student of the School of Software and Engineering at Sun Yat-Sen University. His research interests mainly focus on blockchain. In September 2025, he joined HuangLab to help develop BlockEmulator and BrokerChain Testnet/Wallet.
