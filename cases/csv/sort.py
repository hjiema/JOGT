import pandas as pd

# 读取CSV文件
df = pd.read_csv('norm_openb_pod_list_gpushare20.csv')

# 确保列名是字符串
df_sorted = df.sort_values('V100')

# 将排序后的数据写回到新的 CSV 文件
df_sorted.to_csv('sorted_output1.csv', index=False)

print("Sorted CSV file has been saved as 'sorted_output.csv'")
