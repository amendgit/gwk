package views

import (
	"time"
)

type event_pump_delegate_t interface {
	DoWork() bool
	DoDelayedWork() bool
}

type event_pump_t interface {
	Run()
	Quit()
	ScheduleWork()
	ScheduleDelayedWork(time time.Time)
}
