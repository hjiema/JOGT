import json
import csv

# 读取JSON文件
def load_json(file_path):
    with open(file_path, 'r', encoding='utf-8') as f:
        return json.load(f)

# 提取调度器信息
def extract_scheduler_data(data):
    scheduler_info = []
    
    reports = data.get("reports", {})
    
    for scheduler, scheduler_reports in reports.items():
        for report in scheduler_reports:
            scheduler_name = report.get("scheduler_name")
            case_range = report.get("case_range")
            execution = report.get("execution", {})
            
            # 提取需要的字段
            average_jct_seconds = execution.get("average_jct_seconds")
            average_queue_delay_seconds = execution.get("average_queue_delay_seconds")
            average_ddl_violation_duration_seconds = execution.get("average_ddl_violation_duration_seconds")
            total_ddl_violation_duration_seconds = execution.get("total_ddl_violation_duration_seconds")
            
            scheduler_info.append({
                "scheduler_name": scheduler_name,
                "case_range": case_range,
                "average_jct_seconds": average_jct_seconds,
                "average_queue_delay_seconds": average_queue_delay_seconds,
                "average_ddl_violation_duration_seconds": average_ddl_violation_duration_seconds,
                "total_ddl_violation_duration_seconds": total_ddl_violation_duration_seconds
            })
    
    return scheduler_info

# 将结果保存到 CSV 文件
def save_to_csv(scheduler_info, output_file):
    # CSV 文件的表头
    headers = ["scheduler_name", "case_range", "average_jct_seconds", 
               "average_queue_delay_seconds", "average_ddl_violation_duration_seconds", 
               "total_ddl_violation_duration_seconds"]

    with open(output_file, mode='w', newline='', encoding='utf-8') as f:
        writer = csv.DictWriter(f, fieldnames=headers)

        # 写入表头
        writer.writeheader()

        # 写入数据行
        for info in scheduler_info:
            writer.writerow(info)

# 主函数
if __name__ == "__main__":
    file_path = "[tic_lox_vel_nus_Fgt]_30_ddl_case_range_([0_100]-[0_200])_01-02_09:09:48.json"  # 替换文件路径
    output_file = "scheduler_info11.csv"  # 输出的 CSV 文件名称

    data = load_json(file_path)
    scheduler_info = extract_scheduler_data(data)
    
    # 保存结果到CSV文件
    save_to_csv(scheduler_info, output_file)
    print(f"调度器信息已保存到 {output_file}")
