package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func homePage(ctx *fiber.Ctx) error {
	return ctx.Render("templates/homePage.html", nil)
}

func addDeadline(ctx *fiber.Ctx) error {
	if ctx.Method() != fiber.MethodPost { // перешли на api по адресу отправки данных (с методом get)
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"error": "Метод запроса не поддерживается!",
		})
	}

	datetime := ctx.FormValue("datetime-input") // получение даты дедлайна из формы
	taskName := ctx.FormValue("task-name")
	taskDescription := ctx.FormValue("task-description")

	log.Printf("Добавлен новый дедлайн.\nДата конца: %s\nНазвание: %s\nОписание: %s", taskName, datetime, taskDescription)                         // вывод в консоль
	return ctx.SendString(fmt.Sprintf("Добавлен новый дедлайн.\nДата конца: %s\nНазвание: %s\nОписание: %s", datetime, taskName, taskDescription)) // ответ сервера
}

func setupRoutes(app *fiber.App) {
	app.Static("/static", "./static") // шарим папку static, для передачи js скрипта
	app.Get("/", homePage)
	app.Post("/api/addDeadline", addDeadline) // техническая ссылка для передачи данных post запросами
	log.Println("Сервер запущен на порту 8080")
	app.Listen(":8080")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // настройка логгирования для вывода строки в коде
	app := fiber.New(fiber.Config{
		Prefork:       true,    // включаем предварительное форкование для увеличения производительности на многоядерных процессорах (проще говоря запуск на всех ядрах процессора)
		ServerHeader:  "Fiber", // добавляем заголовок для идентификации сервера
		CaseSensitive: true,    // включаем чувствительность к регистру в URL
		StrictRouting: true,    // включаем строгую маршрутизацию
	})
	setupRoutes(app)
}
