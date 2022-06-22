package taskgroup

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/log"
	"danbing/plugin"
	recordchannel "danbing/recordChannel"
	statistic "danbing/statistics"
	"fmt"
	"sync"
	"time"
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
	return t.Reader != nil && t.Writer != nil && t.Communication != nil
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

func str(i interface{}) string {
	r, ok := i.(string)
	if !ok {
		panic(fmt.Errorf("%v cant change to string", i))
	}
	return r
}

func (t *Task) Run() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func(t *Task) {
		defer wg.Done()

		record := t.Reader.Reader()

		t.Record.PutRecord(record)
	}(t)

	go func(t *Task) {
		defer wg.Done()
		record := t.Record.GetRecord()

		t.Writer.Writer(record)

	}(t)

	wg.Wait()
}

type TaskGroup struct {
	ID            int
	Tasks         []*Task
	Communication *statistic.Communication
	Table         string
	logger        log.Logger
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

func (tg *TaskGroup) SetLogger(logger log.Logger) {
	tg.logger = logger
}

func (tg *TaskGroup) Check() bool {
	return tg.Communication != nil
}

func (tg *TaskGroup) PutTask(task *Task) {
	tg.Tasks = append(tg.Tasks, task)
	tg.registerTaskCommunication(task)
}

func (tg *TaskGroup) Run() {
	t := time.Now()
	for i := 0; i < len(tg.Tasks); i++ {
		t := tg.Tasks[i]
		t.Run()
	}
	consume := fmt.Sprintf("%v", time.Since(t))
	msg := fmt.Sprintf("taskgroup [%d] end", tg.ID)
	tg.logger.Debug(msg, []interface{}{"consume", consume}...)

}

func (tg *TaskGroup) registerTaskCommunication(task *Task) {
	if task.CheckTask() {
		tg.Communication.Build(task.Communication)
	}
}
