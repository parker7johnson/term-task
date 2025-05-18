package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"golang.org/x/term"
)

type Task struct {
	Priority string `json:"priority"`
	Title    string `json:"title"`
	DeadLine string `json:"deadLine,omitempty"`
	Selected bool `json:"selected, omitemtpy"`
}

const (
	FILE_NAME = "tasks.json"
)

const (
	READ = iota
	WRITE
	MARSHAL
	UNMARSHAL
)

var selectedIndex int = 0

func readFile() []byte {
	file, err := os.ReadFile(FILE_NAME)
	if err != nil {
		error(err, READ)
	}
	return file
}

func error(err any, messageType int) {
	
	switch messageType {
	case MARSHAL:
		fmt.Printf("Error marshaling the new task list: %v", err)
	case UNMARSHAL:
		fmt.Printf("Error when unmarshingaling tasks.json: %v", err)
	case READ:
		fmt.Printf("Error when trying to read tasks.json: %v", err)
	case WRITE:
		fmt.Printf("Error when trying to write to tasks.json: %v", err)
	default:
		fmt.Printf("Unknown error: %v", err)
	}
	os.Exit(1)
}
//task add laundry normal tomorrow
//task list
//[] 1. Laundry by tomorrow
func main() {
	add("laundry", "normal", "tomorrow")
	add("dishes", "normal", "tomorrow")
  list()

}

func list() {

	file := readFile()
	var tasks []Task

	err := json.Unmarshal(file, &tasks)

	if err != nil {
		error(err, UNMARSHAL)
	}

	
	displayInteractableList(tasks)
}

func displayInteractableList(tasks []Task) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	//this is the same as the code above but manually listens for the ctrl c keys
	buf := make([]byte, 3)


	redraw(tasks)
	fmt.Print(selectedIndex)

	for {
		os.Stdin.Read(buf)
		fmt.Print(buf)	

		if buf[0] == 3 {
			term.Restore(int(os.Stdin.Fd()), oldState)
		}

		if buf[0] == 27 && buf[1] == 91 {
			if buf[2] == 65 {
				if selectedIndex > 0 {
					selectedIndex--
				}
			} else if buf[2] == 66 {
				if selectedIndex < len(tasks) {
					selectedIndex++
				}
			}
		} else if buf[0] == 13 {
			tasks[selectedIndex].Selected = !tasks[selectedIndex].Selected
		}
		redraw(tasks)
	}
  
	
}

func redraw(tasks []Task) {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	checked := "[]"
	cursor := ""
	for i, task := range tasks {
		
		if i == selectedIndex {
			cursor = "> "
		} else {
			cursor = ""
		}

		if task.Selected {
			checked = "[x]"
		}
		if (i == len(tasks) - 1) {
			fmt.Printf("%s %s %s %s %s", cursor, checked, task.Title, task.Priority, task.DeadLine)
		} else {
			fmt.Printf("%s %s %s %s %s\r\n", cursor, checked, task.Title, task.Priority, task.DeadLine)
		}
	}
}


func add(title, priority, deadline string) {

	file := readFile()

	newTask := Task{priority, title, deadline, false}
	var tasks []Task 

	err := json.Unmarshal(file, &tasks)
	if err != nil {
		error(err, UNMARSHAL)
	}

	tasks = append(tasks, newTask)
	

	data, err := json.Marshal(tasks)	

	if err != nil {
		error(err, MARSHAL)
	}

	err = os.WriteFile(FILE_NAME, data, fs.FileMode(0664))
	if err != nil {
		error(err, WRITE)
	}
}
