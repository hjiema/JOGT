package schedulers

import (
	"DES-go/schedulers/types"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

type FragScheduler struct {
	cluster         types.Cluster
	allArrivedjobs  []types.JobMeta
	allWaitingJobs  []types.Job
	SchedulerRecord []*types.DoScheduleCallRecord
	FinishedJob     []types.Job
	gpuTrees        map[types.GPUID]*TreeNode // 新增，用于存储每个GPU的任务调度树
}

type TreeNode struct {
	Job         types.Job
	GPUType     types.GPUType
	GPUId       types.GPUID
	StartTime   types.Time
	EndTime     types.Time
	GPUMilli    int //当前GPU横向被占用的显存大小
	Frag        float64
	Runtime     types.Time
	FGT         float64
	RemainMilli int //横向GPU剩余显存大小
	Left        *TreeNode
	Right       *TreeNode
}

func NewFragScheduler() *FragScheduler {

	return &FragScheduler{
		SchedulerRecord: make([]*types.DoScheduleCallRecord, 0, 128),
		//SchedulerRecord: &types.SchedulerRecord{
		//	DoScheduleRecords: []*types.DoScheduleCallRecord{},
		//},
		gpuTrees: make(map[types.GPUID]*TreeNode), // 初始化GPU任务调度树

	}
}

func (d *FragScheduler) DoSchedule() {
	start := time.Now()
	d.doSchedule()
	duration := time.Since(start)
	d.SchedulerRecord = append(d.SchedulerRecord, &types.DoScheduleCallRecord{Duration: duration})
	// d.SchedulerRecord.Extra = d.FinishedJob

	//for _, gpuId := range d.cluster.GPUIDs() {
	//	// 获取对应的根 TreeNode
	//	println(gpuId)
	//	printTree(d.gpuTrees[gpuId], 0)
	//	println("=====================")
	//}
}

func (d *FragScheduler) SetCluster(cluster types.Cluster) {
	d.cluster = cluster
	//fmt.Printf("%+v\n", d.gpuCluster.GPUIDs())
	for _, gpuId := range d.cluster.GPUIDs() {
		//println(gpuId)
		d.gpuTrees[gpuId] = nil // 初始化每个GPU的任务树
	}
}

func calculateFragRuntime(remainMilli int) float64 {
	// 简化的Fragment计算函数
	// 修改后是形式函数，无实际意义
	aboveMilliPopularity := 1.0
	return float64(remainMilli) * aboveMilliPopularity
}

func sumFGTOnGPU(node *TreeNode) float64 {
	if node == nil {
		return 0
	}
	return node.FGT + sumFGTOnGPU(node.Left)
}

// 获取某个 GPU 类型对应的执行时间
func getDurationForGPU(job types.Job, gpuType types.GPUType) types.Duration {
	durations := job.JobMeta().Durations()
	return durations[gpuType]
}

const A100 types.GPUType = "A100"

func sortJobsByGpumilli(jobs []types.Job) {
	sort.Slice(jobs, func(i, j int) bool {
		//return jobs[i].Durations()[A100] > jobs[j].Durations()[A100]
		return jobs[i].JobMeta().GPUMilli() < jobs[j].JobMeta().GPUMilli()
	})
}

func sortJobsByA100Duration(jobs []types.Job) {
	sort.Slice(jobs, func(i, j int) bool {
		//return jobs[i].Durations()[A100] > jobs[j].Durations()[A100]
		return jobs[i].JobMeta().Durations()[A100] < jobs[j].JobMeta().Durations()[A100]
	})
}
func sortJobsByddl(jobs []types.Job) {
	sort.Slice(jobs, func(i, j int) bool {
		//return jobs[i].Durations()[A100] > jobs[j].Durations()[A100]
		return jobs[i].JobMeta().DDL() < jobs[j].JobMeta().DDL()
	})
}

func (d *FragScheduler) doSchedule() {
	// 初始化任务等待队列
	var jobs []types.Job
	jobs = d.allWaitingJobs

	// jobFGTList列表，存储每个GPU的FGT值
	jobFGTList := make(map[types.GPUID]float64)
	sortJobsByGpumilli(jobs)
	sortJobsByA100Duration(jobs)
	sortJobsByddl(jobs)
	for gpuid := range d.cluster.GPUIDs() {
		println(gpuid)
	}

	for _, job := range jobs {
		// 调度每个作业，计算FGT值
		jobFGTList = scheduleJob(job, d.cluster.GPUIDs(), d.cluster, d.gpuTrees, jobFGTList)
		// 找到FGT最小的GPU
		bestGPUID := findMinFGTGPU(jobFGTList)
		if bestGPUID != -1 {
			// 将作业调度到最优GPU
			d.gpuTrees = scheduleJobForGPU(job, bestGPUID, d.cluster, d.gpuTrees)
		}
	}

	var finishedallJob []Jobfish

	for _, gpuId := range d.cluster.GPUIDs() {
		// 获取对应的根 TreeNode
		rootTreeNode := d.gpuTrees[gpuId]

		// 如果根节点不为 nil，开始递归遍历树
		if rootTreeNode != nil {
			traverseTree(rootTreeNode, &finishedallJob, d.cluster)
		}
	}
	writeFinishedJobsToFile(finishedallJob)

}

func writeFinishedJobsToFile(finishedallJob []Jobfish) error {
	// 设置文件路径
	filePath := "/hydra/data/finished_jobs.json"

	// 创建或打开文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	// 将数据编码为 JSON 并写入文件
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(finishedallJob); err != nil {
		return fmt.Errorf("无法写入 JSON 数据: %v", err)
	}

	//fmt.Printf("已成功将数据写入 %s\n", filePath)
	return nil
}

func traverseTree(treeNode *TreeNode, finishedallJob *[]Jobfish, cluster types.Cluster) {
	// 如果 treeNode 为 nil，结束递归
	if treeNode == nil {
		return
	}

	// 创建 Jobfish 对象
	jobchange := Jobfish{
		JobName:             treeNode.Job.JobName(),
		GpuID:               treeNode.GPUId,
		GpuType:             treeNode.GPUType,
		Gpu:                 cluster.GPU(treeNode.GPUId), // 假设 d.cluster.GPU 返回 GPU 对象
		FirstExecutionTime:  treeNode.StartTime,
		FinishExecutionTime: treeNode.EndTime,
		RemainingRatio:      0, // 这里根据需要设置剩余比例
		IsRunning:           false,
	}

	*finishedallJob = append(*finishedallJob, jobchange)

	//// 将 jobchange 添加到 finishedallJob 列表中
	//finishedallJob[gpuID] = append(finishedallJob[gpuID], jobchange)

	// 递归遍历左子树
	traverseTree(treeNode.Left, finishedallJob, cluster)

	// 递归遍历右子树
	traverseTree(treeNode.Right, finishedallJob, cluster)
}

type Jobfish struct {
	JobName types.JobName
	//executionDetail     *JobExecutionDetail
	GpuID               types.GPUID
	GpuType             types.GPUType
	Gpu                 types.GPU
	FirstExecutionTime  types.Time
	FinishExecutionTime types.Time
	RemainingRatio      float64 // 任务未执行完的部分的剩余比例。
	IsRunning           bool
}

// copyTree 实现树的深拷贝
func copyTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	newNode := &TreeNode{
		Job:         root.Job,
		GPUType:     root.GPUType,
		GPUId:       root.GPUId,
		StartTime:   root.StartTime,
		EndTime:     root.EndTime,
		GPUMilli:    root.GPUMilli,
		Frag:        root.Frag,
		Runtime:     root.Runtime,
		FGT:         root.FGT,
		RemainMilli: root.RemainMilli,
	}
	newNode.Left = copyTree(root.Left)
	newNode.Right = copyTree(root.Right)
	return newNode
}

