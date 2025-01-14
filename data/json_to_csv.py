import json
import csv

# 读取JSON文件
with open('[vel_tic_nus_c_1_c_9_lox_c_3_c_5_c_7]_norm_openb_pod_list_gpushare20_case_range_([0_1]-[100_101])_07-30_01:56:17.json', 'r') as file:
    data = json.load(file)

# 打开CSV文件以写入
with open('output.csv', 'w', newline='') as csvfile:
    csvwriter = csv.writer(csvfile)

    # 写入CSV文件的表头
    headers = ['Scheduler']
    case_ranges = data['case_ranges'][0]
    for case_range in case_ranges:
        headers.append(f'Case Range {case_range}')
    csvwriter.writerow(headers)

    # 写入每个调度器的数据
    for scheduler, reports in data['reports'].items():
        row = [scheduler]
        for report in reports:
            row.append(report['execution']['average_jct_seconds'])
        csvwriter.writerow(row)

print("数据已成功写入output.csv文件")