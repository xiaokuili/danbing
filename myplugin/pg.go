package myplugin

import (
	"danbing/plugin"
	"danbing/task"
	"encoding/json"

	"github.com/go-pg/pg/v10"
)

type PgReader struct {
	Query   *task.Query
	Connect *task.Connect
	db      *pg.DB
}

func (reader *PgReader) Init(tq *task.Query, tc *task.Connect) {
	db := pg.Connect(&pg.Options{
		Addr:     tc.URL,
		User:     tc.Username,
		Password: tc.Password,
		Database: tc.Database,
	})
	reader.db = db
	reader.Query = tq
}
func (reader *PgReader) Name() string {
	return "pgsqlreader"
}

func (reader *PgReader) Split(taskNum int) []plugin.ReaderPlugin {
	plugins := make([]plugin.ReaderPlugin, 0)
	for i := 0; i < taskNum; i++ {
		reader.Query.Offset = i * reader.Query.Size
		plugins = append(plugins, reader)
	}
	return plugins
}

func (reader *PgReader) Reader() string {
	result := make(map[string]interface{})
	reader.db.Model(result).Query(reader.Query.SQL, "")
	s, _ := json.Marshal(result)
	return string(s)
}

func (reader *PgReader) Close() {
	reader.db.Close()
}
