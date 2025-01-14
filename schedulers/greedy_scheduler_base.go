package schedulers

import (
	"DES-go/schedulers/jobs_util"
	"DES-go/schedulers/types"
	"DES-go/simulator"
	"fmt"
	"math"
	"time"
)

// 调度逻辑是,对于任务,按照在不同的GPU上的运行时间进行排序,如当前场景有三种GPU,那么就有三个切片,每个切片是在一种GPU上对运行时间的排序.当有空的GPU时,首先根据GPU类型,选择在当前GPU上运行时间最短的任务,然后将任务从等待队列中删除.

type SchedulerTemplate struct {
	cluster types.Cluster

	// 等待队列中的所有任务，其分别在每种类型的GPU上，按照RemainingDuration排序。types.GPUType是键值
	sortedWaitingJobs map[types.GPUType][]types.Job

	DoScheduleCalls []*types.DoScheduleCallRecord
	impl            GreedyScheduler
}

type GreedyScheduler interface {
	types.Scheduler
	insertJob2SortedWaitingJobs(job types.Job)
	pickTarget(emptyQueues []types.GPUJobQueue) (types.Job, types.GPUJobQueue)
}

// 创建并初始化一个新的SchedulerTemplate实例。
func NewGreedySchedulerTemplate() *SchedulerTemplate {
	// 返回一个空的DoSchedulerCalls切片，用于存储调度调用记录。
	return &SchedulerTemplate{
		DoScheduleCalls: make([]*types.DoScheduleCallRecord, 0),
	}
}

func (s *SchedulerTemplate) DoSchedule() {
	start := time.Now()
	s.doSchedule()
	duration := time.Since(start)
	s.DoScheduleCalls = append(s.DoScheduleCalls, &types.DoScheduleCallRecord{Duration: duration})
}

func (s *SchedulerTemplate) doSchedule() {
	for s.hasWaitingJob() && s.hasEmptyGPUQueue() {
		// 从waitingJobs中，在全部可能的EmptyGPUSlot上，挑选一个速度最快的。
		emptyQueues := s.getEmptyGPUQueues()
		targetJob, targetQueue := s.impl.pickTarget(emptyQueues)
		// 遍历全部的waiting job，按照gpu type进行分类，在每个waitingJobs上找首个job（即在这个类型上剩余执行时间最短的任务）
		// 遍历结束后，找到一个速度最快的任务。
		if targetJob == nil || targetQueue == nil {
			panic("SchedulerTemplate targetJob == nil || targetQueue == nil")
		}
		s.removeFromSortedWaitingJobs(targetJob)
		targetQueue.SetJobs(targetJob)
	}
}

func (s *SchedulerTemplate) pickTarget(emptyQueues []types.GPUJobQueue) (types.Job, types.GPUJobQueue) {
	panic("SchedulerTemplate pickTarget cannot be called.")
}

func (s *SchedulerTemplate) hasWaitingJob() bool {
	for _, l := range s.sortedWaitingJobs {
		if len(l) > 0 {
			return true
		}
	}
	return false
}

func (s *SchedulerTemplate) insertJob2SortedWaitingJobs(job types.Job) {
	panic("SchedulerTemplate insertJob2SortedWaitingJobs Cannot be called.")
}

func (s *SchedulerTemplate) removeFromSortedWaitingJobs(job types.Job) {
	for _, gpuType := range s.cluster.GPUTypes() {
		ls := s.sortedWaitingJobs[gpuType]
		targetIdx := -1
		for idx, jobInWaitingList := range ls {
			if jobInWaitingList.JobName() == job.JobName() {
				targetIdx = idx
			}
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

func (s *SchedulerTemplate) hasEmptyGPUQueue() bool {
	for _, queue := range s.cluster.GPUJobQueues() {
		if len(queue.Jobs()) == 0 {
			return true
		}
	}
	return false
}

func (s *SchedulerTemplate) getEmptyGPUQueues() []types.GPUJobQueue {
	queues := make([]types.GPUJobQueue, 0, len(s.cluster.GPUJobQueues()))
	for _, queue := range s.cluster.GPUJobQueues() {
		if len(queue.Jobs()) == 0 {
			queues = append(queues, queue)
		}
	}
	return queues
}

func (s *SchedulerTemplate) SetCluster(cluster types.Cluster) {
	s.cluster = cluster
	s.sortedWaitingJobs = make(map[types.GPUType][]types.Job)
	for _, gpuType := range s.cluster.GPUTypes() {
		s.sortedWaitingJobs[gpuType] = make([]types.Job, 0)
	}
}

func (s *SchedulerTemplate) OnScheduleEvent(event types.ScheduleEvent) {
	switch e := event.(type) {
	case *types.ScheduleEventJobsArrived:
		{
			for _, jobMeta := range e.JobMetas() {
				s.impl.insertJob2SortedWaitingJobs(simulator.NewJob(jobMeta.JobName()))
			}
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

func (s *SchedulerTemplate) NextActiveScheduleTime() types.Time {
	return types.Time(math.Inf(1))
}

func (s *SchedulerTemplate) Name() string {
	return fmt.Sprintf("SchedulerTemplate")
}

func (s *SchedulerTemplate) Info() interface{} {
	return s.Name()
}

func (s *SchedulerTemplate) Record() *types.SchedulerRecord {
	return &types.SchedulerRecord{
		DoScheduleRecords: s.DoScheduleCalls,
	}
}
