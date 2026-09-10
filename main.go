
package main

import (
	"todo-tui/jsonmanager"
)

func main() {
	fileName := "artifacts/todo.json"
	jsonmanager.ReadJson(fileName)
}
