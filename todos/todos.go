package todos

import (
	"go-fiber-todos/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Todo struct {
	gorm.Model
	Id        int   `gorm:"primaryKey"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

/*
* @ func GetTodos
* Get all todos
* @param c *fiber.Ctx -- fiber context
*/

func GetTodos(c *fiber.Ctx) error {
	db := database.DBConn
	var todoss []Todo
	db.Find(&todoss)
	return c.Status(fiber.StatusOK).JSON(todoss)
}

/*
* @ func GetTodo
* Get single todo
* @param c *fiber.Ctx -- fiber context
*/

func GetTodo(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return err
	}

	db := database.DBConn

	var todo Todo
	db.Find(&todo, id)

	if int(todo.Id) == id{
	  return ctx.Status(fiber.StatusOK).JSON(todo)
	}

	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "todo not found",
	})
}


/*
* @ func CreateTodo
* Create new todo
* @param c *fiber.Ctx -- fiber context
*/

func CreateTodo(ctx *fiber.Ctx) error {
	db := database.DBConn
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
	id := uuid.New()
	todo := Todo{
		Id:       int(id.ID()),    
		Name:      body.Name,
		Completed: false,
	}
	
	db.Create(&todo)

	return ctx.Status(fiber.StatusOK).JSON(todo)
}


/*
* @ func DeleteTodo
* Delete todo
* @param c *fiber.Ctx -- fiber context
*/

func DeleteTodo(ctx *fiber.Ctx) error {
	db := database.DBConn
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	var todo Todo
	db.First(&todo, id)

	if int(todo.Id) != id {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		     "error": "todo not found",
	     })
		
	}

	db.Delete(&todo)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "todo deleted successfully",
	})
	
	
}

/*
* @ func UpdateTodo
* Update todo
* @param c *fiber.Ctx -- fiber context
*/

func UpdateTodo(ctx *fiber.Ctx) error {
	db := database.DBConn

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

	var todo Todo
	// Check if todo exist, if exist assign it value to todo 
	db.First(&todo, id)
	
	if int(todo.Id) != id {
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

	db.Save(&todo)

	return ctx.Status(fiber.StatusOK).JSON(todo)


}