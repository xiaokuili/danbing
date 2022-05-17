package statistic

import (
	"danbing/cons"
	"fmt"
	"time"
)

// 接口
type Communicator interface {
	Report(*Metric)
	Name() string
}

var Collects map[string]Communicator = make(map[string]Communicator)
var Instance *Communication

func Register(Communicator Communicator) {
	Collects[Communicator.Name()] = Communicator
}

type Communication struct {
	Child  []*Communication
	ID     int
	Metric *Metric
	Table  string
	Name   string
}

var C *Communication = &Communication{}

// 单例模式
func SingletonNew() *Communication {
	if Instance == nil {
		Instance = &Communication{
			Child: []*Communication{},
			Metric: &Metric{
				Counter:   map[string]int{},
				State:     0,
				Throwable: "",
				Message:   map[string]string{},
				Timestamp: time.Time{},
			},
			ID:   -1,
			Name: cons.STAGEJOB,
		}
	}
	return Instance
}

func New(id int, name, table string) *Communication {

	return &Communication{
		ID:    id,
		Name:  name,
		Child: []*Communication{},
		Metric: &Metric{
			Counter:   map[string]int{},
			State:     0,
			Throwable: "",
			Message:   map[string]string{},
			Timestamp: time.Time{},
		},
	}
}

func (c *Communication) Collect() *Metric {
	m := c.Metric
	for i := 0; i < len(c.Child); i++ {
		m.MergeFrom(c.Child[i].Collect())
	}
	return m
}

func (c *Communication) Report() {
	m := c.Collect()

	Collects[c.Name].Report(m)
}

func (c *Communication) Build(newC *Communication) {
	c.Child = append(c.Child, newC)
}

func (c *Communication) IncreaseCounter(key string) {
	c.Metric.IncreaseCounter(key)
}
func (c *Communication) AddCounter(key string, value int) {
	c.Metric.AddCounter(key, value)
}

type TaskGroupCommunicator struct {
}

func (t *TaskGroupCommunicator) Report(m *Metric) {
	fmt.Println(m)
}

func (t *TaskGroupCommunicator) Name() string {
	return cons.STAGETASKGROUP
}

type SchedulerCommunicator struct {
}

func (s *SchedulerCommunicator) Report(m *Metric) {
	fmt.Println(m)
}

func (t *SchedulerCommunicator) Name() string {
	return cons.STAGEJOB
}

type TaskCommunicator struct {
}

func (t *TaskCommunicator) Report(m *Metric) {
	fmt.Println(m)
}

func (t *TaskCommunicator) Name() string {
	return cons.STAGETASK
}

func init() {
	Register(&SchedulerCommunicator{})
	Register(&TaskGroupCommunicator{})
	Register(&TaskCommunicator{})
}
