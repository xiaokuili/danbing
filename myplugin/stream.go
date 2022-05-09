package myplugin

import (
	"danbing/plugin"
	"danbing/task"
	"fmt"
)

type StreamReader struct {
	Query *task.Query
}

func (sr *StreamReader) Init(tq *task.Query, tc *task.Connect) {
	sr.Query = tq
}
func (sr *StreamReader) Name() string {
	return "streamreader"
}

func (sr *StreamReader) Split(taskNum int) []plugin.ReaderPlugin {
	plugins := make([]plugin.ReaderPlugin, 0)
	for i := 0; i < taskNum; i++ {
		plugins = append(plugins, sr)
	}
	return plugins
}

func (sr *StreamReader) Reader() {
	fmt.Printf("[reader]%s \n", sr.Name())

}

type StreamWriter struct {
}

func (sw *StreamWriter) Init(tq *task.Query, tc *task.Connect) {

}

func (sw *StreamWriter) Name() string {
	return "streamwriter"
}

func (sw *StreamWriter) Split(taskNum int) []plugin.WriterPlugin {
	plugins := make([]plugin.WriterPlugin, 0)
	for i := 0; i < taskNum; i++ {
		plugins = append(plugins, sw)
	}
	return plugins
}

func (sw *StreamWriter) Writer() {
	fmt.Printf("[writer]%s\n", sw.Name())

}

// TODO: init必须手动维护
func init() {
	plugin.Register(&StreamReader{})
	plugin.Register(&StreamWriter{})
}
