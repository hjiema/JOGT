package schedulers

import (
	"DES-go/schedulers/jobs_util"
	"DES-go/schedulers/types"
	"DES-go/simulator"
	"fmt"
	"math"
	"sort"
	"time"
)

// GavelScheduler
// 采取了Shortest Job First策略。分别实现了抢占与非抢占。
// 采取调度器内部做任务缓存队列的思想。将所有未在GPU上运行的任务缓存到sortedWaitingJobs中来。
// 这就是说集群中的每个GPUJobQueue最多只有一个任务在队列中。
// 每当有新的任务到来时，将它按照remainingDuration排序的顺序，插入到sortedWaitingJobs中。
// 当发现有空闲的GPU槽，并且存在任务正在等待队列中时，将进行调度。
// 一旦当收到有job在GPU上运行结束的消息，能够确定的是那个GPU一定能够空闲下来。
// SJF具体的算法过程：sortedWaitingJobs是按照GPUType分类的，在迭代中选择下一个要上GPU的任务时，
// 将所有GPUType中最短的任务，与可能的全部GPUQueue做一一对应，在这些两两组合中，选择一个最短的任务。
// 将它放置到GPU队列后，再迭代这个过程。
// 抢占与非抢占的区别就在于，非抢占只遍历那些空闲的GPU，而抢占式，则会先将GPU队列中的全部任务卸下，
// 让它们全部变为空闲的GPU，再按照非抢占的调度算法执行即可。
type GavelScheduler struct {
	*SchedulerTemplate
	//
	//cluster types.Cluster
	//
	//// 等待队列中的所有任务，其分别在每种类型的GPU上，按照RemainingDuration排序。
	//sortedWaitingJobs map[types.GPUType][]types.Job
	//
	//DoScheduleCalls []*types.DoScheduleCallRecord
}

//func NewSJFSchedulerXXX() *GavelScheduler {
//	return &GavelScheduler{
//		DoScheduleCalls: make([]*types.DoScheduleCallRecord, 0),
//	}
//}

func NewGavelScheduler() *GavelScheduler {
	template := NewGreedySchedulerTemplate()
	edf := &GavelScheduler{
		template,
	}
	template.impl = edf
	return edf
}

func (s *GavelScheduler) DoSchedule() {
	start := time.Now()
	s.doSchedule()
	duration := time.Since(start)
	s.DoScheduleCalls = append(s.DoScheduleCalls, &types.DoScheduleCallRecord{Duration: duration})
}

// 从等待队列中挑选作业,将其分配到空闲GPU队列中进行执行.首先尝试在所有空闲GPU队列中找到最合适的作业,并将改作业分配给相应的GPU队列.
func (s *GavelScheduler) doSchedule() {
	for s.hasWaitingJob() && s.hasEmptyGPUQueue() {
		// 获取空闲GPU队列.
		emptyQueues := s.getEmptyGPUQueues()
		// 用于存储选择的目标作业和目标GPU队列.
		var targetJob types.Job = nil
		var targetQueue types.GPUJobQueue = nil
		// 初始化为正无穷大,用于跟踪当前找到的剩余执行时间最短的作业
		leastRemainingDuration := types.Duration(math.Inf(1))
		// 遍历全部的waiting job，按照gpu type进行分类，在每个waitingJobs上找首个job（即在这个类型上剩余执行时间最短的任务）
		// 遍历结束后，找到一个速度最快的任务。
		for gpuType, waitingJobs := range s.sortedWaitingJobs {
			// 当前GPU类型下的等待作业列表,若等待当前GPU类型的作业为空,则跳过该GPU类型.
			if len(waitingJobs) == 0 {
				continue
			}
			// 取出当前GPU类型等待列表中的第一个作业,首个作业是该类型下剩余时间最短的
			firstWaitingJob := waitingJobs[0]
			var candidateQueue types.GPUJobQueue = nil
			// 遍历所有空闲GPU队列,寻找与当前作业所需GPU类型匹配的空闲队列.
			for _, queue := range emptyQueues {
				if queue.GPU().Type() != gpuType {
					continue
				}
				if candidateQueue == nil {
					candidateQueue = queue
					break
				}
			}
			if candidateQueue == nil {
				continue
			}
			if targetJob == nil {
				targetJob, targetQueue = firstWaitingJob, candidateQueue
				leastRemainingDuration = targetJob.RemainingDuration(gpuType)
			} else if rd := firstWaitingJob.RemainingDuration(gpuType); rd < leastRemainingDuration {
				targetJob, targetQueue = firstWaitingJob, candidateQueue
				leastRemainingDuration = rd
			}
		}
		// 没有找到目标作业和队列,则抛出异常
		if targetJob == nil || targetQueue == nil {
			panic("SchedulerTemplate targetJob == nil || targetQueue == nil")
		}
		s.removeFromSortedWaitingJobs(targetJob)
		targetQueue.SetJobs(targetJob)
	}
}

func (s *GavelScheduler) pickTarget(emptyQueues []types.GPUJobQueue) (types.Job, types.GPUJobQueue) {
	var targetJob types.Job = nil
	var targetQueue types.GPUJobQueue = nil
	leastRemainingDuration := types.Duration(math.Inf(1))
	for gpuType, waitingJobs := range s.sortedWaitingJobs {
		if len(waitingJobs) == 0 {
			continue
		}
		firstWaitingJob := waitingJobs[0]
		var candidateQueue types.GPUJobQueue = nil
		for _, queue := range emptyQueues {
			if queue.GPU().Type() != gpuType {
				continue
			}
			if candidateQueue == nil {
				candidateQueue = queue
				break
			}
		}
		if candidateQueue == nil {
			continue
		}
		if targetJob == nil {
			targetJob, targetQueue = firstWaitingJob, candidateQueue
			leastRemainingDuration = targetJob.RemainingDuration(gpuType)
		} else if rd := firstWaitingJob.RemainingDuration(gpuType); rd < leastRemainingDuration {
			targetJob, targetQueue = firstWaitingJob, candidateQueue
			leastRemainingDuration = rd
		}
	}
	return targetJob, targetQueue
}

