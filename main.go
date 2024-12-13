package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type User struct {
	UserId       uint      `gorm:"primaryKey"`
	Nickname     string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	PasswordHash string    `gorm:"not null,default:_"`
	CreateDat    time.Time `gorm:"autoCreateTime"`
	UpdateDat    time.Time `gorm:"autoUpdateTime"`
}

type Claims struct { // для сессий
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
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
	Status       *bool     `gorm:"default:false"`
}

func homePage(ctx *fiber.Ctx) error {
	return ctx.Render("templates/homePage.html", nil)
}

func loginPage(ctx *fiber.Ctx) error {
	return ctx.Render("templates/loginPage.html", nil)
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
	var task_id int = -1
	fmt.Sscan(ctx.FormValue("task-id"), &task_id)
	var priority uint
	fmt.Sscan(ctx.FormValue("priority"), &priority) // приведение строки к uint
	var status bool
	fmt.Sscan(ctx.FormValue("status"), &status)

	log.Printf("Полученные данные из formValue:\n%s | %s | %s | %s | %s", task_id, taskName, taskDescription, deadline, status, priority)

	var result int
	DB.Table("tasks").Where("task_id = ?", task_id).Count(&result) // проверка есть ли такой ID в БД
	// сохранение полученных данных о дедлайне в БД
	if result == 0 { // если такого id нет в бд => добавляем таску
		var lastTaskId Task
		DB.Table("tasks").Order("task_id desc").Last(&lastTaskId) // получение последней записи (обратная сортировка {desk} по id)
		newTask := Task{
			TaskId:       lastTaskId.TaskId + 1,
			UserId:       ctx.Locals("user_id").(uint),
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

func hashPassword(password string) (string, error) {
	// хэширование пароля
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	// проверка пароля и хэша на совпадение
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func loginHandler(ctx *fiber.Ctx) error {
	// метод для проверки наличия пользователя в БД
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	var user User
	err := DB.Table("users").Where("email = ?", email).Find(&user).Error
	if err != nil {
		log.Println("Failed to find user:", err)
	}

	if checkPasswordHash(password, user.PasswordHash) { // сравнение хэшей паролей
		token, err := generateJWT(user.UserId) // генерация jwt токена
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка генерации токена"})
		}
		log.Printf("User %s is logged in. ID: %d", user.Nickname, user.UserId)
		return ctx.JSON(fiber.Map{"token": token})
	}
	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверные учетные данные"}) // ответ для фронта
}

func registerHandler(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	nickname := ctx.FormValue("nickname")

	var result int // проверка на уникальность емайла в бд
	DB.Table("users").Where("email = ?", email).Count(&result)
	if result != 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Ошибка! Аккаунт с такой почтой уже существует."})
	}

	passwordHash, err := hashPassword(password) // генерируем хэш пароля перед сохранением в базе
	if err != nil {
		panic("failed to generate hash password for new user")
	}
	var lastUserId User
	DB.Table("users").Order("user_id desc").Last(&lastUserId) // получение последней записи (обратная сортировка {desk} по id)

	newUser := User{
		UserId:       lastUserId.UserId + 1,
		Nickname:     nickname,
		Email:        email,
		PasswordHash: passwordHash,
		CreateDat:    time.Now(),
		UpdateDat:    time.Now(),
	}
	err = DB.Table("users").Create(&newUser).Error

	if err != nil {
		panic("failed to register new user")
	}
	log.Printf("Был создан новый аккаунт с id:%d", lastUserId.UserId+1)
	token, err := generateJWT(newUser.UserId) // генерация jwt токена
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ошибка генерации токена"})
	}
	return ctx.JSON(fiber.Map{"token": token})
}

func sendTasksFromDB(ctx *fiber.Ctx) error { // /api/getTasks // выдача списка задач в json'e для конкретного юзера
	var tasks []Task
	var taskId string
	var userId = ctx.Locals("user_id").(uint) // получаем айди юзера от которого пришел запрос
	log.Printf("Айди из локали ctx для получения тасок: userId: %d", userId)

	fmt.Sscan(ctx.Params("id"), &taskId) // изъятие параметра айди таски через слеш
	log.Println("taskId: ", taskId)

	if taskId == "all" { // для отдачи всех тасок определенному юзеру
		err := DB.Order("deadline_date").Where("user_id = ?", userId).Find(&tasks).Error // сбор записей таблицы задач из бд в массив
		if err != nil {
			log.Fatal("Failed to fetch records:", err)
		}
	} else { // для прочтения информации о конкретной 1 таске по task_id
		err := DB.Where("task_id = ?", taskId).Find(&tasks).Error // сбор записей таблицы задач из бд в массив
		if err != nil {
			log.Fatal("Failed to fetch records:", err)
		}
	}
	return ctx.JSON(tasks)
}

func generateJWT(userId uint) (string, error) {
	claims := Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Токен действителен 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // метод шифрования HS256
	tokenString, err := token.SignedString([]byte(JWTKey))
	if err != nil {
		return "", err
	}
	log.Printf("tokenString в генерации: %s", tokenString)
	return tokenString, nil
}

func authMiddleware(ctx *fiber.Ctx) error {
	// проверка КАЖДОГО ЗАПРОСА на наличие валидного JWT токена, или простым языком проверка того что юзер вошел в аккаунт
	tokenString := ctx.Get("Authorization") // Получение JWT токена из заголовка запроса Authorization

	if tokenString == "" {
		return ctx.Redirect("/login", fiber.StatusSeeOther)
	}

	// парсинг токена проверяет подпись токена, срок действия и другие утверждения (claims)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTKey), nil // Тот же самый секретный ключ, что и при генерации
	})

	if err != nil {
		log.Printf("Ошибка! токен в обработчике: token:\n%s\nerr=%v", token, err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Недействительный токен"})
	}

	claims, ok := token.Claims.(*Claims)
	log.Printf("Данные текущего запроса в claims: %v", claims, ok)

	if ok && token.Valid {
		ctx.Locals("user_id", claims.UserID) // Токен валиден. Cохраняем айдишник юзера в контексте ТЕКУЩЕГО запроса Fiber
		log.Println("Окей, токен валид, айди userа сохранен в текущем запросе")
		return ctx.Next() // Передача управления следующему обработчику в цепочке.
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Недействительный токен"})
}

func setupRoutes(app *fiber.App) {
	app.Static("/static", "./static") // шарим папку static, для передачи js скрипта

	app.Get("/", homePage)
	app.Get("/login", loginPage)
	app.Post("/authorization", loginHandler) // login
	app.Post("/register", registerHandler)   // reg

	api := app.Group("/api")
	api.Use(authMiddleware) // применяем проверку сессии и JWT токенов для всех маршрутов /api

	app.Get("/api/getTasks/:id", sendTasksFromDB)           // отдача списка задач из бд фронту
	app.Post("/api/addAndEditDeadline", addAndEditDeadline) // техническая ссылка для передачи данных post запросами
	app.Post("/api/deleteDeadline/:id", deleteDeadline)

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
	setupRoutes(app) // установка пар маршрутов и функций
}

var DB *gorm.DB
var JWTKey string

func main() {
	fileConfig, _ := os.ReadFile("config.json") // чтение конф. данных из конфига
	var jsonConfig map[string]interface{}
	json.Unmarshal(fileConfig, &jsonConfig)
	JWTKey = jsonConfig["jwt-key"].(string) // чтение секретного ключа для генерации jwt токенов
	log.Printf("JWTKey: ", JWTKey)

	log.SetFlags(log.LstdFlags | log.Lshortfile) // настройка логирования для вывода строки в коде
	DB = connectDB()
	defer DB.Close()
	DB.AutoMigrate(&Task{})
	DB.AutoMigrate(&User{})

	startServer()
}
