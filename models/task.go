package models

import (
	"database/sql"
	"time"

	"todolist/db"
)

// Task represents a todo task
type Task struct {
	ID              int       `json:"id"`
	Date            string    `json:"date"`              // YYYY-MM-DD format
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Completed       bool      `json:"completed"`
	Notes           string    `json:"notes"`
	CarriedFromDate string    `json:"carried_from_date"` // Empty if not carried over
	CreatedAt       time.Time `json:"created_at"`
}

// Create inserts a new task into the database
func (t *Task) Create() error {
	query := `
		INSERT INTO tasks (date, title, description, completed, notes, carried_from_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	completed := 0
	if t.Completed {
		completed = 1
	}

	carriedFrom := sql.NullString{String: t.CarriedFromDate, Valid: t.CarriedFromDate != ""}

	result, err := db.DB.Exec(query, t.Date, t.Title, t.Description, completed, t.Notes, carriedFrom)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = int(id)
	return nil
}

// GetTaskByID retrieves a task by its ID
func GetTaskByID(id int) (*Task, error) {
	query := `SELECT id, date, title, description, completed, notes, carried_from_date, created_at
	          FROM tasks WHERE id = ?`

	task := &Task{}
	var completed int
	var carriedFrom sql.NullString
	var description sql.NullString
	var notes sql.NullString

	err := db.DB.QueryRow(query, id).Scan(
		&task.ID, &task.Date, &task.Title, &description,
		&completed, &notes, &carriedFrom, &task.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	task.Completed = completed == 1
	if carriedFrom.Valid {
		task.CarriedFromDate = carriedFrom.String
	}
	if description.Valid {
		task.Description = description.String
	}
	if notes.Valid {
		task.Notes = notes.String
	}

	return task, nil
}

// GetTasksByDate retrieves all tasks for a specific date
func GetTasksByDate(date string) ([]Task, error) {
	query := `SELECT id, date, title, description, completed, notes, carried_from_date, created_at
	          FROM tasks WHERE date = ? ORDER BY created_at ASC`

	rows, err := db.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completed int
		var carriedFrom sql.NullString
		var description sql.NullString
		var notes sql.NullString

		err := rows.Scan(
			&task.ID, &task.Date, &task.Title, &description,
			&completed, &notes, &carriedFrom, &task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		task.Completed = completed == 1
		if carriedFrom.Valid {
			task.CarriedFromDate = carriedFrom.String
		}
		if description.Valid {
			task.Description = description.String
		}
		if notes.Valid {
			task.Notes = notes.String
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update updates an existing task
func (t *Task) Update() error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, completed = ?, notes = ?
		WHERE id = ?
	`

	completed := 0
	if t.Completed {
		completed = 1
	}

	_, err := db.DB.Exec(query, t.Title, t.Description, completed, t.Notes, t.ID)
	return err
}

// ToggleTaskComplete toggles the completion status of a task
func ToggleTaskComplete(id int) error {
	query := `UPDATE tasks SET completed = CASE WHEN completed = 0 THEN 1 ELSE 0 END WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	return err
}

// DeleteTask removes a task from the database
func DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	return err
}

// GetUncompletedTasksByDate retrieves uncompleted tasks for carry-over
func GetUncompletedTasksByDate(date string) ([]Task, error) {
	query := `SELECT id, date, title, description, completed, notes, carried_from_date, created_at
	          FROM tasks WHERE date = ? AND completed = 0 ORDER BY created_at ASC`

	rows, err := db.DB.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completed int
		var carriedFrom sql.NullString
		var description sql.NullString
		var notes sql.NullString

		err := rows.Scan(
			&task.ID, &task.Date, &task.Title, &description,
			&completed, &notes, &carriedFrom, &task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		task.Completed = completed == 1
		if carriedFrom.Valid {
			task.CarriedFromDate = carriedFrom.String
		}
		if description.Valid {
			task.Description = description.String
		}
		if notes.Valid {
			task.Notes = notes.String
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
