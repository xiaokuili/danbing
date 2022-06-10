package engine

import (
	"danbing/job"
	_ "danbing/myplugin"
	"danbing/scheduler"
	"time"
)

func Engine(job *job.Job) {
	scheduelr := scheduler.New(job)
	genesis := scheduelr.Init()
	tasks := scheduelr.Split(genesis)
	scheduelr.GroupTasks(tasks)
	scheduelr.Scheduler()
	// scheduelr.Report()
}

func EngineReport(job *job.Job) {
	scheduelr := scheduler.New(job)
	genesis := scheduelr.Init()
	tasks := scheduelr.Split(genesis)
	scheduelr.GroupTasks(tasks)
	go func() {
		for {
			time.Sleep(time.Second * 1)
			scheduelr.Report()
		}
	}()

	scheduelr.Scheduler()
	scheduelr.Report()

}
