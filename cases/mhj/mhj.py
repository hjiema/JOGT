import pandas as pd
import numpy as np

# 读取任务文件和流行度文件
tasks_df = pd.read_csv('input.csv')
popularity_df = pd.read_csv('popularity.csv')

# 初始化GPU资源，包括GPU数量、ID、内存和任务。
gpu_resources = {
    'A100': {'count': 10, 'ids': list(range(0, 10)), 'memory': 1000, 'tasks': {i: None for i in range(10)}},
    'V100': {'count': 20, 'ids': list(range(10, 30)), 'memory': 1000, 'tasks': {i: None for i in range(10, 30)}},
    'GTX2080Ti': {'count': 15, 'ids': list(range(30, 45)), 'memory': 1000, 'tasks': {i: None for i in range(30, 45)}}
}

# 任务排序：按照任务在A100上的完成时间从高到低排序任务
tasks_df = tasks_df.sort_values(by='A100', ascending=False)

# 将任务流行度转换为字典，键值为gpu_milli，值为ratio
popularity_dict = {row['gpu_milli']: row['ratio'] for _, row in popularity_df.iterrows()}

# 定义二叉树节点类。定义TsakNode类，表示任务节点。
class TaskNode:
    def __init__(self, task_id, gpu_milli, runtime):
        self.task_id = task_id
        self.gpu_milli = gpu_milli
        self.runtime = runtime
        self.left = None
        self.right = None

# 计算放置到根节点上后的frag*runtime值
def calculate_frag_runtime(gpu_memory, gpu_milli, runtime, popularity_dict):
    frag = gpu_memory - gpu_milli
    frag_runtime = frag * sum(ratio for mem, ratio in popularity_dict.items() if mem > frag) * runtime
    return frag_runtime

# 计算任务放置在节点P的右子树上的frag*runtime值
def calculate_right_subtree_frag_runtime(node, gpu_milli, runtime, gpu_memory, popularity_dict):
    total_gpu_milli = node.gpu_milli
    current = node.right
    while current:
        total_gpu_milli += current.gpu_milli
        current = current.right
    frag = gpu_memory - total_gpu_milli
    frag_runtime = frag * sum(ratio for mem, ratio in popularity_dict.items() if mem > frag) * runtime
    return frag_runtime

# 放置任务到GPU上
def place_task(gpu_id, gpu_type, node, task_id, gpu_milli, runtime, gpu_resources, popularity_dict):
    gpu_memory = gpu_resources[gpu_type]['memory']
    if node is None:
        node = TaskNode(task_id, gpu_milli, runtime)
        frag_runtime = calculate_frag_runtime(gpu_memory, gpu_milli, runtime, popularity_dict)
    elif node.right is None:
        node.right = TaskNode(task_id, gpu_milli, runtime)
        frag_runtime = calculate_right_subtree_frag_runtime(node, gpu_milli, runtime, gpu_memory, popularity_dict)
    else:
        if can_place_on_right(node, gpu_milli, runtime):
            frag_runtime = calculate_right_subtree_frag_runtime(node, gpu_milli, runtime, gpu_memory, popularity_dict)
            node.right = TaskNode(task_id, gpu_milli, runtime)
        else:
            # 将任务放置在最左叶子节点的左孩子上
            left_most_node = node
            while left_most_node.left:
                left_most_node = left_most_node.left
            left_most_node.left = TaskNode(task_id, gpu_milli, runtime)
            frag_runtime = calculate_frag_runtime(gpu_memory, gpu_milli, runtime, popularity_dict)
    gpu_resources[gpu_type]['tasks'][gpu_id] = node
    return frag_runtime

# 判断任务能否放在节点P的右子树上
def can_place_on_right(node, gpu_milli, runtime):
    total_gpu_milli = node.gpu_milli
    current = node.right
    while current:
        total_gpu_milli += current.gpu_milli
        current = current.right
    return gpu_milli < total_gpu_milli and runtime < node.runtime

# 调度任务
schedule_result = []
for _, task in tasks_df.iterrows():
    task_id = task['job_name']
    gpu_milli = task['gpu_milli']
    best_gpu = None
    best_frag_runtime = float('inf')
    best_node = None

    for gpu_type, resources in gpu_resources.items():
        for gpu_id in resources['ids']:
            node = resources['tasks'][gpu_id]
            frag_runtime = place_task(gpu_id, gpu_type, node, task_id, gpu_milli, task[gpu_type], gpu_resources, popularity_dict)
            if frag_runtime < best_frag_runtime:
                best_frag_runtime = frag_runtime
                best_gpu = (gpu_id, gpu_type)
                best_node = node
                best_runtime = task[gpu_type]

    if best_gpu:
        gpu_id, gpu_type = best_gpu
        place_task(gpu_id, gpu_type, best_node, task_id, gpu_milli, best_runtime, gpu_resources, popularity_dict)
        schedule_result.append({
            'task_id': task_id,
            'gpu_id': gpu_id,
            'gpu_type': gpu_type,
            'start_time': 0,  # 根据具体需求调整
            'end_time': best_runtime,  # 根据具体需求调整
            'frag_runtime': best_frag_runtime
        })

# 输出每个GPU上的任务二叉树及调度结果
for gpu_type, resources in gpu_resources.items():
    for gpu_id, root in resources['tasks'].items():
        print(f"GPU {gpu_id} ({gpu_type}):")
        def print_tree(node, depth=0):
            if node:
                print('  ' * depth + f"Task {node.task_id} - Memory: {node.gpu_milli}, Runtime: {node.runtime}")
                print_tree(node.left, depth + 1)
                print_tree(node.right, depth + 1)
        print_tree(root)

# 输出调度结果
schedule_df = pd.DataFrame(schedule_result)
print(schedule_df)
