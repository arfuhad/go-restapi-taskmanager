package main

import (
	"taskapi/handlers"
	"taskapi/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	defer DB.Close()

	r := gin.Default()

	// Inject DB into handlers
	handlers.SetDB(DB)

	r.POST("/tasks", middleware.ValidateTaskInput(), handlers.CreateTask)
	r.GET("/tasks", handlers.GetTasks)
	r.GET("/tasks/:id", handlers.GetTaskById)
	r.PUT("/tasks/:id/done", handlers.MarkTaskDone)
	r.PUT("/tasks/:id/undone", handlers.MarkTaskUndone)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.Run(":8080")
}
