package main

import (
	"go-fiber-todos/todos"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupV1(app *fiber.App)  {
	v1 := app.Group("/v1")
	setupTodosRoutes(v1)
}

func setupTodosRoutes(grp fiber.Router)  {
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", todos.GetTodos)
	todosRoutes.Get("/:id", todos.GetTodo)
	todosRoutes.Post("/", todos.CreateTodo)
	todosRoutes.Delete("/:id", todos.DeleteTodo)
	todosRoutes.Patch("/:id", todos.UpdateTodo)
}

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
        Format:     "[${ip}]:${port} ${status} - ${method} ${path}\n",
    }))
	
	app.Get("/", func(c *fiber.Ctx) error {
		var msg string = "Hello, World"
		return c.Send([]byte(msg))
	})
 
	setupV1(app)


	// Listen on PORT 3000
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}