// 交换两个节点的信息
func swapNodes(node1, node2 *TreeNode) {
	node1.Job, node2.Job = node2.Job, node1.Job
	node1.GPUType, node2.GPUType = node2.GPUType, node1.GPUType
	node1.GPUId, node2.GPUId = node2.GPUId, node1.GPUId
	node1.StartTime, node2.StartTime = node2.StartTime, node1.StartTime
	node1.EndTime, node2.EndTime = node2.EndTime, node1.EndTime
	node1.GPUMilli, node2.GPUMilli = node2.GPUMilli, node1.GPUMilli
	node1.Frag, node2.Frag = node2.Frag, node1.Frag
	node1.Runtime, node2.Runtime = node2.Runtime, node1.Runtime
	node1.FGT, node2.FGT = node2.FGT, node1.FGT
	node1.RemainMilli, node2.RemainMilli = node2.RemainMilli, node1.RemainMilli
}

// 模拟调度作业到每个GPU并计算FGT值
func scheduleJob(job types.Job, gpus []types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode, jobFGTList map[types.GPUID]float64) map[types.GPUID]float64 {
	if len(jobFGTList) == 0 {
		for _, gpu := range gpus {
			jobFGTList[gpu] = 0
		}
	}

	// 模拟调度：深拷贝gpuTrees，以免影响原始树
	simulatedTrees := make(map[types.GPUID]*TreeNode)
	for gpuID, root := range gpuTrees {
		simulatedTrees[gpuID] = copyTree(root) // 深拷贝树
	}

	for _, gpuId := range gpus {
		// 获取当前GPU类型
		currentGPUType := cluster.GPU(gpuId).Type()
		// 任务在当前GPU类型上的运行时间
		jobRuntime := float64(getDurationForGPU(job, currentGPUType))
		root := simulatedTrees[gpuId]

		if root == nil { // 根节点
			remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
			fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
			fgtOnGPU := fragOnGPU * jobRuntime
			jobFGTList[gpuId] = fgtOnGPU * jobRuntime
		} else {
			if job.JobMeta().GPUMilli() <= root.RemainMilli { // 根的右子树
				root.RemainMilli -= job.JobMeta().GPUMilli()
				current := root
				for current.Right != nil {
					current = current.Right
				}
				// fragOnGPU := calculateFragRuntime(root.RemainMilli)
				// fgtOnGPU := fragOnGPU * jobRuntime // 即fragOnGPU*任务的runtime
				// current.FGT = current.Frag * math.Abs(float64(current.Runtime)-float64(jobRuntime))
				if root.Runtime < types.Time(jobRuntime) {
					root.FGT = root.FGT + 1000*math.Abs(float64(current.Runtime)-float64(jobRuntime)) - float64(job.JobMeta().GPUMilli())*jobRuntime
					// current.FGT = float64(current.GPUMilli) * math.Abs(float64(current.Runtime)-float64(jobRuntime))
					jobFGTList[gpuId] = root.FGT * (float64(current.StartTime) + float64(jobRuntime)) // 为了实现负载均衡，碎片大小乘以总的JCT
					// 不进行更新只模拟调度
					// root.Runtime = types.Time(jobRuntime)
				} else {
					root.FGT = root.FGT - float64(job.JobMeta().GPUMilli())*jobRuntime
					jobFGTList[gpuId] = root.FGT * (float64(current.StartTime) + float64(current.Runtime))
				}
			} else {
				current := root.Left
				for current != nil { // 找到最左侧
					if job.JobMeta().GPUMilli() > current.RemainMilli {
						if current.Left == nil {
							remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
							fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
							fgtOnGPU := fragOnGPU * jobRuntime
							jobFGTList[gpuId] = fgtOnGPU * (float64(current.EndTime) + float64(jobRuntime))
							break
						} else {
							current = current.Left
						}
					} else { // 左子树的右子树
						current.RemainMilli -= job.JobMeta().GPUMilli()
						target := current
						for current.Right != nil {
							current = current.Right
						}
						// fragOnGPU := calculateFragRuntime(current.RemainMilli)
						// fgtOnGPU := fragOnGPU * jobRuntime
						// current.FGT = current.Frag * math.Abs(float64(current.Runtime)-float64(jobRuntime))
						if current.Runtime < types.Time(jobRuntime) {
							target.FGT = target.FGT + 1000*math.Abs(float64(current.Runtime)-float64(jobRuntime)) - float64(job.JobMeta().GPUMilli())*jobRuntime
							// current.FGT = float64(current.GPUMilli) * math.Abs(float64(current.Runtime)-float64(jobRuntime))
							// current.FGT = current.Frag * float64(jobRuntime)
							jobFGTList[gpuId] = target.FGT * (float64(current.StartTime) + float64(jobRuntime))
							// current.Runtime = types.Time(jobRuntime)
						} else {
							target.FGT = target.FGT - float64(job.JobMeta().GPUMilli())*jobRuntime
							// current.FGT = current.Frag * math.Abs(float64(current.Runtime)-float64(jobRuntime))
							jobFGTList[gpuId] = target.FGT * (float64(current.StartTime) + float64(current.Runtime))
						}
						break
					}
				}
				if current == nil { //左子树为空
					remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
					fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
					fgtOnGPU := fragOnGPU * jobRuntime
					jobFGTList[gpuId] = fgtOnGPU * (float64(root.EndTime) + float64(jobRuntime))
				}
			}
		}
	}

	return jobFGTList
}

