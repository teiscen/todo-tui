package jsonmanager

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJson(filePath string) (TaskList, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return TaskList{}, fmt.Errorf("could not read file")
	}

	var tasks TaskList

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return TaskList{}, fmt.Errorf("could not unmarshal %s", filePath)
	}

	return tasks, nil
}

// explanation of the 0644 flag for the third arg of os.WriteFile
// https://projectai.in/projects/dd14b0b6-1d16-4694-af16-e58995181bca/tasks/09557118-cf63-4688-943f-e01ee6664570
func WriteJson(filePath string, tasks TaskList) error {
	b, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("could not marshal tasks")
	}

	err = os.WriteFile(filePath, b, 0o644)
	if err != nil {
		return fmt.Errorf("could not write to file %s", filePath)
	}
	return nil
}
