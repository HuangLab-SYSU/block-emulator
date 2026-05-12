# BlockEmulator 使用帮助

## 一、设置

### 下载数据集

数据集来自xblock发布的 [Ethereum On-chain Data](https://xblock.pro/#/dataset/14)

### *设置系统参数*

打开配置文件 [params/global_config.go](https://github.com/Jianru-Lin/block-emulator-v1/blob/79325c6ddd009c450a00ffbc0e06073a74f3c428/params/global_config.go)，设置参数，具体含义如下：

```Plain
Block_Interval      = 5000                                       // The interval (millisecond) to generate a new block
MaxBlockSize_global = 2000                                       // The maximum number of transactions contained in a block 
InjectSpeed         = 2000                                       // The transaction inject speed
TotalDataSize       = 100000                                     // The total number of transactions
DataWrite_path      = "./result/"                                // Output path of the measurement data
LogWrite_path       = "./log"                                    // Output path of the log
SupervisorAddr      = "127.0.0.1:18800"                          // Supervisor ip address
FileInput           = "../2000000to2999999_BlockTransaction.csv" // The path of the raw BlockTransaction data
```



## 二、使用说明

### 指令说明

```Plain
  -c,                  Running as a supervisor client 
  -g,                  To generate bat file
  -m, --modID int      Choice Committee Method,for example, 0, [CLPA_Broker,CLPA,Broker,Relay]  (default 3)
  -n, --nodeID int     Node ID, for example, 0
  -N, --nodeNum int    Indicate how many nodes of each shard are deployed (default 4)
  -s, --shardID int    Id of the shard to which this node belongs, for example, 0
  -S, --shardNum int   Indicate that how many shards are deployed (default 2)
```

### 启动方法

1. 首先配置 supervisor client 节点，此节点用于整体网络参数配置、状态观察与数据输出。

   使用如下命令启动一个 supervisor client，初始化一个包含 2 个分片，各分片 4 节点，使用 Relay 跨链机制的区块链网络：

```Plain
go run main.go -c -N 4 -S 2 -m 3 
```

2. 启动区块链节点，该节点是仿真区块链中的矿工节点，具备生成区块、发起共识等功能。

   使用如下命令启动一个区块链节点，其节点 ID 为 0，所属分片 ID 为 0，处于具有 2 分片、各分片 4 节点的区块链网络  (采用的是 Relay 跨链机制)：

```Go
go run main.go -n 0 -N 4 -s 0 -S 2 -m 3 
```

### 生成启动批文件

除了逐步手动配置启动区块链网络之外，还可以使用批文件进行快速配置和启动。

使用以下命令快速设置分片数量、分片内节点数量、仿真模式，可生成 Bat文件：

```Go
go run main.go -g -S 2 -N 4 -m 3
```





## 三、使用实例

以启动 Relay 跨链机制区块链网络为例，2 分片，每个分片4 节点。

### 手动启动

启动 Supervisor

```Go
go run main.go -c -N 4 -S 2 -m 3 
```

启动 共识节点

```Go
go run main.go -n 1 -N 4 -s 0 -S 2 -m 3 

go run main.go -n 1 -N 4 -s 1 -S 2 -m 3 

go run main.go -n 2 -N 4 -s 0 -S 2 -m 3 

go run main.go -n 2 -N 4 -s 1 -S 2 -m 3 

go run main.go -n 3 -N 4 -s 0 -S 2 -m 3 

go run main.go -n 3 -N 4 -s 1 -S 2 -m 3 

go run main.go -n 0 -N 4 -s 0 -S 2 -m 3 

go run main.go -n 0 -N 4 -s 1 -S 2 -m 3 
```

### 自动启动

#### Windows 
1.生成 Bat 批处理文件，如

```Go
go run main.go -g -S 2 -N 4 -m 3
```

2.点击项目路径下的 `bat_shardNum=2_NodeNum=4_mod=Relay.bat` 文件运行系统

#### Linux
<font color='red'>注意：如果通过命令行访问 Linux 系统，需要提前做好杀死进程的准备，防止程序运行出现问题导致进程无法关闭。
比如：
```
pkill
```
</font>

1. 生成 shell 批处理文件，如

```Go
go run main.go -g -S 2 -N 4 -m 3
```

2.点击项目路径下的 `bat_shardNum=2_NodeNum=4_mod=Relay.shell` 文件运行系统