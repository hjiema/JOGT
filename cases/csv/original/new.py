import pandas as pd
import numpy as np
import random

def base_process(file_path):
    # 读取CSV文件
    df = pd.read_csv(file_path)

    # 只保留num_gpu为1的行
    df = df[df['num_gpu'] == 1]

    # 删除执行失败的任务
    df = df[df['pod_phase'] == 'Succeeded']

    # 修改列名
    df.rename(columns={'name': 'job_name'}, inplace=True)
    # 删除重复的任务
    df.drop_duplicates(subset=['job_name'], keep='first', inplace=True)

    # 删除不需要的列
    df.drop(columns=['cpu_milli', 'memory_mib', 'num_gpu', 'gpu_milli', 'gpu_spec', 'qos'], inplace=True)

    # 计算jct列
    df['runtime'] = df['deletion_time'] - df['scheduled_time']
    # df['jct'] = df['jct'].dt.total_seconds()

    # 按creation_time列排序
    df.sort_values(by='creation_time', inplace=True)

    # 获取最小的creation_time值
    min_creation_time = df['creation_time'].min()

    # 更新creation_time列的值
    df['creation_time'] = df['creation_time'] - min_creation_time

    # 删除deletion_time和scheduled_time列
    df.drop(columns=['deletion_time', 'scheduled_time'], inplace=True)

    # 修改列名
    df.rename(columns={'creation_time': 'norm_job_submit_time'}, inplace=True)

    # 删除pod_phase列
    df.drop(columns=['pod_phase'], inplace=True)

    df['gpu_type'] = 'A100'

    # 添加索引列
    df.reset_index(drop=True, inplace=True)
    df.index.name = ''
    df.reset_index(inplace=True)

    # 保存处理后的CSV文件
    df.to_csv('norm_openb_pod_list_gpushare20.csv', index=False)

    print("CSV文件处理完成！")
    return df

def gen_ddl_and_gpu_runtimes(df_one_inst, output_path):
    # 手动设置gpu_types
    gpu_types = ['A100', 'GTX2080Ti', 'V100']

    # 不同GPU的运行时间比例，将任务的标准运行时间转换为特定GPU上的运行时间。
    runtimes = {
        'A100': (1, 1),
        'GTX2080Ti': (1.4, 2.),
        'V100': (2.4, 2.66),
    }

    ddl_ratio = 10 # 截止时间比例
    ddl_range = (1.2, 3.0) # 截止时间范围
    submit_together = True # 任务是否可以一起提交
    
    # 用于生成一个任务从from_gpu类型的GPU转到to_gpu类型的GPU时的运行时间。
    def gen_runtime(from_gpu, to_gpu, origin_runtime):
        if from_gpu == to_gpu:
            return origin_runtime
        if from_gpu not in gpu_types:
            print("not in gpu_types:", from_gpu)
        to_rand = random.uniform(*runtimes[to_gpu])
        from_rand = random.uniform(*runtimes[from_gpu])
        return int(origin_runtime * to_rand / from_rand)
    
    # 用于生成一个基于任务提交时间、运行时间和随机因子的截止时间。
    def gen_ddl(norm_submit_time, runtime):
        if random.randint(0, 100) < ddl_ratio:
            return int(norm_submit_time + runtime * random.uniform(*ddl_range))
        return np.inf
    
    # 确保 df_one_inst 包含 gpu_type 和 runtime 列
    if 'gpu_type' not in df_one_inst.columns or 'runtime' not in df_one_inst.columns:
        raise ValueError("DataFrame must contain 'gpu_type' and 'runtime' columns")

    for gpu_type in gpu_types:
        df_one_inst.loc[:, gpu_type] = df_one_inst.apply(lambda x: gen_runtime(x.gpu_type, gpu_type, x.runtime), axis=1)
    
    df_one_inst.loc[:, 'ddl'] = df_one_inst.apply(lambda x: gen_ddl(x.norm_job_submit_time, x.runtime), axis=1)
    
    if submit_together:
        df_one_inst = df_one_inst.iloc[np.random.permutation(len(df_one_inst))]
        df_one_inst.loc[:, 'ddl'] = df_one_inst.apply(lambda x: x['ddl'] - x['norm_job_submit_time'], axis=1)
        df_one_inst.loc[:, 'norm_job_submit_time'] = df_one_inst.apply(lambda x: 0, axis=1)
    
    df_output = df_one_inst[['job_name', 'norm_job_submit_time', 'ddl', 'A100', 'GTX2080Ti', 'V100']]
    df_output.to_csv(output_path, index=False)

def main():
    dfj = base_process('openb_pod_list_gpushare20.csv')
    gen_ddl_and_gpu_runtimes(dfj, '../norm_openb_pod_list_gpushare20.csv')

# 调用 main 函数
if __name__ == "__main__":
    main()
