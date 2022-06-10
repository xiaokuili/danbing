package taskgroup

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/plugin"
	recordchannel "danbing/recordChannel"
	statistic "danbing/statistics"
	"fmt"
	"sync"
)

type Task struct {
	ID            int
	Reader        plugin.ReaderPlugin `json:"reader,omitempty"`
	Writer        plugin.WriterPlugin `json:"writer,omitempty"`
	Record        *recordchannel.Record
	Communication *statistic.Communication
	ReaderParam   *conf.Param
	WriterParam   *conf.Param
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
func (t *Task) SetReaderParam(p *conf.Param) {
	t.ReaderParam = p
}

func (t *Task) SetWriterParam(p *conf.Param) {
	t.WriterParam = p
}

func toString(i interface{}) string {
	r, ok := i.(string)
	if !ok {
		panic(fmt.Errorf("%v cant change to string", i))
	}
	return r
}

func (t *Task) Run() {
	var wg sync.WaitGroup
	fmt.Printf("%d reader begin run\n", t.ID)

	wg.Add(2)
	go func(t *Task) {
		defer wg.Done()

		record := t.Reader.Reader()
		// 在这里基于字段进行收集
		if len(record) > 0 {
			columns := t.ReaderParam.Query.Columns

			r := record[len(record)-1]
			for i := 0; i < len(columns); i++ {
				field := columns[i]
				if field.CollectField {
					name := field.Name
					t.Communication.Metric.SetMessage(name, toString(r[name]))

				}
			}

		}
		t.Record.PutRecord(record)
	}(t)

	go func(t *Task) {
		defer wg.Done()
		record := t.Record.GetRecord()
		// 在这里基于字段进行收集
		t.Writer.Writer(record)
	}(t)
	fmt.Printf("%d writer ending run\n", t.ID)

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
