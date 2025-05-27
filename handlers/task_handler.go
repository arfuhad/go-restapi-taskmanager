package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"taskapi/models"

	"github.com/gin-gonic/gin"
)

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}

// POST /tasks
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	now := time.Now()
	res, err := DB.Exec(`INSERT INTO tasks (title, completed, created_at, updated_at, priority, due_date, tags) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		task.Title, false, now, now, task.Priority, task.DueDate, task.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	id, _ := res.LastInsertId()
	task.ID = int(id)
	task.Completed = false
	task.CreatedAt = now
	task.UpdatedAt = now

	c.JSON(http.StatusCreated, task)
}

// GET /tasks
func GetTasks(c *gin.Context) {
	rows, err := DB.Query(`SELECT id, title, completed, created_at, updated_at, priority, due_date, tags FROM tasks`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var created string
		var updated string
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &created, &updated); err == nil {
			task.CreatedAt, _ = time.Parse(time.RFC3339, created)
			task.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, []models.Task{})
	} else {
		c.JSON(http.StatusOK, tasks)
	}
}

// GET /tasks/:id
func GetTaskById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var task models.Task
	var created, updated string
	err := DB.QueryRow(`SELECT id, title, completed, created_at, updated_at, priority, due_date, tags FROM tasks WHERE id = ?`, id).
		Scan(&task.ID, &task.Title, &task.Completed, &created, &updated, &task.Priority, &task.DueDate, &task.Tags)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}
	task.CreatedAt, _ = time.Parse(time.RFC3339, created)
	task.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	c.JSON(http.StatusOK, task)
}

// PUT /tasks/:id/done
func MarkTaskDone(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	now := time.Now()
	res, err := DB.Exec(`UPDATE tasks SET completed = 1, updated_at = ? WHERE id = ?`, now, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task marked as done"})
}

// PUT /tasks/:id/undone
func MarkTaskUndone(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	now := time.Now()
	res, err := DB.Exec(`UPDATE tasks SET completed = 0, updated_at = ? WHERE id = ?`, now, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task marked as undone"})
}

// DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	res, err := DB.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
