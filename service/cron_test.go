package common

import (
	"fmt"
	"testing"
	"time"
)

const (
	Ten        = 10*time.Second + 50*time.Millisecond
	HalfMinute = 30*time.Second + 50*time.Millisecond
	TwoMinute = 120*time.Second + 50*time.Millisecond
)

type MockTask struct {
	name     string
	schedule string
	count    int
}

func NewMockTask(name, spec string) CronTask {
	return &MockTask{name: name, schedule: spec}
}

func (m *MockTask) GetSchedule() string {
	return m.schedule
}

func (m *MockTask) Run() {
	m.count++
	fmt.Printf("%s run #%d\n", m.name, m.count)
}

var (
	srv   = NewCronServiceImpl()
	task1 = NewMockTask("t1", "@every 1s")
	task2 = NewMockTask("t2", "@every 2s")
	task3 = NewMockTask("t3", "@every 3s")
)

func TestCronServiceImpl_AddTask(t *testing.T) {
	id1, _ := srv.AddTask(task1)
	fmt.Printf("task1 Id: %s\n", id1)
	id2, _ := srv.AddTask(task2)
	fmt.Printf("task2 Id: %s\n", id2)
	id3, _ := srv.AddTask(task3)
	fmt.Printf("task3 Id: %s\n", id3)
	select {
	case <-time.After(HalfMinute):
		fmt.Println("end")
	}
}

func TestCronServiceImpl_CancelTaskById(t *testing.T) {
	id,_  := srv.AddTask(task1)
	select {
	case <-time.After(Ten):
		_ = srv.CancelTask(id)
		fmt.Println("canceled")
	}
	select {
	case <-time.After(Ten):
		fmt.Println("end")
	}
}

func TestCronServiceImpl_SetTaskSchedule(t *testing.T) {
	id,_  := srv.AddTask(task1)
	fmt.Printf("task1 Id: %s\n", id)
	select {
	case <-time.After(Ten):
		_ = srv.SetTaskSchedule(id, "* * * * *") // 分 时 日 月 星期
		fmt.Println("set schedule")
	}
	select {
	case <-time.After(TwoMinute):
		fmt.Println("end")
	}
}

func TestCronServiceImpl_SetTaskSchedule_Wrong_Spec(t *testing.T) {
	id,_  := srv.AddTask(task1)
	fmt.Printf("task1 Id: %s\n", id)
	err := srv.SetTaskSchedule(id, "f**k off")
	if err != nil {
		t.Logf("wrong Spec err: %s\n", err.Error())
	} else {
		t.Fatalf("should be wrong")
	}
}

func TestCronServiceImpl_AddTask2(t *testing.T) {
	ch := make(chan string, 1)
	go func() {
		for i := 0; i < 20; i++ {
			name := fmt.Sprintf("【TASK%d】", i)
			id, _ := srv.AddTask(NewMockTask(name, "@every 2s"))
			fmt.Printf("%s added: %s\n", name, id)
			if i == 10 {
				ch <- id
			}
		}
	}()

	id := <-ch
	_ = srv.SetTaskSchedule(id, "@every 1s")

	select {
	case <-time.After(HalfMinute):
		fmt.Println("end")
	}
}
