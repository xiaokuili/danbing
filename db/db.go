package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Info struct {
	Batch  string
	Table  string
	Uptime string
}

type Database interface {
	Insert(info *Info) error
	Search() *Info
}

func init() {
	os.Remove("./log.db")

	db, err := sql.Open("sqlite3", "./log.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	create table IF NOT EXISTS log(id integer not null primary key, batch varchar(20),table varchar(20), uptime varchar(20));
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return
	}
}

func (i *Info) Insert() error {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(batch, table, uptime) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(i.Batch, i.Table, i.Uptime)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (i *Info) Search(table string) *Info {

	db, err := sql.Open("sqlite3", "./log.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from log ")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	return nil
}
