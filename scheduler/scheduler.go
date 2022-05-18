package scheduler

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/job"
	"danbing/plugin"
	statistic "danbing/statistics"
	"danbing/taskgroup"
	"sync"
)

type Scheduler struct {
	Communication *statistic.Communication
	TaskGroup     []*taskgroup.TaskGroup
	Job           *job.Job
	Table         string
	wg            sync.WaitGroup
}

func New(job *job.Job) *Scheduler {
	communication := statistic.SingletonNew()
	return &Scheduler{
		Communication: communication,
		TaskGroup:     []*taskgroup.TaskGroup{},
		Job:           job,
		Table:         job.Table,
	}
}

func (s *Scheduler) Persister() {}

func InitReader(p *conf.Param) plugin.ReaderPlugin {
	if p.Type != cons.PLUGINREADER {
		panic("please check plugin is reader")
	}
	r := plugin.ReaderPlugins[p.Name]
	r.Init(p.Query, p.Connect)
	return r
}

func InitWriter(p *conf.Param) plugin.WriterPlugin {
	if p.Type != cons.PLUGINWRITER {
		panic("please check plugin is writer")
	}
	w := plugin.WriterPlugins[p.Name]
	w.Init(p.Query, p.Connect)
	return w
}

func (s *Scheduler) Init() *taskgroup.Task {
	param := s.Job.Param
	genesis := &taskgroup.Task{}
	for i := 0; i < len(param); i++ {
		p := param[i]
		switch p.Type {
		case cons.PLUGINREADER:
			genesis.Reader = InitReader(p)
		case cons.PLUGINWRITER:
			genesis.Writer = InitWriter(p)
		}
	}
	return genesis
}

func (s *Scheduler) Split(genesis *taskgroup.Task) []*taskgroup.Task {
	if !genesis.CheckTask() {
		panic("please init task...")
	}
	reader := genesis.Reader.Split(s.Job.Speed.Channel)
	writer := genesis.Writer.Split(len(reader))
	tasks := s.MergeRWTask(reader, writer)
	return tasks
}

func (s *Scheduler) MergeRWTask(r []plugin.ReaderPlugin, w []plugin.WriterPlugin) []*taskgroup.Task {

	tasks := make([]*taskgroup.Task, 0)
	for i := 0; i < len(r); i++ {
		t := taskgroup.NewTask(i, s.Table, r[i], w[i])
		rp, err := conf.ReaderParam(s.Job.Param)
		if err != nil {
			panic("")
		}
		t.SetReaderParam(rp)
		wp, err := conf.WriterParam(s.Job.Param)
		if err != nil {
			panic("")
		}
		t.SetWriterParam(wp)
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *Scheduler) InitTaskGroup() []*taskgroup.TaskGroup {
	threat := s.Job.Speed.Thread

	group := make([]*taskgroup.TaskGroup, threat)
	for i := 0; i < threat; i++ {
		tg := taskgroup.New(i, s.Table)
		group[i] = tg
		s.registerTaskgroupCommunication(tg)
	}
	return group
}

func (s *Scheduler) PutTaskToGroup(tasks []*taskgroup.Task, bucket []*taskgroup.TaskGroup) {
	for i := 0; i < len(tasks); i++ {
		t := tasks[i]
		gid := i % len(bucket)
		bucket[gid].PutTask(t)

	}
	s.TaskGroup = bucket
}

func (s *Scheduler) GroupTasks(tasks []*taskgroup.Task) {
	group := s.InitTaskGroup()
	s.PutTaskToGroup(tasks, group)
}

func (s *Scheduler) Scheduler() {
	group := s.TaskGroup

	s.wg.Add(len(group))
	for i := 0; i < len(group); i++ {
		go func(group *taskgroup.TaskGroup) {
			defer s.wg.Done()
			group.Run()
		}(group[i])
	}
	s.wg.Wait()
}

func (s *Scheduler) Report() {
	s.Communication.Report()
}

func (s *Scheduler) registerTaskgroupCommunication(tg *taskgroup.TaskGroup) {
	if !tg.Check() {
		panic("please check taskgroup")
	}

	s.Communication.Build(tg.Communication)
}
