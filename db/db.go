package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	table     = "log"
	createSQL = `		CREATE table if not exists log (
		id integer not null primary key,
		batch varchar(20),
		name varchar(20),
		begin INTEGER,
		end  	INTEGER
	)`
	insertSQL = "INSERT INTO  log (batch, name, begin, end) values(?, ?, ?, ?)"
	searchSQL = "select batch,name,begin,end   from log where name='%s' order by end desc limit 1;	"
)

type Info struct {
	Batch string
	Name  string
	Begin int64
	End   int64
}

type Database interface {
	Insert(info *Info) error
	Search() *Info
}

type File struct {
	Path string
}

func New(path string) *File {
	f := &File{
		Path: path,
	}
	f.createTable()
	return f
}

func (f *File) initDB() *sql.DB {
	db, err := sql.Open("sqlite3", f.Path)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (f *File) createTable() error {
	db := f.initDB()
	defer db.Close()
	sqlStmt := createSQL
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil

}

func (f *File) Insert(i *Info) error {
	db := f.initDB()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(i.Batch, i.Name, i.Begin, i.End)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (f *File) Search(sql string) []*Info {

	db := f.initDB()
	defer db.Close()
	result := make([]*Info, 0)
	rows, err := db.Query(sql)
	if err != nil {

	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var begin int64
		var end int64
		var batch string
		err = rows.Scan(&batch, &name, &begin, &end)
		if err != nil {

		}
		info := Info{
			Batch: batch,
			Name:  name,
			Begin: begin,
			End:   end,
		}
		result = append(result, &info)
	}
	return result
}

func (f *File) SearchLast(name string) *Info {
	var one *Info
	sql := fmt.Sprintf(searchSQL, name)
	result := f.Search(sql)
	if len(result) == 1 {
		one = result[0]
	}
	return one
}
