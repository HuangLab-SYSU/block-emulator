# BlockEmulator 参数配置

首先下载 BlockEmulator 所需的数据集，数据集来中山大学·软件工程学院·InPlusLab 在 [xblock](https://xblock.pro/#/) 发布的 [Ethereum On-chain Data](https://xblock.pro/#/dataset/14) 数据集。

BlockEmulator 系统中，采用位于 **package params** 中的配置文件来对于系统的各类参数进行设置，每个字段默认值以及具体含义如下：

| 字段                | 类型   | 默认值                                     | 说明                  |
| ------------------- | ------ | ------------------------------------------ | --------------------- |
| Block_Interval      | int    | 5000                                       | 新区块生成间隔        |
| MaxBlockSize_global | int    | 2000                                       | 区块包含交易数量上限  |
| InjectSpeed         | int    | 2000                                       | 交易注入速度          |
| TotalDataSize       | int    | 100000                                     | 交易的总数量          |
| DataWrite_path      | String | "./result/"                                | 评估数据的输出位置    |
| LogWrite_path       | String | "./log"                                    | 日志的输出位置        |
| SupervisorAddr      | String | "127.0.0.1:18800"                          | Supervisor 的 IP 地址 |
| FileInput           | String | "../2000000to2999999_BlockTransaction.csv" | 原始区块交易数据位置  |

同时也包含更细致化的内容配置。

| 字段             | 类型     | 默认值                                                       | 说明                        |
| ---------------- | -------- | ------------------------------------------------------------ | --------------------------- |
| Init_Balance     | big.Int  | 100000000000000000000000000000000000000000000                | 账户初始余额                |
| NodesinShard     | int      | 4                                                            | 每个分片中节点的数量        |
| ShardNum         | int      | 4                                                            | 区块链分片的数量            |
| CommitteeMethod  | []string | {"CLPA_Broker", "CLPA", "Broker", "Relay"}                   | 默认委员会方法              |
| MeasureBrokerMod | []string | {"TPS_Broker", "Latency_Broker", "CrossTxRate_Broker", "TxNumberCount_Broker"} | 测试 BrokerChain 机制的方法 |
| MeasureRelayMod  | []string | {"TPS_Relay", "Latency_Relay", "CrossTxRate_Relay", "TxNumberCount_Relay"} | 测试 Relay 机制的方法       |