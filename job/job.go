package job

import (
	"danbing/conf"
	"danbing/cons"
)

type Job struct {
	Param []*conf.Param `json:"job"`
	Speed *conf.Speed   `json:"speed"`
	Table string
}

func New(table string) *Job {
	return &Job{
		Param: []*conf.Param{},
		Speed: &conf.Speed{},
		Table: table,
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
	if j.Param[0].Type != cons.CONfREADER {
		j.Param = append(j.Param, p)
	}
}

func (j *Job) SetWriterParam(p *conf.Param) {
	if j.setParam(p) {
		return
	}
	if j.Param[0].Type != cons.CONfWRITER {
		j.Param = append(j.Param, p)
	}
}
