package myplugin

import (
	"danbing/conf"
	"danbing/plugin"
	"fmt"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mohae/deepcopy"
)

type PgReader struct {
	Query   *conf.Query
	Connect *conf.Connect
	db      *sql.DB
}

func (reader *PgReader) Init(tq *conf.Query, tc *conf.Connect) {
	var pool *sql.DB
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		tc.Host, tc.Port, tc.Username, tc.Password, tc.Database)
	pool, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Println(err)
	}
	err = pool.Ping()
	if err != nil {
		panic(err)
	}
	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)
	reader.db = pool
	reader.Query = tq

}
func (reader *PgReader) Name() string {
	return "pgsqlreader"
}

func (reader *PgReader) Copy() *PgReader {
	new := deepcopy.Copy(reader)
	p, ok := new.(*PgReader)
	if !ok {
		fmt.Println("")
	}
	p.db = reader.db
	return p
}

func (reader *PgReader) Split(taskNum int) []plugin.ReaderPlugin {
	plugins := make([]plugin.ReaderPlugin, 0)
	sqlbase := reader.Query.SQL

	for i := 0; i < taskNum; i++ {
		new := reader.Copy()

		new.Query.Offset = i * reader.Query.Size
		sql := fmt.Sprintf("%s offset %d limit %d", sqlbase, new.Query.Offset, reader.Query.Size)
		// sql := sqlbase

		new.Query.SQL = sql
		plugins = append(plugins, new)
	}
	return plugins
}

func (reader *PgReader) Reader() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	rows, err := reader.db.Query(reader.Query.SQL)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	if len(cols) > 0 {
		for rows.Next() {
			buff := make([]interface{}, len(cols))
			data := make([][]byte, len(cols)) //数据库中的NULL值可以扫描到字节中
			for i, _ := range buff {
				buff[i] = &data[i]
			}
			rows.Scan(buff...)
			dataKv := make(map[string]interface{})
			for k, col := range data {
				dataKv[cols[k]] = string(col)
			}
			result = append(result, dataKv)
		}
	}
	return result
}

func (reader *PgReader) Close() {
	reader.db.Close()
}

// TODO: init必须手动维护
func init() {
	plugin.Register(&PgReader{})
}
