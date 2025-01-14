import json
import csv
import os

def calculate_average_end_time(data):
    results = {}
    reports = data['reports']

    for scheduler, scheduler_data in reports.items():
        end_times = []
        for report in scheduler_data:
            job_reports = report['job_reports']
            for job in job_reports:
                end_times.append(job['end_time'])

        average_end_time = sum(end_times) / len(end_times) if end_times else 0
        results[scheduler] = average_end_time

    return results

# 读取调度报告的JSON文件
with open('[Fgt_tic_lox_vel_nus]_all_case_case_range_([0_100]-[0_100])_01-09_09:53:52.json', 'r') as file:
    data = json.load(file)

# 计算平均 end_time
average_end_times = calculate_average_end_time(data)

# 定义CSV文件路径
csv_file_path = 'average_end_times.csv'

# 检查CSV文件是否存在，判断是否需要写入标题行
file_exists = os.path.isfile(csv_file_path)

# 写入结果到CSV文件
with open(csv_file_path, 'a', newline='') as csvfile:
    fieldnames = ['scheduler', 'average_end_time']
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

    if not file_exists:
        writer.writeheader()  # 写入标题行

    for scheduler, avg_end_time in average_end_times.items():
        writer.writerow({'scheduler': scheduler, 'average_end_time': avg_end_time})

print("结果已写入到 average_end_times.csv 文件中。")

