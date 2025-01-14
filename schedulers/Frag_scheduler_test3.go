package schedulers

//
//import (
//	"DES-go/schedulers/types"
//	"math"
//	"time"
//)
//
//type FragScheduler struct {
//	gpuCluster      types.Cluster
//	allWaitingJobs  []types.JobMeta
//	SchedulerRecord *types.SchedulerRecord
//	gpuTrees        map[types.GPUID]*TreeNode // 新增，用于存储每个GPU的任务调度树
//}
//
//type TreeNode struct {
//	Job         types.JobMeta
//	GPUType     types.GPUType
//	GPUId       types.GPUID
//	StartTime   float64
//	EndTime     float64
//	GPUMilli    int
//	Frag        float64
//	Runtime     float64
//	FGT         float64
//	RemainMilli int
//	Left        *TreeNode
//	Right       *TreeNode
//}
//
//func NewFragScheduler() *FragScheduler {
//	return &FragScheduler{
//		SchedulerRecord: &types.SchedulerRecord{
//			DoScheduleRecords: []*types.DoScheduleCallRecord{},
//		},
//		gpuTrees: make(map[types.GPUID]*TreeNode), // 初始化GPU任务调度树
//	}
//}
//
//func (d *FragScheduler) DoSchedule() {
//	start := time.Now()
//	d.doSchedule()
//	duration := time.Since(start)
//	d.SchedulerRecord.DoScheduleRecords = append(d.SchedulerRecord.DoScheduleRecords, &types.DoScheduleCallRecord{Duration: duration})
//}
//
//func (d *FragScheduler) SetCluster(cluster types.Cluster) {
//	d.gpuCluster = cluster
//	for _, gpuId := range cluster.GPUIDs() {
//		d.gpuTrees[gpuId] = nil // 初始化每个GPU的任务树
//	}
//}
//
//func calculateFragRuntime(remainMilli int) float64 {
//	// 简化的Fragment计算函数
//	aboveMilliPopularity := 1.0
//	return float64(remainMilli) * aboveMilliPopularity
//}
//
//func sumFGTOnGPU(node *TreeNode) float64 {
//	if node == nil {
//		return 0
//	}
//	return node.FGT + sumFGTOnGPU(node.Left) + sumFGTOnGPU(node.Right)
//}
//
//// 获取某个 GPU 类型对应的执行时间
//func getDurationForGPU(job types.JobMeta, gpuType types.GPUType) types.Duration {
//	durations := job.Durations()
//	return durations[gpuType]
//}
//
//func (d *FragScheduler) doSchedule() {
//	// 初始化任务等待队列
//	var jobs []types.JobMeta
//	jobs = d.allWaitingJobs
//
//	// 存储每个GPU的FGT值
//	jobFGTList := make(map[types.GPUID]float64)
//	for _, job := range jobs {
//		// 调度每个作业，计算FGT值
//		jobFGTList = scheduleJob(job, d.gpuCluster.GPUIDs(), d.gpuCluster, d.gpuTrees, jobFGTList)
//		// 找到FGT最小的GPU
//		bestGPUID := findMinFGTGPU(jobFGTList)
//		if bestGPUID != -1 {
//			// 将作业调度到最优GPU
//			d.gpuTrees = scheduleJobForGPU(job, bestGPUID, d.gpuCluster, d.gpuTrees)
//		}
//	}
//}
//
//// 调度作业到每个GPU并计算FGT值
//func scheduleJob(job types.JobMeta, gpus []types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode, jobFGTList map[types.GPUID]float64) map[types.GPUID]float64 {
//	if len(jobFGTList) == 0 {
//		for _, gpu := range gpus {
//			jobFGTList[gpu] = 0
//		}
//	}
//
//	for _, gpuId := range gpus {
//		currentGPUType := cluster.GPU(gpuId).Type()
//		jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//		root := gpuTrees[gpuId]
//
//		if root == nil {
//			remainMilliOnGPU := 1000 - job.GPUMilli()
//			fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			jobFGTList[gpuId] = fgtOnGPU
//		} else {
//			if job.GPUMilli() <= root.RemainMilli {
//				root.RemainMilli -= job.GPUMilli()
//				current := root
//				for current.Right != nil {
//					current = current.Right
//				}
//				fragOnGPU := calculateFragRuntime(root.RemainMilli)
//				fgtOnGPU := fragOnGPU * jobRuntime
//				current.FGT = current.Frag * (current.Runtime - jobRuntime)
//				jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//			} else {
//				current := root.Left
//				for current != nil {
//					if job.GPUMilli() > current.RemainMilli {
//						if current.Left == nil {
//							remainMilliOnGPU := 1000 - job.GPUMilli()
//							fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//							fgtOnGPU := fragOnGPU * jobRuntime
//							jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//							break
//						} else {
//							current = current.Left
//						}
//					} else {
//						current.RemainMilli -= job.GPUMilli()
//						for current.Right != nil {
//							current = current.Right
//						}
//						fragOnGPU := calculateFragRuntime(current.RemainMilli)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.FGT = current.Frag * (current.Runtime - jobRuntime)
//						jobFGTList[gpuId] = fgtOnGPU
//						break
//					}
//				}
//				if current == nil {
//					remainMilliOnGPU := 1000 - job.GPUMilli()
//					fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//				}
//			}
//		}
//	}
//
//	return jobFGTList
//}
//
//// 将作业调度到最佳GPU
//func scheduleJobForGPU(job types.JobMeta, gpuID types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode) map[types.GPUID]*TreeNode {
//	currentGPUType := cluster.GPU(gpuID).Type()
//	jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//	root := gpuTrees[gpuID]
//
//	if root == nil {
//		remainMilliOnGPU := 1000 - job.GPUMilli()
//		fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//		fgtOnGPU := fragOnGPU * jobRuntime
//		root = &TreeNode{
//			Job:         job,
//			GPUType:     currentGPUType,
//			GPUId:       gpuID,
//			StartTime:   0,
//			EndTime:     jobRuntime,
//			GPUMilli:    job.GPUMilli(),
//			Frag:        fragOnGPU,
//			Runtime:     jobRuntime,
//			FGT:         fgtOnGPU,
//			RemainMilli: remainMilliOnGPU,
//		}
//		gpuTrees[gpuID] = root
//	} else {
//		if job.GPUMilli() <= root.RemainMilli {
//			root.RemainMilli -= job.GPUMilli()
//			current := root
//			for current.Right != nil {
//				current = current.Right
//			}
//			fragOnGPU := calculateFragRuntime(root.RemainMilli)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			current.Right = &TreeNode{
//				Job:         job,
//				GPUType:     currentGPUType,
//				GPUId:       gpuID,
//				StartTime:   current.EndTime,
//				EndTime:     current.EndTime + jobRuntime,
//				GPUMilli:    job.GPUMilli(),
//				Frag:        fragOnGPU,
//				Runtime:     jobRuntime,
//				FGT:         fgtOnGPU,
//				RemainMilli: root.RemainMilli,
//			}
//		} else {
//			// 处理剩余资源不足的情况
//			current := root.Left
//			for current != nil {
//				if job.GPUMilli() > current.RemainMilli {
//					if current.Left == nil {
//						remainMilliOnGPU := 1000 - job.GPUMilli()
//						fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.Left = &TreeNode{
//							Job:         job,
//							GPUType:     currentGPUType,
//							GPUId:       gpuID,
//							StartTime:   current.EndTime,
//							EndTime:     current.EndTime + jobRuntime,
//							GPUMilli:    job.GPUMilli(),
//							Frag:        fragOnGPU,
//							Runtime:     jobRuntime,
//							FGT:         fgtOnGPU,
//							RemainMilli: remainMilliOnGPU,
//						}
//						break
//					} else {
//						current = current.Left
//					}
//				} else {
//					current.RemainMilli -= job.GPUMilli()
//					for current.Right != nil {
//						current = current.Right
//					}
//					fragOnGPU := calculateFragRuntime(current.RemainMilli)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					current.Right = &TreeNode{
//						Job:         job,
//						GPUType:     currentGPUType,
//						GPUId:       gpuID,
//						StartTime:   current.EndTime,
//						EndTime:     current.EndTime + jobRuntime,
//						GPUMilli:    job.GPUMilli(),
//						Frag:        fragOnGPU,
//						Runtime:     jobRuntime,
//						FGT:         fgtOnGPU,
//						RemainMilli: current.RemainMilli,
//					}
//					break
//				}
//			}
//		}
//	}
//
//	gpuTrees[gpuID] = root
//	return gpuTrees
//}
//
//// 找到FGT最小的GPU
//func findMinFGTGPU(jobFGTList map[types.GPUID]float64) types.GPUID {
//	minFGTValue := math.Inf(1)
//	var bestGPUID types.GPUID
//
//	for gpuID, fgtValue := range jobFGTList {
//		if fgtValue < minFGTValue {
//			minFGTValue = fgtValue
//			bestGPUID = gpuID
//		}
//	}
//
//	return bestGPUID
//}
//
//func (d *FragScheduler) OnScheduleEvent(event types.ScheduleEvent) {
//	switch e := event.(type) {
//	case *types.ScheduleEventJobsArrived:
//		{
//			d.allWaitingJobs = e.JobMetas()
//			d.DoSchedule()
//		}
//	case *types.ScheduleEventJobsFinished:
//		{
//			d.DoSchedule()
//		}
//	}
//}
//
//func (d *FragScheduler) Name() string {
//	return "FragScheduler"
//}
//
//func (d *FragScheduler) Info() interface{} {
//	return d.Name()
//}
//
//func (d *FragScheduler) NextActiveScheduleTime() types.Time {
//	return types.Time(math.Inf(1))
//}
//
//func (d *FragScheduler) Record() *types.SchedulerRecord {
//	return d.SchedulerRecord
//}
//=========================================================================================================================

