package main

import (
	_ "danbing/myplugin"
	"danbing/plugin"
	"fmt"
)

func main() {
	fmt.Println(plugin.ReaderPlugins)
	fmt.Println(plugin.WriterPlugins)
}
