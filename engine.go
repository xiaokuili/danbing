package main

import (
	_ "danbing/myplugin"
	"danbing/plugin"
	"danbing/task"
	"fmt"
)

type Job struct {
	Param     []*task.Param `json:"job"`
	Speed     *task.Speed   `json:"speed"`
	Task      *Task         `json:"task,omitempty"`
	Tasks     []*Task       `json:"tasks,omitempty"`
	TaskGroup [][]*Task     `json:"taskgroup,omitempty"`
}

type Task struct {
	Reader plugin.ReaderPlugin `json:"reader,omitempty"`
	Writer plugin.WriterPlugin `json:"writer,omitempty"`
}

// 基于配置文件生成job
func New() *Job {
	// struct
	// type struct 是声明
	// 这里是实例化

	j := &Job{
		Param: []*task.Param{},
		Speed: &task.Speed{},
		Task: &Task{
			Reader: nil,
			Writer: nil,
		},
		Tasks:     []*Task{},
		TaskGroup: [][]*Task{},
	}

	return j
}

func Temple() *Job {

	job := New()
	reader := &task.Param{
		Connect: &task.Connect{},
		Query:   &task.Query{},
		Type:    task.READER,
		Name:    "streamreader",
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
		Channel:          10,
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
	j.TaskGroup = make([][]*Task, j.Speed.Thread)
	for i := 0; i < len(param); i++ {
		p := param[i]
		if p.Type == task.READER {
			j.Task.Reader = plugin.ReaderPlugins[p.Name]
		}
		if p.Type == task.WRITER {
			j.Task.Writer = plugin.WriterPlugins[p.Name]
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
	fmt.Println(Wtask)
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

	for i := 0; i < len(tasks); i++ {
		t := tasks[i]
		gid := i % threat
		j.TaskGroup[gid] = append(j.TaskGroup[gid], t)
	}
	fmt.Println("group task")
}

func (j *Job) Scheduler() {
	group := j.TaskGroup
	for i := 0; i < len(group); i++ {
		tasks := group[i]
		for j := 0; j < len(tasks); j++ {
			t := tasks[j]
			fmt.Printf("%d-%d is running...\n", i, j)
			t.Reader.Reader()
			t.Writer.Writer()
		}
	}
	fmt.Println("scheduler")
}

func main() {

	j := Temple()
	j.Init()
	j.Split()
	j.GroupTask()
	j.Scheduler()

}