//package schedulers
//
//import (
//"DES-go/schedulers/types"
//"math"
//)
//
//type FragScheduler struct {
// cluster types.Cluster // 集群
// // nextScheduleToGPUID int             // 下一个需要调度的GPUID
// // lastScheduleTime    types.Time      // 上次调度的时间
// maxGPUJobQueueID    int             // GPU作业队列的最大ID
// unscheduledJobMetas []types.JobMeta // 未调度的作业元数据的切片
//
// unscheduledJobsCacheLength int // 未调度作业缓存的长度
// maxScheduleInterval        int // 最大调度间隔时间
//}
//
//// 定义一个结构体来存储每个GPU的FGT信息
//type GPUFgtInfo struct {
// GPUID int
// FGT   float64
//}
//
//func NewFragScheduler() *FragScheduler {
// // some casual Frag configs
// unscheduledJobsCacheLength := 1000
// maxScheduleInterval := 1000
// return &FragScheduler{
//	 unscheduledJobMetas:        make([]types.JobMeta, 0, unscheduledJobsCacheLength),
//	 maxScheduleInterval:        maxScheduleInterval,
//	 unscheduledJobsCacheLength: unscheduledJobsCacheLength,
// }
//}
//
//func calculateFragRuntime(remainMilli int) float64 {
// // 简化的Fragment计算函数
// aboveMilliPopularity := 1.0
// return float64(remainMilli) * aboveMilliPopularity
//}
//
//func getDurationForGPU(job types.Job, gpuType types.GPUType) types.Duration {
// durations := job.JobMeta().Durations()
// return durations[gpuType]
//}
//
//// 计算每个GPU的FGT值
//func (d *FragScheduler) CalculateFGTForGPU(gpuID types.GPUID, job types.Job) float64 {
// currentGPUType := d.cluster.GPU(gpuID).Type()
// //gpuQueue := d.cluster.GPUJobQueues()[types.GPUID(gpuID)]
// remainMilli := 1000 - job.JobMeta().GPUMilli()
// jobRuntime := float64(getDurationForGPU(job, currentGPUType))
// fragOnGPU := calculateFragRuntime(remainMilli)
// fgt := fragOnGPU * jobRuntime
// // 假设 job.gpuMilli 为作业的计算需求
// // 根据 remainMilli 和 job.gpuMilli 计算 FGT 值
// //ragOnGPU := remainMilli * 0.1 // 假设碎片率与剩余容量成正比
// // fgt := fragOnGPU * job.Runtime()
// return fgt
//}
//
//// 查找具有最小FGT的GPU
//func (d *FragScheduler) FindBestGPUForJob(job types.Job) int {
// bestGPUID := -1
// minFGT := math.MaxFloat64
//
// // 遍历所有GPU，计算FGT值
// for gpuID := 0; gpuID < d.maxGPUJobQueueID; gpuID++ {
//	 fgt := d.CalculateFGTForGPU(types.GPUID(gpuID), job)
//	 if fgt < minFGT {
//		 minFGT = fgt
//		 bestGPUID = gpuID
//	 }
// }
//
// return bestGPUID
//}
//
//func (d *FragScheduler) DoSchedule() {
// if d.unscheduledJobMetas == nil {
//	 panic("FragScheduler d.unscheduledJobMetas == nil")
// }
//
// jobs := make([]types.Job, 0, len(d.unscheduledJobMetas))
// for _, jobMeta := range d.unscheduledJobMetas {
//	 jobs = append(jobs, d.cluster.InitJob(jobMeta))
// }
//
// var targetJobQueue types.GPUJobQueue
// for _, job := range jobs {
//	 // 调用 FindBestGPUForJob 查找最佳GPU
//	 bestGPUID := d.FindBestGPUForJob(job)
//	 if bestGPUID == -1 {
//		 panic("No valid GPU found for scheduling")
//	 }
//	 targetJobQueue = d.cluster.GPUJobQueues()[types.GPUID(bestGPUID)]
//	 targetJobQueue.SetJobs(job)
// }
//
// // d.lastScheduleTime = d.cluster.Now()
// d.unscheduledJobMetas = d.unscheduledJobMetas[:0]
//}
//
//func (d *FragScheduler) SetCluster(cluster types.Cluster) {
// d.cluster = cluster
// // d.nextScheduleToGPUID = 0
// // d.lastScheduleTime = d.cluster.Now()
//
// d.maxGPUJobQueueID = math.MinInt64
// for _, gpuList := range d.cluster.GPUs() {
//	 for _, gpu := range gpuList {
//		 d.maxGPUJobQueueID = int(math.Max(float64(gpu.ID())+1, float64(d.maxGPUJobQueueID)))
//	 }
// }
//}
//
//func (d *FragScheduler) OnScheduleEvent(event types.ScheduleEvent) {
// switch e := event.(type) {
// //case *types.ScheduleEventDurationPassed:
// //	{
// //		 if d.cluster.Now()-d.lastScheduleTime >= types.Time(d.maxScheduleInterval) {
// //			d.DoSchedule()
// //		}
// //	}
// case *types.ScheduleEventJobsArrived:
//	 {
//		 d.unscheduledJobMetas = append(d.unscheduledJobMetas, e.JobMetas()...)
//		 if len(d.unscheduledJobMetas) > d.unscheduledJobsCacheLength {
//			 d.DoSchedule()
//		 }
//	 }
// }
//}
//
//func (d *FragScheduler) NextActiveScheduleTime() types.Time {
// if len(d.unscheduledJobMetas) > 0 {
//	 return types.Time(d.maxScheduleInterval) - (d.cluster.Now() - d.lastScheduleTime)
// }
// return types.Time(math.Inf(1))
//}
//
//func (d *FragScheduler) Name() string {
// return "FragScheduler"
//}
//
//func (d *FragScheduler) Info() interface{} {
// return d.Name()
//}
//
//func (d *FragScheduler) Record() *types.SchedulerRecord {
// return &types.SchedulerRecord{
//	 DoScheduleRecords: []*types.DoScheduleCallRecord{},
// }
//}

