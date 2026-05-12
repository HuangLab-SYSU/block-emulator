# BlockEmulator 可测量的指标

## BlockEmulator 指标调用

BlockEmulator的指标将由 supervisor 进程来测量。因此，如果要在 BlockEmulator 中调用指标，需要在创建 supervisor 时，指定 supervisor 测量的指标名称：

```Go
spv:= new(supervisor.Supervisor)
spv.NewSupervisor(supervisor_ip, chainConfig, committeeMethod, mearsureModNames...)
```

在 Golang 中，mearsureModNames... 指的是“可变参数”。

## BlockEmulator 支持的指标

目前，blockEmulator 支持的指标有八种（四类）：

| 指标类型                             | 说明                                                         |                                                              |
| ------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| Relay 跨分片交易处理模式             | 吞吐量（TPS_Relay）                                          | 针对 Relay 方式的吞吐量测试，吞吐量是平均下区块链系统每秒处理交易的数目；它是区块链用来评测性能的基本指标之一 |
| 交易确认时延（TCL_Relay）            | 针对 Relay 方式的交易确认时延（transaction confirm latency，TCL）测试，某一笔交易的 TCL 是指该交易从进入交易池到最终确认上链所耗费的时间，而本指标测试的是一堆交易的平均 TCL ；它也是区块链用来评测性能的基本指标之一 |                                                              |
| 跨分片交易占比（CrossTxRate_Relay）  | 针对 Relay 方式的跨分片交易占比，指的是跨分片交易所占区块链系统处理的所有交易的占比；它常被用于评测一个账户重划分算法的性能好坏 |                                                              |
| 交易总数(TxNumberCount_Relay)        | 测量交易总数，可以用来统计最终上链交易的数目。               |                                                              |
| Broker account跨分片交易处理模式     | 吞吐量（TPS_Broker）                                         | 针对 Broker account 方式的吞吐量测试，吞吐量是平均下区块链系统每秒处理交易的数目；它是区块链用来评测性能的基本指标之一 |
| 交易确认时延（TCL_Broker）           | 针对 Broker account 方式的交易确认时延（transaction confirm latency，TCL）测试，某一笔交易的 TCL 是指该交易从进入交易池到最终确认上链所耗费的时间，而本指标测试的是一堆交易的平均 TCL ；它也是区块链用来评测性能的基本指标之一 |                                                              |
| 跨分片交易占比（CrossTxRate_Broker） | 针对 Broker account 方式的跨分片交易占比，指的是跨分片交易所占区块链系统处理的所有交易的占比；它常被用于评测一个账户重划分算法的性能好坏 |                                                              |
| 交易总数(TxNumberCount_Broker)       | 测量交易总数，可以用来统计最终上链交易的数目。               |                                                              |

以上指标会在 supervisor 程序运行结束时输出（在此版本的 BlockEmulator 中，一笔跨分片交易将会被记为两笔0.5笔交易）。这些指标的输出格式是：

```Go
var metricByEpoch []float // output the metric result in an array, the index of the metric result is set by epoch
var metricResultAll float // the metric result throughout the entire running
```

即测量结果会按照 Epoch 逐一置入数组 metricByEpoch 中，并且将整个运行过程的测量结果写入 metricResultAll 中。 

另外，各个指标还会以 .csv 文件的形式 被写入 params.DataWrite_path/supervisor_measureOutput 的文件夹中。
