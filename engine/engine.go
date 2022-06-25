package engine

import (
	"danbing/conf"
	"danbing/db"
	"danbing/job"
	_ "danbing/myplugin"
	"danbing/scheduler"
	"errors"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	BEGIN      = "2020-01-01 01:01:01"
	End        = "2000-01-01 01:01:01"
	TimeFormat = "2006-01-02 15:04:05"

	rhost     = "reader.conn.host"
	rport     = "reader.conn.port"
	ruser     = "reader.conn.user"
	rpassword = "reader.conn.password"
	rdb       = "reader.conn.db"

	rBaseSQL = "reader.query.baseSQL"
	rWhere   = "reader.query.where"
	rPrimary = "reader.query.primary"
	rName    = "reader.name"
	rCount   = "reader.count"

	whost     = "writer.conn.host"
	wport     = "writer.conn.port"
	wuser     = "writer.conn.user"
	wpassword = "writer.conn.password"
	wdb       = "writer.conn.db"

	wBaseSQL = "writer.query.baseSQL"
	wWhere   = "writer.query.where"
	wPrimary = "writer.query.primary"
	wName    = "writer.name"
	wCount   = "writer.count"

	thread     = "speed.thread"
	numPerTask = "speed.num_per_task"

	project = "danbing"
	dbPath  = "/app/danbing/file/danbing.db"
	jobPath = "/app/danbing/file"

	log = "log"
)

func Begin(f *db.File, j *job.Job) string {
	info := f.SearchLast(j.Table)
	if info != nil {
		t := time.Unix(info.End, 0)
		return t.Format(TimeFormat)
	}
	return BEGIN
}

func Ending() string {
	return time.Now().Format(TimeFormat)
}

func Build(path string) *job.Job {

	name := viper.GetString("name")
	j := job.New(name)
	// 创建链接
	rconn := conf.NewConn(viper.GetString(rhost),
		viper.GetInt(rport),
		viper.GetString(ruser),
		viper.GetString(rpassword),
		viper.GetString(rdb),
	)
	wconn := conf.NewConn(viper.GetString(whost),
		viper.GetInt(wport),
		viper.GetString(wuser),
		viper.GetString(wpassword),
		viper.GetString(wdb),
	)
	// 创建查询语句
	rquery := conf.NewQuery(
		viper.GetString(rBaseSQL),
		name,
		viper.GetString(rWhere),
		viper.GetInt(rCount),
		viper.GetStringSlice(rPrimary),
	)
	wquery := conf.NewQuery(
		viper.GetString(wBaseSQL),
		name,
		viper.GetString(wWhere),
		viper.GetInt(wCount),
		viper.GetStringSlice(wPrimary))

	// 创建reader和writer
	j.SetReader(
		conf.NewReader(viper.GetString(rName),
			rconn,
			rquery,
		),
	)
	j.SetWriter(
		conf.NewWriter(viper.GetString(wName),
			wconn,
			wquery,
		),
	)

	// 速度调节
	speed := conf.NewSpeed(viper.GetInt(numPerTask), viper.GetInt(thread))
	j.SetSpeed(speed)
	return j
}

func Run() {
	pflag.String("d", "/app/danbing/file/danbing.db", "dbpath: 数据库存储地址")
	pflag.String("j", "/app/danbing/file", "jobpath: 任务路径")
	pflag.String("n", "config", "name:配置文件名称")
	pflag.String("t", "transform", "type: transform or check")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.AddConfigPath(viper.GetString("j"))
	viper.SetConfigName(viper.GetString("n"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	t := viper.GetString("t")
	switch t {
	case "transform":
		tranform()
	case "check":
		check()
	}

}

func tranform() {
	file := db.New(dbPath)
	job := Build(jobPath)

	reader := job.Reader()
	where := reader.Query.Where
	b, e := BEGIN, End
	if where != "" {
		b = Begin(file, job)
		e = Ending()
		reader.Query.Begin = b
		reader.Query.End = e
	}

	scheduler.Run(job, viper.GetString(log))

	nb, err := time.ParseInLocation(TimeFormat, b, time.Local)
	if err != nil {
	}
	ne, err := time.ParseInLocation(TimeFormat, e, time.Local)
	if err != nil {
	}

	info := &db.Info{
		Name:  job.Table,
		Batch: job.Batch,
		Begin: nb.Unix(),
		End:   ne.Unix(),
	}
	file.Insert(info)
}

func countBeginEnd(f *db.File, j *job.Job) (begin, end string, err error) {
	info := f.SearchLast(j.Name)
	if info != nil {
		begin := time.Unix(info.Begin, 0)
		end := time.Unix(info.End, 0)
		return begin.Format(TimeFormat), end.Format(TimeFormat), nil
	}
	return "", "", errors.New("dont run this job")
}

func check() {
	file := db.New(dbPath)
	job := Build(jobPath)
	reader := job.Reader()
	where := reader.Query.Where
	b, e, err := countBeginEnd(file, job)
	if err != nil {
		panic("")
	}
	if where == "" {
		panic("")
	}

	reader.Query.Begin = b
	reader.Query.End = e

	scheduler.Run(job, viper.GetString(log))

}
