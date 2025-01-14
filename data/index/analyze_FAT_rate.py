import csv
import json

# 读取popularity.csv
popularity = {}
with open('popularity.csv', mode='r', encoding='utf-8') as file:
    reader = csv.reader(file)
    next(reader)  # 跳过表头
    for row in reader:
        gpu_milli, _, ratio = row
        popularity[int(gpu_milli)] = float(ratio)

# 读取all_case.csv
all_cases = {}
with open('all_case.csv', mode='r', encoding='utf-8') as file:
    reader = csv.reader(file)
    next(reader)  # 跳过表头
    for row in reader:
        _, job_name, _, gpu_milli, _, _, _, _ = row
        all_cases[job_name] = int(gpu_milli)

# 读取调度报告json文件
with open('150_rate.json', mode='r', encoding='utf-8') as file:
    report = json.load(file)

# 计算每种调度算法下每个case_range内每个作业的FAT
fat_results = {}
for scheduler in report['reports']:
    for schedule in report['reports'][scheduler]:
        scheduler_name = schedule['scheduler_name']
        case_range = tuple(schedule['case_range'])
        fat_sum = 0
        for job in schedule['job_reports']:
            job_name = job['job_name']
            gpu_milli = all_cases[job_name]
            remaining_gpu = 1000 - gpu_milli
            runtime_on_gpu = job['runtime_on_gpu']
            end_time = job['end_time']

            # 计算流行度总和
            total_popularity = sum(ratio for milli, ratio in popularity.items() if milli > remaining_gpu)

            # 计算FAT
            fat = remaining_gpu * runtime_on_gpu * end_time * total_popularity
            fat_sum += fat

        if (scheduler_name, case_range) not in fat_results:
            fat_results[(scheduler_name, case_range)] = 0
        fat_results[(scheduler_name, case_range)] += fat_sum

# 输出结果到CSV文件
with open('fat_results.csv', mode='a', encoding='utf-8', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(['scheduler_name', 'case_range', 'fat_sum'])
    for (scheduler_name, case_range), fat_sum in fat_results.items():
        writer.writerow([scheduler_name, case_range, fat_sum])

print("计算完成，结果已保存到fat_results.csv文件中。")
