This is the **Help Documents** for BlockEmulator

# 1. Settings
## 1.1 Download Dataset
The dataset is from [Ethereum-On-chain data](https://xblock.pro/#/dataset/14)，which is published by [xblock.pro](https://xblock.pro/#/).

## 1.2 Configuration of the System Parameters

open the configuration file [params/global_config.go](https://github.com/HuangLab-SYSU/block-emulator/blob/main/params/global_config.go), configures the parameters. The meaning of parameters are shown as follows:

```
1 Block_Interval  = 5000   // generate new block interval
2 MaxBlockSize_global = 2000    // the block contains the maximum number of transactions
3 InjectSpeed   = 2000     // the transaction inject speed
4 TotalDataSize  = 100000    // the total number of txs
5 DataWrite_path  = "./result/"     // measurement data result output path
6 LogWrite_path       = "./log"     // log output path
7 SupervisorAddr   = "127.0.0.1:18800"     //supervisor ip address
8 FileInput           = "../2000000to2999999_BlockTransaction.csv" //the raw BlockTransaction data path
```

# 2. Usages Explaination

## 2.1 Command  Explaination
```
 1 -c, --client   whether this node is a client
 2 -g, --gen      generation bat
 3 -m, --modID int      choice Committee Method,for example, 0, [CLPA_Broker,CLPA,Broker,Relay]  (default 3)
 4 -n, --nodeID int     id of this node, for example, 0
 5 -N, --nodeNum int    indicate how many nodes of each shard are deployed (default 4)
 6 -s, --shardID int    id of the shard to which this node belongs, for example, 0
 7 -S, --shardNum int   indicate that how many shards are deployed (default 2)
```

## 2.2 Launch 
1. Set the number of fragment nodes, number of fragments, and simulation mode, and start the supervisor client
   ```
   1 go run main.go -c -N 4 -S 2 -m 3 
   ```
2. Set node ID, number of fragment nodes, owning fragment ID, number of fragment, simulation mode, and start the consensus node

    ```
    1 go run main.go -n 0 -N 4 -s 0 -S 2 -m 3 
    ```


## 2.3 Generates bat files
Set the number of fragments, number of nodes in fragments, and simulation mode to generate Bat files

``` 
1 go run main.go -g -S 2 -N 4 -m 3 
```
# 3. Demos
Take the Relay cross-chain mechanism blockchain as an example, 2 shards, 4 nodes
## 3.1 manual start
1. lanuch **Supervisor**
   ```
   1 go run main.go -c -N 4 -S 2 -m 3 
   ```
2. Lanuch consensus node

```
1 go run main.go -n 1 -N 4 -s 0 -S 2 -m 3 

2 go run main.go -n 1 -N 4 -s 1 -S 2 -m 3 

3 go run main.go -n 2 -N 4 -s 0 -S 2 -m 3 

4 go run main.go -n 2 -N 4 -s 1 -S 2 -m 3 

5 go run main.go -n 3 -N 4 -s 0 -S 2 -m 3 

7 go run main.go -n 3 -N 4 -s 1 -S 2 -m 3 

7 8 go run main.go -n 0 -N 4 -s 0 -S 2 -m 3 

go run main.go -n 0 -N 4 -s 1 -S 2 -m 3 
```
## 3.2 auto start

### Windows 
1. Generate bat files

```Go
go run main.go -g -S 2 -N 4 -m 3
```

2. click the `bat_shardNum=2_NodeNum=4_mod=Relay.bat` that is in the path of project to run the system 


### Linux
<font color='red'>Note: When using the command line in Linux, be ready to terminate processes in advance to avoid complications during program execution.

For example：
```
pkill
```
</font>

1. Generate shell files

```Go
go run main.go -g -S 2 -N 4 -m 3
```

2. click the `bat_shardNum=2_NodeNum=4_mod=Relay.shell` that is in the path of project to run the system

