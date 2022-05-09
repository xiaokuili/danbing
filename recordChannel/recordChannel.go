package recordchannel

import "time"

type Record struct {
	C chan []byte
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
	r.C <- record
}
