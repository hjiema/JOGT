import pandas as pd

# 读取CSV文件
data = pd.read_csv('openb_node_list_gpu_node.csv')

# 删除model列为空的行
data = data.dropna(subset=['model'])

# 仅保留model列值为T4、P100和V100开头的行
data = data[data['model'].str.startswith(('T4', 'P100', 'V100'))]

# 将V100M32和V100M16统一为V100
data['model'] = data['model'].replace({'V100M32': 'V100', 'V100M16': 'V100'})

# 保存处理后的数据到新的CSV文件
data.to_csv('openb_node_list_gpu_node_normalization.csv', index=False)

#print("openb_node_list_gpu_node_normalization.csv")