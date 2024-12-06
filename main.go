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
	User_id      uint      `gorm:"primaryKey"`
	Nickname     string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	Passwordhash string    `gorm:"not null"`
	Createdat    time.Time `gorm:"autoCreateTime"`
	Updatedat    time.Time `gorm:"autoUpdateTime"`
}

type Task struct {
	Task_id      uint      `gorm:"primaryKey"`
	User_id      uint      `gorm:"foreignKey"`
	Title        string    `gorm:"not null"`
	Description  string    `gorm:"column:description"`
	Priority     uint      `gorm:"default:0"`
	Deadlinedate time.Time `gorm:"not null"`
	Createdat    time.Time `gorm:"autoCreateTime"`
	Updatedat    time.Time `gorm:"autoUpdateTime"`
	Status       bool      `gorm:"column:status"`
}

func homePage(ctx *fiber.Ctx) error {
	return ctx.Render("templates/homePage.html", nil)
}

func deleteDeadline(ctx *fiber.Ctx) error {
	db := connectDB()
	defer db.Close()

	var tasks []Task
	var taskId uint
	fmt.Sscan(ctx.Params("id"), &taskId) // изъятие параметра айди таски через слеш

	err := db.Where("Task_id = ?", taskId).Delete(&tasks).Error // сбор записей таблицы задач из бд в массив
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

	log.Println(task_id, taskName, taskDescription, deadline, status, priority)

	// сохранение полученных данных о дедлайне в БД
	db := connectDB()
	defer db.Close() // отложенное отклчючение от бд (пока не вышли из текущей функции)

	if task_id == 0 { // если поле id пустое => добавляем таску
		var lastTaskId Task
		db.Table("tasks").Order("Task_id desc").Last(&lastTaskId) // получение последней записи (обратная сортировка {desk} по id)
		newTask := Task{
			Task_id:      lastTaskId.Task_id + 1,
			User_id:      1,
			Title:        taskName,
			Description:  taskDescription,
			Priority:     priority,
			Deadlinedate: deadline,
			Createdat:    time.Now(),
			Updatedat:    time.Now(),
			Status:       false,
		}
		err := db.Table("tasks").Create(&newTask).Error
		if err != nil {
			panic("failed to create tasks")
		}
		log.Printf("Дедлайн с id:%d был изменён в бд", lastTaskId.Task_id+1)
	} else {
		var editTask Task
		err := db.Table("tasks").Where("Task_id = ?", task_id).Find(&editTask).Limit(1).Error // Находим запись по ID
		if err != nil {
			panic("failed to find tasks for edit") // Обработка ошибки, если запись не найдена
		}
		log.Printf(editTask.Title, editTask.Status, status)
		// изменение значений записи в бд
		editTask.Title = taskName
		editTask.Description = taskDescription
		editTask.Priority = priority
		editTask.Deadlinedate = deadline
		editTask.Updatedat = time.Now()
		editTask.Status = status

		err = db.Table("tasks").Where("Task_id = ?", task_id).Updates(&editTask).Error // сохраняем изменения
		if err != nil {
			panic("failed to find tasks for edit") // Обработка ошибки, если запись не найдена
		}
		// сохраняем статус отдельно (т.к он почему-то не хочет сохраняться сразу со всеми)
		err = db.Table("tasks").Where("Task_id = ?", task_id).UpdateColumn("status", status).Error
		if err != nil {
			panic("failed to find tasks for edit") // Обработка ошибки, если запись не найдена
		}
		// сохраняем приоритет отдельно (т.к он почему-то не хочет сохраняться сразу со всеми)
		err = db.Table("tasks").Where("Task_id = ?", task_id).UpdateColumn("priority", priority).Error
		if err != nil {
			panic("failed to find tasks for edit") // Обработка ошибки, если запись не найдена
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

func getUsersFromDB(db *gorm.DB) {
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

func sendTasksFromDB(ctx *fiber.Ctx) error { // /api/getTasks // выдача списка задач в json'e для конкретного юзера
	db := connectDB()
	defer db.Close()

	var tasks []Task
	var taskId uint
	fmt.Sscan(ctx.Params("id"), &taskId) // изъятие параметра айди таски через слеш
	if taskId == 0 {                     // 0 для отдачи всех тасок
		err := db.Order("deadlinedate").Find(&tasks).Error // сбор записей таблицы задач из бд в массив
		if err != nil {
			log.Fatal("Failed to fetch records:", err)
		}
	} else { // конкретный id таски
		err := db.Where("Task_id = ?", taskId).Find(&tasks).Error // сбор записей таблицы задач из бд в массив
		if err != nil {
			log.Fatal("Failed to fetch records:", err)
		}
	}

	//for i, _ := range tasks { // *ВАЖНО* изменение даты дедлайнов (ставим наш час. пояс) - 3часа
	//	tasks[i].Deadlinedate = tasks[i].Deadlinedate.Add(-3 * time.Hour)
	//}

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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // настройка логирования для вывода строки в коде
	startServer()
}
