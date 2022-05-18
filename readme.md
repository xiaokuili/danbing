## danbing
dangbing内网数据同步工具，同时有简单的数据处理模块

## feature
- 内网数据传输
- 支持不同版本数据库插件式开发
- 基于表传输结果统计
- 自动调节传输速度

## danbing vs datax
1. danbing只需要配置数据传输必要信息，不需要配置速度，danbing会自动调节速度，尽可能缩减传输时间

## qucikstart
cd example
go test 

```
package engine_test

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/engine"
	"danbing/job"
)

// 创建streamreader和streamwriter
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


```

