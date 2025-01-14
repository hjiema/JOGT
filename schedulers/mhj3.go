package schedulers

//
//import (
//	"DES-go/schedulers/types"
//	"encoding/csv"
//	"fmt"
//	"log"
//	"math"
//	"os"
//	"sort"
//	"strconv"
//	"time"
//)
//
//// -----------------算法思想--------------------------------------------------
//// FragScheduler
//// 为保证公平，采取作业先到先服务
//// 作业逐个调度
//// 对于每个作业，计算调度到不同GPU上的fgd*JCT的值
//// 将作业调度到fgd*JCT值最小的GPU上。fgd*JCT的计算方式见frag_jct.py文件。
//// -------------------------------------------------------------------------------
//
//// 包含集群信息、等待调度的作业队列和调度调用记录
//type FragScheduler struct {
//	*SchedulerTemplate
//	cluster           types.Cluster
//	sortedWaitingJobs map[types.GPUType][]types.Job
//	DoScheduleCalls   []*types.DoScheduleCallRecord
//}
//
//// 创建并初始化一个新的FragScheduler实例
//func NewFragScheduler() *FragScheduler {
//	template := NewGreedySchedulerTemplate()
//	fs := &FragScheduler{
//		template,
//		nil,
//		make(map[types.GPUType][]types.Job),
//		make([]*types.DoScheduleCallRecord, 0),
//	}
//	template.impl = fs
//	return fs
//}
//
//func (s *FragScheduler) DoSchedule() {
//	start := time.Now()
//	s.doSchedule()
//	duration := time.Since(start)
//	s.DoScheduleCalls = append(s.DoScheduleCalls, &types.DoScheduleCallRecord{Duration: duration})
//}
//
//// 在有等待作业和空闲GPU队列的情况下,选择一个作业和对应的GPU队列进行调度.
//func (s *FragScheduler) doSchedule() {
//	for s.hasWaitingJob() && s.hasEmptyGPUQueue() {
//		emptyQueues := s.getEmptyGPUQueues()
//		targetJob, targetQueue := s.pickTarget(emptyQueues)
//		if targetJob == nil || targetQueue == nil {
//			panic("SchedulerTemplate targetJob == nil || targetQueue == nil")
//		}
//		s.removeFromSortedWaitingJobs(targetJob)
//		targetQueue.SetJobs(targetJob)
//	}
//}
//
//// 从空闲GPU队列中选择一个最适合当前作业的队列.
//func (s *FragScheduler) pickTarget(emptyQueues []types.GPUJobQueue) (types.Job, types.GPUJobQueue) {
//	var targetJob types.Job
//	var targetQueue types.GPUJobQueue
//	leastMetric := math.Inf(1)
//	for gpuType, waitingJobs := range s.sortedWaitingJobs {
//		if len(waitingJobs) == 0 {
//			continue
//		}
//		firstWaitingJob := waitingJobs[0]
//		var candidateQueue types.GPUJobQueue
//		for _, queue := range emptyQueues {
//			if queue.GPU().Type() != gpuType {
//				continue
//			}
//			if candidateQueue == nil {
//				candidateQueue = queue
//				break
//			}
//		}
//		if candidateQueue == nil {
//			continue
//		}
//		metric := s.calculateMetric(firstWaitingJob, candidateQueue.GPU())
//		if metric < leastMetric {
//			targetJob, targetQueue = firstWaitingJob, candidateQueue
//			leastMetric = metric
//		}
//	}
//	return targetJob, targetQueue
//}
//
//func (s *FragScheduler) hasWaitingJob() bool {
//	for _, l := range s.sortedWaitingJobs {
//		if len(l) > 0 {
//			return true
//		}
//	}
//	return false
//}
//
//func (s *FragScheduler) insertJob2SortedWaitingJobs(job types.Job) {
//	for _, gpuType := range s.cluster.GPUTypes() {
//		ls := s.sortedWaitingJobs[gpuType]
//		target := job.GPUMilli
//		i := sort.Search(len(ls), func(i int) bool {
//			return ls[i].GPUMilli >= target
//		})
//		s.sortedWaitingJobs[gpuType] = append(ls[:i], append([]types.Job{job}, ls[i:]...)...)
//	}
//}
//
//func (s *FragScheduler) removeFromSortedWaitingJobs(job types.Job) {
//	for _, gpuType := range s.cluster.GPUTypes() {
//		ls := s.sortedWaitingJobs[gpuType]
//		target := job.GPUMilli
//		i := sort.Search(len(ls), func(i int) bool {
//			return ls[i].GPUMilli >= target
//		})
//		if ls[i].GPUMilli != target {
//			panic("SchedulerTemplate removeFromSortedWaitingJobs ls[i].GPUMilli != target")
//		}
//		var targetIdx = -1
//		for ls[i].GPUMilli == target {
//			if ls[i].JobName == job.JobName {
//				targetIdx = i
//				break
//			}
//			i++
//		}
//		if targetIdx == -1 {
//			panic("SchedulerTemplate removeFromSortedWaitingJobs targetIdx == -1")
//		}
//		s.sortedWaitingJobs[gpuType] = append(ls[:targetIdx], ls[targetIdx+1:]...)
//	}
//}
//
//func (s *FragScheduler) hasEmptyGPUQueue() bool {
//	for _, queue := range s.cluster.GPUJobQueues() {
//		if len(queue.Jobs()) == 0 {
//			return true
//		}
//	}
//	return false
//}
//
//func (s *FragScheduler) getEmptyGPUQueues() []types.GPUJobQueue {
//	queues := make([]types.GPUJobQueue, 0, len(s.cluster.GPUJobQueues()))
//	for _, queue := range s.cluster.GPUJobQueues() {
//		if len(queue.Jobs()) == 0 {
//			queues = append(queues, queue)
//		}
//	}
//	return queues
//}
//
//func (s *FragScheduler) SetCluster(cluster types.Cluster) {
//	s.cluster = cluster
//	s.sortedWaitingJobs = make(map[types.GPUType][]types.Job)
//	for _, gpuType := range s.cluster.GPUTypes() {
//		s.sortedWaitingJobs[gpuType] = make([]types.Job, 0)
//	}
//}
//
//func (s *FragScheduler) OnScheduleEvent(event types.ScheduleEvent) {
//	switch e := event.(type) {
//	case *types.ScheduleEventJobsArrived:
//		for _, jobMeta := range e.JobMetas() {
//			s.insertJob2SortedWaitingJobs(types.Job{
//				JobName: jobMeta.JobName(),
//				// 其他字段根据需要填充
//			})
//		}
//		s.DoSchedule()
//	case *types.ScheduleEventJobsFinished:
//		if !s.hasEmptyGPUQueue() {
//			panic("!s.hasEmptyGPUQueue() when some jobs finished.")
//		}
//		s.DoSchedule()
//	}
//}
//
//func (s *FragScheduler) NextActiveScheduleTime() types.Time {
//	return types.Time(math.Inf(1))
//}
//
//func (s *FragScheduler) Name() string {
//	return fmt.Sprintf("FragScheduler")
//}
//
//func (s *FragScheduler) Info() interface{} {
//	return s.Name()
//}
//
//func (s *FragScheduler) Record() *types.SchedulerRecord {
//	return &types.SchedulerRecord{
//		DoScheduleRecords: s.DoScheduleCalls,
//	}
//}
//
//func readTasksFromCSV(filename string) []types.Job {
//	file, err := os.Open(filename)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	reader := csv.NewReader(file)
//	records, err := reader.ReadAll()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var tasks []types.Job
//	for _, record := range records[1:] {
//		normJobSubmitTime, _ := strconv.Atoi(record[1])
//		gpuMilli, _ := strconv.Atoi(record[2])
//		ddl, _ := strconv.Atoi(record[3])
//		a100, _ := strconv.ParseFloat(record[4], 64)
//		gtx2080Ti, _ := strconv.ParseFloat(record[5], 64)
//		v100, _ := strconv.ParseFloat(record[6], 64)
//
//		task := types.Job{
//			JobName:           record[0],
//			NormJobSubmitTime: normJobSubmitTime,
//			GPUMilli:          gpuMilli,
//			DDL:               ddl,
//			A100:              a100,
//			GTX2080Ti:         gtx2080Ti,
//			V100:              v100,
//		}
//		tasks = append(tasks, task)
//	}
//	return tasks
//}
//
//func readPopularityFromCSV(filename string) map[int]float64 {
//	file, err := os.Open(filename)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer file.Close()
//
//	reader := csv.NewReader(file)
//	records, err := reader.ReadAll()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	popularity := make(map[int]float64)
//	for _, record := range records[1:] {
//		gpuMilli, _ := strconv.Atoi(record[0])
//		ratio, _ := strconv.ParseFloat(record[1], 64)
//		popularity[gpuMilli] = ratio
//	}
//	return popularity
//}
//
//// 计算作业调度到不同GPU上的fgd*jct值,用于选择最优的GPU队列.
//func (s *FragScheduler) calculateMetric(task types.Job, gpu types.GPU) float64 {
//	fragSize := (1000 - task.GPUMilli) * s.calculatePopularitySum(task, gpu)
//	jct := s.calculateWaitTime(task, gpu) + s.getRuntime(task, gpu.Type())
//	return float64(gpu.ID()) * fragSize * jct
//}
//
//// 计算当前GPU上所有任务的流行度之和
//func (s *FragScheduler) calculatePopularitySum(task types.Job, gpu types.GPU) float64 {
//	sum := 0.0
//	for _, t := range gpu.Tasks() {
//		if t.GPUMilli < (1000 - task.GPUMilli) {
//			sum += s.popularity[t.GPUMilli]
//		}
//	}
//	return sum
//}
//
//// 计算作业在当前GPU上的等待时间
//func (s *FragScheduler) calculateWaitTime(task types.Job, gpu types.GPU) int {
//	waitTime := 0
//	for _, t := range gpu.Tasks() {
//		if gpu.Used()+task.GPUMilli > gpu.Capacity() {
//			waitTime += int(s.getRuntime(t, gpu.Type()))
//		}
//	}
//	return waitTime
//}
//
//// 获取作业的运行时间.
//func (s *FragScheduler) getRuntime(task types.Job, gpuType string) float64 {
//	switch gpuType {
//	case "A100":
//		return task.A100
//	case "GTX2080Ti":
//		return task.GTX2080Ti
//	case "V100":
//		return task.V100
//	default:
//		return math.MaxFloat64
//	}
//}