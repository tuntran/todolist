package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"todolist/models"
)

// CreateTask handles POST /tasks - creates a new task
func CreateTask(c *gin.Context) {
	date := c.PostForm("date")
	title := c.PostForm("title")
	description := c.PostForm("description")

	// Validate title is required
	if title == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Title is required"})
		return
	}

	// Default to today if no date provided
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Check if date is within ±1 day window
	if !isWithinEditWindow(date) {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Can only create tasks within 1 day of today"})
		return
	}

	task := &models.Task{
		Date:        date,
		Title:       title,
		Description: description,
		Completed:   false,
	}

	if err := task.Create(); err != nil {
		log.Printf("Error creating task: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to create task"})
		return
	}

	// Redirect back to day view
	c.Redirect(http.StatusFound, "/day/"+date)
}

// GetTask handles GET /tasks/:id - retrieves a task
func GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := models.GetTaskByID(id)
	if err != nil {
		log.Printf("Error getting task: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// ToggleTask handles POST /tasks/:id/toggle - toggles task completion
func ToggleTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid task ID"})
		return
	}

	// Get task to find its date for redirect
	task, err := models.GetTaskByID(id)
	if err != nil {
		log.Printf("Error getting task: %v", err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Task not found"})
		return
	}

	if err := models.ToggleTaskComplete(id); err != nil {
		log.Printf("Error toggling task: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to toggle task"})
		return
	}

	// Redirect back to day view
	c.Redirect(http.StatusFound, "/day/"+task.Date)
}

// AddNotes handles POST /tasks/:id/notes - appends notes to a task
func AddNotes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid task ID"})
		return
	}

	newNotes := c.PostForm("notes")
	if newNotes == "" {
		// Just redirect back if empty notes
		task, _ := models.GetTaskByID(id)
		if task != nil {
			c.Redirect(http.StatusFound, "/day/"+task.Date)
		} else {
			c.Redirect(http.StatusFound, "/")
		}
		return
	}

	task, err := models.GetTaskByID(id)
	if err != nil {
		log.Printf("Error getting task: %v", err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Task not found"})
		return
	}

	// Append notes with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04")
	if task.Notes != "" {
		task.Notes += "\n\n"
	}
	task.Notes += "[" + timestamp + "] " + newNotes

	if err := task.Update(); err != nil {
		log.Printf("Error updating notes: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to add notes"})
		return
	}

	// Redirect back to day view
	c.Redirect(http.StatusFound, "/day/"+task.Date)
}

// DeleteTask handles POST /tasks/:id/delete - removes a task
func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid task ID"})
		return
	}

	// Get task to find its date for redirect
	task, err := models.GetTaskByID(id)
	if err != nil {
		log.Printf("Error getting task: %v", err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Task not found"})
		return
	}

	date := task.Date

	if err := models.DeleteTask(id); err != nil {
		log.Printf("Error deleting task: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to delete task"})
		return
	}

	// Redirect back to day view
	c.Redirect(http.StatusFound, "/day/"+date)
}

// isWithinEditWindow checks if task date is within ±1 day of today
func isWithinEditWindow(taskDate string) bool {
	today := time.Now().Truncate(24 * time.Hour)
	parsed, err := time.Parse("2006-01-02", taskDate)
	if err != nil {
		return false
	}

	dayBefore := today.AddDate(0, 0, -1)
	dayAfter := today.AddDate(0, 0, 1)

	// Task date must be >= yesterday and <= tomorrow
	return !parsed.Before(dayBefore) && !parsed.After(dayAfter)
}

// UpdateTask handles POST /tasks/:id/edit - updates task title and description
func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := models.GetTaskByID(id)
	if err != nil {
		log.Printf("Error getting task: %v", err)
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Task not found"})
		return
	}

	// Check if task is within edit window (±1 day)
	if !isWithinEditWindow(task.Date) {
		c.HTML(http.StatusForbidden, "error.html", gin.H{"error": "Can only edit tasks within 1 day of today"})
		return
	}

	// Get form values
	newTitle := c.PostForm("title")
	newDescription := c.PostForm("description")

	// Validate title
	if newTitle == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Title is required"})
		return
	}

	// Update task
	task.Title = newTitle
	task.Description = newDescription

	if err := task.Update(); err != nil {
		log.Printf("Error updating task: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to update task"})
		return
	}

	// Redirect back to day view
	c.Redirect(http.StatusFound, "/day/"+task.Date)
}
