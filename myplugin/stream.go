package myplugin

import (
	"danbing/plugin"
	"fmt"
)

type StreamReader struct {
}

func (sr *StreamReader) Name() string {
	return "streamreader"
}

func (sr *StreamReader) Split() {
	fmt.Println("run plugin split")
}

func (sr *StreamReader) Reader() {

}

type StreamWriter struct {
}

func (sw *StreamWriter) Name() string {
	return "streamwriter"
}

func (sw *StreamWriter) Writer() {

}

// TODO: init必须手动维护
func init() {
	plugin.Register(&StreamReader{})
	plugin.Register(&StreamWriter{})
}
