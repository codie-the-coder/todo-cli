package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Task struct {
	Text      string
	Completed bool
}

var configDir, _ = os.UserConfigDir()
var savePath = filepath.Join(configDir, "todo-cli", "todo.json")

const helpText string = `some todolist written by codie the coder as first project in golang
usage:
	todo list - get all tasks
	todo add "text" - add task
	todo remove <task id> - remove task
	todo clear - clear every task
	todo complete <task id> - mark task as completed
	todo uncomplete <task id> - mark task as not completed
`

var taskList []Task // потом используется

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			fmt.Println(helpText)
		case "list":
			fetchAndPrintTaskList()
		case "add":
			addTask(argsAfterCommand())
		case "remove":
			removeTask(argsAfterCommand())
		case "complete":
			completeTask(argsAfterCommand())
		case "uncomplete":
			unCompleteTask(argsAfterCommand())
		default:
			fmt.Println("unknown command:", os.Args[1])
		}
	} else {
		fmt.Println(helpText)
	}
}

func argsAfterCommand() string {
	if len(os.Args) <= 2 {
		return ""
	}
	return strings.Join(os.Args[2:], " ")
}

func fetchAndPrintTaskList() {
	loadJson()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "[ID]\t[Task]\t[Completed]")

	for i, task := range taskList {
		fmt.Fprintf(w, "%v\t%v\t%v\n", i, task.Text, boolToYesNo(task.Completed))
	}

	w.Flush()
}

func addTask(taskText string) {
	loadJson()

	task := Task{
		Text:      taskText,
		Completed: false,
	}

	taskList = append(taskList, task)

	saveJson()
}

func removeTask(taskId string) {
	loadJson()

	num64, err := strconv.ParseUint(taskId, 10, 64)

	if err != nil {
		fmt.Println("negative numbers is not allowed")
		return
	}

	id := int(num64)

	if id >= len(taskList) {
		fmt.Println("there is no task with that id")
	} else {
		taskList = slices.Delete(taskList, id, id+1)
		saveJson()
	}
}

func completeTask(taskId string) {
	loadJson()

	num64, err := strconv.ParseUint(taskId, 10, 64)

	if err != nil {
		fmt.Println("negative numbers is not allowed")
		return
	}

	id := int(num64)

	if id >= len(taskList) {
		fmt.Println("there is no task with that id")
	} else if taskList[id].Completed != true {
		taskList[id].Completed = true
		saveJson()
	} else {
		fmt.Println("task is already completed")
	}
}

func unCompleteTask(taskId string) {
	loadJson()

	num64, err := strconv.ParseUint(taskId, 10, 64)

	if err != nil {
		fmt.Println("negative numbers is not allowed")
		return
	}

	id := int(num64)

	if id >= len(taskList) {
		fmt.Println("there is no task with that id")
	} else if taskList[id].Completed != false {
		taskList[id].Completed = false
		saveJson()
	} else {
		fmt.Println("task is not completed")
	}
}

func loadJson() {
	if checkFile() != true {
		return
	}

	data, _ := os.ReadFile(savePath)
	json.Unmarshal(data, &taskList)
}

func saveJson() {
	if !checkFile() {
		return
	}

	jsonData, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		fmt.Println("json serialization err:", err)
		return
	}

	if err := os.WriteFile(savePath, jsonData, 0644); err != nil {
		fmt.Println("file write err:", err)
	}
}

func checkFile() bool {
	_, err := os.Stat(savePath)

	if err == nil {
		return true
	}

	if !errors.Is(err, os.ErrNotExist) {
		return false
	}

	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false
	}

	file, err := os.Create(savePath)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}