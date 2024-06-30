package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTasks(response http.ResponseWriter, request *http.Request) {
	body, err := json.Marshal(tasks)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "applications/json")
	response.WriteHeader(http.StatusOK)

	_, err = response.Write(body)
	if err != nil {
		log.Printf("response writing error: %s\n", err.Error())
		return
	}
}

func getTask(response http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	task, ok := tasks[id]

	if !ok {
		http.Error(response, "unknown ID", http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(task)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Header().Set("Conten-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	_, err = response.Write(body)
	if err != nil {
		log.Printf("response writing error: %s\n", err.Error())
		return
	}
}

func addTask(response http.ResponseWriter, request *http.Request) {
	buffer := bytes.Buffer{}
	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	var task Task
	err = json.Unmarshal(buffer.Bytes(), &task)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	if _, ok := tasks[task.ID]; ok {
		http.Error(response, fmt.Sprintf("ID {%s} allready exists", task.ID), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	response.WriteHeader(http.StatusCreated)
}

func deleteTask(response http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	if _, ok := tasks[id]; !ok {
		http.Error(response, "unknown ID", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	response.WriteHeader(http.StatusOK)
}

func main() {
	file, err := os.OpenFile("precode.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

	r := chi.NewRouter()
	r.Get("/tasks", getTasks)
	r.Get("/task/{id}", getTask)
	r.Post("/tasks", addTask)
	r.Delete("/tasks/{id}", deleteTask)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
