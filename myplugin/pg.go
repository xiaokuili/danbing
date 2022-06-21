package myplugin

import (
	"danbing/conf"
	"danbing/plugin"
	"fmt"
	"log"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mohae/deepcopy"
)

type PgReader struct {
	Query   *conf.Query
	Connect *conf.Connect
	db      *sql.DB
	Total   int
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

func (reader *PgReader) searchCount(sql string) int {
	rows, err := reader.db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}
	return count

}

func shuffle(total, task int) int {

	partition := total / task

	return partition
}

func (reader *PgReader) Split(taskNum int) []plugin.ReaderPlugin {
	plugins := make([]plugin.ReaderPlugin, 0)

	total := reader.Count()

	partition := shuffle(total, taskNum)

	sqlbase := reader.Query.SQL

	for i := 0; i < taskNum; i++ {
		new := reader.Copy()
		offset := i * partition
		if i == taskNum-1 {
			partition = total - offset + 10
		}

		sql := fmt.Sprintf("%s offset %d limit %d", sqlbase, offset, partition)
		// sql := sqlbase

		new.Query.SQL = sql
		plugins = append(plugins, new)
	}
	return plugins
}

func (reader *PgReader) Count() int {
	if reader.Total == 0 {
		reader.Total = reader.searchCount(reader.Query.Count)
	}

	return reader.Total
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
