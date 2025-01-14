import pandas as pd
import csv
import os

# 读取CSV文件
df = pd.read_csv('all_case.csv')

# 定义区间
bins = [0, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000]

# 将gpu_milli列中的数据分到区间中
df['gpu_milli_bins'] = pd.cut(df['gpu_milli'], bins)

# 统计每个区间的数量
bin_counts = df['gpu_milli_bins'].value_counts().sort_index()

# 计算每个区间的比例
total_count = df.shape[0]
bin_ratios = bin_counts / total_count

# 输出结果
print("---------------------------------------")
output = []
for bin_range, count in bin_counts.items():
    ratio = bin_ratios[bin_range]
    output.append([bin_range, count, f"{ratio:.2%}"])
    # 修改区间格式为（x, y]形式
    print(f"({bin_range.left}, {bin_range.right}] , {count},{ratio:.2%}")

# 指定保存路径
output_dir = "./"
output_file = os.path.join(output_dir, "output.csv")

# 创建目录（如果不存在）
os.makedirs(output_dir, exist_ok=True)

# 将结果保存到CSV文件
with open(output_file, 'w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(['range', 'amount', 'rate'])
    for row in output:
        # 修改区间格式为（x, y]形式
        writer.writerow([f"({row[0].left}, {row[0].right}]", row[1], row[2]])

print(f"文件已保存到：{output_file}")
