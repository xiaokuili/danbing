package recordchannel

import (
	statistic "danbing/statistics"
	"time"
)

type Record struct {
	C chan []byte

	Communication *statistic.Communication
}

func New(communication *statistic.Communication) *Record {
	return &Record{
		C:             make(chan []byte),
		Communication: communication,
	}
}

func (r *Record) GetRecord() []byte {
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

func (r *Record) SetRecord(record []byte) {
	r.Communication.AddCounter("byteSize", len(record))
	r.Communication.IncreaseCounter("recordcount")
	r.C <- record

}
