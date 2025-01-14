import pandas as pd

# 读取CSV文件
df = pd.read_csv('openb_pod_list_gpushare60.csv')

# 统计model列中每种类型的数量
model_counts = df['gpu_milli'].value_counts().sort_index()

total_counts = df.shape[0]
ratios = model_counts/total_counts

# 输出结果
print("--------------------2--------------------")
for bin_range, count in model_counts.items():
    ratio = ratios[bin_range]
    print(f"{bin_range}: {count} ({ratio:.2%})")