# BlockEmulator Handbook
This is the official handbook of **BlockEmulator** (version 1.0).
We also provide a **[Chinese-version handbook](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/ch/readme_ChnVersion.md)**.

# FAQ (Frequently Asked Questions)
See **[FAQ](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/FAQ.md)**. 

# Version Updates
View the updates of **BlockEmulator** here: 
**[ Version Updates ](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/versionUpdate.md)**

# Introduction 

## 1. Background of BlockEmulator


Initiated by **[HuangLab](http://xintelligence.pro/)**  (a research group in the School of Software Engineering, Sun Yat-sen University, China), ***BlockEmulator** is a blockchain testbed that enables researchers to verify their proposed new protocols and mechanisms. It supports multiple consensus protocols, particularly the blockchain sharding mechanism.


The main purpose of this testbed is to help users (researchers, students, etc.) quickly verify their own blockchain consensus protocols and blockchain-sharding protocols. 

**BlockEmulator** is designed as an experimental platform that adopts a lightweight system architecture. It simplifies the implementation of industrial-class blockchains since **BlockEmulator** only implements the core functions of a blockchain, including the transaction pool, block packaging, consensus protocols, and on-chain transaction storage. It also supports common consensus protocols, such as Practical Byzantine Fault Tolerance (PBFT).

In particular, BlockEmulator offers the system-level design and implementation for blockchain-sharding mechanisms. For example, the cross-shard transaction mechanisms implemented by BlockEmulator include the following two representative solutions, i.e., i) **Relay transaction mechanism** proposed by **Monoxide** (NSDI'2019), and the **BrokerChain** protocol (INFOCOM'2022) [PDF](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding).

BlockEmulator is oriented toward blockchain researchers. It provides a blockchain experimental platform for quickly implementing their own algorithms, protocols, and mechanisms. It also offers very helpful functions to help researchers collect experimental data, facilitating their plotting experimental figures.


## 2. Official Technical Paper & Citation

