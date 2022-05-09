package main

import (
	_ "danbing/myplugin"
	"danbing/plugin"
	recordchannel "danbing/recordChannel"
	"danbing/task"
	"fmt"
	"sync"
)

type Job struct {
	Param     []*task.Param `json:"job"`
	Speed     *task.Speed   `json:"speed"`
	Task      *Task         `json:"task,omitempty"`
	Tasks     []*Task       `json:"tasks,omitempty"`
	TaskGroup []*TaskGroup  `json:"taskgroup,omitempty"`
}

type Task struct {
	Reader plugin.ReaderPlugin  `json:"reader,omitempty"`
	Writer plugin.WriterPlugin  `json:"writer,omitempty"`
	Record recordchannel.Record ``
}

type TaskGroup struct {
	Id    int
	Tasks []*Task
}

func (t *Task) Run() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func(t *Task) {
		defer wg.Done()
		record := t.Reader.Reader()
		t.Record.SetRecord([]byte(record))
	}(t)

	go func(t *Task) {
		defer wg.Done()
		record := t.Record.GetRecord()
		t.Writer.Writer(string(record))
	}(t)
	wg.Wait()
}

func (tg *TaskGroup) Run() {
	for i := 0; i < len(tg.Tasks); i++ {
		tg.Tasks[i].Run()
	}
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
		Task:      &Task{},
		Tasks:     []*Task{},
		TaskGroup: []*TaskGroup{},
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
		Channel:          100,
		Thread:           10,
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

	tasks := make([]*Task, 0)
	for i := 0; i < len(r); i++ {
		t := &Task{Reader: r[i], Writer: w[i]}
		tasks = append(tasks, t)
	}
	j.Tasks = tasks
}

func (j *Job) GroupTask() {
	threat := j.Speed.Thread
	tasks := j.Tasks
	group := make([]*TaskGroup, threat)
	for i := 0; i < threat; i++ {

		group[i] = &TaskGroup{
			Id:    i,
			Tasks: []*Task{},
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
	group := j.TaskGroup
	var wg sync.WaitGroup
	wg.Add(len(group))
	for i := 0; i < len(group); i++ {
		go func(group *TaskGroup) {
			defer wg.Done()
			group.Run()
		}(group[i])
	}
	wg.Wait()
	fmt.Println("scheduler")
}

func main() {

	j := Temple()
	j.Init()
	j.Split()
	j.GroupTask()
	j.Scheduler()

}
