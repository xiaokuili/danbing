package main

import (
	_ "danbing/myplugin"
	"danbing/plugin"
	"encoding/json"
	"fmt"
	"os"
)

// 字段
type Column struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// 参数
type Para struct {
	Column []*Column `json:"column"`
	Slice  int       `json:"sliceRecordCount"`
}

// 插件RW
type RWPlugin struct {
	Name string `json:"name"`
	Para *Para  `json:"parameter"`
}

type Job struct {
	Content []map[string]*RWPlugin `json:"content"`
	Setting map[string]interface{} `json:"setting"`
}

type Task struct {
	Reader *plugin.ReaderPlugin
	Writer *plugin.WriterPlugin
}

func New() *Job {
	j := &Job{
		Content: []map[string]*RWPlugin{},
		Setting: map[string]interface{}{},
	}

	return j
}

func (j *Job) Init(path string) {
	config, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	v := make(map[string]*Job)
	err = json.Unmarshal(config, &v)
	if err != nil {
		panic(err)
	}
	j.Content = v["job"].Content
	j.Setting = v["job"].Setting
}

func (j *Job) Split() {
	content := j.Content
	for i := 0; i < len(content); i++ {
		rw := content[i]
		r := rw["reader"]
		plugin.ReaderPlugins[r.Name].Split()

	}
}

func (j *Job) GroupTask() {
	fmt.Println("group task")
}

func (j *Job) Scheduler() {
	fmt.Println("scheduler")
}

func main() {
	path := "/app/danbing/example/stream2steam.json"
	j := New()
	j.Init(path)
	j.Split()
	j.GroupTask()
	j.Scheduler()
}
