func ReceiveFinishedJobs(jobs []Jobfish) {
	// 接收到的 jobs 赋值给 allJob
	allJob := jobs

	// 如果需要，可以进一步操作 allJob
	fmt.Println("Received finished jobs:", allJob)
}