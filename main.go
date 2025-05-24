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
	Selected bool `json:"selected,omitempty"`
	ToDelete bool `json:"toDelete"`
}

const (
	ADD = "add"
	LIST = "list"
	FILE_NAME = "tasks.json"
	UP = 'k' 
	DOWN = 'j' 
	DELETE = 'x' 
	SELECT = 13
	QUIT = 'q' 
)

const (
	READ = iota
	WRITE
	MARSHAL
	UNMARSHAL
)

var selectedIndex int = 0
var cursor = ""
var checked = ""
var tasks []Task = make([]Task, 0)


func main() {
	
	command := os.Args[1]

	switch command {
	case ADD:
		taskTitle := os.Args[2]
		taskPriority := os.Args[3]
		add(taskTitle, taskPriority)
	case LIST:
		list()
	}

}

func list() {

	file := readFile()

	err := json.Unmarshal(file, &tasks)

	if err != nil {
		error(err, UNMARSHAL)
	}

	
	displayInteractableList()
	removeSelectedAndDeletedTasks()
}

func removeSelectedAndDeletedTasks() {
	clonedTasks := make([]Task , 0)
	for _, task := range clonedTasks {

		if task.Selected || task.ToDelete {
			continue
		}
		clonedTasks = append(clonedTasks, task)
	}
	tasks = clonedTasks
	writeTasksToFile()
}

func displayInteractableList() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)


	redraw(tasks)

	for {
		os.Stdin.Read(buf)
		userInput := buf[0]
		if userInput == QUIT {
			break
		}
		switch userInput {
		
		case UP:
			if selectedIndex > 0 {
				selectedIndex--
			}	
		case DOWN:
			if selectedIndex < len(tasks) - 1 {
				selectedIndex++
			}
		case SELECT: 
			tasks[selectedIndex].Selected = !tasks[selectedIndex].Selected
		case DELETE:
			tasks[selectedIndex].ToDelete = !tasks[selectedIndex].ToDelete;
		}

		redraw(tasks)
	}
}

func redraw(tasks []Task) {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	if len(tasks) == 0 {
		fmt.Printf("Task List is currently Emtpy")
	}
	for i, task := range tasks {
		
		if i == selectedIndex {
			cursor = "> "
		} else {
			cursor = ""
		}

		if task.Selected {
			checked = "[x]"
		} else {
			checked = "[]"
		}
		
		displayTask(task, i)
	}
}

func displayTask(task Task, idx int) {
	if task.ToDelete {
		fmt.Printf("\u001b[31m%s   %s   %s   %s\u001b[0m", cursor, checked, task.Title, task.Priority)
	}	else {
		fmt.Printf("%s   %s   %s   %s", cursor, checked, task.Title, task.Priority)
	}

	if (len(tasks) - 1 != idx) {
		fmt.Print("\r\n")
	}
}


func add(title, priority string) {

	file := readFile()

	newTask := Task{priority, title, false, false}

	err := json.Unmarshal(file, &tasks)
	if err != nil {
		error(err, UNMARSHAL)
	}

	tasks = append(tasks, newTask)
	
	fmt.Printf("tasks: %+v", tasks)

	writeTasksToFile()
}


func writeTasksToFile() {
	data, err := json.Marshal(tasks)	
	if err != nil {
		error(err, MARSHAL)
	}

	err = os.WriteFile(FILE_NAME, data, fs.FileMode(0664))
	if err != nil {
		error(err, WRITE)
	}
}

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
