package conf

import "danbing/cons"

// 字段
type Column struct {
	Value        string `json:"value"`
	Name         string `json:"name"`
	WhereField   bool   `json:"where_field"`   // update
	PrimaryField bool   `json:"primary_field"` // upsert
}

type Query struct {
	SQL     string `json:"sql"`
	Size    int
	Offset  int
	Table   string    `json:"table"`
	Columns []*Column `json:"columns"`
}

type Connect struct {
	Host     string `json:"host"` // host: port
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Speed struct {
	Byte             int `json:"byte"`
	BytePerChannel   int `json:"byte_per_channel"`
	Record           int `json:"record"`
	RecordPerChannel int `json:"record_per_channel"`
	Channel          int `json:"channel"` // 拆分任务数量 = 总记录/每个任务数量(byte类似)
	Thread           int `json:"thread"`  // 执行线程数
}

type Param struct {
	Connect *Connect `json:"connect"`
	Query   *Query   `json:"query"`
	Name    string   `json:"name"` // reader or writer name
	Type    string   `json:"type"` // reader type or writer type
}

func NewReader(name string) *Param {
	reader := &Param{
		Connect: &Connect{},
		Query:   &Query{},
		Name:    name,
		Type:    cons.CONfREADER,
	}
	return reader
}

func (p *Param) SetConnect(c *Connect) {
	p.Connect = c
}

func (p *Param) SetQuery(q *Query) {
	p.Query = q
}
