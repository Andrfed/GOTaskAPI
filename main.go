// TODO:
// go get github.com/lib/pq -- драйвер PostgreSQL https://metanit.com/go/tutorial/10.3.php https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/

package main

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"

	"fmt"

	"database/sql"

	_ "github.com/lib/pq"

	"os"
)

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Новый запрос:", c.Request.Method, c.Request.URL.Path)
		c.Next() // Передаем управление следующему обработчику
	}
}

// - DB_HOST=postgres
// - DB_PORT=5433
// - DB_USER=postgres
// - DB_PASSWORD=admin
// - DB_NAME=goAppTasksDB

func databaseConnection() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	fmt.Println("Trying to connect with database...")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	fmt.Println("Info:", psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}

// методы: db.Exec(), db.QueryRow(), db.Query(), результат.Scan()
// db.Database.Exec("CALL mydatabase.mystoredprocedure($1, $2)", param1, param2)

func main() {
	// подключение к базе данных
	db := databaseConnection()

	// проверка соединения
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	// Создаем новый экземпляр роутера
	r := gin.Default()

	// Подключаем middleware
	r.Use(LoggerMiddleware())

	// Определяем маршрут для главной страницы
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Список задач")
	})

	taskRouter := r.Group("/tasks")
	{
		taskRouter.POST("/", func(c *gin.Context) {
			// добавление новой задачи
			var task Task

			if err := c.ShouldBindJSON(&task); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			fmt.Println("add new task:", "title:", task.Title, "text:", task.Text)

			new_task_title := task.Title
			new_task_text := task.Text
			rows, err := db.Query("select * from addTask($1, $2)", new_task_title, new_task_text)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			tasks := []Task{}

			for rows.Next() {
				t := Task{}
				err := rows.Scan(&t.Id, &t.Title, &t.Text)
				if err != nil {
					fmt.Println(err)
					continue
				}
				tasks = append(tasks, t)
			}
			if len(tasks) > 0 {
				c.IndentedJSON(http.StatusOK, tasks[0])
			} else {
				c.IndentedJSON(http.StatusOK, Task{})
			}
		})

		taskRouter.GET("/", func(c *gin.Context) {

			// реализация получения списка задач из базы данных

			rows, err := db.Query("select * from getAllTasks()")
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			tasks := []Task{}

			for rows.Next() {
				t := Task{}
				err := rows.Scan(&t.Id, &t.Title, &t.Text)
				if err != nil {
					fmt.Println(err)
					continue
				}
				tasks = append(tasks, t)
			}

			if len(tasks) > 0 {
				c.IndentedJSON(http.StatusOK, tasks)
			} else {
				c.IndentedJSON(http.StatusOK, []Task{})
			}
		})

		taskRouter.GET("/:id", func(c *gin.Context) {
			// получение конкретной задачи
			task_id := c.Param("id")
			rows, err := db.Query("select * from getTask($1)", task_id)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			tasks := []Task{}

			for rows.Next() {
				t := Task{}
				err := rows.Scan(&t.Id, &t.Title, &t.Text)
				if err != nil {
					fmt.Println(err)
					continue
				}
				tasks = append(tasks, t)
			}

			if len(tasks) > 0 {
				c.IndentedJSON(http.StatusOK, tasks[0])
			} else {
				c.IndentedJSON(http.StatusOK, Task{})
			}

		})

		taskRouter.PUT("/:id", func(c *gin.Context) {
			// обновление конкретной задачи
			var task Task

			if err := c.ShouldBindJSON(&task); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			task_id := c.Param("id")
			new_title := task.Title
			new_text := task.Text
			rows, err := db.Query("select * from updateTask($1, $2, $3)", task_id, new_title, new_text)
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			tasks := []Task{}

			for rows.Next() {
				t := Task{}
				err := rows.Scan(&t.Id, &t.Title, &t.Text)
				if err != nil {
					fmt.Println(err)
					continue
				}
				tasks = append(tasks, t)
			}

			if len(tasks) > 0 {
				c.IndentedJSON(http.StatusOK, tasks[0])
			} else {
				c.IndentedJSON(http.StatusOK, Task{})
			}
		})

		taskRouter.DELETE("/:id", func(c *gin.Context) {
			// удаление конкретной задачи
			task_id := c.Param("id")
			result, err := db.Exec("call deleteTask($1)", task_id)
			if err != nil {
				panic(err)
			}
			fmt.Println(result)
			var task Task
			var t_id int64
			t_id, err = strconv.ParseInt(task_id, 10, 0)
			if err != nil {
				panic(err)
			}
			task.Id = int(t_id)
			task.Title = "deleted_id"
			task.Text = ""
			c.IndentedJSON(http.StatusOK, task)
		})
	}

	r.Run(":8000")
}
