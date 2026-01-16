package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
)

type Task struct {
	Name      string    `json:"name"`
	IsDone    bool      `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
}

var tasks []Task

func main() {
	// Проверяем существования файла tasks.json, если файла нет - создаём

	fileInfo, err := os.Lstat("./tasks.json")
	if fileInfo == nil || err != nil {
		// WriteFile - записывает в файл новые строки, а если файла нет - создаст автоматически.
		// Конкретно тут мы записываем в tasks.json пустой массив, чтобы при десериализации
		// функция json.Unmarshal не ругалась на то, что в файле не существует json формата

		os.WriteFile("./tasks.json", []byte("[]"), 0644)
	}

	// Читаем содержимое tasks.json в байтах

	jsonBytes, err := os.ReadFile("./tasks.json")
	if err != nil {
		log.Fatal(err)
	}

	// Десериализируем байты из нашего файла в структуру Task и добавляем последние задачи из файла в наш слайс тасков

	var lastDataJson []Task
	err = json.Unmarshal(jsonBytes, &lastDataJson)
	if err != nil {
		log.Fatal(err)
	}
	tasks = append(tasks, lastDataJson...)

	// Основная логика выбора команды

	flags := os.Args
	if len(flags) == 1 {
		log.Fatal("You didn't enter a command. To see all commands, type 'help'")
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) == 3 {
			Create(os.Args[2])
		} else {
			log.Fatal("To create a task, use: create <name>")
		}
	case "delete":
		if len(os.Args) == 3 {
			Delete(os.Args[2])
		} else {
			log.Fatal("To delete a task, use: delete <task ID | all>")
		}
	case "done":
		if len(os.Args) == 3 {
			Done(os.Args[2])
		} else {
			log.Fatal("To mark a task as done, use: done <task ID>")
		}
	case "undone":
		if len(os.Args) == 3 {
			Undone(os.Args[2])
		} else {
			log.Fatal("To unmark a task as done, use: undone <task ID>")
		}
	case "list":
		if len(tasks) == 0 {
			log.Println("You don't have any tasks yet. Add one!")
		} else {
			List()
		}
	case "help":
		fmt.Println("Create a task: create <name>\nMark a task as done: done <task ID>\nUnmark a task as done: undone <task ID>\nDelete a task: delete <task ID | all>\nList all tasks: list\nShow this message: help")
	default:
		log.Fatal("Unknown command! To see available commands, type 'help'")
	}
}

func Create(name string) {
	// Создаём новую задачу по структуре Task

	newTask := Task{
		Name:      name,       // То что мы передали в функцию название (os.Args[2])
		IsDone:    false,      // Выполнена ли задача, по дефолту - false
		CreatedAt: time.Now(), // Время создания - сейчас
	}

	// Пополняем слайс новой задачей, а после сериализируем наш слайс задач в json формат и записываем через os.WriteFile

	tasks = append(tasks, newTask)

	tasksBytes, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./tasks.json", tasksBytes, 0755)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task '%s' added!", name)
}

func Delete(idStr string) {
	// Если пользователь вместо ID задачи написал all, то удаляем все задачи

	if idStr == "all" {
		tasks = nil

		tasksBytes, err := json.Marshal(tasks)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile("./tasks.json", tasksBytes, 0755)
		if err != nil {
			log.Fatal(err)
		}

		log.Print("All tasks deleted!")

		return
	}

	// Переводим строку в число, после чего заранее присваиваем название задачи которую мы удаляем

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	taskName := tasks[id].Name

	// Удаление задачи (берём все что идёт до этого ID задачи и после этого ID задачи и соединяем - получаем слайс без этой задачи)

	tasks = append(tasks[:id], tasks[id+1:]...)

	tasksBytes, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./tasks.json", tasksBytes, 0755)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task '%s' deleted!", taskName) // Говорим что удалили конкретную задачу, а не ID, так UX выглядит как-будто лучше.
}

func Done(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	taskName := tasks[id].Name

	tasks[id].IsDone = true

	tasksBytes, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./tasks.json", tasksBytes, 0755)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task '%s' marked as done!", taskName)
}

func Undone(idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	taskName := tasks[id].Name

	tasks[id].IsDone = false

	tasksBytes, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./tasks.json", tasksBytes, 0755)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task '%s' unmarked as done!", taskName)
}

func List() {
	for i, task := range tasks {
		done := "❌"
		if task.IsDone == true {
			done = "✅"
		}
		fmt.Printf("%d) [%s] %s | %v\n", i, done, task.Name, humanize.Time(task.CreatedAt)) // humanize - переводит время в человеческий формат (3 weeks ago, 1 hours ago и т.д.)
	}
}
