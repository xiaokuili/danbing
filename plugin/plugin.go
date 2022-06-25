package plugin

import (
	"danbing/conf"
	"fmt"
)

const (
	RPlugin = "reader"
	WPlugin = "writer"
)

type ReaderPlugin interface {
	Name() string
	Init(*conf.Query, *conf.Connect) int
	Split(taskNum int) []ReaderPlugin
	Reader() []map[string]interface{}
	Close()
}

type WriterPlugin interface {
	Name() string
	Init(*conf.Query, *conf.Connect)
	Split(taskNum int) []WriterPlugin
	Writer([]map[string]interface{})
	Close()
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
