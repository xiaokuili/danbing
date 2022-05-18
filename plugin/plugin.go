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
	Init(*conf.Query, *conf.Connect)
	Name() string
	Split(taskNum int) []ReaderPlugin
	Reader() string
	Close()
}

type WriterPlugin interface {
	Name() string
	Init(*conf.Query, *conf.Connect)
	Split(taskNum int) []WriterPlugin
	Writer(s string)
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
