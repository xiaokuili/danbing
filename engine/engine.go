package engine

import (
	"danbing/db"
	"danbing/job"
	_ "danbing/myplugin"
	"danbing/scheduler"
	"time"
)

const (
	BEGIN      = "2022-01-01 01:01:01"
	TimeFormat = "2006-01-02 15:04:05"
)

func Begin(j *job.Job) string {
	info := db.SearchLast(j.Table)
	if info != nil {
		t := time.Unix(info.Uptime, 0)
		return t.Format(TimeFormat)
	}
	return BEGIN
}

func Ending() string {
	return time.Now().Format(TimeFormat)
}

func Run(job *job.Job) {
	b := Begin(job)
	e := Ending()

	uptime, err := time.Parse(e, TimeFormat)
	if err != nil {

	}
	job.SetBeginTime(b)
	job.SetEndTime(e)
	job.Refresh()

	scheduler.Run(job)
	if job.Stream {
		info := db.Info{
			Name:   job.Table,
			Batch:  job.Batch,
			Uptime: uptime.Unix(),
		}
		info.Insert()
	}

}
