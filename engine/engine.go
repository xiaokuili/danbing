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
	scheduler.Run(job)

}
