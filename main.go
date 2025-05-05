package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
)

type Task struct {
    ID          int    `json:"id"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

func main() {
    taskDesc := flag.String("desc", "", "Task description")
    complete := flag.Int("complete", -1, "Task ID to mark as completed")
    list := flag.Bool("list", false, "List all tasks")

    flag.Parse()

    tasks, err := loadTasks("tasks.json")
    if err != nil {
        fmt.Println("Error loading tasks:", err)
        os.Exit(1)
    }

    if *taskDesc != "" {
        tasks = append(tasks, Task{ID: len(tasks) + 1, Description: *taskDesc})
    } else if *complete != -1 {
        for i := range tasks {
            if tasks[i].ID == *complete {
                tasks[i].Completed = true
                break
            }
        }
    } else if *list {
        for _, task := range tasks {
            status := "pending"
            if task.Completed {
                status = "completed"
            }
            fmt.Printf("[%d] %s - %s\n", task.ID, task.Description, status)
        }
    }

    err = saveTasks("tasks.json", tasks)
    if err != nil {
        fmt.Println("Error saving tasks:", err)
        os.Exit(1)
    }
}

func loadTasks(filename string) ([]Task, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil && !os.IsNotExist(err) {
        return nil, err
    }

    var tasks []Task
    if len(data) > 0 {
        err = json.Unmarshal(data, &tasks)
        if err != nil {
            return nil, err
        }
    }

    return tasks, nil
}

func saveTasks(filename string, tasks []Task) error {
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return err
    }

    return ioutil.WriteFile(filename, data, 0644)
}
