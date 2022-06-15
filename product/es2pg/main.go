package main

import (
	"danbing/conf"
	"danbing/cons"
	"danbing/engine"
	"danbing/job"
)

// 	select
// 	ENCODE(dm.id::bytea, 'hex') as key,
// 	jsonb_build_object('data', dm.*, 'id', dm.id, 's_time', cast(extract(EPOCH from CURRENT_TIMESTAMP )* 1000 as int8), 'savetime', to_char(current_timestamp, 'YYYY-MM-DD HH24:MI:SS'), 'table', 'dm_ecoi_edb_value') as value
// from
// dm_ecoi_edb_value as dm

func pg2esJob() *job.Job {
	sql := `select id from dm_lget_company_addissue`
	job := job.New("danbing")
	c := make([]*conf.Column, 0)
	c = append(c, &conf.Column{
		Name:         "id",
		PrimaryField: true, // 收集这个字段的最后一条数据
	})
	// c = append(c, &conf.Column{
	// 	Name:       "value",
	// 	FieldStype: "object",
	// })
	reader := &conf.Param{
		Connect: &conf.Connect{
			Host:     "192.168.200.200",
			Port:     5432,
			Username: "dm",
			Password: "ZybsHt2oY4l2V200",
			Database: "dm",
		},
		Query: &conf.Query{
			SQL:   sql,
			Count: "select count(*) from dm_lget_company_addissue",
		},
		Type: cons.PLUGINREADER,
		Name: "pgsqlreader",
	}
	job.SetReaderParam(reader)

	writer := &conf.Param{
		Connect: &conf.Connect{
			Host:     "192.168.216.124",
			Port:     18054,
			Username: "elastic",
			Password: "HGeMa7BMi6CLjNbPmONZ",
			Database: "",
		},
		Query: &conf.Query{

			Table:   "danbingtest",
			Columns: c,
		},
		Type: cons.PLUGINWRITER,
		Name: "eswriter",
	}
	job.SetWriterParam(writer)
	speed := &conf.Speed{
		Byte:             0,
		BytePerChannel:   0,
		Record:           0,
		RecordPerChannel: 0,
		Channel:          80, // task 数量
		Thread:           50, // threat group数量
	}
	job.SetSpeed(speed)
	return job
}
func main() {

	job := pg2esJob()
	engine.EngineReport(job)
}
