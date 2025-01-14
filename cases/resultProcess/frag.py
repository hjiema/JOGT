import pandas as pd


# 数据处理第二步，frag处理第一步。输入combined_result.csv文件，输出计算后的各种指标，保存到flie_with_frag_size.csv文件中。

# 读取文件1和文件2
file1 = pd.read_csv('combined_result.csv')
file2 = pd.read_csv('popularity.csv')

# 定义计算frag_size和frag_size_runtime的函数
def calculate_frag_size(gpu_milli, runtime_on_gpu, job_jct):
    # frag_size = (1000 - gpu_milli) * file2[file2['gpu_milli'] > (1000 - gpu_milli)]['ratio'].sum()
    frag_size = (1000 - gpu_milli)
    frag_size_runtime = frag_size * runtime_on_gpu
    fgt = frag_size_runtime * job_jct
    return frag_size, frag_size_runtime, fgt

# 为文件1添加新的列，并将类型设置为float64
file1['frag_size'] = 0.0
file1['frag_size_runtime'] = 0.0

# 计算并添加frag_size和frag_size_runtime
for index, row in file1.iterrows():
    frag_size, frag_size_runtime, fgt = calculate_frag_size(row['gpu_milli'], row['runtime_on_gpu'], row['job_jct'])
    file1.at[index, 'frag_size'] = frag_size
    file1.at[index, 'frag_size_runtime'] = frag_size_runtime
    file1.at[index, 'frag_size_jct'] = fgt

# 保存结果到新的CSV文件
file1.to_csv('file2_with_frag_size.csv', index=False)
