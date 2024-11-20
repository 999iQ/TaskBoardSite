package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type user struct {
	name  string
	email string
}

func home_page(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home_page.html"))
	tmpl.ExecuteTemplate(w, "home_page.html", nil)
}

func add_deadline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // например сразу перешли по адресу отправки данных кнопкой при нажатии, но без нажатия
		http.Error(w, "Метод запроса не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	datetime := r.FormValue("datetime-input") // получение даты дедлайна из формы
	taskName := r.FormValue("task-name")
	taskDescription := r.FormValue("task-description")

	w.Write([]byte(fmt.Sprintf("Добавлен новый дедлайн.\nДата конца: %s\nНазвание: %s\nОписание: %s", datetime, taskName, taskDescription))) // ответ сервера
	log.Printf("Добавлен новый дедлайн.\nДата конца: %s\nНазвание: %s\nОписание: %s", taskName, datetime, taskDescription)                   // вывод в консоль
}

func handle_request() {
	// шарим папку static, чтобы передать js скрипт в этой папке
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", home_page)
	http.HandleFunc("/process", add_deadline)
	log.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // настройка логгирования для вывода строки в коде
	handle_request()
}
