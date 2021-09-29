package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{
	{Id: 1, Name: "Walking the dog", Completed: false},
	{Id: 2, Name: "Walking the cat", Completed: false},
}


func main() {
	app := fiber.New()


	app.Get("/", func(c *fiber.Ctx) error {
		var msg string = "Hello, World"
		return c.Send([]byte(msg))
	})

	app.Get("/todos", GetTodos)
	app.Get("/todos/:id", GetTodo)
	app.Post("/todos", CreateTodo)
	app.Delete("/todos/:id", DeleteTodo)
	app.Patch("/todos/:id", UpdateTodo)


	// Listen on PORT 3000
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func GetTodos(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(todos)
}

func CreateTodo (ctx *fiber.Ctx) error {
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

	todo := Todo{
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
	ctx.BodyParser(&body)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Cannot parse body",
		})
	}

	var todo Todo

	for _, t := range todos {
		if t.Id == id  {
			todo = t
			break
		}
	}

	if todo.Id == 0 {
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

	return ctx.Status(fiber.StatusOK).JSON(todos)

}