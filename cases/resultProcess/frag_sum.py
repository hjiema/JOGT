import pandas as pd

# 读取CSV文件
df = pd.read_csv('frag_hydrababwithheuristic_3_[0-10]_result.csv')

# 计算frag_size列的总和
frag_size_sum = df['frag_size'].sum()

# 计算frag_size_runtime列的总和
frag_size_runtime_sum = df['frag_size_runtime'].sum()

# 输出结果
print(f"frag_hydrababwithheuristic3_[0-10]_result的frag_size列的总和: {frag_size_sum}")
print(f"frag_hydrababwithheuristic3_[0-1]_result的frag_size_runtime列的总和: {frag_size_runtime_sum}")

# # 将结果写入新的CSV文件
# result_df = pd.DataFrame({
#     'metric': ['frag_size_sum', 'frag_size_runtime_sum'],
#     'value': [frag_size_sum, frag_size_runtime_sum]
# })
#
# result_df.to_csv('output_summary.csv', index=False)
#
# print("结果已写入output_summary.csv文件")
