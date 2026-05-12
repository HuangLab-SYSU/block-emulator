# blockEmulator .log files

The blockEmulator's Worker Nodes and supervisor will output and save the relevant run logs during the run, and the logs are saved in ***params. LogWrite_path*** path. Log files will be organized according to the following file directory structure (two shards, two nodes per shard, as an example):

```
1 // assume params.LogWrite_path is viewed as "./log"
2
3 log：
4 +---S0
5 |       N0.log
6 |       N1.log
7 |
8 \---S1
9        N0.log
10       N1.log
11 \---Supervisor
12        supervisor.log
```


# blockEmulator .csv files
In addition to .log files, the **Worker Nodes** of blockEmulator will also output .csv files when the blockchain is on the chain, and the files are saved in *params.DataWrite_path* path. It records the content related to the block within the nodes in the blockchain (such as the number of transactions in the block, the number of transactions across shards, etc.).
