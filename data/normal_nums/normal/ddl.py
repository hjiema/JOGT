import json
import csv
from collections import defaultdict

# 读取JSON文件
with open('num_all_case.json', 'r') as f:
    data = json.load(f)

# 读取CSV文件
ddl_requirements = {}
with open('all_case.csv', 'r') as f:
    reader = csv.DictReader(f)
    for row in reader:
        job_name = row['job_name']
        ddl = row['ddl']
        ddl_requirements[job_name] = float(ddl) if ddl != 'inf' else float('inf')

# 初始化结果字典
violation_rate = defaultdict(lambda: defaultdict(float))
average_violation_time = defaultdict(lambda: defaultdict(float))

# 计算每种调度算法在每种case_range下的违规情况
for scheduler_name, reports in data['reports'].items():
    for report in reports:
        total_jobs = 0
        total_violation_time = 0
        violated_jobs_count = 0
        case_range = tuple(report['case_range'])

        for job in report['job_reports']:
            job_name = job['job_name']
            end_time = job['end_time']
            ddl = ddl_requirements.get(job_name, float('inf'))
            
            if end_time > ddl:
                violated_jobs_count += 1
                total_violation_time += (end_time - ddl)
        
        total_jobs = len(report['job_reports'])
        
        violation_rate[scheduler_name][case_range] = violated_jobs_count / total_jobs
        average_violation_time[scheduler_name][case_range] = total_violation_time / total_jobs

# 输出违规率CSV文件
with open('violation_rate.csv', 'w', newline='') as f:
    writer = csv.writer(f)
    writer.writerow(['Scheduler', 'Case Range', 'Violation Rate'])
    for scheduler, case_ranges in violation_rate.items():
        for case_range, rate in case_ranges.items():
            writer.writerow([scheduler, case_range, rate])

# 输出平均违规时间CSV文件
with open('average_violation_time.csv', 'w', newline='') as f:
    writer = csv.writer(f)
    writer.writerow(['Scheduler', 'Case Range', 'Average Violation Time'])
    for scheduler, case_ranges in average_violation_time.items():
        for case_range, avg_time in case_ranges.items():
            writer.writerow([scheduler, case_range, avg_time])

print("计算完成并输出结果到CSV文件。")
