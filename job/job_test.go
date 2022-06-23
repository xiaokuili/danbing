package job

import (
	"danbing/conf"
	"danbing/cons"
	"testing"
)

const (
	begin = "2022-02-19 01:03:05"
	end   = "2022-02-20 01:03:05"
	sql   = "select * from danbing where update_time > '2022-02-19 01:03:05' and update_time <= '2022-02-20 01:03:05'"
	count = "select count(*) from danbing where update_time > '2022-02-19 01:03:05' and update_time <= '2022-02-20 01:03:05'"
)

func TJobRefresh(name string) *Job {
	job := New(name)
	c := make([]*conf.Column, 0)
	c = append(c, &conf.Column{
		Name:       "update_time",
		WhereField: true,
	})

	reader := &conf.Param{
		Connect: &conf.Connect{},
		Query: &conf.Query{
			SQL:     "select * from danbing",
			Columns: c,
			Count:   "select count(*) from danbing",
		},
		Type: cons.PLUGINREADER,
		Name: "streamreader",
	}

	job.SetReader(reader)

	writer := &conf.Param{
		Connect: &conf.Connect{},
		Query:   &conf.Query{},
		Type:    cons.PLUGINWRITER,
		Name:    "streamwriter",
	}
	job.SetWriter(writer)
	speed := &conf.Speed{
		Byte:             0,
		BytePerChannel:   0,
		Record:           0,
		RecordPerChannel: 0,
		TaskRecordsNum:   10, // task 数量
		Thread:           10, // threat group数量
	}
	job.SetSpeed(speed)
	return job
}

func TestJob_Refresh(t *testing.T) {
	name := "danbing_refresh_test"
	j := TJobRefresh(name)
	tests := []struct {
		name string
		job  *Job
	}{
		// TODO: Add test cases.
		{
			name: name,
			job:  j,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := tt.job
			j.SetBeginTime("2022-02-19 01:03:05")
			j.SetEndTime("2022-02-20 01:03:05")
			j.Refresh()
			s := j.reader().Query.SQL
			c := j.reader().Query.Count
			if s != sql {
				t.Errorf("sql err; %s-%s", s, sql)
			}
			if c != count {
				t.Errorf("count err; %s-%s", c, count)
			}
		})
	}
}
