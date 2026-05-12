# 常见问题

## 我该在哪下载数据集？
我们使用的数据集是来自 **[xblock 收集的以太坊历史交易](https://xblock.pro/xblock-eth.html)**。

## 下载完数据集后，我该如何运行 blockEmulator？
具体步骤可参见 **[blockEmulator 使用帮助](https://github.com/HuangLab-SYSU/block-emulator/blob/main/docs/ch/blockEmlator_manual.md)**。

如果想要更快捷地测试是否 blockEmulator 能正常运行，可以在上述使用帮助的 **三、使用实例** 中选择自动启动（即 批处理方式启动）。

## 为什么我运行 blockEmulator 时，节点的 PBFT 无法达成共识 / 共识停滞不前？
原因可能有如下几种：
1. 启动 blockEmulator 前，未删除上一次运行时产生的数据或日志。

    **解决方案**：删除 *./log* *./record* *./result* 文件夹
2. 节点中止 / 崩溃。

    这种情况产生的原因可能是 *节点端口被占用*，*设备性能不足（一般是一台机器超过 20 个分片时）*

## 为什么运行 blockEmulator 之后，测试结果的数据全是 0？
有可能是 Supervisor 节点没有启动 / 启动失败，这导致交易没有发送给矿工们所在的分片。

可以查看 Supervisor 节点所对应的终端是否出现了报错信息。