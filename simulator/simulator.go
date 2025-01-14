package simulator

import (
	"DES-go/schedulers/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
)

type Simulator struct {
	// 模拟器的选项·

	opts *Options
	// 调度器实例，负责调度任务
	scheduler types.Scheduler
	// 集群实例
	cluster *Cluster
	// 日志记录器
	logger *logger
	// 记录已完成的作业
	recordedFinishedJobs []types.Job
}

// 创建并初始化新的模拟器实例  setOpts是可变参数，用于设置模拟器的选项。
func NewSimulator(scheduler types.Scheduler, setOpts ...SetOption) *Simulator {
	opts := defaultOptions

	for _, setOpt := range setOpts {
		setOpt(opts)
	}
	if opts.dataSourceCSVPath != "" {
		initDataSource(opts.dataSourceCSVPath, opts.dataSourceRange)
	}

	logger := NewLogger(opts.logEnabled, opts.logDirPath)
	return &Simulator{
		scheduler:            scheduler,
		opts:                 opts,
		cluster:              NewCluster(opts.gpuType2Count),
		logger:               logger,
		recordedFinishedJobs: make([]types.Job, 0),
	}
}

// 模拟器的核心函数，用于运行模拟并返回模拟结果
func (s *Simulator) Run() *types.Record {
	// 启动集群服务,使其开始运行.
	s.cluster.startServe()
	// 将传入的集群参数,设置为调度器的集群.
	s.scheduler.SetCluster(s.cluster)
	//println(s.scheduler.Name())
	getDataSource().IterBySubmitTime(func(indices []int, metas []types.JobMeta) {
		// 按照提交时间迭代数据源,并检查所有任务的提交时间是否相同.如果不相同则抛出异常.
		submitTime := metas[0].SubmitTime()
		for _, meta := range metas {
			if meta.SubmitTime() != submitTime {
				panic("getDataSource().IterBySubmitTime metas' submit times are different.")
			}
		}
		// 当前时间以及超过任务提交时间,并且超出了运行的最小间隔,触发panic.
		if float64(submitTime-s.cluster.Now()) < -float64(s.opts.minDurationPassInterval) {
			panic(fmt.Sprintf("meta.submitTime() = %v - s.cluster.Now() = %v) >= -float64(s.opts.minDurationPassInterval = %v)", submitTime, s.cluster.Now(), s.opts.minDurationPassInterval))
		}
		// 将模拟器的当前时间推进到任务的提交时间.
		for s.cluster.Now() < submitTime {
			passDuration := submitTime - s.cluster.Now()
			s.passDuration(types.Duration(passDuration), false)
		}
		// 触发任务到达事件,通知调度器有新任务到达.
		s.emitEvent(types.NewScheduleEventJobsArrived(metas))
	})
	// 将模拟器的时间推进到所有任务完成.
	s.passDuration(0, true)
	//println("*******************")
	//println(s.scheduler.Name()

	if s.scheduler.Name() == "FragScheduler" {
		//	println("*******************")
		// 从 JSON 文件中读取 Job 数据
		filePath := "/hydra/data/finished_jobs.json" // 设置 JSON 文件路径
		jobs, err := readJobsFromJSON(filePath)
		if err != nil {
			fmt.Println(err)
		}
		var jobList []types.Job
		for _, job := range jobs {
			gpuThere, err := s.cluster.GetGPUByID(job.GpuID)
			if err != nil {
				fmt.Println(err)
			}
			// 创建 JobImpl 实例，并将 Jobfish 的数据赋值
			startTime := job.FirstExecutionTime
			endTime := job.FinishExecutionTime
			timeRange := &TimeRange{
				start: startTime, // 大写 Start
				end:   endTime,   // 大写 End
			}
			executionRange := &JobExecutionRange{
				gpu:               gpuThere, // 大写 Gpu
				jobName:           job.JobName,
				timeRange:         timeRange,
				completenessRatio: job.RemainingRatio,
			}

			executionDetail := &JobExecutionDetail{
				jobName: job.JobName,
				executionRanges: map[types.GPU][]types.JobExecutionRange{ // 确保类型匹配
					gpuThere: {executionRange},
				},
			}

			var jobImpl = &Job{
				jobName:             job.JobName,
				executionDetail:     executionDetail,
				firstExecutionTime:  job.FirstExecutionTime,
				finishExecutionTime: job.FinishExecutionTime,
				remainingRatio:      job.RemainingRatio,
				isRunning:           job.IsRunning,
			} // 添加到 Job 接口类型的列表
			jobList = append(jobList, jobImpl)
		}

		//if err != nil {
		//	fmt.Printf("读取 JSON 文件出错: %v\n", err)
		//}
		//// 将 Jobfish 列表转换为 Job 接口类型的列表
		//jobList := convertToJobList(jobs)
		s.recordedFinishedJobs = jobList
		//s.cluster.gpus

		return &types.Record{
			SchedulerName:   s.scheduler.Name(),
			SchedulerInfo:   s.scheduler.Info(),
			GPUs:            s.cluster.GPUs(),
			FinishedJobs:    s.recordedFinishedJobs,
			SchedulerRecord: s.scheduler.Record(),
		}
	}

	//for _, job := range s.recordedFinishedJobs {
	//	println(job.JobName())
	//	println(job.Violation())
	//}
	return &types.Record{
		SchedulerName:   s.scheduler.Name(),
		SchedulerInfo:   s.scheduler.Info(),
		GPUs:            s.cluster.GPUs(),
		FinishedJobs:    s.recordedFinishedJobs,
		SchedulerRecord: s.scheduler.Record(),
	}
}

