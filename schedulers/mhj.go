package schedulers

//
//import (
//	"encoding/csv"
//	"fmt"
//	"log"
//	"math"
//	"os"
//	"sort"
//	"strconv"
//)
//
//// FragScheduler
//// 为保证公平，采取作业先到先服务
//// 作业逐个调度
//// 对于每个作业，计算调度到不同GPU上的fgd*JCT的值
//// 将作业调度到fgd*JCT值最小的GPU上。fgd*JCT的计算方式见frag_jct.py文件。
//
//type Task struct {
//	JobName           string
//	NormJobSubmitTime int
//	GPUMilli          int
//	DDL               int
//	A100              float64
//	GTX2080Ti         float64
//	V100              float64
//}
//
//type GPU struct {
//	ID       int
//	Type     string
//	Capacity int
//	Used     int
//	Tasks    []Task
//}
//
//func FragScheduler(taskFile string, popularityFile string, clusterConfig map[string]int) {
//	tasks := readTasksFromCSV(taskFile)
//	popularity := readPopularityFromCSV(popularityFile)
//	gpus := initializeGPUs(clusterConfig)
//
//	for _, task := range tasks {
//		assignTaskToGPU(task, gpus, popularity)
//	}
//}
//
//func readTasksFromCSV(filename string) []Task {
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
//	var tasks []Task
//	for _, record := range records[1:] { // Skip header
//		normJobSubmitTime, _ := strconv.Atoi(record[1])
//		gpuMilli, _ := strconv.Atoi(record[2])
//		ddl, _ := strconv.Atoi(record[3])
//		a100, _ := strconv.ParseFloat(record[4], 64)
//		gtx2080Ti, _ := strconv.ParseFloat(record[5], 64)
//		v100, _ := strconv.ParseFloat(record[6], 64)
//
//		task := Task{
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
//	for _, record := range records[1:] { // Skip header
//		gpuMilli, _ := strconv.Atoi(record[0])
//		ratio, _ := strconv.ParseFloat(record[1], 64)
//		popularity[gpuMilli] = ratio
//	}
//	return popularity
//}
//
//func initializeGPUs(clusterConfig map[string]int) []GPU {
//	var gpus []GPU
//	id := 0
//	for gpuType, count := range clusterConfig {
//		for i := 0; i < count; i++ {
//			gpus = append(gpus, GPU{ID: id, Type: gpuType, Capacity: 1000, Used: 0})
//			id++
//		}
//	}
//	return gpus
//}
//
//func assignTaskToGPU(task Task, gpus []GPU, popularity map[int]float64) {
//	sort.Slice(gpus, func(i, j int) bool {
//		return calculateMetric(task, gpus[i], popularity) < calculateMetric(task, gpus[j], popularity)
//	})
//
//	for i := range gpus {
//		if gpus[i].Used+task.GPUMilli <= gpus[i].Capacity {
//			waitTime := calculateWaitTime(task, gpus[i])
//			runtime := getRuntime(task, gpus[i].Type)
//			gpus[i].Used += task.GPUMilli
//			gpus[i].Tasks = append(gpus[i].Tasks, task)
//			fmt.Printf(`"job_name": "%s", "selected_gpu_id": %d, "selected_gpu_type": "%s", "runtime_on_gpu": %.1f, "start_time": %d, "end_time": %.1f\n`,
//				task.JobName, gpus[i].ID, gpus[i].Type, runtime, task.NormJobSubmitTime+waitTime, float64(task.NormJobSubmitTime+waitTime)+runtime)
//			break
//		}
//	}
//}
//
//func calculateMetric(task Task, gpu GPU, popularity map[int]float64) float64 {
//	fragSize := (1000 - task.GPUMilli) * calculatePopularitySum(task, gpu, popularity)
//	jct := calculateWaitTime(task, gpu) + getRuntime(task, gpu.Type)
//	return float64(gpu.ID) * fragSize * jct
//}
//
//func calculatePopularitySum(task Task, gpu GPU, popularity map[int]float64) float64 {
//	sum := 0.0
//	for _, t := range gpu.Tasks {
//		if t.GPUMilli < (1000 - task.GPUMilli) {
//			sum += popularity[t.GPUMilli]
//		}
//	}
//	return sum
//}
//
//func calculateWaitTime(task Task, gpu GPU) int {
//	waitTime := 0
//	for _, t := range gpu.Tasks {
//		if gpu.Used+task.GPUMilli > gpu.Capacity {
//			waitTime += int(getRuntime(t, gpu.Type))
//		}
//	}
//	return waitTime
//}
//
//func getRuntime(task Task, gpuType string) float64 {
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
