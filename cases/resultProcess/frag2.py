import pandas as pd

# 数据处理第三步，frag计算统计第二步，输入file_with_frag_size.csv，统计结果输出到

# 读取输入文件
file = pd.read_csv('file2_with_frag_size.csv')

# 按filename分组并计算总和
grouped = file.groupby('filename').agg({
    'frag_size': 'sum',
    'frag_size_runtime': 'sum',
    'frag_size_jct': 'sum'
}).reset_index()

# 保存结果到新的CSV文件
grouped.to_csv('analyzed_file2.csv', index=False)
