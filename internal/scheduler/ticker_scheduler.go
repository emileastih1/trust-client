package scheduler

import (
	"time"

	"go.uber.org/zap"
)

type TickerScheduler struct {
	logger *zap.Logger
	tasks  []*ScheduledTask
}

func NewTickerScheduler(logger *zap.Logger, tasks []*ScheduledTask) *TickerScheduler {
	return &TickerScheduler{
		tasks:  tasks,
		logger: logger,
	}
}

// Start the execution of the tasks in their specified schedule.
func (s TickerScheduler) Start() error {
	for _, task := range s.tasks {
		ticker := time.NewTicker(task.Schedule)
		// TODO find a way to clean up the routines in case some have a completion
		go func(t *ScheduledTask) {
			for range ticker.C {
				t.Task()
			}
		}(task)
	}
	return nil
}
