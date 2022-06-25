package scheduler

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/job"
	"danbing/log"
	"danbing/plugin"
	statistic "danbing/statistics"
	"danbing/taskgroup"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Scheduler struct {
	Communication *statistic.Communication
	TaskGroup     []*taskgroup.TaskGroup
	Job           *job.Job
	Table         string
	Total         int
	wg            sync.WaitGroup
	logger        log.Logger
	Reader        plugin.ReaderPlugin
	Writer        plugin.WriterPlugin
}

func New(job *job.Job, logger log.Logger) *Scheduler {
	communication := statistic.SingletonNew()
	return &Scheduler{
		Communication: communication,
		TaskGroup:     []*taskgroup.TaskGroup{},
		Job:           job,
		Table:         job.Table,
		logger:        logger,
	}
}

func (s *Scheduler) Init() {
	param := s.Job.Param
	var reader plugin.ReaderPlugin
	var writer plugin.WriterPlugin
	for i := 0; i < len(param); i++ {
		p := param[i]
		switch p.Type {
		case cons.PLUGINREADER:
			reader = initReader(p)
			s.Total = reader.Init(p.Query, p.Connect)
		case cons.PLUGINWRITER:
			writer = initWriter(p)
			writer.Init(p.Query, p.Connect)
		}
	}
	s.Reader = reader
	s.Writer = writer

}

func (s *Scheduler) Split() []*taskgroup.Task {

	reader := s.Reader.Split(s.Job.Speed.NumPerTask)
	writer := s.Writer.Split(len(reader))
	tasks := s.mergeRW(reader, writer)
	return tasks
}

func (s *Scheduler) Close() {
	s.Reader.Close()
	s.Writer.Close()
}

func (s *Scheduler) AssignTasks(tasks []*taskgroup.Task) {
	group := s.initTaskGroup()
	s.assignTask(tasks, group)
}

func (s *Scheduler) Scheduler() {
	t := time.Now()
	group := s.TaskGroup
	s.wg.Add(len(group))
	for i := 0; i < len(group); i++ {
		go func(group *taskgroup.TaskGroup) {
			defer s.wg.Done()
			group.Run()
		}(group[i])
	}
	s.wg.Wait()

	consume := fmt.Sprintf("%v", time.Since(t))
	s.logger.Debug("scheduler  end", []interface{}{"consume", consume}...)
}

func Run(job *job.Job, level string) {
	run(job, log.MustNewDefaultLogger(log.LogFormatJSON, level))
}

func run(job *job.Job, logger log.Logger) {
	t := time.Now()
	s := New(job, logger)
	s.Init()
	tasks := s.Split()
	s.AssignTasks(tasks)
	go func(time.Time) {
		for {
			s.Report(t)
			time.Sleep(time.Second * 1)
		}
	}(t)

	s.Scheduler()
	s.Report(t)
	s.Close()
}

func (s *Scheduler) Report(t time.Time) {
	deal := s.Communication.Report()

	total := s.Total
	d, err := strconv.Atoi(deal)
	if err != nil {
		d = 0
	}
	percent := percentage(d, total)
	et := 0
	if percent == 0 {
		et = s.Total / 20000
	} else {
		et = (100 - percent) * int(time.Since(t).Seconds()) / percent
	}

	s.logger.Info("", "deal", deal, "total", s.Total, "percent", percent, "estimate", et)
}

func initReader(p *conf.Param) plugin.ReaderPlugin {
	if p.Type != cons.PLUGINREADER {
		panic("please check plugin is reader")
	}
	r := plugin.ReaderPlugins[p.Name]

	return r
}

func initWriter(p *conf.Param) plugin.WriterPlugin {
	if p.Type != cons.PLUGINWRITER {
		panic("please check plugin is writer")
	}
	w := plugin.WriterPlugins[p.Name]

	return w
}

func (s *Scheduler) mergeRW(r []plugin.ReaderPlugin, w []plugin.WriterPlugin) []*taskgroup.Task {
	tasks := make([]*taskgroup.Task, 0)
	for i := 0; i < len(r); i++ {
		t := taskgroup.NewTask(i, s.Table, r[i], w[i])
		rp, err := conf.Reader(s.Job.Param)
		if err != nil {
			panic("")
		}
		t.SetReaderParam(rp)
		wp, err := conf.Writer(s.Job.Param)
		if err != nil {
			panic("")
		}
		t.SetWriterParam(wp)
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *Scheduler) initTaskGroup() []*taskgroup.TaskGroup {
	threat := s.Job.Speed.Thread

	group := make([]*taskgroup.TaskGroup, threat)
	for i := 0; i < threat; i++ {
		tg := taskgroup.New(i, s.Table)
		tg.SetLogger(s.logger)
		group[i] = tg
		s.registerTaskgroupCommunication(tg)
	}
	return group
}

func (s *Scheduler) assignTask(tasks []*taskgroup.Task, bucket []*taskgroup.TaskGroup) {
	for i := 0; i < len(tasks); i++ {
		t := tasks[i]
		gid := i % len(bucket)
		bucket[gid].PutTask(t)

	}
	s.TaskGroup = bucket
}

func (s *Scheduler) registerTaskgroupCommunication(tg *taskgroup.TaskGroup) {
	if !tg.Check() {
		panic("please check taskgroup")
	}

	s.Communication.Build(tg.Communication)
}

func percentage(part, total int) int {

	delta := (float64(part) / float64(total)) * 100
	return int(delta)
}
