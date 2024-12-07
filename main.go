package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"time"
)

type User struct {
	UserId       uint      `gorm:"primaryKey"`
	Nickname     string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	PasswordDash string    `gorm:"not null"`
	CreateDat    time.Time `gorm:"autoCreateTime"`
	UpdateDat    time.Time `gorm:"autoUpdateTime"`
}

type Task struct {
	TaskId       uint      `gorm:"primary_key,column:task_id,auto_increment"`
	UserId       uint      `gorm:"foreign_key"`
	Title        string    `gorm:"not null"`
	Description  string    `gorm:"column:description"`
	Priority     *uint     `gorm:"default:0"`
	DeadlineDate time.Time `gorm:"not null,default:CURRENT_TIMESTAMP()"`
	CreateDat    time.Time `gorm:"autoCreateTime"`
	UpdateDat    time.Time `gorm:"autoUpdateTime"`
	Status       *bool
}

func homePage(ctx *fiber.Ctx) error {
	return ctx.Render("templates/homePage.html", nil)
}

func deleteDeadline(ctx *fiber.Ctx) error {
	var tasks []Task
	var taskId uint
	fmt.Sscan(ctx.Params("id"), &taskId) // изъятие параметра айди таски через слеш

	err := DB.Where("task_id = ?", taskId).Delete(&tasks).Error // сбор записей таблицы задач из бд в массив
	if err != nil {
		log.Fatal("Failed to fetch records:", err)
	}

	return ctx.SendString(fmt.Sprintf("Дедлайн с id:%d удалён.", taskId))
}

func addAndEditDeadline(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodPost { // перешли на api по адресу отправки данных (с методом get)
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "Метод запроса не поддерживается!",
		})
	}

	taskName := ctx.FormValue("task-name")
	taskDescription := ctx.FormValue("task-description")
	deadline, _ := time.Parse("2006-01-02T15:04", ctx.FormValue("datetime-input")) // получение даты дедлайна из формы // time.RFC3339 используется для парсинга строк в формате ISO 8601
	var task_id uint
	fmt.Sscan(ctx.FormValue("task-id"), &task_id)
	var priority uint
	fmt.Sscan(ctx.FormValue("priority"), &priority) // приведение строки к uint
	var status bool
	fmt.Sscan(ctx.FormValue("status"), &status)

	log.Printf("Полученные данные из formValue:\n%s | %s | %s | %s | %s", task_id, taskName, taskDescription, deadline, status, priority)

	// сохранение полученных данных о дедлайне в БД
	if task_id == 0 { // если поле id пустое => добавляем таску
		var lastTaskId Task
		DB.Table("tasks").Order("task_id desc").Last(&lastTaskId) // получение последней записи (обратная сортировка {desk} по id)
		newTask := Task{
			UserId:       1,
			Title:        taskName,
			Description:  taskDescription,
			Priority:     &priority,
			DeadlineDate: deadline,
			CreateDat:    time.Now(),
			UpdateDat:    time.Now(),
		}
		err := DB.Table("tasks").Create(&newTask).Error
		if err != nil {
			panic("failed to create tasks")
		}
		log.Printf("Дедлайн с id:%d был изменён в бд", lastTaskId.TaskId+1)
	} else {
		var editTask Task
		err := DB.Table("tasks").Where("task_id = ?", task_id).First(&editTask).Error // Находим запись по ID
		if err != nil {
			panic("failed to find tasks for edit") // Обработка ошибки, если запись не найдена
		}
		// изменение значений записи в бд
		editTask.Title = taskName
		editTask.Description = taskDescription
		editTask.Priority = &priority
		editTask.DeadlineDate = deadline
		editTask.UpdateDat = time.Now()
		editTask.Status = &status

		err = DB.Table("tasks").Where("task_id = ?", task_id).Updates(&editTask).Error // сохраняем изменения
		if err != nil {
			panic("failed to update task") // Обработка ошибки, если запись не найдена
		}

		log.Printf("Дедлайн с id:%d был изменён в бд", task_id)
		return err
	}
	return ctx.SendString(fmt.Sprintf("Дедлайн с id:%d был изменён в бд", task_id))
}

func connectDB() *gorm.DB {

	fileConfig, err := os.ReadFile("config.json") // чтение конф. данных из конфига
	var jsonConfig map[string]interface{}
	json.Unmarshal(fileConfig, &jsonConfig)

	// Подключение к базе данных postgresql
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		jsonConfig["host"],
		jsonConfig["user"],
		jsonConfig["password"],
		jsonConfig["dbname"],
		jsonConfig["port"],
	)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.LogMode(true) // логирование sql запросов
	return db
}

func getUsersFromDB(DB *gorm.DB) {
	// поиск записей в бд
	var users []User
	err := DB.Find(&users).Error
	if err != nil {
		log.Fatal("Failed to fetch records:", err)
	}
	// вывод записей из бд
	for _, user := range users {
		log.Println(user)
		//log.Println(fmt.Sprintf("id: %d, nickname: %s, email: %s, passwordhash: %s", user.User_id, user.Nickname, user.Email, user.Passwordhash))
	}
}

func sendTasksFromDB(ctx *fiber.Ctx) error { // /api/getTasks // выдача списка задач в json'e для конкретного юзера
	var tasks []Task
	var taskId uint
	fmt.Sscan(ctx.Params("id"), &taskId) // изъятие параметра айди таски через слеш
	if taskId == 0 {                     // 0 для отдачи всех тасок
		err := DB.Order("deadline_date").Find(&tasks).Error // сбор записей таблицы задач из бд в массив
		if err != nil {
			log.Fatal("Failed to fetch records:", err)
		}
	} else { // конкретный id таски
		err := DB.Where("task_id = ?", taskId).Find(&tasks).Error // сбор записей таблицы задач из бд в массив
		if err != nil {
			log.Fatal("Failed to fetch records:", err)
		}
	}
	return ctx.JSON(tasks)
}

func setupRoutes(app *fiber.App) {
	app.Static("/static", "./static") // шарим папку static, для передачи js скрипта
	app.Get("/", homePage)
	app.Post("/api/addAndEditDeadline", addAndEditDeadline) // техническая ссылка для передачи данных post запросами
	app.Post("/api/deleteDeadline/:id", deleteDeadline)
	app.Get("/api/getTasks/:id", sendTasksFromDB) // отдача списка задач из бд фронту
	log.Println("Сервер запущен на порту 8080")
	app.Listen(":8080")
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

var DB *gorm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // настройка логирования для вывода строки в коде
	DB = connectDB()
	defer DB.Close()
	DB.AutoMigrate(&Task{})
	startServer()
}
