import csv

# 文件路径
file_path = 'all_case.csv'

# 目标值
gpu_milli_threshold = 45000

def find_threshold_row(file_path, gpu_milli_threshold):
    cumulative_gpu_milli = 0  # 累积的 gpu_milli 值

    # 打开文件并读取
    with open(file_path, mode='r') as file:
        reader = csv.DictReader(file)

        for i, row in enumerate(reader, start=1):
            try:
                # 当前行的 gpu_milli 值
                gpu_milli = float(row['gpu_milli'])

                # 累加 gpu_milli
                cumulative_gpu_milli += gpu_milli

                # 检查累积值是否达到阈值
                if cumulative_gpu_milli >= gpu_milli_threshold:
                    return i, cumulative_gpu_milli  # 返回行号和累积值
            except ValueError:
                print(f"行 {i} 数据有误，跳过：{row}")

    return None, cumulative_gpu_milli  # 如果未找到符合条件的行

# 调用函数
n, total_gpu_milli = find_threshold_row(file_path, gpu_milli_threshold)

if n:
    print(f"满足条件的行号为：{n}")
    print(f"前 {n-1} 行的 gpu_milli 累积值小于 {gpu_milli_threshold}，前 {n} 行的 gpu_milli 累积值为 {total_gpu_milli}")
else:
    print("未找到满足条件的行")