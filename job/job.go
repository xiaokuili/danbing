package job

import (
	"danbing/conf"
	"danbing/cons"
	"fmt"
	"time"
)

type Job struct {
	Param []*conf.Param `json:"job"`
	Speed *conf.Speed   `json:"speed"`
	Table string
	Name  string
	Batch string
}

func New(table string) *Job {
	return &Job{
		Param: []*conf.Param{},
		Speed: &conf.Speed{},
		Table: table,
		Name:  table,
		Batch: time.Now().Format(time.RFC3339),
	}
}

func (j *Job) SetSpeed(s *conf.Speed) {
	j.Speed = s
}

// -1 -> must not need
// 1 -> must need
// 0 -> can't find
func (j *Job) CheckNeedParam() int {
	l := j.Param
	if len(l) >= 2 {
		return -1
	}
	if len(l) == 0 {
		return 1
	}
	return 0
}

func (j *Job) setParam(p *conf.Param) bool {
	need := j.CheckNeedParam()
	if need == -1 {
		return true
	}
	if need == 1 {
		j.Param = append(j.Param, p)
		return true
	}
	return false
}

func (j *Job) SetReaderParam(p *conf.Param) {
	if j.setParam(p) {
		return
	}
	if j.Param[0].Type != cons.PLUGINREADER {
		j.Param = append(j.Param, p)
	}
}

func (j *Job) SetWriterParam(p *conf.Param) {
	if j.setParam(p) {
		return
	}
	if j.Param[0].Type != cons.PLUGINWRITER {
		j.Param = append(j.Param, p)
	}
}

func (j *Job) reader() *conf.Param {
	var r *conf.Param
	for i := 0; i < len(j.Param); i++ {
		if j.Param[i].Type == cons.PLUGINREADER {
			r = j.Param[i]
		}
	}
	return r
}

func (j *Job) writer() *conf.Param {
	var w *conf.Param
	for i := 0; i < len(j.Param); i++ {
		if j.Param[i].Type == cons.PLUGINWRITER {
			w = j.Param[i]
		}
	}
	return w
}

func (j *Job) SetBeginTime(b string) {
	r := j.reader()
	if j.isUpdate() {
		r.Query.BeginTime = b
	}
}

func (j *Job) SetEndTime(b string) {
	r := j.reader()
	if j.isUpdate() {
		r.Query.EndTime = b
	}
}

func (j *Job) isUpdate() bool {
	return j.where() != ""
}

func (j *Job) where() string {
	var column []*conf.Column
	for i := 0; i < len(j.Param); i++ {
		if j.Param[i].Type == cons.PLUGINREADER {
			column = j.Param[i].Query.Columns
		}
	}
	for i := 0; i < len(column); i++ {
		c := column[i]
		if c.WhereField {
			return c.Name
		}
	}
	return ""
}

func (j *Job) Refresh() {
	if j.isUpdate() {
		reader := j.reader()
		where := j.where()
		whereSQL := fmt.Sprintf(" where %s > '%s' and %s <= '%s'", where, reader.Query.BeginTime, where, reader.Query.EndTime)
		sql := reader.Query.SQL
		sql = sql + whereSQL
		reader.Query.SQL = sql
		count := reader.Query.Count + whereSQL
		reader.Query.Count = count

	}

}
