import pandas as pd

def filter_and_count(csv_file):
    # 读取CSV文件
    df = pd.read_csv(csv_file)

    # 筛选出pod_phase列为succeed状态的行
    succeed_df = df[df['pod_phase'] == 'Succeeded']

    # 统计gpu_milli列的值及其数量
    gpu_milli_counts = succeed_df['gpu_milli'].value_counts().sort_index()

    # 输出结果
    print("------------------succeed gpu_milli counts------------------")
    print(gpu_milli_counts)

# 使用示例
filter_and_count('openb_pod_list_gpuspec33.csv')