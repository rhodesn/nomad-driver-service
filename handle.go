package main

import (
	"context"
	"sync"
	"time"

	systemd "github.com/coreos/go-systemd/v22/dbus"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/plugins/drivers"
)

// taskHandle should store all relevant runtime information
// such as process ID if this is a local task or other meta
// data if this driver deals with external APIs
type taskHandle struct {
	// conn
	conn *systemd.Conn

	// unitName
	unitName string

	// collectionInterval
	collectionInterval time.Duration

	// stateLock syncs access to all fields below
	stateLock sync.RWMutex

	logger hclog.Logger

	taskConfig  *drivers.TaskConfig
	procState   drivers.TaskState
	startedAt   time.Time
	completedAt time.Time
	exitResult  *drivers.ExitResult
}

func (h *taskHandle) TaskStatus() *drivers.TaskStatus {
	h.stateLock.RLock()
	defer h.stateLock.RUnlock()

	return &drivers.TaskStatus{
		ID:               h.taskConfig.ID,
		Name:             h.taskConfig.Name,
		State:            h.procState,
		StartedAt:        h.startedAt,
		CompletedAt:      h.completedAt,
		ExitResult:       h.exitResult,
		DriverAttributes: map[string]string{},
	}
}

func (h *taskHandle) IsRunning() bool {
	h.stateLock.RLock()
	defer h.stateLock.RUnlock()
	return h.procState == drivers.TaskStateRunning
}

func (h *taskHandle) monitor() {
	h.stateLock.Lock()
	if h.exitResult == nil {
		h.exitResult = &drivers.ExitResult{}
	}
	h.stateLock.Unlock()

	h.logger.Debug("monitoring unit", "name", h.unitName)
	var state drivers.TaskState
	for {
		timerChan := time.After(h.collectionInterval)
		// We can't use SubscribeUnitsCustom because it only returns either
		// active or enabled units, not inactive ones.
		units, err := h.conn.ListUnitsByNamesContext(context.TODO(), []string{h.unitName})
		if err != nil {
			h.logger.Error("encountered error monitoring unit", err)
			state = drivers.TaskStateUnknown
			break
		}

		if units[0].ActiveState == "inactive" {
			h.logger.Debug("unit inactive", "name", h.unitName)
			state = drivers.TaskStateExited
			break
		}

		<-timerChan
	}

	h.stateLock.Lock()
	defer h.stateLock.Unlock()
	h.procState = state
	h.completedAt = time.Now().Round(time.Millisecond)
}
