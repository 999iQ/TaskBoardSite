package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

type User struct {
	User_id      uint      `gorm:"primaryKey"`
	Nickname     string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	Passwordhash string    `gorm:"not null"`
	Createdat    time.Time `gorm:"autoCreateTime"`
	Updatedat    time.Time `gorm:"autoUpdateTime"`
}

type Task struct {
	Task_id      uint   `gorm:"primaryKey"`
	User_id      uint   `gorm:"foreignKey"`
	Title        string `gorm:"not null"`
	Description  string
	Priority     uint      `gorm:"default:0"`
	Deadlinedate time.Time `gorm:"not null"`
	Createdat    time.Time `gorm:"autoCreateTime"`
	Updatedat    time.Time `gorm:"autoUpdateTime"`
}

func homePage(ctx *fiber.Ctx) error {
	return ctx.Render("templates/homePage.html", nil)
}

func addDeadline(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodPost { // перешли на api по адресу отправки данных (с методом get)
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "Метод запроса не поддерживается!",
		})
	}

	taskName := ctx.FormValue("task-name")
	taskDescription := ctx.FormValue("task-description")
	deadline, _ := time.Parse("2006-01-02T15:04", ctx.FormValue("datetime-input")) // получение даты дедлайна из формы // time.RFC3339 используется для парсинга строк в формате ISO 8601
	var priority uint
	fmt.Sscan(ctx.FormValue("priority"), &priority) // приведение строки к uint

	// сохранение полученных данных о дедлайне в БД
	db := connectDB()
	defer db.Close() // отложенное отклчючение от бд (пока не вышли из текущей функции)

	var lastTask Task
	db.Table("tasks").Last(&lastTask) // получение последней записи (обратная сортировка {desk} по id)
	newUser := Task{
		Task_id:      lastTask.Task_id + 1,
		User_id:      1,
		Title:        taskName,
		Description:  taskDescription,
		Priority:     priority,
		Deadlinedate: deadline,
		Createdat:    time.Now(),
		Updatedat:    time.Now(),
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		panic("failed to create tasks")
	}
	printDBtableTasks(db)

	log.Printf("Добавлен новый дедлайн."+
		"\nНазвание: %s"+
		"\nДата конца: %s"+
		"\nОписание: %s"+
		"\nПриоритет: %d", taskName, deadline, taskDescription, priority) // вывод в консоль

	return ctx.SendString(fmt.Sprintf("Добавлен новый дедлайн."+
		"\nНазвание: %s"+
		"\nДата конца: %s"+
		"\nОписание: %s"+
		"\nПриоритет: %d", taskName, deadline, taskDescription, priority))
}

func setupRoutes(app *fiber.App) {
	app.Static("/static", "./static") // шарим папку static, для передачи js скрипта
	app.Get("/", homePage)
	app.Post("/api/addDeadline", addDeadline) // техническая ссылка для передачи данных post запросами
	log.Println("Сервер запущен на порту 8080")
	app.Listen(":8080")
}

func connectDB() *gorm.DB {
	// Подключение к базе данных postgresql
	dsn := "host=localhost user=postgres password=1234 dbname=task_board_site port=5432 sslmode=disable"
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.LogMode(true) // логирование sql запросов
	return db
}

func printDBtableUsers(db *gorm.DB) {
	// поиск записей в бд
	var users []User
	err := db.Find(&users).Error
	if err != nil {
		log.Fatal("Failed to fetch records:", err)
	}
	// вывод записей из бд
	for _, user := range users {
		log.Println(user)
		//log.Println(fmt.Sprintf("id: %d, nickname: %s, email: %s, passwordhash: %s", user.User_id, user.Nickname, user.Email, user.Passwordhash))
	}
}

func printDBtableTasks(db *gorm.DB) {
	// поиск записей в бд
	var tasks []Task
	err := db.Find(&tasks).Error
	if err != nil {
		log.Fatal("Failed to fetch records:", err)
	}
	// вывод записей из бд
	for _, user := range tasks {
		log.Println(user)
		//log.Println(fmt.Sprintf("id: %d, nickname: %s, email: %s, passwordhash: %s", user.User_id, user.Nickname, user.Email, user.Passwordhash))
	}
}

func startServer() {
	app := fiber.New(fiber.Config{
		Prefork:       false,   // включаем предварительное форкование для увеличения производительности на многоядерных процессорах (проще говоря запуск на всех ядрах процессора)
		ServerHeader:  "Fiber", // добавляем заголовок для идентификации сервера
		CaseSensitive: true,    // включаем чувствительность к регистру в URL
		StrictRouting: true,    // включаем строгую маршрутизацию
	})
	setupRoutes(app)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // настройка логирования для вывода строки в коде
	startServer()

}
