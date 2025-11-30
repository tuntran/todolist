package main

import (
	"html/template"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"todolist/db"
	"todolist/handlers"
)

func main() {
	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create Gin router
	r := gin.Default()

	// Register custom template functions
	r.SetFuncMap(template.FuncMap{
		"split": strings.Split,
	})

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve static files
	r.Static("/static", "./static")

	// Routes
	r.GET("/", handlers.TodayRedirect)
	r.GET("/day/:date", handlers.DayView)
	r.POST("/day/:date/prepare-next", handlers.PrepareNextDay)

	// Task routes
	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks/:id", handlers.GetTask)
	r.POST("/tasks/:id/toggle", handlers.ToggleTask)
	r.POST("/tasks/:id/notes", handlers.AddNotes)
	r.POST("/tasks/:id/edit", handlers.UpdateTask)
	r.POST("/tasks/:id/delete", handlers.DeleteTask)

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
