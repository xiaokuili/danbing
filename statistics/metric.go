package statistic

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	WAITING State = 0
	RUNNING State = 1

	KILLED    State = 2
	FAILED    State = 3
	SUCCEEDED State = 4
)

type Metric struct {
	Counter    map[string]int
	State      State
	Throwable  string
	Message    map[string]string
	Timestamp  time.Time
	sync.Mutex // 直接通过匿名调用
}

func NewMetric() *Metric {
	return &Metric{
		Counter:   map[string]int{},
		State:     0,
		Throwable: "",
		Message:   map[string]string{},
		Timestamp: time.Time{},
		Mutex:     sync.Mutex{},
	}
}

func (c *Metric) GetTimestamp() time.Time {
	c.Lock()
	defer c.Unlock()
	return c.Timestamp
}

func (c *Metric) SetTimestamp(t time.Time) {
	c.Lock()
	defer c.Unlock()
	c.Timestamp = t
}

func (c *Metric) GetCounter(key string) (int, error) {
	c.Lock()
	defer c.Unlock()
	r, ok := c.Counter[key]
	if ok {
		return r, nil
	}
	return 0, errors.New("cant find key value")
}

func (c *Metric) GetMessage(key string) (string, error) {
	c.Lock()
	defer c.Unlock()
	r, ok := c.Message[key]
	if ok {
		return r, nil
	}
	return "", errors.New("cant find msg value")
}

func (c *Metric) IncreaseCounter(key string) {
	c.Lock()
	defer c.Unlock()
	c.Counter[key] = c.Counter[key] + 1
}

func (c *Metric) AddCounter(key string, value int) {
	c.Lock()
	defer c.Unlock()
	c.Counter[key] = c.Counter[key] + value
}

func (c *Metric) MergeFrom(final *Metric) {
	c.Lock()
	defer c.Unlock()
	c.Throwable = final.Throwable
	c.Timestamp = time.Now()

	for k, v := range final.Counter {
		c.Counter[k] = v + c.Counter[k]
	}

	for k, v := range final.Message {
		c.Message[k] = v
	}
}
