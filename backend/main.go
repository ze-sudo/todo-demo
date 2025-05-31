package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/cors"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var todos = []Todo{
	{ID: 1, Task: "サンプルタスク", Done: false},
}
var nextID = 2

// getTodos は全てのToDoを取得するハンドラ
func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// addTodo は新しいToDoを追加するハンドラ
func addTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	todo.ID = nextID
	nextID++
	todos = append(todos, todo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// updateTodo は既存のToDoを更新するハンドラ
func updateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// クエリパラメータからIDを取得（例: /todos?id=1）
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	// リクエストボディをデコード
	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 該当するToDoを探して更新
	updated := false
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Task = updatedTodo.Task
			todos[i].Done = updatedTodo.Done
			updated = true
			break
		}
	}

	if !updated {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// クエリパラメータからIDを取得（例: /todos?id=1）
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	// ToDoを探して削除
	deleted := false
	for i, todo := range todos {
		if todo.ID == id {
			// 削除: i番目を取り除く
			todos = append(todos[:i], todos[i+1:]...)
			deleted = true
			break
		}
	}

	if !deleted {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// todosHandler はToDoのCRUD操作を処理するハンドラ
func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTodos(w, r)
	case "POST":
		addTodo(w, r)
	case "PUT":
		updateTodo(w, r)
	case "DELETE":
		deleteTodo(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


// main はHTTPサーバを起動する
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/todos", todosHandler)
	handler := cors.AllowAll().Handler(mux)
	http.ListenAndServe(":8080", handler)
}