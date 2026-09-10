package jsonmanager

import (
	"os"
	"encoding/json"
	"fmt"
)

func ReadJson(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic("Error reading file!")
		return err
	}

	var tasks TaskList
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		panic("Error unmarshalling file")
		return err
	}

	fmt.Printf(tasks.ToString())
	return nil
}
