package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/enescakir/emoji"
)

type Task struct {
	Title       string
	Completed   bool
	Created     time.Time
	CompletedAt time.Time
}

type TaskList []Task

func main() {
	if len(os.Args) <= 1 {
		panic("Please enter a flag")
	}

	add := flag.Bool("add", false, "Add a task")
	list := flag.Bool("list", false, "List all tasks")
	delete := flag.Int("delete", 0, "Delete a task")
	complete := flag.Int("complete", 0, "Mark a task as completed")

	flag.Parse()

	file, _ := os.ReadFile("data.json")

	var result TaskList

	json.Unmarshal(file, &result)

	tasks := result

	switch {
	case *add:
		taskName := getInputs(flag.Args())
		tasks.Add(taskName)
		tasks.DrawList()
	case *list:
		tasks.DrawList()
	case *delete > 0:
		tasks.Delete(*delete)
	case *complete > 0:
		tasks.Complete(*complete)
		tasks.DrawList()
	}

	writeFile(&tasks)
}

func (t *TaskList) Add(title string) {
	todo := Task{
		title,
		false,
		time.Now(),
		time.Time{},
	}
	*t = append(*t, todo)
}

func (t *TaskList) Complete(index int) {
	if index > len(*t) {
		fmt.Println("Task at that index does not exist")
		os.Exit(2)
	}
	list := *t
	// Handle index not existing
	list[index-1].Completed = true
	list[index-1].CompletedAt = time.Now()
}

func (t *TaskList) Delete(index int) {
	if index > len(*t) {
		fmt.Println("Task at that index does not exist")
		os.Exit(2)
	}
	list := *t
	s := fmt.Sprintf("Task - %s has been deleted", list[index-1].Title)
	fmt.Println(s)
	// Handle index not existing
	*t = append(list[:index-1], list[index:]...)
}

func (t *TaskList) DrawList() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Time Created"},
			{Align: simpletable.AlignCenter, Text: "Completed"},
			{Align: simpletable.AlignCenter, Text: "Time Completed"},
		},
	}
	var cells [][]*simpletable.Cell

	for index, task := range *t {
		// Increase index so it starts at 1
		index++

		var completedIcon emoji.Emoji

		if task.Completed {
			completedIcon = emoji.CheckMarkButton
		} else {
			completedIcon = emoji.CrossMark
		}

		var timeCompleted string
		notCompletedTime := time.Time{}
		if task.CompletedAt == notCompletedTime {
			timeCompleted = "N/A"

		} else {
			timeCompleted = task.CompletedAt.Format(time.RFC822)
		}
		cells = append(cells, *&[]*simpletable.Cell{

			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", index)},
			{Text: task.Title},
			{Align: simpletable.AlignCenter, Text: task.Created.Format(time.RFC822)},
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%v", completedIcon)},
			{Align: simpletable.AlignCenter, Text: timeCompleted},
		})
	}
	table.Body = &simpletable.Body{Cells: cells}
	table.SetStyle(simpletable.StyleUnicode)
	fmt.Println(table.String())

}

func getInputs(args []string) string {
	if len(args) <= 0 {
		panic("Please provide task name")
	}
	taskName := strings.Join(args, " ")

	return taskName
}

func writeFile(tasks *TaskList) {
	json, err := json.Marshal(*tasks)

	if _, err := os.Stat("data.json"); errors.Is(err, os.ErrNotExist) {
		os.Create("data.json")
	}
	if err != nil {
		panic("Error saving tasks")
	}
	os.WriteFile("data.json", json, 0644)

}
