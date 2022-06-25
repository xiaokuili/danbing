package myplugin

import (
	"testing"

	_ "github.com/lib/pq"
)

func Test_shuffle(t *testing.T) {
	type args struct {
		total      int
		numPerTask int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			args: args{
				total:      100,
				numPerTask: 10,
			},
			want: 10,
		},
		{
			args: args{
				total:      100,
				numPerTask: 11,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shuffle(tt.args.total, tt.args.numPerTask); got != tt.want {
				t.Errorf("shuffle() = %v, want %v", got, tt.want)
			}
		})
	}
}