func (s *Simulator) passDuration(duration types.Duration, noMoreNewSubmits bool) {
	currTime := s.cluster.Now()
	targetTime := currTime + types.Time(duration)
	if noMoreNewSubmits {
		targetTime = 1e38
	}
	for currTime < targetTime || noMoreNewSubmits {
		closestTimeToFinishAnyJob := s.cluster.ClosestTimeToFinishAnyJob()
		nextActiveScheduleTime := s.scheduler.NextActiveScheduleTime()
		// 如果调度器将不会进行主动调度，并且将来没有任务要完成，并且指定不会再有新的任务提交了，那么此时认为模拟结束了。
		if math.IsInf(float64(nextActiveScheduleTime), 1) &&
			math.IsInf(float64(closestTimeToFinishAnyJob), 1) &&
			noMoreNewSubmits {
			// All jobs done
			return
		}
		// calculate partial time.
		// in case some jobs finish very closely, use max() to specify a min interval.
		// targetTime - currTime is the upper limit.
		possibleNextEventTime := math.Min(float64(s.scheduler.NextActiveScheduleTime()), float64(closestTimeToFinishAnyJob))
		partialDuration := types.Duration(math.Min(math.Max(possibleNextEventTime, float64(s.opts.minDurationPassInterval)), float64(targetTime-currTime)))
		finishedJobs := make([]*Job, 0)
		//println("*******************")
		finishedJobs = s.cluster.passDuration(partialDuration)
		s.logTimePassed(partialDuration)
		currTime += types.Time(partialDuration)
		for _, job := range finishedJobs {
			s.recordedFinishedJobs = append(s.recordedFinishedJobs, job)
		}
		s.emitEvent(types.NewScheduleEventDurationPassed(partialDuration))
		if len(finishedJobs) > 0 {
			s.emitEvent(types.NewScheduleEventJobsFinished(s.transformJobs(finishedJobs)))
		}

		//finishedJobs = s.cluster.passDuration(partialDuration)
		//// fmt.Printf("finishedJobs len=[%d], all Finished len=[%d]", len(finishedJobs), len(s.recordedFinishedJobs))
		//s.logTimePassed(partialDuration)
		//currTime += types.Time(partialDuration)
		//for _, job := range finishedJobs {
		//	s.recordedFinishedJobs = append(s.recordedFinishedJobs, job)
		//}
		//s.emitEvent(types.NewScheduleEventDurationPassed(partialDuration))
		//if len(finishedJobs) > 0 {
		//	s.emitEvent(types.NewScheduleEventJobsFinished(s.transformJobs(finishedJobs)))
		//}
	}
}

func (s *Simulator) transformJobs(jobs []*Job) []types.Job {
	res := make([]types.Job, 0, len(jobs))
	for _, job := range jobs {
		res = append(res, job)
	}
	return res
}

func (s *Simulator) logTimePassed(duration types.Duration) {
	if s.opts.formatPrintLevel == AllFormatPrint {
		allInfo := fmt.Sprintf("\nTime Passed: %f seconds, finished jobs count: %d. \ncluster info: \n%v.\n", float64(duration), len(s.recordedFinishedJobs), s.cluster)
		log.Printf(allInfo)
	} else if s.opts.formatPrintLevel == ShortMsgPrint {
		log.Printf("\nTime Passed: %f seconds, finished jobs count: %d.\n", float64(duration), len(s.recordedFinishedJobs))
	} else if s.opts.formatPrintLevel == NoPrint {
		// pass.
	}
}