func (s *GavelScheduler) hasWaitingJob() bool {
	for _, l := range s.sortedWaitingJobs {
		if len(l) > 0 {
			return true
		}
	}
	return false
}

// 将一个新的作业插入到调度器的等待队列中,并确保等待队列按照排序规则保持有序.
// 这个函数根据作业在不同GPU类型上的剩余执行时间,将作业插入到对应的等待队列中.
func (s *GavelScheduler) insertJob2SortedWaitingJobs(job types.Job) {
	for _, gpuType := range s.cluster.GPUTypes() {
		// 获取对应GPU类型的等待队列,ls代表当前GPU类型下已经存在的等待作业列表.
		ls := s.sortedWaitingJobs[gpuType]
		// 计算作业的剩余时间.传入当前GPU类型,计算作业在该GPU类型上的剩余执行时间.
		target := job.RemainingDuration(gpuType)
		// 标准库中的sort.Search函数,用于在有序切片中找到一个位置,使得插入元素后仍保持有序.返回第一个满足条件的索引i.
		i := sort.Search(len(ls), func(i int) bool {
			return ls[i].RemainingDuration(gpuType) >= target
		})
		// 调用一个工具方法InsertJobsSlice,将作业插入到列表ls的第i个位置.
		s.sortedWaitingJobs[gpuType] = jobs_util.GetJobsSliceUtil().InsertJobsSlice(job, i, ls)
	}
}

// 在调度器的等待作业队列中,找到并移除特定的作业. 使用二分查找和线性扫描想结合的方式,确保在sortedWaitingJobs中精确定位并删除目标作业.
func (s *GavelScheduler) removeFromSortedWaitingJobs(job types.Job) {
	for _, gpuType := range s.cluster.GPUTypes() {
		// 遍历GPU类型,对cluster中的每种GPU类型执行相同操作.获取当前GPU类型对应的等待作业列表.
		ls := s.sortedWaitingJobs[gpuType]
		// 获取目标作业在当前GPU类型下的剩余时间,用于二分查找.
		target := job.RemainingDuration(gpuType)
		i := sort.Search(len(ls), func(i int) bool {
			return ls[i].RemainingDuration(gpuType) >= target
		})
		if ls[i].RemainingDuration(gpuType) != target {
			panic("SchedulerTemplate removeFromSortedWaitingJobs ls[i].RemainingDuration(gpuType) != target")
		}
		var targetIdx = -1
		for ls[i].RemainingDuration(gpuType) == target {
			if ls[i].JobName() == job.JobName() {
				targetIdx = i
				break
			}
			i++
		}
		if targetIdx == -1 {
			panic("SchedulerTemplate removeFromSortedWaitingJobs targetIdx == -1")
		}
		var removed types.Job
		removed, s.sortedWaitingJobs[gpuType] = jobs_util.GetJobsSliceUtil().RemoveJobsSlice(targetIdx, ls)
		if removed != job {
			panic("SchedulerTemplate removeFromSortedWaitingJobs removed != job")
		}
	}
}

func (s *GavelScheduler) hasEmptyGPUQueue() bool {
	for _, queue := range s.cluster.GPUJobQueues() {
		if len(queue.Jobs()) == 0 {
			return true
		}
	}
	return false
}

func (s *GavelScheduler) getEmptyGPUQueues() []types.GPUJobQueue {
	queues := make([]types.GPUJobQueue, 0, len(s.cluster.GPUJobQueues()))
	for _, queue := range s.cluster.GPUJobQueues() {
		if len(queue.Jobs()) == 0 {
			queues = append(queues, queue)
		}
	}
	return queues
}

func (s *GavelScheduler) SetCluster(cluster types.Cluster) {
	s.cluster = cluster
	s.sortedWaitingJobs = make(map[types.GPUType][]types.Job)
	for _, gpuType := range s.cluster.GPUTypes() {
		s.sortedWaitingJobs[gpuType] = make([]types.Job, 0)
	}
}

// 调度开始的地方
func (s *GavelScheduler) OnScheduleEvent(event types.ScheduleEvent) {
	switch e := event.(type) {
	case *types.ScheduleEventJobsArrived:
		{
			for _, jobMeta := range e.JobMetas() {
				// 将新作业插入到一个排序的等待队列中.
				s.insertJob2SortedWaitingJobs(simulator.NewJob(jobMeta.JobName()))
			}

			// 插入所有的新作业后,调用方法进行调度.
			s.DoSchedule()
		}
	case *types.ScheduleEventJobsFinished:
		{
			if !s.hasEmptyGPUQueue() {
				panic("!s.hasEmptyGPUQueue() when some jobs finished.")
			}
			s.DoSchedule()
		}
	}
}

func (s *GavelScheduler) NextActiveScheduleTime() types.Time {
	return types.Time(math.Inf(1))
}

func (s *GavelScheduler) Name() string {
	return fmt.Sprintf("GavelScheduler")
}

func (s *GavelScheduler) Info() interface{} {
	return s.Name()
}

func (s *GavelScheduler) Record() *types.SchedulerRecord {
	return &types.SchedulerRecord{
		DoScheduleRecords: s.DoScheduleCalls,
	}
}