// 实际调度
func scheduleJobForGPU(job types.Job, gpuID types.GPUID, cluster types.Cluster, gpuTrees map[types.GPUID]*TreeNode) map[types.GPUID]*TreeNode {
	currentGPUType := cluster.GPU(gpuID).Type()
	jobRuntime := float64(getDurationForGPU(job, currentGPUType))
	root := gpuTrees[gpuID]

	if root == nil {
		remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
		fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
		fgtOnGPU := fragOnGPU * jobRuntime
		root = &TreeNode{
			Job:         job,
			GPUType:     currentGPUType,
			GPUId:       gpuID,
			StartTime:   0,
			EndTime:     types.Time(jobRuntime),
			GPUMilli:    job.JobMeta().GPUMilli(),
			Frag:        fragOnGPU,
			Runtime:     types.Time(jobRuntime),
			FGT:         fgtOnGPU,
			RemainMilli: remainMilliOnGPU,
		}
		gpuTrees[gpuID] = root
	} else {
		if job.JobMeta().GPUMilli() <= root.RemainMilli {
			root.RemainMilli -= job.JobMeta().GPUMilli()
			current := root
			for current.Right != nil {
				current = current.Right
			}
			fragOnGPU := calculateFragRuntime(root.RemainMilli)
			fgtOnGPU := fragOnGPU * jobRuntime
			current.Right = &TreeNode{
				Job:         job,
				GPUType:     currentGPUType,
				GPUId:       gpuID,
				StartTime:   types.Time(root.StartTime),
				EndTime:     root.StartTime + types.Time(jobRuntime),
				GPUMilli:    job.JobMeta().GPUMilli(),
				Frag:        fragOnGPU,
				Runtime:     types.Time(jobRuntime),
				FGT:         fgtOnGPU,
				RemainMilli: root.RemainMilli,
			}
			if root.Runtime < types.Time(jobRuntime) {
				// root.FGT = root.FGT - fgtOnGPU
				root.FGT = root.FGT + 1000*math.Abs(float64(current.Runtime)-float64(jobRuntime)) - float64(job.JobMeta().GPUMilli())*jobRuntime
				// root.Runtime = types.Time(jobRuntime)
				// root.EndTime = root.StartTime + types.Time(jobRuntime)
				swapNodes(root, current.Right)
			} else {
				root.FGT = root.FGT - float64(job.JobMeta().GPUMilli())*jobRuntime
			}
		} else {
			current := root.Left
			for current != nil {
				if job.JobMeta().GPUMilli() > current.RemainMilli {
					if current.Left == nil {
						remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
						fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
						fgtOnGPU := fragOnGPU * (jobRuntime + float64(current.EndTime))
						current.Left = &TreeNode{
							Job:         job,
							GPUType:     currentGPUType,
							GPUId:       gpuID,
							StartTime:   types.Time(current.EndTime),
							EndTime:     current.EndTime + types.Time(jobRuntime),
							GPUMilli:    job.JobMeta().GPUMilli(),
							Frag:        fragOnGPU,
							Runtime:     types.Time(jobRuntime),
							FGT:         fgtOnGPU,
							RemainMilli: remainMilliOnGPU,
						}
						break
					} else {
						current = current.Left
					}
				} else {
					current.RemainMilli -= job.JobMeta().GPUMilli()
					target := current
					for current.Right != nil {
						current = current.Right
					}
					fragOnGPU := calculateFragRuntime(current.RemainMilli)
					fgtOnGPU := fragOnGPU * (jobRuntime + float64(current.StartTime))
					current.Right = &TreeNode{
						Job:         job,
						GPUType:     currentGPUType,
						GPUId:       gpuID,
						StartTime:   types.Time(target.StartTime),
						EndTime:     target.StartTime + types.Time(jobRuntime),
						GPUMilli:    job.JobMeta().GPUMilli(),
						Frag:        fragOnGPU,
						Runtime:     types.Time(jobRuntime),
						FGT:         fgtOnGPU,
						RemainMilli: current.RemainMilli,
					}
					if target.Runtime < types.Time(jobRuntime) {
						target.FGT = target.FGT + 1000*math.Abs(float64(current.Runtime)-float64(jobRuntime)) - float64(job.JobMeta().GPUMilli())*jobRuntime
						//target.Runtime = types.Time(jobRuntime)
						//target.EndTime = target.StartTime + types.Time(jobRuntime)
						swapNodes(target, current.Right)
					} else {
						target.FGT = target.FGT - float64(job.JobMeta().GPUMilli())*jobRuntime
					}
					break
				}
			}
			if current == nil {
				remainMilliOnGPU := 1000 - job.JobMeta().GPUMilli()
				fragOnGPU := calculateFragRuntime(remainMilliOnGPU)
				fgtOnGPU := fragOnGPU * (jobRuntime + float64(root.EndTime))
				root.Left = &TreeNode{
					Job:         job,
					GPUType:     currentGPUType,
					GPUId:       gpuID,
					StartTime:   types.Time(root.EndTime),
					EndTime:     root.EndTime + types.Time(jobRuntime),
					GPUMilli:    job.JobMeta().GPUMilli(),
					Frag:        fragOnGPU,
					Runtime:     types.Time(jobRuntime),
					FGT:         fgtOnGPU,
					RemainMilli: remainMilliOnGPU,
				}
			}
		}
	}
	gpuTrees[gpuID] = root
	return gpuTrees
}

