import csv

# 初始化总和变量
gpu_sum = 0

# 打开CSV文件
with open('openb_node_list_gpu_node_normalization.csv', mode='r', encoding='utf-8') as file:
    csv_reader = csv.DictReader(file)
    
    # 遍历每一行并累加gpu列的值
    for row in csv_reader:
        gpu_sum += int(row['gpu'])

print(f"GPU列的总和是: {gpu_sum}")