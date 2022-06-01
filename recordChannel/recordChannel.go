package recordchannel

import (
	statistic "danbing/statistics"
	"time"
)

type Record struct {
	C chan []map[string]interface{}

	Communication *statistic.Communication
}

func New(communication *statistic.Communication) *Record {
	return &Record{
		C:             make(chan []map[string]interface{}),
		Communication: communication,
	}
}

func (r *Record) GetRecord() []map[string]interface{} {
	aliveTime := time.Hour * 3
	t := time.NewTicker(aliveTime)
	for {
		select {
		case record := <-r.C:
			t.Reset(aliveTime)
			return record
		case <-t.C:
			panic("")
		}
	}
}

func (r *Record) PutRecord(record []map[string]interface{}) {

	r.Communication.AddCounter(statistic.RecordCount, len(record))
	r.C <- record

}
