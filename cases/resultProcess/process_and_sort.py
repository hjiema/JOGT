import json
import csv

# 文件路径
result_json_file = 'result2.json'
pod_list_csv_file = 'norm_openb_pod_list_gpushare20.csv'
output_csv_file = 'combined_result2.csv'
sorted_csv_file = 'sorted_result2.csv'

# 第一步：提取调度结果并生成combined_result.csv
# 读取文件1
with open(result_json_file, 'r') as f1:
    data1 = json.load(f1)

# 读取文件2
file2_data = {}
with open(pod_list_csv_file, 'r') as f2:
    reader = csv.DictReader(f2)
    for row in reader:
        file2_data[row['job_name']] = row

# 生成combined_result.csv文件
with open(output_csv_file, 'w', newline='') as csvfile:
    fieldnames = ['scheduler', 'job_name', 'selected_gpu_id', 'selected_gpu_type', 'runtime_on_gpu', 'start_time', 'job_jct', 'gpu_milli', 'ddl']
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
    writer.writeheader()

    for algorithm, reports in data1['reports'].items():
        for report in reports:
            case_range = report['case_range']
            case_range_str = f"[{case_range[0]}-{case_range[1]}]"
            scheduler = f"{algorithm.lower()}_{case_range_str}"

            for job_report in report['job_reports']:
                job_name = job_report['job_name']
                if job_name in file2_data:
                    row = {
                        'scheduler': scheduler,
                        'job_name': job_report['job_name'],
                        'selected_gpu_id': job_report['selected_gpu_id'],
                        'selected_gpu_type': job_report['selected_gpu_type'],
                        'runtime_on_gpu': job_report['runtime_on_gpu'],
                        'start_time': job_report['start_time'],
                        'job_jct': job_report['end_time'],
                        'gpu_milli': file2_data[job_name]['gpu_milli'],
                        'ddl': file2_data[job_name]['ddl']
                    }
                    writer.writerow(row)

print(f"数据已合并并写入 {output_csv_file} 文件")

# 第二步：对combined_result.csv进行排序并生成sorted_result.csv
# 读取combined_result.csv
with open(output_csv_file, 'r') as csvfile:
    reader = csv.DictReader(csvfile)
    rows = list(reader)

# 排序：按scheduler、selected_gpu_id、start_time依次排序
rows.sort(key=lambda x: (x['scheduler'], int(x['selected_gpu_id']), float(x['start_time'])))

# 写入排序后的CSV文件
with open(sorted_csv_file, 'w', newline='') as csvfile:
    fieldnames = ['scheduler', 'job_name', 'selected_gpu_id', 'selected_gpu_type', 'runtime_on_gpu', 'start_time', 'job_jct', 'gpu_milli', 'ddl']
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
    writer.writeheader()
    writer.writerows(rows)

print(f"排序后的数据已写入 {sorted_csv_file} 文件")
