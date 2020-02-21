package common

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"sync"
)

type CronTask interface {
	Run()
	GetSchedule() string
}

type CronService interface {
	AddTask(CronTask) (string, error)
	SetTaskSchedule(id string, spec string) error
	CancelTask(id string) error
}

type BaseCronTaskInfo struct {
	Id         string
	Spec       string
	LastUpdate string
	Done       chan struct{} `json:"-"`
}

func NewBaseCronTaskInfo(id, spec string) *BaseCronTaskInfo {
	info := &BaseCronTaskInfo{
		Id:         id,
		Spec:       spec,
		LastUpdate: "",
		Done:       make(chan struct{}, 1),
	}
	info.Done <- struct{}{}
	return info
}

func NewCronServiceImpl() *CronServiceImpl {
	s := new(CronServiceImpl)
	s.c = cron.New()
	s.c.Start()
	log.Println("CronService[robfig/cron/v3] started")
	return s
}

type CronServiceImpl struct {
	c     *cron.Cron
	tasks sync.Map
}

func (s *CronServiceImpl) AddTask(task CronTask) (string, error) {
	entryID, err := s.c.AddJob(task.GetSchedule(), task)
	if err != nil {
		return "", errors.New("CronSrv add task failed: " + err.Error())
	}
	id := uuid.NewV4().String()
	s.tasks.Store(id, entryID)
	return id, nil
}

func (s *CronServiceImpl) SetTaskSchedule(id string, spec string) error {
	v, ok := s.tasks.Load(id)
	if !ok {
		return errors.New("CronSrv set task failed, " +
			"cannot load by Id: " + id)
	}
	entryID := v.(cron.EntryID)
	job := s.c.Entry(entryID).Job // snapshot
	s.c.Remove(entryID)
	entryID, err := s.c.AddJob(spec, job)
	if err != nil {
		return errors.New("CronSrv set task failed, add job err: " + err.Error())
	}
	s.tasks.Store(id, entryID)
	return nil
}

func (s *CronServiceImpl) CancelTask(id string) error {
	v, ok := s.tasks.Load(id)
	if !ok {
		return errors.New("CronSrv cancel task failed, " +
			"cannot load by Id: " + id)
	}
	entryID := v.(cron.EntryID)
	s.c.Remove(entryID)
	s.tasks.Delete(id)
	return nil
}

type SimpleCronTask struct {
	run      func()
	schedule string
}

func NewSimpleCronTask(spec string, run func()) *SimpleCronTask {
	return &SimpleCronTask{
		run:      run,
		schedule: spec,
	}
}

func (s *SimpleCronTask) GetSchedule() string {
	return s.schedule
}

func (s *SimpleCronTask) Run() {
	s.run()
}
