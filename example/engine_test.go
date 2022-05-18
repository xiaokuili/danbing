package engine_test

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/engine"
	"danbing/job"
)

func Example_Engine() {
	job := job.New("danbing")
	reader := &conf.Param{
		Connect: &conf.Connect{},
		Query: &conf.Query{
			SQL: "hello world",
		},
		Type: cons.CONfREADER,
		Name: "streamreader",
	}
	job.SetReaderParam(reader)

	writer := &conf.Param{
		Connect: &conf.Connect{},
		Query:   &conf.Query{},
		Type:    cons.CONfWRITER,
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