To provide an official handbook for BlockEmulator, we have written a technical paper titled "BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols" [arXiv page](https://arxiv.org/abs/2311.03612). Please cite our paper if you use BlockEmulator as an experiment tool in your own paper.

```
@article{huang2023blockemulator,
   title={BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols},
   author={Huang, Huawei and Ye, Guang and Chen, Qinde and Yin, Zhaokang and Luo, Xiaofei and Lin, Jianru and Li, Taotao and Yang, Qinglin and Zheng, Zibin},
   journal={arXiv preprint arXiv:2311.03612},
   year={2023}
   }
```


## 3. Related Work

The following papers from HuangLab's publications have adopted **BlockEmulator** as an experimental tool.

1. **BrokerChain**: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding (published at INFOCOM 2022) [PDF](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding)
   
2. Achieving Scalability and Load Balance across Blockchain Shards for State Sharding (published at SRDS 2022) [PDF](https://ieeexplore.ieee.org/document/9996899)
   
3. **tMPT**: Reconfiguration across Blockchain Shards via Trimmed Merkle Patricia Trie (published at IWQoS 2023) [PDF](https://www.researchgate.net/publication/370633426_tMPT_Reconfiguration_across_Blockchain_Shards_via_Trimmed_Merkle_Patricia_Trie)
   
4. **MVCom**: Scheduling Most Valuable Committees for the Large-Scale Sharded Blockchain (published at ICDCS 2021) [PDF](https://ieeexplore.ieee.org/document/9546408)


## 4. Highlights of BlockEmulator

1. **Lightweight**. BlockEmulator is a lightweight testbed platform for blockchain experiments.
2. **Fast Configuration**. BlockEmulator enables users to set up their environments quickly and supports remote deployment on the Cloud.
3. **Customization**. BlockEmulator is implemented using goland language, which supports users' customization and modification.
4. **Easy to Conduct Experiments**. BlockEmulator supports the replay of historical transactions of mainstream blockchains (such as Ethereum). It can automatically yield experimental log files. Using those log files, researchers can interpret plenty of metrics such as system throughput, confirmation latency of transactions, the queueing of the transaction pool, etc. This function is very useful for researchers and students to facilitate their experimental data collection and plotting of experimental charts.

## 5. Getting Started to Use BlockEmulator

Quickly get started with BlockEmulator through the following document. Please refer to:
[BlockEmulator help document](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/Help%20Documents.md). More details of the user guidebook can be found at [Handbook of BlockEmulator (Eng.)](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/en/readme.md) or [Handbook of BlockEmulator (Chn.)](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/ch/readme_ChnVersion.md).

## Reference
  - [NSDI'2019 Monoxide] J. Wang and H. Wang, “Monoxide: Scale out blockchains with asynchronous consensus zones,” in 16th USENIX Symposium on Networked Systems Design and Implementation (NSDI 19), 2019, pp. 95–112.
  
  - [INFOCOM'2022 BrokerChain] H. Huang, X. Peng, J. Zhan, S. Zhang, Y. Lin, Z. Zheng, and S. Guo, “Brokerchain: A cross-shard blockchain protocol for account/balance-based state sharding,” in Proc. of IEEE Conference on Computer Communications (INFOCOM’22), 2022, pp. 1–10. 
  
  - [SRDS'2022 Achieving] C. Li, H. Huang, Y. Zhao, X. Peng, R. Yang, Z. Zheng, and S. Guo, “Achieving scalability and load balance across blockchain shards for state sharding,” in Proc. of 2022 41st International Symposium on Reliable Distributed Systems (SRDS’22), 2022, pp. 284–294.


# Contributors

## Huawei Huang

Huawei Huang is an associate professor at Sun Yat-sen University. He is an IEEE senior member. He received his Ph.D. in Computer Science and Engineering from The University of Aizu, Japan. His research interests include blockchain, Web3, metaverse, and distributed computing/protocols. He has more than 100 papers published in top journals or conferences such as IEEE/ACM ToN, JSAC, TPDS, TDSC, TMC, INFOCOM, ICDCS, SRDS, IWQoS, etc. He is a PI or Co-PI of several blockchain-involved research projects funded by the National Key Research & Development Program of China, National Natural Science Foundation of China (NSFC), etc. He has published two blockchain books titled "From Blockchain to Web3 & Metaverse" and "Blockchain Scalability", both published by Springer in 2023.

## Zibin Zheng

Zheng Zibin (IEEE Fellow) received his Ph.D. degree from the Chinese University of Hong Kong, Hong Kong, in 2012. He is a Professor at the School of Software Engineering, Sun Yat-Sen University, China. His current research interests include service computing, blockchain, and cloud computing.

## Jianru Lin

Lin Jianru, a research scientist engineer, has rich experience in the design and implementation of decentralized systems, smart contract languages, and virtual machines. He is the translator of the Chinese edition of Highly Scalable Systems. He is also the technical mentor of HuangLab, a blockchain laboratory at Sun Yat-sen University. Visit url{ https://github.com/Jianru-Lin/ } to learn about his open-source projects.

## Taotao Li

Taotao Li is a postdoctoral researcher at the School of Software Engineering, Sun Yat-sen University. He received his Ph.D. in cyberspace security from the Institute of Information Engineering, Chinese Academy of Sciences in 2022. The main research interests include blockchain theory and technology application, and the application of cryptography, including sidechain technology, cross-chain protocol, lightweight blockchain, and the application of cryptographic tools in blockchain. He participated in the development of BlockEmulator.


## Xiaofei Luo

Xiaofei Luo is currently a postdoctoral researcher at Sun Yat-sen University. He received his Ph.D. degree in Computer Science and Engineering from the University of Aizu in March 2023. His current research interests include blockchain, payment channel networks, and reinforcement learning. His research has been published in IEEE JSAC and other well-known international journals and conferences. He participated in the development of BlockEmulator.


##  Qinglin Yang

Qinglin Yang, Ph.D., research fellow, and IEEE member, received his Ph.D. in Computer Science and Engineering from the University of Aizu in March 2021. His research interests include intelligent edge cloud computing, federated learning privacy protection, and Web3. In recent years, he has published nearly 20 papers in international academic journals/conferences. He is also a guest editorial board member of the IEEE Open Journal of the Computer Society (OJ-CS). He participated in the development of BlockEmulator. He is the co-author of the blockchain book titled "From Blockchain to Web3 & Metaverse", published by Springer, 2023.

## Qinde Chen

Qinde is a Ph.D. student at the School of Software Engineering, Sun Yat-sen University. His research interests mainly include blockchain. In August 2022, he joined HuangLab and participated in the development of BlockEmulator.

## Guang Ye

Guang is currently a student pursuing his Master's degree at the School of Software and Engineering, at Sun Yat-Sen University. His research interests mainly include Blockchain.  Since joining HuangLab as an undergraduate intern in August 2021, he has participated in the development of BlockEmulator.

## Zhaokang Yin

Zhaokang is currently a student pursuing his Master's degree at the School of Software and Engineering. His research interest mainly focuses on blockchain. In October 2022, he joined HuangLab to participate in the development of BlockEmulator.

## Xiaowen Peng 

Xiaowen received his M.Sc. degree from the School of Computer Science and Engineering in 2023, Sun Yat-Sen University, China. His research interests mainly include blockchain. In August 2021, he participated in the development of BlockEmulator.

## Yue Lin

Yue is currently a student pursuing his master's degree at the School of Computer Science and Engineering, at Sun Yat-Sen University. His research interests mainly include Blockchain. Since being a member of Huanglab in May 2021, he has been involved in two papers published in well-known conferences and journals. He is an early developer of the BlockEmulator.

## Miaoyong Xu

Miaoyong is currently a student pursuing his master's degree at the School of Computer Science and Engineering, Sun Yat-sen University. His research interests mainly include blockchain. he joined HuangLab in October 2021 and participated in the development of BlockEmulator.


## Junhao Wu

Junhao is an undergraduate student at the School of Software Engineering, Sun Yat-sen University. In September 2022, he joined HuangLab as an intern and participated in the development of BlockEmulator.
