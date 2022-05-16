package main

import (
	"danbing/cons"
	_ "danbing/myplugin"
	"danbing/plugin"
	recordchannel "danbing/recordChannel"
	statistic "danbing/statistics"
	"danbing/task"
	"danbing/taskgroup"
	"fmt"
	"sync"
)

type Job struct {
	Param     []*task.Param          `json:"job"`
	Speed     *task.Speed            `json:"speed"`
	Task      *taskgroup.Task        `json:"task,omitempty"`
	Tasks     []*taskgroup.Task      `json:"tasks,omitempty"`
	TaskGroup []*taskgroup.TaskGroup `json:"taskgroup,omitempty"`
}

// 基于配置文件生成job
func New() *Job {
	// struct
	// type struct 是声明
	// 这里是实例化
	// group := make([]*TaskGroup, 0)
	j := &Job{
		Param:     []*task.Param{},
		Speed:     &task.Speed{},
		Task:      &taskgroup.Task{},
		Tasks:     []*taskgroup.Task{},
		TaskGroup: []*taskgroup.TaskGroup{},
	}
	return j
}

func Temple() *Job {

	job := New()
	reader := &task.Param{
		Connect: &task.Connect{},
		Query: &task.Query{
			SQL: "hello world",
		},
		Type: task.READER,
		Name: "streamreader",
	}
	job.Param = append(job.Param, reader)

	writer := &task.Param{
		Connect: &task.Connect{},
		Query:   &task.Query{},
		Type:    task.WRITER,
		Name:    "streamwriter",
	}
	job.Param = append(job.Param, writer)

	job.Speed = &task.Speed{
		Byte:             0,
		BytePerChannel:   0,
		Record:           0,
		RecordPerChannel: 0,
		Channel:          100, // task 数量
		Thread:           10,  // threat group数量
	}
	return job
}

func PG2ESTemple() *Job {

	job := New()
	reader := &task.Param{
		Connect: &task.Connect{
			URL:      "127.0.0.1:5432",
			Username: "postgres",
			Password: "postgres",
			Database: "postgres",
		},
		Query: &task.Query{
			SQL: "select * from danbing",
		},
		Type: task.READER,
		Name: "pgsqlreader",
	}
	job.Param = append(job.Param, reader)

	writer := &task.Param{
		Connect: &task.Connect{},
		Query:   &task.Query{},
		Type:    task.WRITER,
		Name:    "streamwriter",
	}
	job.Param = append(job.Param, writer)

	job.Speed = &task.Speed{
		Byte:             0,
		BytePerChannel:   0,
		Record:           0,
		RecordPerChannel: 0,
		Channel:          100, // task 数量
		Thread:           10,  // threat group数量
	}
	return job
}

// 不存在返回false
func (j *Job) CheckTaskExist() bool {
	t := j.Task
	if t != nil {
		return !(t.Reader == nil || t.Writer == nil)
	}
	return false
}

func (j *Job) Init() {
	param := j.Param
	for i := 0; i < len(param); i++ {
		p := param[i]
		if p.Type == task.READER {
			r := plugin.ReaderPlugins[p.Name]
			r.Init(p.Query, p.Connect)
			j.Task.Reader = r
		}
		if p.Type == task.WRITER {
			w := plugin.WriterPlugins[p.Name]
			w.Init(p.Query, p.Connect)
			j.Task.Writer = w
		}

	}

}

func (j *Job) Split() {
	if !j.CheckTaskExist() {
		return
	}
	t := j.Task

	Rtask := t.Reader.Split(j.Speed.Channel)
	Wtask := t.Writer.Split(len(Rtask))
	j.MergeRWTask(Rtask, Wtask)

}

func (j *Job) MergeRWTask(r []plugin.ReaderPlugin, w []plugin.WriterPlugin) {

	tasks := make([]*taskgroup.Task, 0)
	for i := 0; i < len(r); i++ {
		t := &taskgroup.Task{Reader: r[i], Writer: w[i]}
		tasks = append(tasks, t)
	}
	j.Tasks = tasks
}

func (j *Job) GroupTask() {
	threat := j.Speed.Thread
	tasks := j.Tasks
	group := make([]*taskgroup.TaskGroup, threat)
	for i := 0; i < threat; i++ {

		group[i] = &taskgroup.TaskGroup{
			ID:    i,
			Tasks: []*taskgroup.Task{},
		}
	}

	for i := 0; i < len(tasks); i++ {
		t := tasks[i]
		t.Record = recordchannel.Record{
			C: make(chan []byte),
		}
		gid := i % threat

		group[gid].Tasks = append(group[gid].Tasks, t)

	}
	j.TaskGroup = group
	fmt.Println("group task")
}

func (j *Job) Scheduler() {
	communication := statistic.SingletonNew()

	group := j.TaskGroup
	var wg sync.WaitGroup
	wg.Add(len(group))
	for i := 0; i < len(group); i++ {
		gtask := group[i]
		tgCommunication := statistic.New(gtask.ID, cons.STAGETASKGROUP)
		tgCommunication.Metric.IncreaseCounter("taskgroup_count")
		gtask.Communication = tgCommunication
		communication.Build(tgCommunication)

		go func(group *taskgroup.TaskGroup) {
			defer wg.Done()
			group.Run()
		}(group[i])
	}
	wg.Wait()
	communication.Report()

}

func main() {

	j := PG2ESTemple()
	j.Init()
	j.Split()
	j.GroupTask()
	j.Scheduler()

}
