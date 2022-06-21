package engine

import (
	"danbing/job"
	_ "danbing/myplugin"
	"danbing/scheduler"
)

func Engine(job *job.Job) {

	scheduler.Run(job)
}

func EngineReport(job *job.Job) {
	job.SetBeginTime("2022-02-03")
	job.SetEndTime("2022-06-04")
	job.Refresh()
	scheduler.Run(job)

}
