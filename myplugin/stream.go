package myplugin

import "danbing/plugin"

type StreamReader struct {
}

func (sr *StreamReader) Name() string {
	return "streamReader"
}

func (sr *StreamReader) Split() {

}

func (sr *StreamReader) Reader() {

}

type StreamWriter struct {
}

func (sw *StreamWriter) Name() string {
	return "streamWriter"
}

func (sw *StreamWriter) Writer() {

}

// TODO: init必须手动维护
func init() {
	plugin.Register(&StreamReader{})
	plugin.Register(&StreamWriter{})
}
