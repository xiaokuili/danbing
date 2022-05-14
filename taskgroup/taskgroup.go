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
	Reader        plugin.ReaderPlugin  `json:"reader,omitempty"`
	Writer        plugin.WriterPlugin  `json:"writer,omitempty"`
	Record        recordchannel.Record ``
	Communication *statistic.Communication
}

type TaskGroup struct {
	ID            int
	Tasks         []*Task
	Communication *statistic.Communication
}

func (t *Task) Run() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func(t *Task) {
		defer wg.Done()
		record := t.Reader.Reader()
		t.Communication.Metric.IncreaseCounter("taskReader")
		t.Record.SetRecord([]byte(record))
	}(t)

	go func(t *Task) {
		defer wg.Done()
		record := t.Record.GetRecord()
		t.Communication.Metric.IncreaseCounter("taskWriter")
		t.Writer.Writer(string(record))
	}(t)
	wg.Wait()
}

func (tg *TaskGroup) Run() {
	for i := 0; i < len(tg.Tasks); i++ {
		communication := statistic.New(i, cons.STAGETASK)
		tg.Communication.Build(communication)
		t := tg.Tasks[i]
		t.Communication = communication
		t.Run()
	}
}
