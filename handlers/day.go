package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"todolist/models"
)

// DayView handles GET /day/:date - shows all tasks for a specific date
func DayView(c *gin.Context) {
	dateParam := c.Param("date")

	// Parse date or default to today
	var date string
	if dateParam == "" || dateParam == "today" {
		date = time.Now().Format("2006-01-02")
	} else {
		// Validate date format
		_, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"error": "Invalid date format. Use YYYY-MM-DD",
			})
			return
		}
		date = dateParam
	}

	// Get tasks for this date
	tasks, err := models.GetTasksByDate(date)
	if err != nil {
		log.Printf("Error getting tasks for date %s: %v", date, err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load tasks",
		})
		return
	}

	// Calculate previous and next dates
	currentDate, _ := time.Parse("2006-01-02", date)
	prevDate := currentDate.AddDate(0, 0, -1).Format("2006-01-02")
	nextDate := currentDate.AddDate(0, 0, 1).Format("2006-01-02")

	// Check if any tasks are uncompleted (to show/hide "Prepare Tomorrow" button)
	hasUncompleted := false
	completedCount := 0
	for _, task := range tasks {
		if !task.Completed {
			hasUncompleted = true
		} else {
			completedCount++
		}
	}

	// Format date for display
	displayDate := currentDate.Format("Monday, January 2")

	// Check if this date is within edit window (±1 day from today)
	today := time.Now().Truncate(24 * time.Hour)
	dayBefore := today.AddDate(0, 0, -1)
	dayAfter := today.AddDate(0, 0, 1)
	isEditable := !currentDate.Before(dayBefore) && !currentDate.After(dayAfter)

	c.HTML(http.StatusOK, "day.html", gin.H{
		"date":           date,
		"displayDate":    displayDate,
		"prevDate":       prevDate,
		"nextDate":       nextDate,
		"tasks":          tasks,
		"taskCount":      len(tasks),
		"completedCount": completedCount,
		"hasUncompleted": hasUncompleted,
		"isToday":        date == time.Now().Format("2006-01-02"),
		"isEditable":     isEditable,
	})
}

// PrepareNextDay handles POST /day/:date/prepare-next - carries over uncompleted tasks
func PrepareNextDay(c *gin.Context) {
	dateParam := c.Param("date")

	// Validate date
	currentDate, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid date format"})
		return
	}

	// Get uncompleted tasks for this date
	tasks, err := models.GetUncompletedTasksByDate(dateParam)
	if err != nil {
		log.Printf("Error getting uncompleted tasks: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Failed to get tasks"})
		return
	}

	// Calculate next date
	nextDate := currentDate.AddDate(0, 0, 1).Format("2006-01-02")

	// Copy tasks to next day
	copiedCount := 0
	for _, task := range tasks {
		newTask := &models.Task{
			Date:            nextDate,
			Title:           task.Title,
			Description:     task.Description,
			Completed:       false, // Reset completion status
			Notes:           "",    // Don't carry over notes
			CarriedFromDate: dateParam,
		}

		if err := newTask.Create(); err != nil {
			log.Printf("Error creating carried-over task: %v", err)
			continue
		}
		copiedCount++
	}

	log.Printf("Carried over %d tasks from %s to %s", copiedCount, dateParam, nextDate)

	// Redirect to next day
	c.Redirect(http.StatusFound, "/day/"+nextDate)
}

// TodayRedirect handles GET / - redirects to today's day view
func TodayRedirect(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	c.Redirect(http.StatusFound, "/day/"+today)
}
