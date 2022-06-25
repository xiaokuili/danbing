package conf

import (
	"danbing/cons"
	"fmt"
)

type Query struct {
	BaseSQL string `json:"sql"`
	Table   string `json:"table"`
	Where   string
	Primary []string
	Begin   string
	End     string
	Count   int
}

func NewQuery(sql, t, w string, c int, p []string) *Query {
	return &Query{
		BaseSQL: sql,
		Table:   t,
		Where:   w,
		Primary: p,
		Count:   c,
	}
}

type Connect struct {
	Host     string `json:"host"` // host: port
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func NewConn(host string, port int, user, pass, db string) *Connect {
	return &Connect{
		Host:     host,
		Port:     port,
		Username: user,
		Password: pass,
		Database: db,
	}
}

type Speed struct {
	NumPerTask int `json:"num_per_task"` // 每个任务的数据条数，总任务数=总数量/每个任务数据数量
	Thread     int `json:"thread"`       // 执行线程数

}

// NewSpeed need RecordPerTask and Thread
// task num = total / RecordPerTask
// then group by thread
func NewSpeed(r, t int) *Speed {
	return &Speed{
		NumPerTask: r,
		Thread:     t,
	}
}

type Param struct {
	Connect *Connect `json:"connect"`
	Query   *Query   `json:"query"`
	Name    string   `json:"name"` // reader or writer name
	Type    string   `json:"type"` // reader type or writer type
}

func NewReader(name string, conn *Connect, query *Query) *Param {
	reader := &Param{
		Connect: conn,
		Query:   query,
		Name:    name,
		Type:    cons.PLUGINREADER,
	}
	return reader
}

func NewWriter(name string, conn *Connect, query *Query) *Param {
	reader := &Param{
		Connect: conn,
		Query:   query,
		Name:    name,
		Type:    cons.PLUGINWRITER,
	}
	return reader
}

func Reader(ps []*Param) (*Param, error) {
	return getParam(ps, cons.PLUGINREADER)
}

func Writer(ps []*Param) (*Param, error) {
	return getParam(ps, cons.PLUGINWRITER)
}

func getParam(ps []*Param, t string) (*Param, error) {
	for i := 0; i < len(ps); i++ {
		p := ps[i]

		if p.Type == t {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not exist %s param ", t)
}
