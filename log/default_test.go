package log

import (
	"testing"
)

func Test_defaultLogger_Info(t *testing.T) {

	type args struct {
		msg     string
		keyVals []interface{}
	}
	logger := MustNewDefaultLogger(LogFormatText, LogLevelDebug)
	kv := []interface{}{"key", "value"}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "logtest",
			args: args{
				msg:     "msg",
				keyVals: kv,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Debug(tt.args.msg, tt.args.keyVals...)
			logger.Info(tt.args.msg, "aaa", "bbb")
		})
	}
}
