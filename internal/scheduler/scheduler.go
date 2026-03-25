// Scheduler is responsible for starting the execution of tasks at
// their specified schedule. It uses a cron job package to achieve it.
package scheduler

import "time"

type Scheduler interface {
	Start()
}

type taskFunction func()

type ScheduledTask struct {
	Schedule time.Duration
	Task     taskFunction
}
