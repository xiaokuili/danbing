package engine_test

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/engine"
	"danbing/job"
)

func streamJob() *job.Job {
	job := job.New("danbing")
	c := make([]*conf.Column, 0)
	c = append(c, &conf.Column{
		Name:         "out",
		CollectField: true, // 收集这个字段的最后一条数据
	})
	reader := &conf.Param{
		Connect: &conf.Connect{},
		Query: &conf.Query{
			SQL:     "hello world",
			Columns: c,
		},
		Type: cons.PLUGINREADER,
		Name: "streamreader",
	}
	job.SetReaderParam(reader)

	writer := &conf.Param{
		Connect: &conf.Connect{},
		Query:   &conf.Query{},
		Type:    cons.PLUGINWRITER,
		Name:    "streamwriter",
	}
	job.SetWriterParam(writer)
	speed := &conf.Speed{
		Byte:             0,
		BytePerChannel:   0,
		Record:           0,
		RecordPerChannel: 0,
		Channel:          10, // task 数量
		Thread:           10, // threat group数量
	}
	job.SetSpeed(speed)
	return job
}

func Example_Engine() {
	job := streamJob()
	engine.Engine(job)

	// Output:
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world

}

func Example_Engine_Report() {
	job := streamJob()
	engine.EngineReport(job)

	// Output:
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// hello world
	// map[byteSize:10 recordcount:10]
	// map[out:hello world]
}
