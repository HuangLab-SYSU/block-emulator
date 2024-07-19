import os
import pandas as pd
import matplotlib.pyplot as plt

# 定义存储CSV文件的目录
directory = './expTest/result/pbft_shardNum=4'

# 获取目录中所有的CSV文件
csv_files = [f for f in os.listdir(directory) if f.endswith('.csv')]

# 创建一个空的图表
plt.figure(figsize=(10, 6))

# 存储文件名（去除共同后缀）
legend_labels = []

# 指定从倒数第几个位置开始删除
x = 1

# 遍历每个CSV文件并绘制折线图
for csv_file in csv_files:
    file_path = os.path.join(directory, csv_file)
    
    # 从CSV文件中读取数据
    df = pd.read_csv(file_path)
    
    # 提取需要的列
    block_height = df['Block Height']
    tx_pool_size = df['TxPool Size']
    
    # 获取不带共同后缀的文件名
    parts = csv_file.split('_')
    if len(parts) > x:
        file_name = '_'.join(parts[:-x])
    else:
        file_name = csv_file
    
    # 绘制折线图
    plt.plot(block_height, tx_pool_size, label=file_name)
    
    # 记录用于图例的文件名
    legend_labels.append(file_name)

# 添加图例和标签
plt.xlabel('Block Height')
plt.ylabel('TxPool Size')
plt.title('TxPool Size vs Block Height')

# 设置图例
plt.legend()

# 显示图表
plt.grid(True)
plt.tight_layout()
plt.show()