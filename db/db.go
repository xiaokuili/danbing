package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	table     = "log"
	insertSQL = "INSERT INTO  log (batch, name, uptime) values(?, ?, ?)"
	searchSQL = "select batch,name,uptime   from log where name='%s' order by uptime desc limit 1;	"
)

type Info struct {
	Batch  string
	Name   string
	Uptime int64
}

type Database interface {
	Insert(info *Info) error
	Search() *Info
}

func init() {
	createTable()
}

func createDB() *sql.DB {
	locat := fmt.Sprintf("./%s.db", table)
	db, err := sql.Open("sqlite3", locat)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createTable() error {
	db := createDB()
	defer db.Close()
	sqlStmt := `
		CREATE table if not exists log (
			id integer not null primary key,
			batch varchar(20),
			name varchar(20),
			uptime INTEGER 	
		)
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil

}

func (i *Info) Insert() error {
	db := createDB()
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
	_, err = stmt.Exec(i.Batch, i.Name, i.Uptime)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Search(sql string) []*Info {

	db := createDB()
	defer db.Close()
	result := make([]*Info, 0)
	rows, err := db.Query(sql)
	if err != nil {

	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var uptime int64
		var batch string
		err = rows.Scan(&batch, &name, &uptime)
		if err != nil {

		}
		info := Info{
			Batch:  batch,
			Name:   name,
			Uptime: uptime,
		}
		result = append(result, &info)
	}
	return result
}

func SearchLast(name string) *Info {
	var one *Info
	sql := fmt.Sprintf(searchSQL, name)
	result := Search(sql)
	if len(result) == 1 {
		one = result[0]
	}
	return one
}
