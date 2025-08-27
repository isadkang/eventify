package controllers

import (
	"context"
	"eventify/config"
	"eventify/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	uid, ok := c.Get("user_id")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "please login first"})
		return
	}

	var u models.User
	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, name, email, role, created_at FROM users WHERE id=$1`,
		uid,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": u})
}

func GetAllUser(c *gin.Context) {
	rows, err := config.DB.Query(
		context.Background(),
		`SELECT id, name, email, role, created_at FROM users`,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user"})
			return 
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	var u models.User

	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, name, email, role, created_at WHERE id=$1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": u})
}