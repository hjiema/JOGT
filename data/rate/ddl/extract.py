import json
import csv

# 读取JSON文件
with open('ddl_rate_200.json', 'r') as f:
    data = json.load(f)

# 初始化结果列表
results = []

# 提取每种调度算法在不同case_range下的数据
for scheduler_name, reports in data['reports'].items():
    for report in reports:
        case_range = tuple(report['case_range'])
        execution = report['execution']
        result = {
            'Scheduler': scheduler_name,
            'Case Range': case_range,
            'Average JCT Seconds': execution['average_jct_seconds'],
            'Average Queue Delay Seconds': execution['average_queue_delay_seconds'],
            'Average DDL Violation Duration Seconds': execution['average_ddl_violation_duration_seconds'],
            'Total DDL Violation Duration Seconds': execution['total_ddl_violation_duration_seconds'],
            'DDL Violated Jobs Count': execution['ddl_violated_jobs_count']
        }
        results.append(result)

# 输出为CSV文件
with open('ddl_scheduler_report.csv', 'a', newline='') as f:
    writer = csv.DictWriter(f, fieldnames=['Scheduler', 'Case Range', 'Average JCT Seconds', 'Average Queue Delay Seconds', 'Average DDL Violation Duration Seconds', 'Total DDL Violation Duration Seconds', 'DDL Violated Jobs Count'])
    writer.writeheader()
    for result in results:
        writer.writerow(result)

print("数据提取完成并输出到CSV文件。")
