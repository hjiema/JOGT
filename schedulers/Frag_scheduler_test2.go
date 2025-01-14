package schedulers

//
//import (
//	"DES-go/schedulers/types"
//	"math"
//	"time"
//)
//
//type FragScheduler struct {
//	gpuCluster types.Cluster
//	// waitingJobs 按照任意一个GPU上的速度排序，它的顺序决定了进行算法时的迭代（KMeans聚类时算法）的顺序。（这个顺序会造成怎样的影响，有待商榷）
//	allWaitingJobs []types.JobMeta
//	//sortedWaitingJobs     []types.Jobs
//	SchedulerRecord *types.SchedulerRecord
//}
//
//func NewFragScheduler() *FragScheduler {
//	//scheduler := &FragScheduler{
//	//	DoScheduleCallRecords: make([]*types.DoScheduleCallRecord, 0, 128),
//	//}
//	return &FragScheduler{
//		SchedulerRecord: &types.SchedulerRecord{
//			DoScheduleRecords: []*types.DoScheduleCallRecord{},
//		},
//	}
//}
//
//func (d *FragScheduler) DoSchedule() {
//	start := time.Now()
//	d.doSchedule()
//	duration := time.Since(start)
//	//d.DoScheduleCallRecords = append(d.DoScheduleCallRecords, &types.DoScheduleCallRecord{Duration: duration})
//	d.SchedulerRecord.DoScheduleRecords = append(d.SchedulerRecord.DoScheduleRecords, &types.DoScheduleCallRecord{Duration: duration})
//}
//func (d *FragScheduler) SetCluster(cluster types.Cluster) {
//	d.gpuCluster = cluster
//}
//
//func (d *FragScheduler) doSchedule() {
//	gpuTrees := make(map[types.GPUID]*TreeNode)
//	//println("hers")
//	//fmt.Printf("%+v\n", s.gpuCluster)
//	for _, gpuId := range d.gpuCluster.GPUIDs() {
//		// println(gpuId)
//		gpuTrees[gpuId] = nil
//	}
//
//	// jobs, popularity := ReadInputFiles("input.csv", "popularity.csv")
//	var jobs []types.JobMeta
//	jobs = d.allWaitingJobs
//
//	//for _, job := range s.allWaitingJobs {
//	//	fmt.Printf("Name: %s, GPUMilli: %d", job.JobName(), job.GPUMilli())
//	//	for gpuType, duration := range job.Durations() {
//	//		fmt.Printf("GPU Type: %s, Duration: %.2f", gpuType, duration)
//	//	}
//	//	print("\n")
//	//}
//
//	jobFGTList := make(map[types.GPUID]float64)
//	var gpus []types.GPUID
//	gpus = d.gpuCluster.GPUIDs()
//	//for _, gpuID := range gpus {
//	//	fmt.Printf("GPUID: %d\n", gpuID)
//	//}
//	var Gpus []types.GPU
//	for _, gpuID := range gpus {
//		var m types.GPU
//		m = d.gpuCluster.GPU(gpuID)
//		Gpus = append(Gpus, m)
//		// fmt.Printf("GPU ID: %d, GPU Type: %s\n", m.ID(), m.Type())
//	}
//
//	for _, job := range jobs {
//		jobFGTList = ScheduleJob(job, gpus, Gpus, gpuTrees, jobFGTList)
//		bestGPUID := FindMinFGTGPU(jobFGTList)
//		if bestGPUID != -1 {
//			gpuTrees = ScheduleJobForGPU(job, bestGPUID, Gpus, gpuTrees)
//		}
//	}
//}
//
////func GetGPUType(gpuid types.GPUID, Gpus []types.GPU) types.GPUType {
////	for _, gpu := range Gpus {
////		if gpu.ID() == gpuid {
////			return gpu.Type()
////			break
////		}
////	}
////}
//
//func GetGPUType(gpuId types.GPUID, Gpus []types.GPU) types.GPUType {
//	for _, gpu := range Gpus {
//		if gpu.ID() == gpuId {
//			return gpu.Type()
//		}
//	}
//	// 添加一个默认的返回值
//	return types.GPUType("Unknown")
//}
//
////func (s *FragScheduler) insertJob2SortedWaitingJobs(job types.Job) {
////	for _, gpuType := range s.gpuCluster.GPUTypes() {
////		// 获取对应GPU类型的等待队列,ls代表当前GPU类型下已经存在的等待作业列表.
////		ls := s.sortedWaitingJobs[]
////		// 计算作业在A100上的剩余时间.
////		target := job.RemainingDuration("A100")
////		// 标准库中的sort.Search函数,用于在有序切片中找到一个位置,使得插入元素后仍保持有序.返回第一个满足条件的索引i.
////		i := sort.Search(len(ls), func(i int) bool {
////			return ls[i].RemainingDuration("A100") >= target
////		})
////		// 调用一个工具方法InsertJobsSlice,将作业插入到列表ls的第i个位置.
////		s.sortedWaitingJobs[] = jobs_util.GetJobsSliceUtil().InsertJobsSlice(job, i, ls)
////	}
////}
//
//func (d *FragScheduler) OnScheduleEvent(event types.ScheduleEvent) {
//	switch e := event.(type) {
//	case *types.ScheduleEventJobsArrived:
//		{
//			d.allWaitingJobs = e.JobMetas()
//
//			// 插入所有的新作业后,调用方法进行调度.
//			d.DoSchedule()
//		}
//	case *types.ScheduleEventJobsFinished:
//		{
//			//if !s.hasEmptyGPUQueue() {
//			//	panic("!s.hasEmptyGPUQueue() when some jobs finished.")
//			//}
//			d.DoSchedule()
//		}
//	}
//}
//
////type Job struct {
////	Name      string
////	GPUMilli  int
////	A100      float64
////	GTX2080Ti float64
////	V100      float64
////}
//
////func (j *Job) GetTimeForGPU(gpuType string) float64 {
////	switch gpuType {
////	case "A100":
////		return j.A100
////	case "GTX2080Ti":
////		return j.GTX2080Ti
////	case "V100":
////		return j.V100
////	default:
////		return 0
////	}
////}
//
////func (m *JobMeta) Duration(gpu types.GPU) types.Duration {
////	return m.durations[gpu.Type()]
////}
//
////type GPU struct {
////	Type string
////	ID   int
////}
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
////func ReadInputFiles(jobFile, popularityFile string) ([]types.JobMeta, []map[string]float64) {
////	jobs := []types.JobMeta{}
////	popularity := []map[string]float64{}
////
////	jobFile, err := os.Open(jobFile)
////	if err != nil {
////		panic(err)
////	}
////	defer jobFile.Close()
////
////	reader := csv.NewReader(jobFile)
////	records, err := reader.ReadAll()
////	if err != nil {
////		panic(err)
////	}
////
////	for _, record := range records[1:] {
////		gpuMilli, _ := strconv.Atoi(record[1])
////		a100, _ := strconv.ParseFloat(record[2], 64)
////		gtx2080Ti, _ := strconv.ParseFloat(record[3], 64)
////		v100, _ := strconv.ParseFloat(record[4], 64)
////		job := &Job{
////			Name:      record[0],
////			GPUMilli:  gpuMilli,
////			A100:      a100,
////			GTX2080Ti: gtx2080Ti,
////			V100:      v100,
////		}
////		jobs = append(jobs, job)
////	}
////
////	popFile, err := os.Open(popularityFile)
////	if err != nil {
////		panic(err)
////	}
////	defer popFile.Close()
////
////	popReader := csv.NewReader(popFile)
////	popRecords, err := popReader.ReadAll()
////	if err != nil {
////		panic(err)
////	}
////
////	for _, record := range popRecords[1:] {
////		gpuMilli, _ := strconv.Atoi(record[0])
////		ratio, _ := strconv.ParseFloat(record[1], 64)
////		popularity = append(popularity, map[string]float64{
////			"gpu_milli": float64(gpuMilli),
////			"ratio":     ratio,
////		})
////	}
////
////	return jobs, popularity
////}
//
//func CalculateFragRuntime(remainMilli int) float64 {
//	aboveMilliPopularity := 1.0
//	frag := float64(remainMilli) * aboveMilliPopularity
//	return frag
//}
//
//func SumFGTOnGPU(node *TreeNode) float64 {
//	if node == nil {
//		return 0
//	}
//	leftSum := SumFGTOnGPU(node.Left)
//	rightSum := SumFGTOnGPU(node.Right)
//	return node.FGT + leftSum + rightSum
//}
//
//// 获取某个 GPU 类型对应的执行时间
//func getDurationForGPU(job types.JobMeta, gpuType types.GPUType) types.Duration {
//	durations := job.Durations()   // 获取 job 的 GPU 执行时间映射
//	duration := durations[gpuType] // 查找对应的 GPU 类型的执行时间
//	return duration
//}
//
//func ScheduleJob(job types.JobMeta, gpus []types.GPUID, Gpus []types.GPU, gpuTrees map[types.GPUID]*TreeNode, jobFGTList map[types.GPUID]float64) map[types.GPUID]float64 {
//	if len(jobFGTList) == 0 {
//		for _, gpu := range gpus {
//			jobFGTList[gpu] = 0
//		}
//	}
//	for _, gpuId := range gpus {
//		currentGPUType := GetGPUType(gpuId, Gpus)
//		jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//		root := gpuTrees[gpuId]
//
//		if root == nil {
//			remainMilliOnGPU := 1000 - job.GPUMilli()
//			fragOnGPU := CalculateFragRuntime(remainMilliOnGPU)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			jobFGTList[gpuId] = fgtOnGPU
//		} else {
//			if job.GPUMilli() <= root.RemainMilli {
//				root.RemainMilli -= job.GPUMilli()
//				current := root
//				for current.Right != nil {
//					current = current.Right
//				}
//				fragOnGPU := CalculateFragRuntime(root.RemainMilli)
//				fgtOnGPU := fragOnGPU * jobRuntime
//				current.FGT = current.Frag * (current.Runtime - jobRuntime)
//				jobFGTList[gpuId] = SumFGTOnGPU(root) + fgtOnGPU
//			} else {
//				current := root.Left
//				for current != nil {
//					if job.GPUMilli() > current.RemainMilli {
//						if current.Left == nil {
//							remainMilliOnGPU := 1000 - job.GPUMilli()
//							fragOnGPU := CalculateFragRuntime(remainMilliOnGPU)
//							fgtOnGPU := fragOnGPU * jobRuntime
//							jobFGTList[gpuId] = SumFGTOnGPU(root) + fgtOnGPU
//							break
//						} else {
//							current = current.Left
//						}
//					} else {
//						current.RemainMilli -= job.GPUMilli()
//						for current.Right != nil {
//							current = current.Right
//						}
//						fragOnGPU := CalculateFragRuntime(current.RemainMilli)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.FGT = current.Frag * (current.Runtime - jobRuntime)
//						jobFGTList[gpuId] = fgtOnGPU
//						break
//					}
//				}
//				if current == nil {
//					remainMilliOnGPU := 1000 - job.GPUMilli()
//					fragOnGPU := CalculateFragRuntime(remainMilliOnGPU)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					jobFGTList[gpuId] = SumFGTOnGPU(root) + fgtOnGPU
//				}
//			}
//		}
//	}
//	return jobFGTList
//}
//
//func ScheduleJobForGPU(job types.JobMeta, gpuID types.GPUID, Gpus []types.GPU, gpuTrees map[types.GPUID]*TreeNode) map[types.GPUID]*TreeNode {
//	var gpu types.GPU
//	for _, g := range Gpus {
//		if g.ID() == gpuID {
//			gpu = g
//			break
//		}
//	}
//	if gpu == nil {
//		panic("GPU with ID not found")
//	}
//	currentGPUType := GetGPUType(gpuID, Gpus)
//	jobRuntime := float64(getDurationForGPU(job, currentGPUType))
//	root := gpuTrees[gpuID]
//
//	if root == nil {
//		remainMilliOnGPU := 1000 - job.GPUMilli()
//		fragOnGPU := CalculateFragRuntime(remainMilliOnGPU)
//		fgtOnGPU := fragOnGPU * jobRuntime
//		root = &TreeNode{
//			Job:         job,
//			GPUType:     gpu.Type(),
//			GPUId:       gpu.ID(),
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
//			fragOnGPU := CalculateFragRuntime(root.RemainMilli)
//			fgtOnGPU := fragOnGPU * jobRuntime
//			current.FGT = current.Frag * (current.Runtime - jobRuntime)
//			current.Right = &TreeNode{
//				Job:         job,
//				GPUType:     gpu.Type(),
//				GPUId:       gpu.ID(),
//				StartTime:   current.EndTime,
//				EndTime:     current.EndTime + jobRuntime,
//				GPUMilli:    job.GPUMilli(),
//				Frag:        fragOnGPU,
//				Runtime:     jobRuntime,
//				FGT:         fgtOnGPU,
//				RemainMilli: root.RemainMilli,
//			}
//		} else {
//			current := root.Left
//			for current != nil {
//				if job.GPUMilli() > current.RemainMilli {
//					if current.Left == nil {
//						remainMilliOnGPU := 1000 - job.GPUMilli()
//						fragOnGPU := CalculateFragRuntime(remainMilliOnGPU)
//						fgtOnGPU := fragOnGPU * jobRuntime
//						current.Left = &TreeNode{
//							Job:         job,
//							GPUType:     gpu.Type(),
//							GPUId:       gpu.ID(),
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
//					fragOnGPU := CalculateFragRuntime(current.RemainMilli)
//					fgtOnGPU := fragOnGPU * jobRuntime
//					current.FGT = current.Frag * (current.Runtime - jobRuntime)
//					current.Right = &TreeNode{
//						Job:         job,
//						GPUType:     gpu.Type(),
//						GPUId:       gpu.ID(),
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
//			if current == nil {
//				remainMilliOnGPU := 1000 - job.GPUMilli()
//				fragOnGPU := CalculateFragRuntime(remainMilliOnGPU)
//				fgtOnGPU := fragOnGPU * jobRuntime
//				root.Left = &TreeNode{
//					Job:         job,
//					GPUType:     gpu.Type(),
//					GPUId:       gpu.ID(),
//					StartTime:   root.EndTime,
//					EndTime:     root.EndTime + jobRuntime,
//					GPUMilli:    job.GPUMilli(),
//					Frag:        fragOnGPU,
//					Runtime:     jobRuntime,
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
//func FindMinFGTGPU(jobFGTList map[types.GPUID]float64) types.GPUID {
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
////func PrintTree(node *TreeNode, level int) {
////	if node != nil {
////		PrintTree(node.Right, level+1)
////		fmt.Printf("%s-> %s\n", strings.Repeat(" ", 4*level), node.Job.Name)
////		PrintTree(node.Left, level+1)
////	}
////}
//
////func main() {
////	gpus := []*GPU{
////		{Type: "A100", ID: 1},
////		{Type: "V100", ID: 2},
////		{Type: "V100", ID: 3},
////		{Type: "GTX2080Ti", ID: 4},
////		{Type: "GTX2080Ti", ID: 5},
////	}
////
////	gpuTrees := make(map[int]*TreeNode)
////	for _, gpu := range gpus {
////		gpuTrees[gpu.ID] = nil
////	}
////
////	jobs, popularity := ReadInputFiles("input.csv", "popularity.csv")
////	jobFGTList := make(map[int]float64)
////
////	for _, job := range jobs {
////		jobFGTList = ScheduleJob(job, gpus, gpuTrees, jobFGTList, popularity)
////		bestGPUID := FindMinFGTGPU(jobFGTList)
////		if bestGPUID != -1 {
////			gpuTrees = ScheduleJobForGPU(job, bestGPUID, gpus, gpuTrees, popularity)
////		}
////	}
////}
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
////func (d *FragScheduler) Record() *types.SchedulerRecord {
////	return &types.SchedulerRecord{
////		DoScheduleRecords: d.DoScheduleCallRecords,
////	}
////}
//func (d *FragScheduler) Record() *types.SchedulerRecord {
//	return d.SchedulerRecord
//}
