package myplugin

import (
	"danbing/conf"
	"danbing/plugin"
	"fmt"
)

type StreamReader struct {
	Query *conf.Query
}

func (sr *StreamReader) Init(tq *conf.Query, tc *conf.Connect) {
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

func (sr *StreamReader) Close() {

}
func (sr *StreamReader) Reader() []map[string]interface{} {
	rst := make([]map[string]interface{}, 0)
	m := make(map[string]interface{})
	m["out"] = sr.Query.SQL
	rst = append(rst, m)
	return rst
}

type StreamWriter struct {
}

func (sw *StreamWriter) Init(tq *conf.Query, tc *conf.Connect) {

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

func (sw *StreamWriter) Close() {

}

func (sw *StreamWriter) Writer(ss []map[string]interface{}) {
	for i := 0; i < len(ss); i++ {
		s := ss[i]
		fmt.Printf("%s\n", s["out"])

	}
}

// TODO: init必须手动维护
func init() {
	plugin.Register(&StreamReader{})
	plugin.Register(&StreamWriter{})
}
