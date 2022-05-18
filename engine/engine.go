package engine

import (
	"danbing/job"
	_ "danbing/myplugin"
	"danbing/scheduler"
)

func Engine(job *job.Job) {
	scheduelr := scheduler.New(job)
	genesis := scheduelr.Init()
	tasks := scheduelr.Split(genesis)
	scheduelr.GroupTasks(tasks)
	scheduelr.Scheduler()
	// scheduelr.Report()
}