// 找到FGT最小的GPU
func findMinFGTGPU(jobFGTList map[types.GPUID]float64) types.GPUID {
	minFGTValue := math.Inf(1)
	var bestGPUID types.GPUID

	for gpuID, fgtValue := range jobFGTList {
		if fgtValue < minFGTValue {
			minFGTValue = fgtValue
			bestGPUID = gpuID
		}
	}

	return bestGPUID
}

func (d *FragScheduler) OnScheduleEvent(event types.ScheduleEvent) {
	switch e := event.(type) {
	case *types.ScheduleEventJobsArrived:
		{
			d.allArrivedjobs = e.JobMetas()
			newJobs := make([]types.Job, 0, len(e.JobMetas()))
			for _, jobMeta := range e.JobMetas() {
				newJobs = append(newJobs, d.cluster.InitJob(jobMeta))
			}
			// d.allWaitingJobs = e.JobMetas()
			d.allWaitingJobs = newJobs
			d.DoSchedule()
		}
	case *types.ScheduleEventJobsFinished:
		{
			d.DoSchedule()
		}
	}
}

func (d *FragScheduler) Name() string {
	return "FragScheduler"
}

func (d *FragScheduler) Info() interface{} {
	return d.Name()
}

func (d *FragScheduler) NextActiveScheduleTime() types.Time {
	return types.Time(math.Inf(1))
}

func (d *FragScheduler) Record() *types.SchedulerRecord {
	//return d.SchedulerRecord

	return &types.SchedulerRecord{
		DoScheduleRecords: d.SchedulerRecord,
		//Extra:             d.RecordExtra(),
	}
}

func printTree(node *TreeNode, level int) {
	if node != nil {
		printTree(node.Right, level+1)
		//fmt.Printf("%s-> %s\n", getIndent(level), node.Job.JobName())
		fmt.Printf("%s-> %s\n", strings.Repeat(" ", 4*level), node.Job.JobName())

		printTree(node.Left, level+1)
	}
}
