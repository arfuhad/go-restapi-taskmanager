package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"taskapi/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory db: %v", err)
	}

	createTable := `
	CREATE TABLE tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		priority TEXT,
		due_date TEXT,
		tags TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Insert test data
	db.Exec(`INSERT INTO tasks (title, completed, created_at, updated_at, priority, due_date, tags) VALUES
		("Task 1", 0, "2025-05-01T12:00:00Z", "2025-05-01T12:00:00Z", "high", "2025-06-01", "api,test"),
		("Task 2", 1, "2025-05-02T12:00:00Z", "2025-05-02T12:00:00Z", "low", "2025-06-10", "golang,unit")`)

	SetDB(db)
	return db
}

func TestGetTasks_FilterByPriority(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/tasks", GetTasks)

	req, _ := http.NewRequest("GET", "/tasks?priority=high", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	if !strings.Contains(resp.Body.String(), `"title":"Task 1"`) {
		t.Errorf("Expected response to include Task 1: %s", resp.Body.String())
	}
	if strings.Contains(resp.Body.String(), `"title":"Task 2"`) {
		t.Errorf("Did not expect Task 2: %s", resp.Body.String())
	}
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.POST("/tasks", middleware.ValidateTaskInput(), CreateTask)

	payload := `{
		"title": "New Test Task",
		"priority": "medium",
		"due_date": "2025-06-15",
		"tags": "test,api"
	}`

	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), `"title":"New Test Task"`) {
		t.Errorf("Response doesn't contain task title: %s", resp.Body.String())
	}
}

func TestMarkTaskDone(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.PUT("/tasks/:id/done", MarkTaskDone)

	req := httptest.NewRequest("PUT", "/tasks/1/done", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), `"message":"Task marked as done"`) {
		t.Errorf("Expected success message, got %s", resp.Body.String())
	}
}

func TestMarkTaskUndone(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.PUT("/tasks/:id/undone", MarkTaskUndone)

	req := httptest.NewRequest("PUT", "/tasks/1/undone", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), `"message":"Task marked as undone"`) {
		t.Errorf("Expected success message, got %s", resp.Body.String())
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.DELETE("/tasks/:id", DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp.Code)
	}
	if !strings.Contains(resp.Body.String(), `"message":"Task deleted"`) {
		t.Errorf("Expected success message, got %s", resp.Body.String())
	}
}

func TestCreateTask_MissingTitle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.POST("/tasks", middleware.ValidateTaskInput(), CreateTask)

	payload := `{"priority":"high"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest && resp.Code != http.StatusInternalServerError {
		t.Errorf("Expected 400 or 500, got %d", resp.Code)
	}
}

func TestMarkTaskDone_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.PUT("/tasks/:id/done", MarkTaskDone)

	req := httptest.NewRequest("PUT", "/tasks/999/done", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound && resp.Code != http.StatusOK {
		t.Errorf("Expected 404 or 200, got %d", resp.Code)
	}
}

func TestDeleteTask_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.DELETE("/tasks/:id", DeleteTask)

	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound && resp.Code != http.StatusOK {
		t.Errorf("Expected 404 or 200, got %d", resp.Code)
	}
}

func TestGetTasks_InvalidFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := gin.Default()
	router.GET("/tasks", GetTasks)

	req := httptest.NewRequest("GET", "/tasks?priority=superhard", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.Code)
	}

	if !strings.Contains(resp.Body.String(), "[]") {
		t.Errorf("Expected empty result, got %s", resp.Body.String())
	}
}
