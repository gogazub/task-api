package model

type ExecutionResult struct {
	TaskID   string `json:"task_id"`
	Status   string `json:"status"` // pending/running/completed/failed
	Output   string `json:"output"`
	Error    string `json:"error"`
	ExitCode int    `json:"exit_code"`
	Duration int    `json:"duration_ms"`
}
