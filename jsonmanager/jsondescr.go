package jsonmanager

import (
	"encoding/json"
	"fmt"
	"time"
)

type JsonPrintable interface {
	ToString() string
}

type Step struct {
	StepName string `json:"step_name"`
	Step     string `json:"step"`
	IsDone   bool   `json:"is_done"`
}

type LongTask struct {
	Id          int       `json:"id"`
	IsDone      bool      `json:"is_finished"`
	Name        string    `json:"name"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	FinishAt    time.Time `json:"finish_at"`
	Tags        []string  `json:"tags"`
	Description string    `json:"description"`
	Steps       []Step    `json:"steps"`
}

type Short struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsFinished  bool      `json:"is_finished"`
	CreatedAt   time.Time `json:"created_at"`
	FinishAt    time.Time `json:"finsh_at"`
}

type TaskList struct {
	Long  []LongTask `json:"long"`
	Short []Short    `json:"short"`
}

func (j TaskList) ToString() (string, error) {
	debugInfo, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return "", fmt.Errorf("could not MarshalIndent TaskList")
	}

	return string(debugInfo), nil
}
