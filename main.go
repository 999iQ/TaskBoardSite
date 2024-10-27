package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Name                 string
	Age                  uint16
	Money                int16
	AvgGrades, Happiness float64
	Hobbies              []string
}

func (u *User) getAllInfo() string { // метод структуры для вывода информации об объекте структуры
	return fmt.Sprintf("User name is: %s.\nHis age is: %d.\nAnd he has money equal: %d", u.Name, u.Age, u.Money)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}

func home_page(w http.ResponseWriter, r *http.Request) {
	// ResponseWriter - обращение к страничке
	// Request - отслеживание состояния
	bobik := User{"Bobik", 25, 10000, 4.5, 0.9, []string{"skate", "box", "yoga", "coding"}} // []string{"skate"}: Срез строк
	bobik.setNewName("Ne Bobik")
	tmpl, _ := template.ParseFiles("templates/home_page.html") // встроенный HTML шаблонизатор
	tmpl.Execute(w, bobik)
}

func login_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>The login page :D</>`)

}

func HandleRequests() {
	fmt.Println("Start up!")
	http.HandleFunc("/", home_page) // отслеживание перехода на URL адрес и вызов метода
	http.HandleFunc("/login/", login_page)
	http.ListenAndServe(":8080", nil) // отслеживание порта и передача настроек сервера

}

func main() {
	HandleRequests()

}
