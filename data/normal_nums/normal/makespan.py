import json
import csv

# 读取JSON文件
with open('num_all_case.json', 'r') as f:
    data = json.load(f)

# 初始化结果列表
results = []

# 提取每种调度算法在不同case_range下作业"end_time"的最大值
for scheduler_name, reports in data['reports'].items():
    for report in reports:
        case_range = tuple(report['case_range'])
        max_end_time = max(job['end_time'] for job in report['job_reports'])
        result = {
            'Scheduler': scheduler_name,
            'Case Range': case_range,
            'Max End Time': max_end_time
        }
        results.append(result)

# 输出为CSV文件
with open('makespan.csv', 'w', newline='') as f:
    writer = csv.DictWriter(f, fieldnames=['Scheduler', 'Case Range', 'Max End Time'])
    writer.writeheader()
    for result in results:
        writer.writerow(result)

print("数据提取完成并输出到CSV文件。")
