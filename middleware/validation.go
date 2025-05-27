package middleware

import (
	"fmt"
	"strings"
	"time"

	"taskapi/models"

	"github.com/gin-gonic/gin"
)

func ValidateTaskInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			fmt.Printf("ValidateTaskInput ShouldBindJSON failed: %v", err)
			c.JSON(400, gin.H{"error": "Invalid JSON body"})
			c.Abort()
			return
		}

		if strings.TrimSpace(task.Title) == "" {
			fmt.Printf("ValidateTaskInput Title is required")
			c.JSON(400, gin.H{"error": "Title is required"})
			c.Abort()
			return
		}

		// Validate priority
		validPriorities := map[string]bool{"low": true, "medium": true, "high": true, "": true}
		if !validPriorities[strings.ToLower(task.Priority)] {
			fmt.Printf("ValidateTaskInput Invalid priority. Use low, medium, or high")
			c.JSON(400, gin.H{"error": "Invalid priority. Use low, medium, or high"})
			c.Abort()
			return
		}

		// Optional: validate due_date (if provided)
		if task.DueDate != "" {
			_, err := time.Parse("2006-01-02", task.DueDate)
			if err != nil {
				fmt.Printf("ValidateTaskInput Invalid due_date format. Use YYYY-MM-DD")
				c.JSON(400, gin.H{"error": "Invalid due_date format. Use YYYY-MM-DD"})
				c.Abort()
				return
			}
		}

		// Save validated task in context for handler to use
		c.Set("task", task)

		c.Next()
	}
}
