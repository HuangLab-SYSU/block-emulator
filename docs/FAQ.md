# Frequently Asked Questions (English Version)

## Where can I download the dataset?
- (Updated 2024/05/20) To help users run blockEmulator more quickly, we have generated the `selectedTxs_300K.csv` file, which contains three hundred thousand selected historical transactions.
- The dataset we use is from **[Ethereum Historical Transactions Collected by xblock](https://xblock.pro/xblock-eth.html)**.

## How do I run blockEmulator after downloading the dataset?
For detailed steps, please refer to the **[blockEmulator User Manual](https://github.com/HuangLab-SYSU/block-emulator/blob/main/20240906-blockEmulator%E4%BD%BF%E7%94%A8%E6%89%8B%E5%86%8C-%E5%BC%80%E6%BA%90%E7%89%88%E6%9C%AC-v2.pdf)**.

If you want to quickly test whether blockEmulator can run normally, you can choose the automatic startup method (i.e., batch processing startup) in **Section III: Usage Examples** of the aforementioned user manual.

## Why does the PBFT consensus of nodes fail to be reached or stall when I run blockEmulator?
This could be due to the following reasons:
1.  Data or logs generated from the previous run were not deleted before starting blockEmulator.
    In the update on 2024/09/01, users can specify `ExpDataRootDir` in the `./paramsConfig.json` file as the output directory. As long as the `ExpDataRootDir` is different for two consecutive experiments, the results will not cause file conflicts.

    **Solution**: Alternatively, you can delete the `./log`, `./record`, and `./result` folders under the `ExpDataRootDir`.

2.  Node termination/crash.

    This might be caused by *node port being occupied* or *insufficient device performance (typically when running over 20 shards on a single machine)*.

    **Port Occupied**: You can modify the IP of the occupied node in the `./ipTable.json` file.

## Why are all the test result data zeros after running blockEmulator?
It is possible that the Supervisor node did not start or failed to start, which prevents transactions from being sent to the shards where the miners are located.
Please check the terminal corresponding to the Supervisor node for any error messages.



## _Node_ is different from _Account_ in BlockEmulator
1. (A real question from a user, 2025Nov7) In the Ethereum transaction dataset (e.g., selected_300k.csv), we have 300,000 transactions with fields such as from and sender. My question is: when we select, for example, 8 shards and 4 nodes per shard (i.e., 32 nodes in total), while there are 10,000 unique senders in the dataset, how are these senders mapped to the 32 nodes? Are some of the transactions discarded in this process? I would really appreciate it if you could explain how this mapping is performed.  

**Response:** These sender accounts are not mapped to consensus nodes, but rather to shards. Also, no transactions have been discarded. In BlockEmulator, accounts and nodes are distinct concepts. The sender (account) of a transaction is not mapped to a specific node but is mapped to a particular shard (The static sharding rule is: the sender’s shard's ID = sender account_address % (the number of shards)).

When the supervisor node in BlockEmulator reads transactions from the dataset, it distributes them to different shards. The nodes within these shards will then execute the received transactions. By default, transactions are sent to the [sender's corresponding shard].



=======#################################################################################=======



# 常见问题（中文版本）

## 我该在哪下载数据集？
- (updated in 2024/05/20) 为了帮助用户更快运行 blockEmulator，我们生成了 ''selectedTxs_300K.csv'' 文件，其中包含三十万条挑选后的历史交易。
- 我们使用的数据集是来自 **[xblock 收集的以太坊历史交易](https://xblock.pro/xblock-eth.html)**。

## 下载完数据集后，我该如何运行 blockEmulator？
具体步骤可参见 **[blockEmulator 用户手册](https://github.com/HuangLab-SYSU/block-emulator/blob/main/20240906-blockEmulator%E4%BD%BF%E7%94%A8%E6%89%8B%E5%86%8C-%E5%BC%80%E6%BA%90%E7%89%88%E6%9C%AC-v2.pdf)**。

如果想要更快捷地测试是否 blockEmulator 能正常运行，可以在上述使用帮助的 **三、使用实例** 中选择自动启动（即 批处理方式启动）。

## 为什么我运行 blockEmulator 时，节点的 PBFT 无法达成共识 / 共识停滞不前？
原因可能有如下几种：
1. 启动 blockEmulator 前，未删除上一次运行时产生的数据或日志。
    在 2024/09/01 的更新中，用户可以在 `./paramsConfig.json` 文件中指定 `ExpDataRootDir` 作为输出文件的目录。只需要确保两次实验的 `ExpDataRootDir` 不一样，那么实验的运行结果便不会造成文件冲突。

    **解决方案**：或者也可以删除 `ExpDataRootDir` 下的 `./log` `./record` `./result` 文件夹
    
2. 节点中止 / 崩溃。

    这种情况产生的原因可能是 *节点端口被占用*，*设备性能不足（一般是一台机器超过 20 个分片时）*

    **端口占用**：可以在 `./ipTable.json` 文件中修改被占用节点的 IP。

## 为什么运行 blockEmulator 之后，测试结果的数据全是 0？
有可能是 Supervisor 节点没有启动 / 启动失败，这导致交易没有发送给矿工们所在的分片。

可以查看 Supervisor 节点所对应的终端是否出现了报错信息。

## _Node_ is different from _Account_ in BlockEmulator
1. (A real question from an international user, 2025Nov7) In the Ethereum transaction dataset (e.g., selected_300k.csv), we have 300,000 transactions with fields such as from and sender. My question is: when we select, for example, 8 shards and 4 nodes per shard (i.e., 32 nodes in total), while there are 10,000 unique senders in the dataset, how are these senders mapped to the 32 nodes? Are some of the transactions discarded in this process? I would really appreciate it if you could explain how this mapping is performed.  

**回复:** 先说结论：这些发送方账户（sender account）并不会被映射到共识节点（nodes），它们是被映射到分片（shard）的；而且并没有交易被 discard。

**解释：** 在 BlockEmulator 中，账户（account）和节点（node）是两个不同的概念，交易的 sender（account） 并不会被映射到某一个节点（node），但会被映射到某一个分片（静态分片规则: 发送方账户所属分片 = 发送方账户地址 % 分片数目）。

当 BlockEmulator 的 supervisor 节点从数据集中读取交易后，supervisor 节点会将这些交易发送到不同的分片；相关分片内的节点会执行接收到的交易。默认情况下，交易会被发送到【发送方的所属分片】。