//   重新开始改
//package schedulers
//
//import (
//"DES-go/schedulers/types"
//"fmt"
//"math"
//"strings"
//"time"
//)
//
//type FragScheduler struct {
//	gpuCluster      types.Cluster
//	allWaitingJobs  []types.JobMeta
//	SchedulerRecord *types.SchedulerRecord
//	gpuTrees        map[types.GPUID]*TreeNode // 新增，用于存储每个GPU的任务调度树
//}
//
//type TreeNode struct {
//	Job         types.JobMeta
//	GPUType     types.GPUType
//	GPUId       types.GPUID
//	StartTime   float64
//	EndTime     float64
//	GPUMilli    int
//	Frag        float64
//	Runtime     float64
//	FGT         float64
//	RemainMilli int
//	Left        *TreeNode
//	Right       *TreeNode
//}
//
//func NewFragScheduler() *FragScheduler {
//	return &FragScheduler{
//		SchedulerRecord: &types.SchedulerRecord{
//			DoScheduleRecords: []*types.DoScheduleCallRecord{},
//		},
//		gpuTrees: make(map[types.GPUID]*TreeNode), // 初始化GPU任务调度树
//	}
//}
//
//func (d *FragScheduler) DoSchedule() {
//	start := time.Now()
//	d.doSchedule()
//	duration := time.Since(start)
//	d.SchedulerRecord.DoScheduleRecords = append(d.SchedulerRecord.DoScheduleRecords, &types.DoScheduleCallRecord{Duration: duration})
//}
//
//func (d *FragScheduler) SetCluster(cluster types.Cluster) {
//	d.gpuCluster = cluster
//	for _, gpuId := range cluster.GPUIDs() {
//		d.gpuTrees[gpuId] = nil // 初始化每个GPU的任务树
//	}
//}
//
//func calculateFragRuntime(remainMilli int) float64 {
//	// 简化的Fragment计算函数
//	aboveMilliPopularity := 1.0
//	return float64(remainMilli) * aboveMilliPopularity
//}
//
//func sumFGTOnGPU(node *TreeNode) float64 {
//	if node == nil {
//		return 0
//	}
//	return node.FGT + sumFGTOnGPU(node.Left) + sumFGTOnGPU(node.Right)
//}
//
//// 获取某个 GPU 类型对应的执行时间
//func getDurationForGPU(job types.JobMeta, gpuType types.GPUType) types.Duration {
//	durations := job.Durations()
//	return durations[gpuType]
//}
//
//func (d *FragScheduler) doSchedule() {
//	// 初始化任务等待队列
//	var jobs []types.JobMeta
//	jobs = d.allWaitingJobs
//	//for _, job := range jobs {
//	//	println(job.JobName())
//	//}
//	// 存储每个GPU的FGT值
//	jobFGTList := make(map[types.GPUID]float64)
//
//	for _, job := range jobs {
//		// 调度每个作业，计算FGT值
//		jobFGTList = scheduleJob(job, d.gpuCluster.GPUIDs(), d.gpuCluster, d.gpuTrees, jobFGTList)
//		// 找到FGT最小的GPU
//		bestGPUID := findMinFGTGPU(jobFGTList)
//
//		if bestGPUID != -1 {
//			// 将作业调度到最优GPU
//			d.gpuTrees = scheduleJobForGPU(job, bestGPUID, d.gpuCluster, d.gpuTrees)
//		}
//	}
//
//	//for gpuID, treeNode := range d.gpuTrees {
//	//	fmt.Printf("GPU ID: %d\n", gpuID)
//	//	printTree(treeNode, 0)
//	//}
//}
//func printTree(node *TreeNode, level int) {
//	if node != nil {
//		printTree(node.Right, level+1)
//		//fmt.Printf("%s-> %s\n", getIndent(level), node.Job.JobName())
//		fmt.Printf("%s-> %s\n", strings.Repeat(" ", 4*level), node.Job.JobName())
//
//		printTree(node.Left, level+1)
//	}
//}
//
////func getIndent(level int) string {
////	return fmt.Sprintf("%s", string(' ', level*4))
////}
//
//// 调度作业到每个GPU并计算FGT值
//func scheduleJob(job types.JobMeta, gpus []types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode, jobFGTList map[types.GPUID]float64) map[types.GPUID]float64 {
//	if len(jobFGTList) == 0 {
//		for _, gpu := range gpus {
//			jobFGTList[gpu] = 0
//		}
//	}
//
//	for _, gpuId := range gpus {
//		currentGPUType := cluster.GPU(gpuId).Type()
//		jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//		root := gpuTrees[gpuId]
//
//		if root == nil {
//			remainMilliOnGPU := 1000 - job.GPUMilli()
//			fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			jobFGTList[gpuId] = fgtOnGPU
//		} else {
//			if job.GPUMilli() <= root.RemainMilli {
//				root.RemainMilli -= job.GPUMilli()
//				current := root
//				for current.Right != nil {
//					current = current.Right
//				}
//				fragOnGPU := calculateFragRuntime(root.RemainMilli)
//				fgtOnGPU := fragOnGPU * jobRuntime
//				current.FGT = current.Frag * (current.Runtime - jobRuntime)
//				jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//			} else {
//				current := root.Left
//				for current != nil {
//					if job.GPUMilli() > current.RemainMilli {
//						if current.Left == nil {
//							remainMilliOnGPU := 1000 - job.GPUMilli()
//							fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//							fgtOnGPU := fragOnGPU * jobRuntime
//							jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//							break
//						} else {
//							current = current.Left
//						}
//					} else {
//						current.RemainMilli -= job.GPUMilli()
//						for current.Right != nil {
//							current = current.Right
//						}
//						fragOnGPU := calculateFragRuntime(current.RemainMilli)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.FGT = current.Frag * (current.Runtime - jobRuntime)
//						jobFGTList[gpuId] = fgtOnGPU
//						break
//					}
//				}
//				if current == nil {
//					remainMilliOnGPU := 1000 - job.GPUMilli()
//					fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//				}
//			}
//		}
//	}
//
//	return jobFGTList
//}
//
//// 将作业调度到最佳GPU
//func scheduleJobForGPU(job types.JobMeta, gpuID types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode) map[types.GPUID]*TreeNode {
//	currentGPUType := cluster.GPU(gpuID).Type()
//	jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//	root := gpuTrees[gpuID]
//
//	if root == nil {
//		remainMilliOnGPU := 1000 - job.GPUMilli()
//		fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//		fgtOnGPU := fragOnGPU * jobRuntime
//		root = &TreeNode{
//			Job:         job,
//			GPUType:     currentGPUType,
//			GPUId:       gpuID,
//			StartTime:   0,
//			EndTime:     jobRuntime,
//			GPUMilli:    job.GPUMilli(),
//			Frag:        fragOnGPU,
//			Runtime:     jobRuntime,
//			FGT:         fgtOnGPU,
//			RemainMilli: remainMilliOnGPU,
//		}
//		gpuTrees[gpuID] = root
//	} else {
//		if job.GPUMilli() <= root.RemainMilli {
//			root.RemainMilli -= job.GPUMilli()
//			current := root
//			for current.Right != nil {
//				current = current.Right
//			}
//			fragOnGPU := calculateFragRuntime(root.RemainMilli)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			current.Right = &TreeNode{
//				Job:         job,
//				GPUType:     currentGPUType,
//				GPUId:       gpuID,
//				StartTime:   current.EndTime,
//				EndTime:     current.EndTime + jobRuntime,
//				GPUMilli:    job.GPUMilli(),
//				Frag:        fragOnGPU,
//				Runtime:     jobRuntime,
//				FGT:         fgtOnGPU,
//				RemainMilli: root.RemainMilli,
//			}
//		} else {
//			// 处理剩余资源不足的情况
//			current := root.Left
//			for current != nil {
//				if job.GPUMilli() > current.RemainMilli {
//					if current.Left == nil {
//						remainMilliOnGPU := 1000 - job.GPUMilli()
//						fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.Left = &TreeNode{
//							Job:         job,
//							GPUType:     currentGPUType,
//							GPUId:       gpuID,
//							StartTime:   current.EndTime,
//							EndTime:     current.EndTime + jobRuntime,
//							GPUMilli:    job.GPUMilli(),
//							Frag:        fragOnGPU,
//							Runtime:     jobRuntime,
//							FGT:         fgtOnGPU,
//							RemainMilli: remainMilliOnGPU,
//						}
//						break
//					} else {
//						current = current.Left
//					}
//				} else {
//					current.RemainMilli -= job.GPUMilli()
//					for current.Right != nil {
//						current = current.Right
//					}
//					fragOnGPU := calculateFragRuntime(current.RemainMilli)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					current.Right = &TreeNode{
//						Job:         job,
//						GPUType:     currentGPUType,
//						GPUId:       gpuID,
//						StartTime:   current.EndTime,
//						EndTime:     current.EndTime + jobRuntime,
//						GPUMilli:    job.GPUMilli(),
//						Frag:        fragOnGPU,
//						Runtime:     jobRuntime,
//						FGT:         fgtOnGPU,
//						RemainMilli: current.RemainMilli,
//					}
//					break
//				}
//			}
//		}
//	}
//
//	gpuTrees[gpuID] = root
//	return gpuTrees
//}
//
//// 找到FGT最小的GPU
//func findMinFGTGPU(jobFGTList map[types.GPUID]float64) types.GPUID {
//	minFGTValue := math.Inf(1)
//	var bestGPUID types.GPUID
//
//	for gpuID, fgtValue := range jobFGTList {
//		if fgtValue < minFGTValue {
//			minFGTValue = fgtValue
//			bestGPUID = gpuID
//		}
//	}
//
//	return bestGPUID
//}
//
//func (d *FragScheduler) OnScheduleEvent(event types.ScheduleEvent) {
//	switch e := event.(type) {
//	case *types.ScheduleEventJobsArrived:
//		{
//			d.allWaitingJobs = e.JobMetas()
//			d.DoSchedule()
//		}
//	case *types.ScheduleEventJobsFinished:
//		{
//			d.DoSchedule()
//		}
//	}
//}
//
//func (d *FragScheduler) Name() string {
//	return "FragScheduler"
//}
//
//func (d *FragScheduler) Info() interface{} {
//	return d.Name()
//}
//
//func (d *FragScheduler) NextActiveScheduleTime() types.Time {
//	return types.Time(math.Inf(1))
//}
//
//func (d *FragScheduler) Record() *types.SchedulerRecord {
//	return d.SchedulerRecord
//}

