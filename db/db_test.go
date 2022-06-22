package db

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	testname = "test"
)

func TestInfo_Insert(t *testing.T) {
	type fields struct {
		Batch  string
		Name   string
		Uptime int64
	}
	one := time.Now().Add(time.Second * 10).Unix()
	two := time.Now().Add(time.Second * 20).Unix()
	three := time.Now().Add(time.Second * 30).Unix()
	fmt.Println(one, two, three)
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "one",
			fields: fields{
				Batch:  "20201",
				Name:   testname,
				Uptime: one,
			},
		},
		{
			name: "two",
			fields: fields{
				Batch:  "20202",
				Name:   testname,
				Uptime: two,
			},
		},
		{
			name: "three",
			fields: fields{
				Batch:  "20203",
				Name:   testname,
				Uptime: three,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Info{
				Batch:  tt.fields.Batch,
				Name:   tt.fields.Name,
				Uptime: tt.fields.Uptime,
			}

			if err := i.Insert(); (err != nil) != tt.wantErr {
				t.Errorf("Info.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	result := SearchLast(testname)
	u := result.Uptime
	if result.Uptime != three {
		t.Errorf("Info.SearchLast() want = %v, get %v", u, three)
	}
}