func (s *Simulator) logJobFinished(finishedJobs []types.Job) {
	if s.opts.formatPrintLevel == AllFormatPrint || s.opts.formatPrintLevel == ShortMsgPrint {
		log.Printf("finishedJobs len=[%d], all Finished len=[%d]\n", len(finishedJobs), len(s.recordedFinishedJobs))
	} else if s.opts.formatPrintLevel == NoPrint {
		// pass.
	}
}

//  将ScheduleEvent类型的事件传递给调度器进行处理,调度器通过调用OnScheduleEvent方法处理事件.
func (s *Simulator) emitEvent(event types.ScheduleEvent) {
	s.scheduler.OnScheduleEvent(event)
}

//func (s *Simulator) ReceiveJobs(jobs []schedulers.Jobfish) {
//	//fmt.Println("Received jobs:")
//	finishedJobs := make([]*Job, 0)
//
//	for _, jobre := range jobs {
//
//		//timeRange := &TimeRange{
//		//	start: jobre.FirstExecutionTime,
//		//	end:   jobre.FinishExecutionTime,
//		//}
//		//var timeRange types.TimeRange
//		//timeRange.start = jobre.FirstExecutionTime
//		//timeRange.end = jobre.FinishExecutionTime
//		//timeRange.Runtime() =
//
//		//var jobExecutionRange types.JobExecutionRange
//		//jobExecutionRange.TimeRange() = timeRange
//		//var executionRanges map[types.GPU][]types.JobExecutionRange
//
//		//jobchange := &Job{
//		//	jobName:             jobre.JobName,
//		//	executionDetail:     newJobExecutionDetail(jobre.JobName),
//		//	firstExecutionTime:  jobre.FirstExecutionTime,
//		//	finishExecutionTime: jobre.FinishExecutionTime,
//		//	remainingRatio:      0,
//		//	isRunning:           false,
//		//}
//
//		finishedJobs = append(finishedJobs, jobre)
//	}
//}

// 读取 JSON 文件并解析为 Jobfish 列表
func readJobsFromJSON(filePath string) ([]Jobfish, error) {
	// 打开 JSON 文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 解析 JSON 数据为 Jobfish 列表
	var jobs []Jobfish
	if err := json.Unmarshal(byteValue, &jobs); err != nil {
		return nil, fmt.Errorf("无法解析 JSON 数据: %v", err)
	}

	return jobs, nil
}

// 将 Jobfish 列表转换为 Job 接口类型的列表
//func convertToJobList(jobs []Jobfish) []types.Job {
//	var jobList []types.Job
//	for _, job := range jobs {
//		// 创建 JobImpl 实例，并将 Jobfish 的数据赋值
//		startTime := job.FirstExecutionTime
//		endTime := job.FinishExecutionTime
//		timeRange := &TimeRange{
//			start: startTime, // 大写 Start
//			end:   endTime,   // 大写 End
//		}
//		executionRange := &JobExecutionRange{
//			gpu:               job.GpuID, // 大写 Gpu
//			jobName:           job.JobName,
//			timeRange:         timeRange,
//			completenessRatio: job.RemainingRatio,
//		}
//
//		executionDetail := &JobExecutionDetail{
//			jobName: job.JobName,
//			executionRanges: map[types.GPU][]types.JobExecutionRange{ // 确保类型匹配
//				job.GpuType: {executionRange},
//			},
//		}
//
//		var jobImpl = &Job{
//			jobName:             job.JobName,
//			executionDetail:     executionDetail,
//			firstExecutionTime:  job.FirstExecutionTime,
//			finishExecutionTime: job.FinishExecutionTime,
//			remainingRatio:      job.RemainingRatio,
//			isRunning:           job.IsRunning,
//		} // 添加到 Job 接口类型的列表
//		jobList = append(jobList, jobImpl)
//	}
//
//	return jobList
//}

type Jobfish struct {
	JobName             types.JobName `json:"JobName"`
	GpuID               types.GPUID   `json:"GpuID"`
	GpuType             types.GPUType `json:"GpuType"`
	FirstExecutionTime  types.Time    `json:"FirstExecutionTime"`
	FinishExecutionTime types.Time    `json:"FinishExecutionTime"`
	RemainingRatio      float64       `json:"RemainingRatio"`
	IsRunning           bool          `json:"IsRunning"`
}

func (c *Cluster) GetGPUByID(gpuID types.GPUID) (types.GPU, error) {
	for _, gpuList := range c.gpus {
		for _, gpu := range gpuList {
			if gpu.ID() == gpuID {
				return gpu, nil
			}
		}
	}
	return nil, fmt.Errorf("GPUID %s not found", gpuID)
}
