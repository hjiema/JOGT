import csv
import re

# 定义输入文件和输出文件
input_file = "num_all_case.txt"  # 日志文件名
output_file = "runtime_results.csv"  # 输出 CSV 文件名

# 定义正则表达式匹配调度算法、CaseRange 和 Runtime
scheduler_pattern = re.compile(r"Starting Simulation For Scheduler (\w+), StartTime:")
case_range_pattern = re.compile(r"CaseRange: \[(\d+) (\d+)\], StartTime:")
runtime_pattern = re.compile(r"RunTime: ([\d.]+)")

# 存储提取的结果
results = []

# 读取日志文件并提取数据
with open(input_file, "r") as f:
    current_scheduler = None
    current_case_range = None

    for line in f:
        # 匹配调度算法
        scheduler_match = scheduler_pattern.search(line)
        if scheduler_match:
            current_scheduler = scheduler_match.group(1)

        # 匹配 CaseRange
        case_range_match = case_range_pattern.search(line)
        if case_range_match:
            current_case_range = f"[{case_range_match.group(1)} {case_range_match.group(2)}]"

        # 匹配 Runtime
        runtime_match = runtime_pattern.search(line)
        if runtime_match and current_scheduler and current_case_range:
            runtime = runtime_match.group(1)
            # 保存当前记录
            results.append([current_scheduler, current_case_range, runtime])
            current_case_range = None  # 重置 CaseRange

# 将结果写入 CSV 文件
with open(output_file, "w", newline="") as csvfile:
    writer = csv.writer(csvfile)
    # 写入表头
    writer.writerow(["Scheduler", "CaseRange", "Runtime"])
    # 写入数据
    writer.writerows(results)

print(f"提取完成！结果已保存到 {output_file}")
