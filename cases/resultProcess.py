import json
import csv

# 处理调度结果json文件的第一步，提取json文件关键信息，输入到combined_result.csv文件中。


# 读取文件1
with open('result.json', 'r') as f1:
    data1 = json.load(f1)

# 读取文件2
file2_data = {}
with open('norm_openb_pod_list_gpushare20.csv', 'r') as f2:
    reader = csv.DictReader(f2)
    for row in reader:
        file2_data[row['job_name']] = row

# 生成CSV文件
with open('combined_result.csv', 'w', newline='') as csvfile:
    fieldnames = ['filename', 'job_name', 'selected_gpu_id', 'selected_gpu_type', 'runtime_on_gpu', 'job_jct', 'gpu_milli', 'A100', 'GTX2080Ti', 'V100', 'ddl']
    writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
    writer.writeheader()

    for algorithm, reports in data1['reports'].items():
        for report in reports:
            case_range = report['case_range']
            case_range_str = f"[{case_range[0]}-{case_range[1]}]"
            filename = f"{algorithm.lower()}_{case_range_str}"

            for job_report in report['job_reports']:
                job_name = job_report['job_name']
                if job_name in file2_data:
                    row = {
                        'filename': filename,
                        'job_name': job_report['job_name'],
                        'selected_gpu_id': job_report['selected_gpu_id'],
                        'selected_gpu_type': job_report['selected_gpu_type'],
                        'runtime_on_gpu': job_report['runtime_on_gpu'],
                        'job_jct': job_report['end_time'],
                        'gpu_milli': file2_data[job_name]['gpu_milli'],
                        'A100': file2_data[job_name]['A100'],
                        'GTX2080Ti': file2_data[job_name]['GTX2080Ti'],
                        'V100': file2_data[job_name]['V100'],
                        'ddl': file2_data[job_name]['ddl']
                    }
                    writer.writerow(row)

print("数据已合并并写入combined_result.csv文件")
