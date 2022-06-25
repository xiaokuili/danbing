package myplugin

import (
	"danbing/conf"
	"danbing/plugin"
	"fmt"
	"log"
	"math"
	"strings"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mohae/deepcopy"
)

type PgReader struct {
	Query    *conf.Query
	Connect  *conf.Connect
	db       *sql.DB
	Total    int
	countSQL string
	sql      string
}

func (reader *PgReader) Init(tq *conf.Query, tc *conf.Connect) int {
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

	reader.sql = reader.streamSQL(
		tq.Where,
		tq.Begin,
		tq.End,
	)
	reader.countSQL = fmt.Sprintf("select count(*) from (%s) b", reader.sql)
	return reader.count()
}

func (reader *PgReader) Split(numPerTask int) []plugin.ReaderPlugin {
	plugins := make([]plugin.ReaderPlugin, 0)

	total := reader.count()

	partition := shuffle(total, numPerTask)

	sqlbase := reader.sql

	for i := 0; i < partition; i++ {
		new := reader.copy()
		offset := i * numPerTask
		if i == partition-1 {
			numPerTask = total - offset + 10
		}

		sql := fmt.Sprintf("%s offset %d limit %d", sqlbase, offset, numPerTask)
		// sql := sqlbase

		new.sql = sql
		plugins = append(plugins, new)
	}
	return plugins
}

func (reader *PgReader) Name() string {
	// s := "%s where %s"
	return "pgsqlreader"
}

func (reader *PgReader) Reader() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	rows, err := reader.db.Query(reader.sql)
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

func (reader *PgReader) streamSQL(where, begin, end string) string {
	base := reader.Query.BaseSQL
	if base == "" {
		panic("baseSQL don't exist...")
	}
	if where == "" {
		return base
	}
	var s string

	if strings.Contains(base, "where") {
		s = fmt.Sprintf(
			"%s and %s > '%s' and %s <= '%s' ",
			base,
			where,
			begin,
			where,
			end,
		)
	} else {
		s = fmt.Sprintf(
			"%s where %s > '%s' and %s <= '%s' ",
			base,
			where,
			begin,
			where,
			end,
		)
	}

	if !strings.Contains(s, "count") {
		s = fmt.Sprintf(
			" %s order by %s ",
			s,
			where,
		)
	}
	return s
}

func (reader *PgReader) copy() *PgReader {
	new := deepcopy.Copy(reader)
	p, ok := new.(*PgReader)
	if !ok {
		fmt.Println("")
	}
	p.db = reader.db
	return p
}

func (reader *PgReader) searchCount(sql string) int {
	c, err := searchCount(sql)
	if err == nil {
		return c
	}

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

func shuffle(total, numPerTask int) int {
	if total < numPerTask {
		return 1
	}
	partition := int(math.Ceil(float64(total) / float64(numPerTask)))

	return partition
}

func (reader *PgReader) count() int {
	if reader.Total == 0 {
		reader.Total = reader.searchCount(reader.countSQL)
	}

	return reader.Total
}

// TODO: init必须手动维护
func init() {
	plugin.Register(&PgReader{})
}
