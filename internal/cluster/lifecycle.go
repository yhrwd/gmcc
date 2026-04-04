package cluster

type StartTrigger string

const (
	StartTriggerManualStart   StartTrigger = "manual_start"
	StartTriggerManualRestart StartTrigger = "manual_restart"
	StartTriggerAutoReconnect StartTrigger = "auto_reconnect"
)

type ExitCategory string

const (
	ExitCategoryManualStop        ExitCategory = "manual_stop"
	ExitCategoryNetworkDisconnect ExitCategory = "network_disconnect"
	ExitCategoryStartupTimeout    ExitCategory = "startup_timeout"
	ExitCategoryAuthFailed        ExitCategory = "auth_failed"
	ExitCategoryUnknown           ExitCategory = "unknown"
)

func canTransition(from, to InstanceStatus) bool {
	allowed := map[InstanceStatus]map[InstanceStatus]bool{
		StatusPending:      {StatusStarting: true},
		StatusStarting:     {StatusRunning: true, StatusError: true, StatusStopped: true},
		StatusRunning:      {StatusReconnecting: true, StatusError: true, StatusStopped: true},
		StatusReconnecting: {StatusStarting: true, StatusError: true, StatusStopped: true},
		StatusError:        {StatusStarting: true, StatusStopped: true},
		StatusStopped:      {StatusStarting: true},
	}

	return allowed[from][to]
}
