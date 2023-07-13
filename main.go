package main

import (
	
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Items struct {
	detail string
	urgent bool
	completed bool
}



func main() {

	// connStr := "postgres://todo_database_k1lt_user:4EOjxkAIA3RvWxyFfAQ6TyNZBNFxbHsb@dpg-cio5v9t9aq06u3msmef0-a.oregon-postgres.render.com/todo_database_k1lt"
	// // Connect to database
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }


	itemList := []Items{
		{detail: "wash dishes", urgent:  true, completed: true},
		{detail: "cook dinner", urgent: false, completed: true},
		{detail: "sleep", urgent: false, completed: false},
	}

	app := fiber.New()

	app.Get("/", hello)
	app.Get("/list", func(c *fiber.Ctx) error {
		return c.JSON(itemList)
	})

	app.Get("/env", func(c *fiber.Ctx) error {
		return c.SendString("Hello, ENV! " + os.Getenv("TEST_ENV"))
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, idiot!")
}

