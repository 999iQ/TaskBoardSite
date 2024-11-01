package main

import (
	"fmt"
	"html/template"
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

func addDeadlineButtonPress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // например сразу перешли по адресу отправки данных кнопкой при нажатии, но без нажатия
		http.Error(w, "Метод запроса не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	datetime := r.FormValue("datetime-input") // получение данных из формы
	w.Write([]byte(fmt.Sprintf("Кнопка была нажата! %s", datetime)))
	fmt.Println("Кнопка была нажата! Дата и время:", datetime)
}

func handleRequest() {
	// шарим папку static, чтобы передать js скрипт
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", home_page)
	http.HandleFunc("/process", addDeadlineButtonPress)
	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
}
