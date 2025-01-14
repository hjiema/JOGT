import pandas as pd

def process_csv(input_path, output_path):
    # 读取CSV文件
    df = pd.read_csv(input_path)

    # 计算 h2 列和 h3 列
    df['h2'] = df['GTX2080Ti'] / df['A100']
    df['h3'] = df['V100'] / df['A100']

    # 保存处理后的CSV文件
    df.to_csv(output_path, index=False)

    print("CSV文件处理完成！")

# 调用函数
process_csv('../norm_openb_pod_list_gpushare20.csv', 'processed_openb_pod_list_gpushare20.csv')
