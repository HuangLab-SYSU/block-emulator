# blockEmulator-Broker2Earn

> Broker2Earn (B2E) —— 分片区块链上面向 Broker 的可持续激励机制，基于 [BlockEmulator](https://github.com/HuangLab-SYSU/block-emulator) 与 [BrokerChain](https://ieeexplore.ieee.org/document/9796859) 实现。
>
> 论文：H. Huang, Q. Chen, et al., *"Broker2Earn: Towards Maximizing Brokers' Revenue and Cross-Shard Transactions Throughput in Sharded Blockchains"*, IEEE INFOCOM 2024.

本仓库为 Broker2Earn 协议在 BlockEmulator 上的开源实现，基于 BrokerChain 会议版代码扩展而来。配合 BlockEmulator 开源项目使用手册的第 13 章阅读，可一步步复现论文中的实验。

---

## 目录

- [研究背景](#研究背景)
- [B2E 协议简介](#b2e-协议简介)
- [本仓库相对 BrokerChain 的新增功能](#本仓库相对-brokerchain-的新增功能)
- [代码结构](#代码结构)
- [快速开始](#快速开始)
- [参数配置](#参数配置)
- [节点 IP 配置](#节点-ip-配置)
- [启动实验](#启动实验)
- [实验结果文件说明](#实验结果文件说明)
- [实验图绘制示例](#实验图绘制示例)
- [注意事项](#注意事项)
- [引用](#引用)

---

## 研究背景

分片技术是保持区块链去中心化的同时提升可扩展性的可行路线。但在状态分片下，一笔涉及不同分片账户的交易（**跨分片交易，CTX**）处理代价远高于片内交易，过高的 CTX 比例会显著拖累整体吞吐。

[BrokerChain (INFOCOM '22)](https://ieeexplore.ieee.org/document/9796859) 通过引入 **broker 账户** 把一笔跨分片交易拆解为两笔片内交易，缓解了 CTX 瓶颈，但留下了一个开放问题：

> **如何吸引足够多的用户主动质押通证、自愿成为 broker，并把这些流动性合理分配到各分片中？**

Broker2Earn 正是为回答这一问题而设计的激励机制。

## B2E 协议简介

Broker2Earn 的运作流程：

1. **注册 Broker**：自愿者 (Volunteer) 通过智能合约注册申请成为 broker 账户；
2. **质押通证**：broker 将一定数量的通证质押到合约中，为分片区块链提供流动性；
3. **B2E 算法分配**：B2E 把各 broker 账户分配到不同分片，以最优化流动性使用；
4. **处理跨分片交易**：跨分片交易 `CTX⟨S_i, S_j, vol, fee⟩` 经由 broker 中转，被拆解为两笔片内交易；
5. **broker 获得收益**：每笔成功中转的 CTX 中，broker 按佣金比例 (默认 10%) 收取手续费分成。

B2E 在算法层面将上述招募-分配问题建模为一个最大化问题，证明其为 **NP-hard**，并提出基于 **Relax-and-Rounding** 的 **在线近似算法 (Online Approximation Algorithm)**：

- **Relax**：将 0/1 整数约束放松为线性规划 (LP)，多项式时间内精确求解；
- **Rounding**：对 LP 连续解做随机舍入，并给出严格的近似比保证；
- **Online**：随着交易流的实时到达进行决策，无需提前知道全局交易信息。

算法细节与近似比证明详见 [Broker2Earn 论文 (INFOCOM '24)](#引用)。

## 本仓库相对 BrokerChain 的新增功能

本仓库基于 BrokerChain 会议版代码（`blockEmulator-broker`），新增了以下 **四项** 功能以适配 B2E 协议：

1. **阻塞式注入机制 (Blocking Injection)**：每批交易注入后，Supervisor 会等待该批次的所有交易全部在链上确认后再注入下一批。当累计确认数（片内交易 + 跨分片交易对）达到 `params.InjectSpeed` 时，一个 epoch 结束。每一轮 epoch 的边界清晰可控，B2E 算法始终在完整的一批交易范围内做决策，与论文中的理论模型严格对应。
2. **未分配 CTX 重试机制 (Retry Pool)**：首次运行 B2E 时，若某些 CTX 因 broker 余额不足而未能分配，会先进入重试池 `restBrokerRawMegPool`。随着 Tx1 上链、broker 余额得到补充，系统会自动重试，直到池清空或连续 5 秒无进展。
3. **交易价值过滤 (Transaction Filtering)**：读取数据集时，若某笔 CTX 的转账金额超过所有 broker 的最大可用余额上限 (`Init_broker_Balance × ShardNum`，约 20 ETH × 分片数)，则直接过滤丢弃。这类交易任何 broker 都无法承接，留在系统中会长期阻塞重试池。
4. **逐 Epoch 统计记录 (Epoch Stats)**：每完成一个 epoch，自动将 B2E 分配结果、交易统计、broker 累计收益写入 CSV 文件，便于后续绘图分析。

## 代码结构

```
blockEmulator-Broker2Earn/
├── block-emulator/
│   ├── params/
│   │   └── global_config.go          # 全局参数配置文件
│   ├── supervisor/
│   │   ├── committee/
│   │   │   └── committee_broker_b2e.go  # B2E Supervisor 核心逻辑
│   │   └── Broker2Earn/
│   │       ├── B2E.go                # Relax-and-Rounding 主流程
│   │       └── URFA.go               # URFA 线性松弛分配算法
│   ├── data/
│   │   └── selectedTxs_300K.csv      # 内置以太坊历史交易数据集 (30 万笔)
│   └── ipTable.json                  # 节点 IP 配置文件
```

## 快速开始

```bash
# 1. 克隆仓库
git clone <YOUR_GITHUB_URL> blockEmulator-Broker2Earn
cd blockEmulator-Broker2Earn/block-emulator

# 2. 编译
go build -o blockEmulator .

# 3. 生成启动脚本 (4 分片 × 4 节点，启用 B2E)
#    Windows
./blockEmulator.exe -g -S 4 -N 4 -m 4
#    Linux / macOS
./blockEmulator -g -S 4 -N 4 -m 4

# 4. 运行
#    Windows: 双击生成的 .bat 文件
#    Linux / macOS:
bash run_IpAddr=127_0_0_1.sh
```

实验结束后，结果自动写入 `./result/` 目录。

## 参数配置

Broker2Earn 的全部系统参数集中在 `params/global_config.go`，用 Go 变量直接定义，无需额外的 JSON 配置文件。用任意文本编辑器打开该文件并按需修改即可：

```go
var (
    Block_Interval      = 5000    // 出块间隔 (ms)
    MaxBlockSize_global = 2000    // 每块最多交易数
    InjectSpeed         = 5000    // 每轮注入交易条数，同时是 epoch 结束阈值
    TotalDataSize       = 500000  // 实验总交易数 (默认 50 万)
    BatchSize           = 5000    // 每批发送量，建议与 InjectSpeed 一致
    BrokerNum           = 50      // 系统中 broker 账户数量
    NodesInShard        = 4       // 每分片共识节点数
    ShardNum            = 4       // 分片总数
    IterNum_B2E         = 5       // B2E 算法每轮迭代次数
    Brokerage           = 0.1     // broker 佣金比例 (0.1 = 10%)
    DataWrite_path      = "./result/"
    FileInput           = "./data/selectedTxs_300K.csv"
)
```

| 参数 | 含义 |
| --- | --- |
| `ShardNum` | 分片总数 |
| `NodesInShard` | 每分片共识节点数 |
| `BrokerNum` | 系统中初始参与的 broker 账户数量，程序启动时从 `broker/` 目录下按顺序读取前 `BrokerNum` 个地址 |
| `InjectSpeed` | Supervisor 每轮注入的交易数，同时为每个 epoch 的结束阈值 |
| `IterNum_B2E` | B2E 算法每个 epoch 内执行 URFA 贪心分配的迭代次数，值越大分配越精细但开销越高 |
| `Brokerage` | broker 的佣金比例 |
| `TotalDataSize` | 参与实验的交易总条数，总 epoch 数 = `TotalDataSize / InjectSpeed` |
| `MaxBlockSize_global` | 每个区块的最大交易数 |
| `Block_Interval` | 出块时间间隔 (ms) |
| `FileInput` | 输入交易数据集路径 |

### 算法模式开关 (`-m` 参数)

| `-m` | 算法模式 |
| :---: | --- |
| 0 | CLPA_Broker |
| 1 | CLPA |
| 2 | Broker (原始 BrokerChain，无激励) |
| 3 | Relay (默认值) |
| **4** | **Broker_b2e (Broker2Earn)** |

启动时必须显式指定 `-m 4` 才会启用 B2E 算法。

## 节点 IP 配置

若修改了 `ShardNum` 或 `NodesInShard`，需要同步更新根目录下的 `ipTable.json`。以 2 分片 × 2 节点为例：

```json
{
    "0": {"0": "127.0.0.1:32217", "1": "127.0.0.1:32227"},
    "1": {"0": "127.0.0.1:32317", "1": "127.0.0.1:32327"},
    "2147483647": {"0": "127.0.0.1:38800"}
}
```

其中键 `"2147483647"` 代表 Supervisor 节点 (shard ID 十六进制为 `0x7fffffff`)。

## 启动实验

完成参数配置后，在 `block-emulator/` 目录下：

```bash
# 1. 编译 (首次或修改代码后)
go build -o blockEmulator .

# 2. 生成批处理启动脚本
./blockEmulator -g -S 4 -N 4 -m 4
#   -S : 分片数      -N : 每分片节点数      -m 4 : 启用 B2E

# 3. 启动所有节点
# Windows : 双击生成的 .bat
# Linux / macOS:
bash run_IpAddr=127_0_0_1.sh
```

系统会依次注入交易、调用 B2E 算法分配 broker、等待链上确认后再注入下一批，直至全部 `TotalDataSize` 笔交易处理完毕，自动将结果写入 `./result/`。

## 实验结果文件说明

```
./result/
├── supervisor_measureOutput/
│   ├── Average_TPS.csv                  # 每 epoch 系统平均 TPS
│   ├── Transaction_Confirm_Latency.csv  # 每 epoch 交易平均确认时延
│   └── CrossTransaction_ratio.csv       # 每 epoch 跨分片交易比例
├── brokerRsult/
│   ├── <brokerAddr>_brokerBalance.csv   # 各 broker 各分片自由余额变化
│   ├── <brokerAddr>_lockBalance.csv     # 各 broker 各分片锁定余额变化
│   └── <brokerAddr>_Profit.csv          # 各 broker 各分片累计手续费收益
└── epoch_stats/
    └── epoch_stats_Broker2Earn_inject5000_shard4.csv  # 逐 epoch B2E 分配统计
```

`epoch_stats` 文件各列含义：

| 列名 | 含义 |
| --- | --- |
| `Epoch` | epoch 序号 (从 0 开始) |
| `Injected_Tx` | 本 epoch 注入的交易总数 (= `InjectSpeed`) |
| `B2E_Allocated_Tx` | B2E 成功分配给 broker 的跨分片交易数 |
| `B2E_Unallocated_Tx` | 因 broker 余额不足未能分配的跨分片交易数 |
| `Broker_Service_Tx` | 本 epoch 链上确认的 broker 中转交易数 (Tx1 + Tx2) |
| `Inner_Tx` | 本 epoch 链上确认的片内交易数 |
| `Total_Profit_ETH` | 截至本 epoch 所有 broker 累计收益 (ETH) |
| `Active_Broker_Count` | 当前活跃 broker 账户数 |
| `Block_Count` | 本 epoch 产生的区块总数 |

> 注：跨分片交易在系统中被拆为两笔片内交易处理，因此 `Average_TPS.csv` 中跨分片交易按 0.5 笔计入 TPS。

## 实验图绘制示例

下面的 Python 脚本读取上述结果文件，绘制三张核心分析图：系统 TPS、B2E 每 epoch 分配情况、各 broker 累计收益。

```python
import os
import pandas as pd
import matplotlib.pyplot as plt

result_dir = './result/'

# ---- 图 1：系统平均 TPS ----
tps_df = pd.read_csv(
    result_dir + 'supervisor_measureOutput/Average_TPS.csv',
    header=None, names=['TPS']
)
plt.figure(figsize=(8, 4))
plt.plot(tps_df['TPS'].values, marker='o', linewidth=2,
         color='steelblue', label='Broker2Earn')
plt.xlabel('Epoch'); plt.ylabel('Average TPS')
plt.title('System Throughput (Broker2Earn)')
plt.legend(); plt.grid(True, linestyle='--', alpha=0.6)
plt.tight_layout(); plt.savefig('B2E_TPS.pdf', dpi=300); plt.show()

# ---- 图 2：B2E 每 epoch 分配情况 ----
epoch_df = pd.read_csv(
    result_dir + 'epoch_stats/'
    'epoch_stats_Broker2Earn_inject5000_shard4.csv'
)
plt.figure(figsize=(8, 4))
plt.plot(epoch_df['B2E_Allocated_Tx'].values, linewidth=2,
         color='forestgreen', label='Allocated CTX')
plt.plot(epoch_df['B2E_Unallocated_Tx'].values, linewidth=2,
         color='tomato', linestyle='--', label='Unallocated CTX')
plt.xlabel('Epoch'); plt.ylabel('Number of Cross-Shard Txs')
plt.title('B2E Allocation per Epoch')
plt.legend(); plt.grid(True, linestyle='--', alpha=0.6)
plt.tight_layout(); plt.savefig('B2E_Allocation.pdf', dpi=300); plt.show()

# ---- 图 3：各 Broker 累计收益 ----
broker_dir = result_dir + 'brokerRsult/'
profit_files = [f for f in os.listdir(broker_dir) if f.endswith('_Profit.csv')]
plt.figure(figsize=(8, 4))
for fname in profit_files[:3]:
    df = pd.read_csv(os.path.join(broker_dir, fname))
    total = df.sum(axis=1)
    plt.plot(total.values, linewidth=2,
             label=fname.replace('_Profit.csv', '')[:12])
plt.xlabel('Epoch'); plt.ylabel('Cumulative Profit (wei)')
plt.title('Broker Profit over Time')
plt.legend(fontsize=9); plt.grid(True, linestyle='--', alpha=0.6)
plt.tight_layout(); plt.savefig('B2E_BrokerProfit.pdf', dpi=300); plt.show()
```

## 数据集

实验使用的输入数据来自 [xBlock — Ethereum On-chain Data](https://xblock.pro/#/dataset/14)，每条记录包含三个字段：

- `from`：付款方账户地址
- `to`：收款方账户地址
- `value`：转账金额 (Wei)

仓库已内置一份 30 万条交易的数据集 `./data/selectedTxs_300K.csv`，**克隆即可直接使用**。如需更大规模数据，可前往 xBlock 下载，并将 `FileInput` 与 `TotalDataSize` 改为相应值。

## 注意事项

- Linux/macOS 下若提示权限不足，请先执行 `chmod +x blockEmulator`。
- 实验异常中断但后台进程未停止时，可执行 `pkill -9 -f blockEmulator` 强制终止。
- 阻塞式注入要求每批交易**全部上链**后才能继续，因此实验耗时与 `Block_Interval` 直接相关。希望加快进度可适当减小 `TotalDataSize` 或 `Block_Interval`。
- 系统在读取数据时会自动过滤掉转账金额超过 `Init_broker_Balance × ShardNum` (约 20 ETH × 分片数) 的跨分片交易，这些交易超出所有 broker 承接能力，过滤后不计入 `TotalDataSize`，属于正常现象。

## 引用

如果本仓库对您的研究有帮助，欢迎引用 Broker2Earn 与 BrokerChain：

```bibtex
@inproceedings{chen2024broker2earn,
  title     = {Broker2Earn: Towards Maximizing Brokers' Revenue and Cross-Shard Transactions Throughput in Sharded Blockchains},
  author    = {Chen, Qinde and Huang, Huawei and others},
  booktitle = {IEEE INFOCOM 2024},
  year      = {2024}
}

@inproceedings{huang2022brokerchain,
  title     = {BrokerChain: A Cross-Shard Blockchain Protocol for Account/Balance-based State Sharding},
  author    = {Huang, Huawei and others},
  booktitle = {IEEE INFOCOM 2022},
  year      = {2022}
}
```

并请同时引用 BlockEmulator 项目本身。详细使用说明可参考《BlockEmulator 开源项目使用手册》**第 13 章 Broker2Earn 算法**。

## License

参见根目录 [LICENSE](./LICENSE) 文件。
