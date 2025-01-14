import pandas as pd
import numpy as np

# 读取任务文件和流行度文件
tasks = pd.read_csv('input.csv')
popularity = pd.read_csv('popularity.csv')

# GPU资源配置
gpu_resources = {
    'A100': {'count': 10, 'memory': 1000, 'ids': list(range(0, 10))},
    'V100': {'count': 20, 'memory': 1000, 'ids': list(range(10, 30))},
    'GTX2080Ti': {'count': 15, 'memory': 1000, 'ids': list(range(30, 45))}
}

# 任务排序：按照任务在A100上的完成时间从高到低排序任务
tasks = tasks.sort_values(by='A100', ascending=False)

# 初始化GPU状态
gpu_status = {gpu_type: {'tasks': [], 'frag_runtime': 0} for gpu_type in gpu_resources.keys()}

# 计算frag*runtime值的函数
def calculate_frag_runtime(gpu_type, task, gpu_status, popularity):
    gpu_memory = gpu_resources[gpu_type]['memory']
    frag_runtime = (gpu_memory - task['gpu_milli']) * \
                   popularity[popularity['gpu_milli'] > (gpu_memory - task['gpu_milli'])]['ratio'].sum() * \
                   task[gpu_type]
    return frag_runtime

# 调度任务
for _, task in tasks.iterrows():
    best_gpu = None
    best_frag_runtime = float('inf')
    
    for gpu_type in gpu_resources.keys():
        frag_runtime = calculate_frag_runtime(gpu_type, task, gpu_status, popularity)
        
        if frag_runtime < best_frag_runtime:
            best_frag_runtime = frag_runtime
            best_gpu = gpu_type
    
    # 更新GPU状态
    gpu_status[best_gpu]['tasks'].append(task['job_name'])
    gpu_status[best_gpu]['frag_runtime'] += best_frag_runtime

# 输出每个GPU上的任务二叉树及调度结果
for gpu_type, status in gpu_status.items():
    print(f"GPU类型: {gpu_type}")
    print(f"任务列表: {status['tasks']}")
    print(f"总的frag*runtime值: {status['frag_runtime']}\n")
