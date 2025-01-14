import pandas as pd

def merge_csv_files(output_file, *input_files):
    # 读取所有CSV文件并存储在一个列表中，忽略原有的索引列
    dataframes = [pd.read_csv(file, index_col=0) for file in input_files]

    # 合并所有数据框
    merged_df = pd.concat(dataframes, ignore_index=True)

    # 删除重复的任务
    merged_df.drop_duplicates(subset=['job_name'], keep='first', inplace=True)

    # 重新设置索引
    merged_df.reset_index(drop=True, inplace=True)

    merged_df.reset_index(drop=True, inplace=True)
    merged_df.index.name = ''
    merged_df.reset_index(inplace=True)

    # 将合并后的数据框写入新的CSV文件
    merged_df.to_csv(output_file, index=False)


input_files = ['norm_openb_pod_list_gpushare20.csv', 'norm_openb_pod_list_gpushare40.csv', 'norm_openb_pod_list_gpushare60.csv', 'norm_openb_pod_list_gpushare80.csv', 'norm_openb_pod_list_gpushare100.csv']
merge_csv_files('all_cases1.csv', *input_files)