//
//
//========================================================================
// 下面是正常版本，但是不能输出finishedJob，其余正常的代码
//package schedulers
//
//import (
//"DES-go/schedulers/types"
//"fmt"
//"math"
//"sort"
//"strings"
//"time"
//)
//
//type FragScheduler struct {
//	cluster         types.Cluster
//	allArrivedjobs  []types.JobMeta
//	allWaitingJobs  []types.Job
//	SchedulerRecord []*types.DoScheduleCallRecord
//	FinishedJob     []types.Job
//	gpuTrees        map[types.GPUID]*TreeNode // 新增，用于存储每个GPU的任务调度树
//}
//
//type TreeNode struct {
//	Job         types.Job
//	GPUType     types.GPUType
//	GPUId       types.GPUID
//	StartTime   types.Time
//	EndTime     types.Time
//	GPUMilli    int
//	Frag        float64
//	Runtime     types.Time
//	FGT         float64
//	RemainMilli int
//	Left        *TreeNode
//	Right       *TreeNode
//}
//
//func NewFragScheduler() *FragScheduler {
//
//	return &FragScheduler{
//		SchedulerRecord: make([]*types.DoScheduleCallRecord, 0, 128),
//		//SchedulerRecord: &types.SchedulerRecord{
//		//	DoScheduleRecords: []*types.DoScheduleCallRecord{},
//		//},
//		gpuTrees: make(map[types.GPUID]*TreeNode), // 初始化GPU任务调度树
//
//	}
//}
//
//func (d *FragScheduler) DoSchedule() {
//	start := time.Now()
//	d.doSchedule()
//	duration := time.Since(start)
//	d.SchedulerRecord = append(d.SchedulerRecord, &types.DoScheduleCallRecord{Duration: duration})
//	// d.SchedulerRecord.Extra = d.FinishedJob
//}
//
//func (d *FragScheduler) SetCluster(cluster types.Cluster) {
//	d.cluster = cluster
//	//fmt.Printf("%+v\n", d.gpuCluster.GPUIDs())
//	for _, gpuId := range d.cluster.GPUIDs() {
//		//println("11111111111111111")
//		//println(gpuId)
//		d.gpuTrees[gpuId] = nil // 初始化每个GPU的任务树
//	}
//}
//
//func calculateFragRuntime(remainMilli int) float64 {
//	// 简化的Fragment计算函数
//	aboveMilliPopularity := 1.0
//	return float64(remainMilli) * aboveMilliPopularity
//}
//
//func sumFGTOnGPU(node *TreeNode) float64 {
//	if node == nil {
//		return 0
//	}
//	return node.FGT + sumFGTOnGPU(node.Left) + sumFGTOnGPU(node.Right)
//}
//
//// 获取某个 GPU 类型对应的执行时间
//func getDurationForGPU(job types.Job, gpuType types.GPUType) types.Duration {
//	durations := job.JobMeta().Durations()
//	return durations[gpuType]
//}
//
//const A100 types.GPUType = "A100"
//
//func sortJobsByA100Duration(jobs []types.Job) {
//	sort.Slice(jobs, func(i, j int) bool {
//		//return jobs[i].Durations()[A100] > jobs[j].Durations()[A100]
//		return jobs[i].JobMeta().Durations()[A100] > jobs[j].JobMeta().Durations()[A100]
//	})
//}
//
//func (d *FragScheduler) doSchedule() {
//	// 初始化任务等待队列
//	var jobs []types.Job
//	jobs = d.allWaitingJobs
//
//	// 存储每个GPU的FGT值
//	jobFGTList := make(map[types.GPUID]float64)
//	sortJobsByA100Duration(jobs)
//
//	for _, job := range jobs {
//		// 调度每个作业，计算FGT值
//		jobFGTList = scheduleJob(job, d.cluster.GPUIDs(), d.cluster, d.gpuTrees, jobFGTList)
//		// 找到FGT最小的GPU
//		bestGPUID := findMinFGTGPU(jobFGTList)
//		if bestGPUID != -1 {
//			// 将作业调度到最优GPU
//			d.gpuTrees = scheduleJobForGPU(job, bestGPUID, d.cluster, d.gpuTrees)
//		}
//	}
//
//	var finishedallJob []Jobfish
//
//	for _, gpuId := range d.cluster.GPUIDs() {
//		// 获取对应的根 TreeNode
//		rootTreeNode := d.gpuTrees[gpuId]
//
//		// 如果根节点不为 nil，开始递归遍历树
//		if rootTreeNode != nil {
//			traverseTree(rootTreeNode, finishedallJob, d.cluster)
//		}
//	}
//
//	for _, job := range finishedallJob {
//		fmt.Printf("Job Name: %s\n", job.JobName)
//		fmt.Printf("GPU ID: %s\n", job.GpuID)
//		fmt.Printf("GPU Type: %s\n", job.GpuType)
//		fmt.Printf("GPU: %s\n", job.Gpu)
//		fmt.Printf("First Execution Time: %s\n", job.FirstExecutionTime)
//		fmt.Printf("Finish Execution Time: %s\n", job.FinishExecutionTime)
//		fmt.Printf("Remaining Ratio: %.2f\n", job.RemainingRatio)
//		fmt.Printf("Is Running: %t\n", job.IsRunning)
//		fmt.Println() // 空行区分每个 Jobfish
//	}
//
//}
//
//func traverseTree(treeNode *TreeNode, finishedallJob []Jobfish, cluster types.Cluster) {
//	// 如果 treeNode 为 nil，结束递归
//	if treeNode == nil {
//		return
//	}
//
//	// 创建 Jobfish 对象
//	jobchange := Jobfish{
//		JobName:             treeNode.Job.JobName(),
//		GpuID:               treeNode.GPUId,
//		GpuType:             treeNode.GPUType,
//		Gpu:                 cluster.GPU(treeNode.GPUId), // 假设 d.cluster.GPU 返回 GPU 对象
//		FirstExecutionTime:  treeNode.StartTime,
//		FinishExecutionTime: treeNode.EndTime,
//		RemainingRatio:      0, // 这里根据需要设置剩余比例
//		IsRunning:           false,
//	}
//
//	finishedallJob = append(finishedallJob, jobchange)
//
//	//// 将 jobchange 添加到 finishedallJob 列表中
//	//finishedallJob[gpuID] = append(finishedallJob[gpuID], jobchange)
//
//	// 递归遍历左子树
//	traverseTree(treeNode.Left, finishedallJob, cluster)
//
//	// 递归遍历右子树
//	traverseTree(treeNode.Right, finishedallJob, cluster)
//}
//
//type Jobfish struct {
//	JobName types.JobName
//	//executionDetail     *JobExecutionDetail
//	GpuID               types.GPUID
//	GpuType             types.GPUType
//	Gpu                 types.GPU
//	FirstExecutionTime  types.Time
//	FinishExecutionTime types.Time
//	RemainingRatio      float64 // 任务未执行完的部分的剩余比例。
//	IsRunning           bool
//}
//
//// 模拟调度作业到每个GPU并计算FGT值
//func scheduleJob(job types.Job, gpus []types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode, jobFGTList map[types.GPUID]float64) map[types.GPUID]float64 {
//	if len(jobFGTList) == 0 {
//		for _, gpu := range gpus {
//			jobFGTList[gpu] = 0
//		}
//	}
//
//	// 模拟调度：深拷贝gpuTrees，以免影响原始树
//	simulatedTrees := make(map[types.GPUID]*TreeNode)
//	for gpuID, root := range gpuTrees {
//		simulatedTrees[gpuID] = copyTree(root) // 深拷贝树
//	}
//
//	for _, gpuId := range gpus {
//		currentGPUType := cluster.GPU(gpuId).Type()
//		jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//		root := simulatedTrees[gpuId]
//
//		if root == nil {
//			remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
//			fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			jobFGTList[gpuId] = fgtOnGPU
//		} else {
//			if job.JobMeta().GPUMilli() <= root.RemainMilli {
//				root.RemainMilli -= job.JobMeta().GPUMilli()
//				current := root
//				for current.Right != nil {
//					current = current.Right
//				}
//				fragOnGPU := calculateFragRuntime(root.RemainMilli)
//				fgtOnGPU := fragOnGPU * jobRuntime
//				current.FGT = current.Frag * (float64(current.Runtime) - float64(jobRuntime))
//				jobFGTList[gpuId] = sumFGTOnGPU(root) - fgtOnGPU
//			} else {
//				current := root.Left
//				for current != nil {
//					if job.JobMeta().GPUMilli() > current.RemainMilli {
//						if current.Left == nil {
//							remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
//							fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//							fgtOnGPU := fragOnGPU * jobRuntime
//							jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//							break
//						} else {
//							current = current.Left
//						}
//					} else {
//						current.RemainMilli -= job.JobMeta().GPUMilli()
//						for current.Right != nil {
//							current = current.Right
//						}
//						fragOnGPU := calculateFragRuntime(current.RemainMilli)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.FGT = current.Frag * (float64(current.Runtime) - float64(jobRuntime))
//						jobFGTList[gpuId] = fgtOnGPU
//						break
//					}
//				}
//				if current == nil {
//					remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
//					fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					jobFGTList[gpuId] = sumFGTOnGPU(root) + fgtOnGPU
//				}
//			}
//		}
//	}
//
//	return jobFGTList
//}
//
//// copyTree 实现树的深拷贝
//func copyTree(root *TreeNode) *TreeNode {
//	if root == nil {
//		return nil
//	}
//	newNode := &TreeNode{
//		Job:         root.Job,
//		GPUType:     root.GPUType,
//		GPUId:       root.GPUId,
//		StartTime:   root.StartTime,
//		EndTime:     root.EndTime,
//		GPUMilli:    root.GPUMilli,
//		Frag:        root.Frag,
//		Runtime:     root.Runtime,
//		FGT:         root.FGT,
//		RemainMilli: root.RemainMilli,
//	}
//	newNode.Left = copyTree(root.Left)
//	newNode.Right = copyTree(root.Right)
//	return newNode
//}
//
//// 实际调度
//func scheduleJobForGPU(job types.Job, gpuID types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode) map[types.GPUID]*TreeNode {
//	currentGPUType := cluster.GPU(gpuID).Type()
//	jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//	root := gpuTrees[gpuID]
//
//	if root == nil {
//		remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
//		fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//		fgtOnGPU := fragOnGPU * jobRuntime
//		root = &TreeNode{
//			Job:         job,
//			GPUType:     currentGPUType,
//			GPUId:       gpuID,
//			StartTime:   types.Time(0),
//			EndTime:     types.Time(jobRuntime),
//			GPUMilli:    job.JobMeta().GPUMilli(),
//			Frag:        fragOnGPU,
//			Runtime:     types.Time(jobRuntime),
//			FGT:         fgtOnGPU,
//			RemainMilli: remainMilliOnGPU,
//		}
//		gpuTrees[gpuID] = root
//	} else {
//		if job.JobMeta().GPUMilli() <= root.RemainMilli {
//			root.RemainMilli -= job.JobMeta().GPUMilli()
//			current := root
//			for current.Right != nil {
//				current = current.Right
//			}
//			fragOnGPU := calculateFragRuntime(root.RemainMilli)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			current.Right = &TreeNode{
//				Job:         job,
//				GPUType:     currentGPUType,
//				GPUId:       gpuID,
//				StartTime:   types.Time(current.EndTime),
//				EndTime:     current.EndTime + types.Time(jobRuntime),
//				GPUMilli:    job.JobMeta().GPUMilli(),
//				Frag:        fragOnGPU,
//				Runtime:     types.Time(jobRuntime),
//				FGT:         fgtOnGPU,
//				RemainMilli: root.RemainMilli,
//			}
//		} else {
//			current := root.Left
//			for current != nil {
//				if job.JobMeta().GPUMilli() > current.RemainMilli {
//					if current.Left == nil {
//						remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
//						fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.Left = &TreeNode{
//							Job:         job,
//							GPUType:     currentGPUType,
//							GPUId:       gpuID,
//							StartTime:   types.Time(current.EndTime),
//							EndTime:     current.EndTime + types.Time(jobRuntime),
//							GPUMilli:    job.JobMeta().GPUMilli(),
//							Frag:        fragOnGPU,
//							Runtime:     types.Time(jobRuntime),
//							FGT:         fgtOnGPU,
//							RemainMilli: remainMilliOnGPU,
//						}
//						break
//					} else {
//						current = current.Left
//					}
//				} else {
//					current.RemainMilli -= job.JobMeta().GPUMilli()
//					for current.Right != nil {
//						current = current.Right
//					}
//					fragOnGPU := calculateFragRuntime(current.RemainMilli)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					current.Right = &TreeNode{
//						Job:         job,
//						GPUType:     currentGPUType,
//						GPUId:       gpuID,
//						StartTime:   types.Time(current.EndTime),
//						EndTime:     current.EndTime + types.Time(jobRuntime),
//						GPUMilli:    job.JobMeta().GPUMilli(),
//						Frag:        fragOnGPU,
//						Runtime:     types.Time(jobRuntime),
//						FGT:         fgtOnGPU,
//						RemainMilli: current.RemainMilli,
//					}
//					break
//				}
//			}
//			if current == nil {
//				remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
//				fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
//				fgtOnGPU := fragOnGPU * jobRuntime
//				root.Left = &TreeNode{
//					Job:         job,
//					GPUType:     currentGPUType,
//					GPUId:       gpuID,
//					StartTime:   types.Time(root.EndTime),
//					EndTime:     root.EndTime + types.Time(jobRuntime),
//					GPUMilli:    job.JobMeta().GPUMilli(),
//					Frag:        fragOnGPU,
//					Runtime:     types.Time(jobRuntime),
//					FGT:         fgtOnGPU,
//					RemainMilli: remainMilliOnGPU,
//				}
//			}
//		}
//	}
//	gpuTrees[gpuID] = root
//	return gpuTrees
//}
//
//// 找到FGT最小的GPU
//func findMinFGTGPU(jobFGTList map[types.GPUID]float64) types.GPUID {
//	minFGTValue := math.Inf(1)
//	var bestGPUID types.GPUID
//
//	for gpuID, fgtValue := range jobFGTList {
//		if fgtValue < minFGTValue {
//			minFGTValue = fgtValue
//			bestGPUID = gpuID
//		}
//	}
//
//	return bestGPUID
//}
//
//func (d *FragScheduler) OnScheduleEvent(event types.ScheduleEvent) {
//	switch e := event.(type) {
//	case *types.ScheduleEventJobsArrived:
//		{
//			d.allArrivedjobs = e.JobMetas()
//			newJobs := make([]types.Job, 0, len(e.JobMetas()))
//			for _, jobMeta := range e.JobMetas() {
//				newJobs = append(newJobs, d.cluster.InitJob(jobMeta))
//			}
//			// d.allWaitingJobs = e.JobMetas()
//			d.allWaitingJobs = newJobs
//			d.DoSchedule()
//		}
//	case *types.ScheduleEventJobsFinished:
//		{
//			d.DoSchedule()
//		}
//	}
//}
//
//func (d *FragScheduler) Name() string {
//	return "FragScheduler"
//}
//
//func (d *FragScheduler) Info() interface{} {
//	return d.Name()
//}
//
//func (d *FragScheduler) NextActiveScheduleTime() types.Time {
//	return types.Time(math.Inf(1))
//}
//
//func (d *FragScheduler) Record() *types.SchedulerRecord {
//	//return d.SchedulerRecord
//
//	return &types.SchedulerRecord{
//		DoScheduleRecords: d.SchedulerRecord,
//		//Extra:             d.RecordExtra(),
//	}
//}
//
////func (d *FragScheduler) RecordExtra() interface{} {
////	return d.FinishedJob
////}
//
//func printTree(node *TreeNode, level int) {
//	if node != nil {
//		printTree(node.Right, level+1)
//		//fmt.Printf("%s-> %s\n", getIndent(level), node.Job.JobName())
//		fmt.Printf("%s-> %s\n", strings.Repeat(" ", 4*level), node.Job.JobName())
//
//		printTree(node.Left, level+1)
//	}
//}
