package db

import (
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInfo_Search(t *testing.T) {
	type fields struct {
		Batch  string
		Table  string
		Uptime string
	}
	type args struct {
		table string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Info
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				Batch:  "1",
				Table:  "test",
				Uptime: "2020-02-03",
			},
			args: args{},
			want: &Info{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Info{
				Batch:  tt.fields.Batch,
				Table:  tt.fields.Table,
				Uptime: tt.fields.Uptime,
			}
			i.Insert()
			if got := i.Search(tt.args.table); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Info.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
