package engine_test

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/engine"
	"danbing/job"
	"os/exec"
	"time"
)

func pgtostreamJob() *job.Job {
	job := job.New("danbing")
	// c := make([]*conf.Column, 0)
	//  c = append(c, &conf.Column{
	// 	Name:         "out",
	// 	CollectField: true, // 收集这个字段的最后一条数据
	// })
	reader := &conf.Param{
		Connect: &conf.Connect{
			Host:     "127.0.0.1",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Database: "postgres",
		},
		Query: &conf.Query{
			SQL:   "select * from danbing",
			Count: "select count(*) from danbing",
		},
		Type: cons.PLUGINREADER,
		Name: "pgsqlreader",
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

func runShell(name string, arg ...string) {
	cmd := exec.Command(name, arg...)

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func Example_PG2STREAM() {
	c := "docker-compose"

	runShell(c, "-f", "../docker/docker-compose.yml", "down")
	runShell(c, "-f", "../docker/docker-compose.yml", "up", "-d")
	time.Sleep(time.Second * 5)
	defer runShell(c, "-f", "../docker/docker-compose.yml", "down")

	job := pgtostreamJob()
	engine.Engine(job)

	// Output:
}
