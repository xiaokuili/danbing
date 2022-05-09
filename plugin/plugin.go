package plugin

import (
	"danbing/task"
	"fmt"
)

const (
	RPlugin = "reader"
	WPlugin = "writer"
)

type ReaderPlugin interface {
	Init(*task.Query, *task.Connect)
	Name() string
	Split(taskNum int) []ReaderPlugin
	Reader()
}

type WriterPlugin interface {
	Name() string
	Init(*task.Query, *task.Connect)
	Split(taskNum int) []WriterPlugin
	Writer()
}

var ReaderPlugins map[string]ReaderPlugin = make(map[string]ReaderPlugin)
var WriterPlugins map[string]WriterPlugin = make(map[string]WriterPlugin)

func Register(plugin interface{}) {
	switch plugin := plugin.(type) {
	case WriterPlugin:
		WriterPlugins[plugin.Name()] = plugin
	case ReaderPlugin:
		ReaderPlugins[plugin.Name()] = plugin

	default:
		fmt.Printf("%s can't register\n", plugin)
	}

}
