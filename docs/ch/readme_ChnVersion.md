# BlockEmulator 介绍文档

**常见问题**见 [FAQ](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/FAQ.md)。

**版本更新**见 [Version Update](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/versionUpdate.md)。

这里是 BlockEmulator 官方中文版，本文档根据当前 BlockEmulator 1.0 版本编写。

<h2>一、项目介绍</h2>

<h3>1. BlockEmulator 诞生的背景、开源初衷</h3>

BlockEmulator 是由中山大学·软件工程学院·黄华威研究组（  [HuangLab](http://www.xintelligence.pro) ）发起的可支持**多种共识协议**与**跨分片机制**的区块链协议验证平台。



此实验平台开源的主要目的是为了帮助用户（研究者、学生）快速验证他们提出的新型区块链共识协议和分片机制。BlockEmulator 被设计为采用轻量化区块链系统架构的实验平台。它简化了工业级区块链系统的实验流程，这是因为 blockEmualtor 仅仅实现了区块链核心功能，比如交易池、区块打包、区块共识、交易上链等核心环节，并且支持常见的几种主流共识协议，如拜占庭容错 ( Practical Byzantine Fault Tolerance, PBFT ) 协议与工作量证明机制。



特别地，blockEmualtor 对主流的“区块链分片机制”进行了系统底层级别的设计与实现 。其中，blockEmualtor 实现的“跨分片交易”机制包含以下两个具有代表性的分片协议：Monoxide (NSDI'2019) 方案中提出的 “Relay 交易机制”，以及 BrokerChain (INFOCOM'2022) 中的 “broker 机制”。



该实验平台主要面向区块链研究人员，当他们需要对提出的新型区块链共识协议、新型跨分片机制进行验证时，可以帮助进行快速搭建一个轻量化的区块链底层协议的实验平台，并对实验数据进行收集，方便绘制科研论文所需的实验图。



本介绍文档旨在为 BlockEmulator 提供一个易上手、足够详细的用户使用手册。


<h3>2. 官方技术论文 </h3>
为了提供一个严谨的官方文档，我们近期撰写了一篇技术论文，题目为 "BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols"，已经上传到了 arXiv. 如果您使用了 BlockEmulator 当做实验工具，请别忘记引用我们这篇技术论文，bib 信息如下：

```
@article{huang2023blockemulator,
   title={BlockEmulator: An Emulator Enabling to Test Blockchain Sharding Protocols},
   author={Huang, Huawei and Ye, Guang and Chen, Qinde and Yin, Zhaokang and Luo, Xiaofei and Lin, Jianru and Li, Taotao and Yang, Qinglin and Zheng, Zibin},
   journal={arXiv preprint arXiv:2311.03612},
   year={2023}
   }
```

<h3>3. 相关工作 -- BlockEmulator 产出的论文</h3>

- 如下几篇 HuangLab 产出的区块链论文使用了 BlockEmulator 作为实验平台。欢迎感兴趣的同行了解：

- - BrokerChain: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding **(INFOCOM 2022)** 【[PDF](https://www.researchgate.net/publication/356789473_BrokerChain_A_Cross-Shard_Blockchain_Protocol_for_AccountBalance-based_State_Sharding)】【[公众号介绍文章](https://mp.weixin.qq.com/s/5MelID6kVMQeM1LAET-37w)】
  - Achieving Scalability and Load Balance across Blockchain Shards for State Sharding **(SRDS 2022)**【[PDF](https://ieeexplore.ieee.org/document/9996899)】【[公众号介绍文章](https://mp.weixin.qq.com/s/UMsQ7VIzyU-nzW5MzUXPfg)】
  - tMPT: Reconfiguration across Blockchain Shards via Trimmed Merkle Patricia Trie **(IWQoS 2023)**【[PDF](https://www.researchgate.net/publication/370633426_tMPT_Reconfiguration_across_Blockchain_Shards_via_Trimmed_Merkle_Patricia_Trie)】【[公众号介绍文章](https://mp.weixin.qq.com/s/M7KxjWTyheygrgKgDfz97g)】
  - MVCom: Scheduling Most Valuable Committees for the Large-Scale Sharded Blockchain **(ICDCS 2021) 【[PDF](https://ieeexplore.ieee.org/document/9546408)】【[公众号介绍文章](https://mp.weixin.qq.com/s/3EUO5Pt5-4hwGWtCJaFnIQ)】



<h3>4. BlockEmulator 亮点</h3>

1. **轻量化：** BlockEmulator 是一个轻量化的区块链实验平台。

2. **快速搭建：** 方便用户进行快速搭建区块链实验平台，并且支持远程部署到云端运行。

3. **可定制化：** BlockEmulator 是基于 Golang 语言实现的区块链实验平台，支持用户定制化二次开发。

4. **易于实验：** BlockEmulator 支持主流区块链（如以太坊）的历史交易数据的回放，可以自动输出、保存各项区块链实验指标（如系统吞吐量、交易确认时延、交易池拥塞程度等等）以及系统运行的日志。这些功能非常便于科研人员与学生进行实验数据的收集以及实验图的绘制。

   

<h3>5. 快速上手</h3>

- 如果你想快速上手使用 BlockEmulator，请参考：

- - BlockEmulator一个快速上手的例子：[BlockEmulator 使用帮助](./blockEmlator_manual.md)



<h2>二、项目使用文档</h2>

<h3>1. 环境配置</h3>

<h4>1.1 代码的下载部署</h4>

BlockEmulator 使用的实验数据来自中山大学·软件工程学院·InPlusLab 在 [xblock.pro](https://xblock.pro/#/) 发布的 [Ethereum On-chain Data](https://xblock.pro/#/dataset/14) 数据集。

BlockEmulator 项目代码：[ [Github Project ](https://github.com/HuangLab-SYSU/block-emulator)]



<h4>1.2 BlockEmulator 相关环境的配置。</h4>

Golang 最低版本为 1.18。

中国境内使用时，可以在命令行中输入如下指令配置镜像：

```Shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```



<h3>2. BlockEmulator 的架构设计</h3>



<h4>2.1. BlockEmulator 的系统架构</h4>

[系统架构](./system_arch.md)

 

<h4>2.2. BlockEmulator 的参数配置</h4>

[BlockEmulator 参数配置](./params.md) 



<h4>2.3. BlockEmulator 的数据结构设计</h4>

[数据结构设计](./data.md) 



<h4>2.4. BlockEmulator 的运行流程</h4>

[运行周期](./flow.md) 



<h4>2.5 内置的共识协议和跨分片交易机制</h4>

<h5>2.5.1 Monoxide 提出的 Relay Transaction 接力交易的跨分片机制 [NSDI'2019 Monoxide] </h4>

[Relay 跨分片交易技术文档](./relay.md)



<h5>2.5.2 BrokerChain 提出的“中间人”跨分片机制 [INFOCOM'2022 BrokerChain] </h5>

[BrokerChain 技术文档](./brokerchain.md) 



<h5>2.5.3 SRDS 论文提出的“账户划分”机制 [SRDS'2022 Achieving] </h5>

[CLPA 算法详解](./CLPA.md) 



<h4>2.6 BlockEmulator 可测量的指标</h4>

[BlockEmulator 可测量的指标](./metirc.md) 





<h4>2.7 BlockEmulator 的运行日志</h4>

[BlockEmulator 的运行日志](./log.md)

 


<h2> References </h2>

- [NSDI'2019 Monoxide] J. Wang and H. Wang, “Monoxide: Scale out blockchains with asynchronous consensus zones,” in 16th USENIX Symposium on Networked Systems Design and Implementation (NSDI 19), 2019, pp. 95–112.
- [INFOCOM'2022 BrokerChain] H. Huang, X. Peng, J. Zhan, S. Zhang, Y. Lin, Z. Zheng, and S. Guo, “Brokerchain: A cross-shard blockchain protocol for account/balance-based state sharding,” in Proc. of IEEE Conference on Computer Communications (INFOCOM’22), 2022, pp. 1–10. 
- [SRDS'2022 Achieving] C. Li, H. Huang, Y. Zhao, X. Peng, R. Yang, Z. Zheng, and S. Guo, “Achieving scalability and load balance across blockchain shards for state sharding,” in Proc. of 2022 41st International Symposium on Reliable Distributed Systems (SRDS’22), 2022, pp. 284–294.



<h2>三、BlockEmulator 的贡献者</h2>

# 黄华威

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/Huawei.PNG" width=200>

中山大学副教授，博士生导师，IEEE 高级会员，中国计算机学会 (CCF) 高级会员，CCF 区块链专委会执行委员、CCF 分布式与并行计算专委会执行委员。2016年取得日本会津大学“计算机科学与工程”博士学位；曾先后担任日本学术振兴会特别研究员、香港理工大学访问学者、日本京都大学特任助理教授。研究方向包括区块链底层机制、分布式系统与协议、Web3 与元宇宙底层关键技术。研究成果发表在 CCF-A 类推荐期刊 IEEE Journal on Selected Areas in Communications（JSAC），IEEE Transactions on Parallel and Distributed Systems (TPDS), IEEE Transactions on Dependable and Secure Computing (TDSC), IEEE Transactions on Mobile Computing (TMC) 与 IEEE Transactions on Computers (TC)等，以及 CCF 推荐 A / B 类国际学术会议 INFOCOM、ICDCS、SRDS、IWQoS等。论文谷歌引用 4700+，H-index 30。出版区块链教材《从区块链到 Web3》人民邮电出版社出版，与区块链学术专著《From Blockchain to Web3 & Metaverse》Springer, 2023。

# 郑子彬

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/zibin.PNG" width=200>

中山大学计算机学院教授，软件工程学院副院长、国家优秀青年科学基金获得者、IEEE Fellow、ACM 杰出科学家、国家数字家庭工程技术研究主心副主任、区块链与智能金融研究中心主任。出版 Springer 英文学术专著 2 部、发表论文 200 余篇，论文谷歌学术引用超过 33000 次。获得教育部自然科学二等奖、吴文俊人工智能自然科学二等奖、青年珠江学者、 ACM 中国新星提名奖、国际软件工程大会（ICSE）ACM SIGSOFT Distinguished Paper Award、国际 Web 服务大会（ICWS）最佳学生论文奖等奖项；担任数十个国际学术会议的程序委员会主席。

# 林建入

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/Jianru.png" width=200>

林建入，资深程序员，在去中心化系统、智能合约语言与虚拟机的设计与实现方面拥有丰富经验。《高伸缩性系统》中文版译者。唬米科技创始人，同时也是中山大学区块链实验室 HuangLab 的技术导师。访问 \url{ https://github.com/Jianru-Lin/} 可以了解他的开源项目。设计并领导开发了 blockEmulator 以及其后续高阶版本。参与编著区块链教材《从区块链到 Web3》人民邮电出版社出版。

# 李涛涛

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/taotao.PNG" width=200>

李涛涛，博士，博士后研究员，于2022年毕业于中国科学院信息工程研究所网络空间安全专业获工学博士学位，现就职于中山大学软件工程学院进行博士后研究工作。主要研究方向为区块链理论与技术应用，应用密码，具体包括侧链技术、跨链协议、轻量级区块链、密码工具在区块链中的应用。参与多项国家/省部级重点研发计划课题，近年来在CCF推荐国际学术会议和期刊上发表论文多篇。参与编写区块链教材《从区块链到 Web3》人民邮电出版社出版。

# 罗肖飞

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/Luo.png" width=200>

罗肖飞，中山大学博士后研究员，于2023年3月获得日本会津大学“计算机科学与工程”博士学位，曾参与日本文部省 RFID 项目的研发等相关工作。目前的研究方向为区块链、支付通道网络，强化学习等等。相关研究成果发表在 CCF A 类期刊及其他知名国际期刊与会议。参与多项科技部和广东省重点研发计划项目。参与编写区块链教材《从区块链到 Web3》人民邮电出版社出版。


# 杨青林

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/qinglin.PNG" width=200>

杨青林，博士，助理研究员， IEEE 会员（IEEE member），于2021年3月取得日本会津大学“计算机科学与工程”博士学位。主要研究方向包括智能边缘云计算，联邦学习隐私保护，Web3。近年来在国际学术期刊/会议发表论文近 20 篇。担任 IEEE Open Journal of Computer Science (OJ-CS) 专刊客座编委成员。参与多项国家重点研发计划课题、国家自然科学基金面上项目的研发工作。参与编著区块链教材《从区块链到 Web3》人民邮电出版社出版，与区块链学术专著《From Blockchain to Web3 & Metaverse》Springer, 2023。

# 陈钦德

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/qinde.png" width=200>


陈钦德，中山大学软件工程学院 2023 级博士生，目前的研究方向为区块链。自 2022 年 8 月进入 HuangLab 研究学习， 参与 BlockEmulator 的开发。

# 叶光

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/yeguang.png" width=200>

叶光，中山大学软件工程学院 2023 级研究生，目前的研究方向为区块链。自 2021 年 8 月以本科实习生身份加入 HuangLab，先后参与了一篇 CCF-A 期刊论文和 BlockEmulator 的开发。

# 殷昭伉

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/zhaokang.png" width=200>

殷昭伉，中山大学软件工程学院2023级研究生，目前的研究方向为区块链。2022年10月，通过研究生推免加入HuangLab，参与 BlockEmulator的开发。


# 彭肖文

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/Xiaowen.jpg" width=200>

彭肖文，中山大学计算机学院 2020 级研究生，已毕业。研究方向为区块链。自 2019 年 6 月加入HuangLab，先后参与两篇 CCF-A 类论文研究，是 BlockEmulator 的早期开发者。

# 林岳

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/Linyue.png" width=200>

林岳，中山大学计算机学院 2021 级研究生，目前的研究方向为区块链。自 2021 年 5 月加入 HuangLab，先后参与了一篇 CCF-A 会议论文和一篇 CCF-A 期刊论文，是 BlockEmulator 的早期开发者。

# 许淼泳

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/MiaoYong.png" width=200>

许淼泳，中山大学计算机学院 2022 级研究生，目前的研究方向为区块链。自 2021 年 10 月进入 HuangLab 研究学习， 参与 BlockEmulator 的开发。

# 吴均豪

<img src ="https://github.com/HuangLab-SYSU/block-emulator/blob/main/bios/junhao.png" width=200>


吴均豪，中山大学软件工程学院2021级本科生。2022年9月加入 HuangLab 实习，参与 BlockEmulator 的开发。
