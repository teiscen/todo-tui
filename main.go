package main

import (
	"fmt"
	"log"

	"todo-tui/jsonmanager"
	"todo-tui/testing"
)

func testJSON() {
	filePathRead := "artifacts/todoread.json"
	tasks, err := jsonmanager.ReadJson(filePathRead)
	if err != nil {
		log.Fatal(err)
	}

	taskStr, err := tasks.ToString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(taskStr)

	filePathWrite := "artifacts/todowrite.json"
	err = jsonmanager.WriteJson(filePathWrite, tasks)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// testJSON()
	// tui.TestTUI()
	testing.TestingMain()
	// calendar.TestCalendar()
}
