package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID        int
	Detail     string 
	Completed bool
	Urgent bool	
}



// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "toolman1"
// 	dbname   = "mydatabase"
// )

// const (
// 	host     = "dpg-cio5v9t9aq06u3msmef0-a"
// 	port     = 5432
// 	user     = "todo_database_k1lt_user"
// 	password = "4EOjxkAIA3RvWxyFfAQ6TyNZBNFxbHsb"
// 	dbname   = "todo_database_k1lt"
// )

var db *sql.DB
var err error
// connecting to DB
func init() {
	// connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	connStr := "postgres://todo_database_k1lt_user:4EOjxkAIA3RvWxyFfAQ6TyNZBNFxbHsb@dpg-cio5v9t9aq06u3msmef0-a.oregon-postgres.render.com/todo_database_k1lt"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}


func main() {

	defer db.Close()

	// creating initial table IF table doesnt exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			detail TEXT NOT NULL,
			completed BOOLEAN NOT NULL,
			urgent BOOLEAN NOT NULL
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	
	
	// Routes
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/todos", getTodos)
	router.POST("/todos", createTodo)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	
	// creating port for railway app
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	router.Run("0.0.0.0:" + port)
}



// handles for routes
// handles for routes
// handles for routes


// fetch for all list items 
func getTodos(c *gin.Context) {

	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		todo := Todo{}
		rows.Scan(&todo.ID, &todo.Detail, &todo.Completed, &todo.Urgent)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

// creating new list item and returning all items
func createTodo(c *gin.Context) {

	todo := Todo{}
	err := c.BindJSON(&todo)

	if err != nil {
		log.Fatal(err)
	}
	
	_, err2 := db.Exec("INSERT INTO todos (detail, completed, urgent) VALUES ($1, $2, $3)", todo.Detail, todo.Completed, todo.Urgent)
	
	if err2 != nil {
		log.Fatal(err2)
	}

	rows, err3 := db.Query("SELECT * FROM todos")
	if err3 != nil {
		log.Fatal(err3)
	}

	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		todo := Todo{}
		rows.Scan(&todo.ID, &todo.Detail, &todo.Completed, &todo.Urgent)
		todos = append(todos, todo)
	}

	c.JSON(http.StatusCreated, todos)
}

// updating list item
func updateTodo(c *gin.Context) {
	id := c.Param("id")

	todo := Todo{}
	err := c.BindJSON(&todo)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE todos SET detail = $1, completed = $2, urgent = $3 WHERE id = $4", todo.Detail, todo.Completed, todo.Urgent, id)
	if err != nil {
		log.Fatal(err)
	}

	c.Status(http.StatusOK)
}

// deleting list item
func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	c.Status(http.StatusOK)
}