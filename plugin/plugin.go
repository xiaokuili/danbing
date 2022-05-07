package plugin

import "fmt"

type ReaderPlugin interface {
	Name() string
	Split()
	Reader()
}

type WriterPlugin interface {
	Name() string
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
