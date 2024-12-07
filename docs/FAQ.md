# 常见问题

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
