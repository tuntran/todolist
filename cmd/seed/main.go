package main

import (
	"log"
	"time"

	"todolist/db"
	"todolist/models"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	defer db.Close()

	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	tasks := []models.Task{
		// Today
		{Date: today, Title: "Review pull requests", Completed: false},
		{Date: today, Title: "Deploy to staging", Completed: false, Notes: "[2025-11-30 10:30] Check CI/CD pipeline"},
		{Date: today, Title: "Team standup meeting", Completed: true, Notes: "[2025-11-30 09:00] Discussed sprint goals"},
		{Date: today, Title: "Write documentation", Completed: false},
		{Date: today, Title: "Fix critical bug #1234", Completed: true, Notes: "[2025-11-30 11:00] Root cause: race condition"},
		// Tomorrow
		{Date: tomorrow, Title: "Sprint planning", Completed: false},
		{Date: tomorrow, Title: "Code review session", Completed: false},
		{Date: tomorrow, Title: "Update dependencies", Completed: false},
		// Yesterday
		{Date: yesterday, Title: "Setup project structure", Completed: true, Notes: "[2025-11-29 14:00] Created Go/Gin app"},
		{Date: yesterday, Title: "Design database schema", Completed: true},
		{Date: yesterday, Title: "Implement UI mockup", Completed: true, Notes: "[2025-11-29 16:30] Cyberpunk style applied"},
	}

	for _, t := range tasks {
		task := &models.Task{
			Date:        t.Date,
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
			Notes:       t.Notes,
		}
		if err := task.Create(); err != nil {
			log.Printf("Failed to create task '%s': %v", t.Title, err)
		} else {
			log.Printf("Created: %s", t.Title)
		}
	}

	log.Println("Demo data seeded successfully!")
}
