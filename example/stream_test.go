package engine_test

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/engine"
	"danbing/job"
	"time"
)

func streamJob() *job.Job {
	job := job.New("danbing_stream_test")
	c := make([]*conf.Column, 0)
	c = append(c, &conf.Column{
		Name: "out",
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
		TaskRecordsNum:   10, // task 数量
		Thread:           10, // threat group数量
	}
	job.SetSpeed(speed)

	job.SetStream()
	return job
}

func Example_Stream_Job() {
	job := streamJob()
	engine.Run(job)
	time.Sleep(time.Second * 5)
	engine.Run(job)
	// Output:
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[out:hello world]
	// map[byteSize:10 recordcount:10]
	// map[out:hello world]
}
