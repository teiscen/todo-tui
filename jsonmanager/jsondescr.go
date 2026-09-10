package jsonmanager

import (
	"fmt"
	"strings"
	"time"
)

type Step struct {
	StepName string `json:"step_name"`
	Step     string `json:"step"`
	IsDone   bool   `json:"is_done"`
}

func (j Step) ToString() (str string) { // Named return output added
	str += "\t{\n"
	str += fmt.Sprintf("\t\tStepName: %s\n", j.StepName)
	str += fmt.Sprintf("\t\tStep: %s\n", j.Step)
	str += fmt.Sprintf("\t\tIsDone: %t\n", j.IsDone)
	str += "\t}\n"
	return str
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

func (j LongTask) ToString() (str string) { // Named return output added
	str += "{\n"
	str += fmt.Sprintf("\tId: %d\n", j.Id)
	str += fmt.Sprintf("\tIsDone: %t\n", j.IsDone)
	str += fmt.Sprintf("\tName: %s\n", j.Name)
	str += fmt.Sprintf("\tPriority: %d\n", j.Priority)
	str += fmt.Sprintf("\tCreatedAt: %v\n", j.CreatedAt)
	str += fmt.Sprintf("\tFinishAt: %v\n", j.FinishAt)
	str += fmt.Sprintf("\tDescription: %s\n", j.Description)
	
	str += "\tSteps: [\n"
	for _, step := range j.Steps {
		str += step.ToString()
	}
	str += "\t]\n"

	str += fmt.Sprintf("\tTags: [%s]\n", strings.Join(j.Tags, ", "))
	str += "}\n"
	return str
}

type Short struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsFinished  bool      `json:"is_finished"`
	CreatedAt   time.Time `json:"created_at"`
	FinishAt    time.Time `json:"finsh_at"`
}

func (j Short) ToString() (str string) { // Named return output added
	str += "{\n"
	str += fmt.Sprintf("\tId: %d\n", j.Id)
	str += fmt.Sprintf("\tName: %s\n", j.Name)
	str += fmt.Sprintf("\tDescription: %s\n", j.Description)
	str += fmt.Sprintf("\tIsFinished: %t\n", j.IsFinished)
	str += fmt.Sprintf("\tCreatedAt: %v\n", j.CreatedAt)
	str += fmt.Sprintf("\tFinishAt: %v\n", j.FinishAt)
	str += "}\n"
	return str
}

type TaskList struct {
	Long  []LongTask `json:"long"`
	Short []Short    `json:"short"`
}

func (j TaskList) ToString() (str string) { // Named return output added
	for _, long := range j.Long {
		str += long.ToString()
	}
	for _, short := range j.Short {
		str += short.ToString()
	}
	return str
}








