import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

# 读取CSV文件
file_path = './expTest/result/supervisor_measureOutput/Tx_Details.csv'  # 替换为你的CSV文件路径
df = pd.read_csv(file_path)

# 提取 "Confirmed latency of this tx (ms)" 列
latency_column = 'Confirmed latency of this tx (ms)'
latency_data_1 = df[latency_column]

# 提取第二组数据（'Broker1 Tx commit timestamp (not a broker tx -> nil)' 列不为空）
timestamp_column = 'Broker1 Tx commit timestamp (not a broker tx -> nil)'
latency_data_2 = df[df[timestamp_column].notna()][latency_column]

# 提取第三组数据（'Broker1 Tx commit timestamp (not a broker tx -> nil)' 列为空）
latency_data_3 = df[df[timestamp_column].isna()][latency_column]

# 设置绘图风格
sns.set(style='whitegrid')

# 创建一个绘图对象
plt.figure(figsize=(10, 6))

# 绘制核密度估计曲线
density_kwargs = {'lw': 2, 'alpha': 0.7}

# 所有交易
sns.kdeplot(latency_data_1, color='blue', label='All Txs', **density_kwargs)

# Broker交易
sns.kdeplot(latency_data_2, color='orange', label='Broker Txs (CTXs)', **density_kwargs)

# InnerShard交易
sns.kdeplot(latency_data_3, color='green', label='InnerShard Txs', **density_kwargs)

# 添加图例和标题
plt.legend()
plt.title('Distribution of Confirmed Latency of This Tx (ms)')
plt.xlabel('Confirmed Latency (ms)')
plt.ylabel('Density')

# 显示图表
plt.show()
