# BlockEmulator · Broker2Earn 分支

> Broker2Earn (B2E) —— 分片区块链上面向 Broker 的可持续激励机制，基于 [BlockEmulator](https://github.com/HuangLab-SYSU/block-emulator) 与 BrokerChain 实现。
>
> 论文：H. Huang, Q. Chen, et al., *"Broker2Earn: Towards Maximizing Brokers' Revenue and Cross-Shard Transactions Throughput in Sharded Blockchains"*, **IEEE INFOCOM 2024**.
>
> 仓库地址：<https://github.com/HuangLab-SYSU/block-emulator/tree/broker2earn>

本分支为 Broker2Earn 协议在 BlockEmulator 上的开源实现，由 BrokerChain 的 INFOCOM 2022 会议版代码修改而来。配合《BlockEmulator 开源项目使用手册》第 13 章一起阅读，可一步步复现论文中的实验。

---

## 目录

- [研究背景](#研究背景)
- [B2E 协议简介](#b2e-协议简介)
- [本分支相对 BrokerChain 的新增功能](#本分支相对-brokerchain-的新增功能)
- [代码结构](#代码结构)
- [快速开始](#快速开始)
- [参数配置](#参数配置)
- [节点 IP 配置](#节点-ip-配置)
- [编译与启动实验](#编译与启动实验)
- [实验数据集](#实验数据集)
- [实验结果文件说明](#实验结果文件说明)
- [实验图绘制示例](#实验图绘制示例)
- [注意事项](#注意事项)
- [引用](#引用)

---

## 研究背景

分片技术是保持区块链去中心化、同时提升其可扩展性的可行路线。但在状态分片下，一笔涉及不同分片账户的交易（**跨分片交易，CTX**）处理代价远高于片内交易，过高的 CTX 比例会显著拖累整体吞吐。

BrokerChain (INFOCOM '22) 通过引入 **broker 账户**，把一笔跨分片交易拆解为两笔片内交易，缓解了 CTX 瓶颈，但留下了一个开放问题：

> **如何吸引足够多的志愿者主动质押 Token、自愿成为 broker，并把这些流动性合理分配到各分片中？**

Broker2Earn 正是为回答这一问题而设计的激励机制。

## B2E 协议简介

Broker2Earn 的运作流程：

1. **注册 Broker**：志愿者 (Volunteer) 通过智能合约注册申请成为 broker 账户；
2. **质押 Token**：broker 将一定数量的 Token 质押到合约中，为分片区块链提供流动性；
3. **B2E 算法分配**：B2E 把各 broker 账户分配到不同分片，以最优化流动性使用；
4. **处理跨分片交易**：跨分片交易 `CTX⟨S_i, S_j, vol, fee⟩` 经由 broker 中转，被拆解为两笔片内交易；
5. **broker 获得收益**：每笔成功中转的 CTX 中，broker 按佣金比例（默认 10%）收取手续费分成。

B2E 在算法层面将上述招募—分配问题建模为一个 utility 最大化问题，证明其为 **NP-hard**，并提出基于 **Relax-and-Rounding** 的 **在线近似算法 (Online Approximation Algorithm)**：

- **Relax**：将 0/1 整数变量松弛为连续变量，转为线性规划 (LP)，多项式时间内精确求解；
- **Rounding**：对 LP 分数解做随机舍入，并给出严格的近似比保证；
- **Online**：随着交易流的实时到达进行决策，无需提前知道全局交易信息。

算法细节与近似比证明详见 [Broker2Earn 论文](#引用)。

## 本分支相对 BrokerChain 的新增功能

本分支基于 BrokerChain 的 INFOCOM 2022 会议版代码，保留原 broker 机制的基础上新增了 **四项** 功能：

1. **阻塞式注入机制 (Blocking Injection)**：每批交易注入后，Supervisor 会等待该批次所有交易在链上完成确认后再注入下一批。当累计确认数（片内交易数 + 跨分片交易对数）达到 `params.InjectSpeed` 时，当前 epoch 终止并触发下一批注入。每个 epoch 的边界明确且可控，B2E 算法始终在完整批次的交易范围内进行分配决策，与论文中的理论模型严格一致。
2. **未分配 CTX 重试机制 (Retry Pool)**：B2E 算法首次运行时，若某些 CTX 因 broker 账户余额不足而未能分配，这些交易会进入重试池 `restBrokerRawMegPool`。随着对应 Tx1 完成上链、broker 账户余额逐步恢复，系统会自动对重试池中的 CTX 重新执行 B2E，直到池清空或连续 5 秒无新分配。
3. **交易金额过滤 (Transaction Filtering)**：读取数据集时，若某笔 CTX 的转账金额超过所有 broker 账户的最大可用余额上限（`Init_broker_Balance × ShardNum`，约 **2 ETH × 分片数**），系统会直接过滤该交易。此类交易超出任何 broker 的承接能力，若不予过滤将长期占据重试池，导致实验无法正常终止。
4. **逐 Epoch 统计记录 (Epoch Stats)**：每个 epoch 结束后，系统自动将该 epoch 的 B2E 分配结果、交易统计与 broker 累计收益等数据写入 CSV 文件，便于后续数据分析与可视化。

## 代码结构

```
blockEmulator-Broker2Earn/
├── block-emulator/
│   ├── params/
│   │   └── global_config.go             # 全局参数配置文件
│   ├── supervisor/
│   │   ├── committee/
│   │   │   └── committee_broker_b2e.go  # B2E Supervisor 核心逻辑
│   │   └── Broker2Earn/
│   │       ├── B2E.go                   # Relax-and-Rounding 主流程
│   │       └── URFA.go                  # URFA 线性松弛分配算法
│   ├── data/
│   │   └── selectedTxs_300K.csv         # 内置以太坊历史交易数据集
│   └── ipTable.json                     # 节点 IP 配置文件
```

## 快速开始

```bash
# 1. clone 并切到 broker2earn 分支
git clone -b broker2earn https://github.com/HuangLab-SYSU/block-emulator.git
cd block-emulator/block-emulator

# 2. 生成启动脚本（4 分片 × 4 节点，启用 B2E）
go run main.go -g -S 4 -N 4 -m 4

# 3. 启动所有节点
#    Windows : 双击生成的 .bat 文件，或在命令行直接运行
#    Linux/macOS:
bash run_IpAddr=127_0_0_1.sh
```

实验结束后，所有结果自动写入 `./result/` 目录。

## 参数配置

B2E 的全部系统参数集中在 `params/global_config.go`，以 Go 变量形式直接声明，无需 JSON 配置文件。用任意文本编辑器打开该文件并按需修改即可：

```go
var (
    Block_Interval      = 5000    // 出块间隔 (ms)
    MaxBlockSize_global = 2000    // 每块最多交易数
    InjectSpeed         = 5000    // 每轮注入交易条数 (= 每个 epoch 的交易数)
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
| `BrokerNum` | 系统初始参与的 broker 账户数量，程序启动时从 `broker/` 目录按顺序读取前 `BrokerNum` 个地址 |
| `InjectSpeed` | Supervisor 每轮注入的交易数，同时为每个 epoch 的结束阈值 |
| `IterNum_B2E` | B2E 算法每个 epoch 内执行 URFA 贪心分配的迭代次数，值越大分配越精细，开销也越高 |
| `Brokerage` | broker 的佣金比例 |
| `TotalDataSize` | 参与实验的交易总数。总 epoch 数 = `TotalDataSize / InjectSpeed` |
| `MaxBlockSize_global` | 每个区块所容纳的最大交易数 |
| `Block_Interval` | 出块时间间隔 (ms) |
| `FileInput` | 输入交易数据集的文件路径 |

### 通过 `-m` 参数选择算法

BlockEmulator 内置多种交易处理算法。启动时通过命令行参数 `-m <编号>` 指定本次实验使用的算法：

| `-m` 值 | 算法模式 |
| :---: | --- |
| 0 | CLPA_Broker |
| 1 | CLPA |
| 2 | Broker（原始 BrokerChain，无激励） |
| 3 | Relay（默认值） |
| **4** | **Broker_b2e (Broker2Earn)** |

启动时必须显式指定 **`-m 4`** 才会启用 B2E 算法。

## 节点 IP 配置

若修改了 `ShardNum` 或 `NodesInShard`，需要同步更新根目录下的 `ipTable.json`。以 2 分片 × 2 节点为例：

```json
{
    "0": {"0": "127.0.0.1:32217", "1": "127.0.0.1:32227"},
    "1": {"0": "127.0.0.1:32317", "1": "127.0.0.1:32327"},
    "2147483647": {"0": "127.0.0.1:38800"}
}
```

其中键 `"2147483647"` 代表 Supervisor 节点 ID（其十六进制为 `0x7fffffff`）。

## 编译与启动实验

完成参数配置后，在 `block-emulator/` 根目录下按以下步骤启动（以 4 分片 × 4 节点为例）：

1. **生成批处理启动脚本。** 该命令会基于 `main.go` 即时编译并生成启动所有节点的批处理脚本：

   ```bash
   go run main.go -g -S 4 -N 4 -m 4
   ```

   - `-S 4`：分片总数 4
   - `-N 4`：每分片节点数 4
   - `-m 4`：启用 Broker2Earn 算法

   执行完毕后，终端会打印：
   ```
   [BlockEmulator] Generating launch scripts: shardNum=4, nodesInShard=4, mode=4 (Broker_b2e)
   [BlockEmulator] Done. Batch files ... have been written ...
   ```
   并在当前目录生成批处理文件（Windows 为 `.bat`，Linux/macOS 为 `.sh`）。

2. **运行批处理脚本以启动所有节点。**

   - Windows：双击或命令行运行生成的 `.bat`；
   - Linux / macOS：
     ```bash
     bash run_IpAddr=127_0_0_1.sh
     ```

   脚本会针对每个节点依次执行：
   ```
   go run main.go -n <node-id> -N <每片节点数> -s <shard-id> -S <分片总数> -m 4
   ```
   从而拉起所有共识节点与 Supervisor。

3. **等待实验自动结束。** 系统依次注入交易、调用 B2E 算法分配 broker、等待该批上链确认后再注入下一批，直至 `TotalDataSize` 笔交易全部处理完毕，自动将结果写入 `./result/`。

## 实验数据集

实验所用数据集来自 **以太坊主网区块高度 19000000–19499999** 区间的历史交易记录，已按如下规则进行预处理：

- 仅保留跨分片交易，剔除片内交易（B2E 关注 broker 机制对 CTX 的处理效果）；
- 剔除智能合约调用类交易（即收方或发方为合约地址的交易）；
- 去除自转账与字段异常等无效记录。

经上述处理后，最终保留 **600,000** 笔交易，与 Broker2Earn 原始论文实验规模一致。每条记录包含三个字段：

- `from`：付款方账户地址
- `to`：收款方账户地址
- `value`：转账金额（单位：Wei）

### 下载链接

> **百度网盘** （文件名：`Tx-dataset_forB2E_600K-TXs-CHEN-Qinde-2026May12.csv`）
>
> 链接：<https://pan.baidu.com/s/1Ap8S2njayTOqTj9lB6yBzw?pwd=1234>　提取码：`1234`

下载后将 `params/global_config.go` 中的 `FileInput` 指向该文件，并将 `TotalDataSize` 设为 `600000` 即可。

> 仓库 `./data/` 下另内置了一份小规模数据 `selectedTxs_300K.csv`，便于快速验证流程。

## 实验结果文件说明

实验结束后，所有输出存放于 `./result/`，主要包括三类 CSV 文件：

```
./result/
├── supervisor_measureOutput/
│   ├── Average_TPS.csv                  # 每 epoch 平均吞吐量
│   ├── Transaction_Confirm_Latency.csv  # 每 epoch 交易平均确认时延
│   └── CrossTransaction_ratio.csv       # 每 epoch 跨分片交易比例
├── brokerResult/
│   ├── <brokerAddr>_brokerBalance.csv   # 各 broker 各分片自由余额
│   ├── <brokerAddr>_lockBalance.csv     # 各 broker 各分片锁定余额
│   └── <brokerAddr>_Profit.csv          # 各 broker 各分片累计佣金收益
└── epoch_stats/
    └── epoch_stats_Broker2Earn_inject5000_shard4.csv  # 逐 epoch B2E 统计
```

`epoch_stats` 文件各列含义：

| 列名 | 含义 |
| --- | --- |
| `Epoch` | epoch 序号（从 0 开始） |
| `Tx_injection_speed` | 本 epoch 内的交易注入速度 |
| `CTXs_served_by_Alg` | 由算法（如 B2E）成功分配给 broker 的跨分片交易数 |
| `CTXs_unserved_by_Alg` | 由算法（如 B2E）未能成功分配给 broker 的跨分片交易数 |
| `CTXs_served_by_Broker` | 本 epoch 内实际上链的、经 broker 中转的 CTX 数量（用于校验算法分配结果是否正确上链） |
| `ITX_count_handled_thisEpoch` | 本 epoch 内上链确认的片内交易数 |
| `Total_Profit_ETH` | 截至本 epoch 全部 broker 的累计收益（ETH） |
| `Active_Broker_Count` | 当前活跃 broker 账户数 |
| `Block_Count_thisEpoch` | 本 epoch 内产生的区块总数 |

## 实验图绘制示例

以下脚本读取两份 `epoch_stats` 文件（B2E 与 BrokerChain），绘制三联箱线图对比两个算法在每个 epoch 上的分布特征：

```python
import pandas as pd
import matplotlib.pyplot as plt

# 两份 epoch_stats CSV 路径
B2E_CSV = './epoch_stats_Broker2Earn_inject5000_shard4.csv'
BC_CSV  = './epoch_stats_Brokerchain_inject5000_shard4.csv'

df_b2e = pd.read_csv(B2E_CSV)
df_bc  = pd.read_csv(BC_CSV)

# 三个对比指标
metrics = [
    ('CTXs_served_by_Alg',          'CTXs served'),
    ('CTXs_unserved_by_Alg',        'CTXs unserved'),
    ('ITX_count_handled_thisEpoch', 'ITXs served'),
]
# 配色: B2E 浅蓝 + 深蓝边; BrokerChain 浅橙 + 深橙边
palette = [('#A8C8E1', '#2E5F9C'), ('#F2C28C', '#C9701B')]

fig, axes = plt.subplots(1, len(metrics), figsize=(3.4 * len(metrics), 5.6))
for ax, (col, label) in zip(axes, metrics):
    data = [df_b2e[col].dropna().values, df_bc[col].dropna().values]
    bp = ax.boxplot(
        data, patch_artist=True, widths=0.6,
        medianprops=dict(color='#C8423A', linewidth=2.4),
        flierprops=dict(marker='o', markersize=6, alpha=0.6,
                        markerfacecolor='#888888', markeredgecolor='#888888'),
    )
    for i, (face, edge) in enumerate(palette):
        bp['boxes'][i].set_facecolor(face)
        bp['boxes'][i].set_edgecolor(edge)
        bp['boxes'][i].set_linewidth(2.0)
        for line in (bp['whiskers'][2*i], bp['whiskers'][2*i+1],
                     bp['caps'][2*i],     bp['caps'][2*i+1]):
            line.set_color(edge); line.set_linewidth(1.8)
    ax.set_xticks([1, 2])
    ax.set_xticklabels(['B2E', 'BrokerChain'], fontsize=17)
    ax.set_title(label, fontsize=17, fontweight='bold')
    ax.tick_params(axis='y', labelsize=17)
    ax.grid(True, axis='y', linestyle='--', alpha=0.45)

fig.suptitle('B2E vs BrokerChain  ·  Per-Epoch Distribution',
             fontsize=20, fontweight='bold')
plt.tight_layout(rect=[0, 0, 1, 0.91])
plt.savefig('B2E_vs_BrokerChain_boxplot.pdf', dpi=300, bbox_inches='tight')
plt.show()
```

典型结果：B2E 服务的 CTX 数量显著高于 BrokerChain；且 B2E 能将注入的跨分片交易基本全部完成；由于两组实验使用同一份数据，ITX 数量基本一致。

## 注意事项

- Linux/macOS 下若执行 `.sh` 提示「权限不足」，可改用 `bash run_*.sh` 直接调用解释器执行，无需为脚本本身赋予可执行权限。
- 若实验异常中断但后台节点进程未自动终止，可在 Linux/macOS 下执行 `pkill -9 -f "go run main.go"` 强制终止所有相关进程；Windows 下可在任务管理器中结束相关 `go.exe` 与节点进程。
- 阻塞式注入机制要求每批交易**全部上链**后方可继续，因此实验总耗时与 `Block_Interval` 直接相关。希望加快进度可适当减小 `TotalDataSize` 或缩短 `Block_Interval`。
- 程序在读取交易数据时会自动过滤掉转账金额超过 `Init_broker_Balance × ShardNum`（即约 **2 ETH × 分片数**）的跨分片交易，被过滤后不计入 `TotalDataSize`。

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

并请同时引用 BlockEmulator 项目本身。完整使用说明可参考《BlockEmulator 开源项目使用手册》**第 13 章 Broker2Earn 算法**。

## License

参见根目录 [LICENSE](./LICENSE) 文件。
