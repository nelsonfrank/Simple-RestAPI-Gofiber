package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{Id: 1, Name: "Walking the dog", Completed: false},
	{Id: 2, Name: "Walking the cat", Completed: false},
}


func main() {
	app := fiber.New()

	app.Use(requestid.New())

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

func setupV1(app *fiber.App)  {
	v1 := app.Group("/v1")
	setupTodosRoutes(v1)
}
func setupTodosRoutes(grp fiber.Router)  {
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", GetTodos)
	todosRoutes.Get("/:id", GetTodo)
	todosRoutes.Post("/", CreateTodo)
	todosRoutes.Delete("/:id", DeleteTodo)
	todosRoutes.Patch("/:id", UpdateTodo)
}

func GetTodos(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(todos)
}

func CreateTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}

	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return err
	}

	todo := &Todo{
		Id:        len(todos)+1,    
		Name:      body.Name,
		Completed: false,
	}
	todos = append(todos, todo)

	return ctx.Status(fiber.StatusOK).JSON(todos)
}

func GetTodo(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return err
	}

	for _, todo := range todos {
		if todo.Id == id  {
			return ctx.Status(fiber.StatusOK).JSON(todo)
		}
	}

	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "todo not found",
	})
}

func DeleteTodo(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	for i, todo := range todos{
		if todo.Id == id {
			todos = append(todos[0:i], todos[i+1:]... )
			
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "todo deleted successfully",
			})
		}
	}
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "todo not found",
	})
}

func UpdateTodo(ctx *fiber.Ctx) error {

	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	var body request

	err = ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Cannot parse body",
		})
	}

	var todo *Todo 

	for _, t := range todos {
		if t.Id == id  {
			todo = t
			break
		}
	}

	if todo == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "todo not found",
		})
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	return ctx.Status(fiber.StatusOK).JSON(todo)


}