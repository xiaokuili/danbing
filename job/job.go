package job

import (
	"danbing/conf"
	"danbing/cons"
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

func (j *Job) SetReader(p *conf.Param) {
	j.setParam(cons.PLUGINREADER, p)
}

func (j *Job) SetWriter(p *conf.Param) {
	j.setParam(cons.PLUGINWRITER, p)

}
func (j *Job) Reader() *conf.Param {
	var r *conf.Param
	r, err := conf.Reader(j.Param)
	if err != nil {
		panic("reader don't exist")
	}
	return r
}

func (j *Job) Writer() *conf.Param {
	var r *conf.Param
	conf.Writer(j.Param)
	return r
}

func (j *Job) setParam(t string, p *conf.Param) {
	if j.existPlugin(t) {
		return
	}
	if j.overPlugin() {
		return
	}
	j.Param = append(j.Param, p)
}

func (j *Job) existPlugin(t string) bool {
	ps := j.Param
	for i := 0; i < len(ps); i++ {
		p := ps[i]
		if p.Type == t {
			return true
		}
	}
	return false
}

func (j *Job) overPlugin() bool {
	return len(j.Param) > 2
}
