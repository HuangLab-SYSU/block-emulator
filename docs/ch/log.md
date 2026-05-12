# BlockEmulator 的运行日志

## blockEmulator .log 文件

blockEmulator 的 Worker Nodes 和 supervisor 会在运行中输出并保存相关的运行日志，日志保存在 params.LogWrite_path 路径下。

日志文件将会按照如下的文件目录结构组织（以两个分片、每个分片两个节点为例）：

```Plain
// 假设 params.LogWrite_path 被设定为 "./log"

log：
+---S0
|       N0.log
|       N1.log
|
\---S1
        N0.log
        N1.log
\---Supervisor
        supervisor.log
```

## blockEmulator .csv 文件

blockEmulator 的 Worker Nodes 除了 .log 文件之外，其中的 Leader 节点 还会在区块链上链时，输出 .csv 文件，文件保存在 params.DataWrite_path 路径下。它记录的是区块链中节点内部与区块相关的内容（如该区块的交易数目，跨分片交易数目等）。