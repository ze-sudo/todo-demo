'use client';
import { useEffect, useState } from "react";
import "./page.css";

type Todo = { id: number; task: string; done: boolean };

export default function Home() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTask, setNewTask] = useState("");
  const [editId, setEditId] = useState<number | null>(null);
  const [editTask, setEditTask] = useState("");
  const [editDone, setEditDone] = useState(false);

  // ToDo一覧取得
  const fetchTodos = () => {
    fetch("http://localhost:8080/todos")
      .then(res => res.json())
      .then(setTodos);
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  // 追加
  const addTodo = async () => {
    if (!newTask) return;
    await fetch("http://localhost:8080/todos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ task: newTask, done: false }),
    });
    setNewTask("");
    fetchTodos();
  };

  // 更新
  const updateTodo = async (id: number) => {
    await fetch(`http://localhost:8080/todos?id=${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ task: editTask, done: editDone }),
    });
    setEditId(null);
    fetchTodos();
  };

  // 削除
  const deleteTodo = async (id: number) => {
    await fetch(`http://localhost:8080/todos?id=${id}`, { method: "DELETE" });
    fetchTodos();
  };


  return (
    <div>
      <h1>ToDoリスト(demo)</h1>
      <input
        value={newTask}
        onChange={e => setNewTask(e.target.value)}
        placeholder="新しいタスク"
      />
      <button onClick={addTodo}>追加</button>
      <ul>
        {todos.map(todo =>
          editId === todo.id ? (
            <li key={todo.id}>
              <input
                value={editTask}
                onChange={e => setEditTask(e.target.value)}
              />
              <label>
                <input
                  type="checkbox"
                  checked={editDone}
                  onChange={e => setEditDone(e.target.checked)}
                />
                完了
              </label>
              <button onClick={() => updateTodo(todo.id)}>保存</button>
              <button onClick={() => setEditId(null)}>キャンセル</button>
            </li>
          ) : (

            <li key={todo.id}>
              {todo.task} {todo.done ? "✔️" : ""}
              <button
                onClick={() => {
                  setEditId(todo.id);
                  setEditTask(todo.task);
                  setEditDone(todo.done);
                }}
              >
                編集
              </button>
              <button onClick={() => deleteTodo(todo.id)}>削除</button>
            </li>
          )
        )}
      </ul>
    </div>
  );
}