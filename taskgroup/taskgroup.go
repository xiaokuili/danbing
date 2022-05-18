package taskgroup

import (
	"danbing/cons"
	"danbing/plugin"
	recordchannel "danbing/recordChannel"
	statistic "danbing/statistics"
	"sync"
)

type Task struct {
	ID            int
	Reader        plugin.ReaderPlugin   `json:"reader,omitempty"`
	Writer        plugin.WriterPlugin   `json:"writer,omitempty"`
	Record        *recordchannel.Record ``
	Communication *statistic.Communication
}

func (t *Task) CheckTask() bool {
	return t.Reader != nil && t.Writer != nil
}

func NewTask(id int, table string, reader plugin.ReaderPlugin, writer plugin.WriterPlugin) *Task {
	communication := statistic.New(id, cons.STAGETASK, table)
	return &Task{
		ID:            id,
		Reader:        reader,
		Writer:        writer,
		Record:        recordchannel.New(communication),
		Communication: communication,
	}
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

type TaskGroup struct {
	ID            int
	Tasks         []*Task
	Communication *statistic.Communication
	Table         string
}

func New(id int, table string) *TaskGroup {
	communication := statistic.New(id, cons.STAGETASKGROUP, table)
	tg := &TaskGroup{
		ID:            id,
		Tasks:         []*Task{},
		Communication: communication,
		Table:         table,
	}
	return tg
}
func (tg *TaskGroup) Check() bool {
	return tg.Communication != nil
}

func (tg *TaskGroup) PutTask(task *Task) {
	tg.Tasks = append(tg.Tasks, task)
	tg.registerTaskCommunication(task)
}

func (tg *TaskGroup) registerTaskCommunication(task *Task) {

	tg.Communication.Build(task.Communication)
}

func (tg *TaskGroup) Run() {
	for i := 0; i < len(tg.Tasks); i++ {
		t := tg.Tasks[i]
		t.Run()
	}
}
