package controllers

import (
	"context"
	"eventify/config"
	"eventify/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	var in models.Event

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(
		context.Background(),
		`INSERT INTO events (title, description, date, location, quota)
		VALUES ($1, $2, $3, $4, $5)`,
		in.Title, in.Description, in.Date, in.Location, in.Quota,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event created!"})
}

func ListEvents(c *gin.Context) {
	rows, err := config.DB.Query(
		context.Background(),
		`SELECT id, title, description, date, location, quota, created_at FROM events`,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch event"})
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event

		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date, &e.Location, &e.Quota, &e.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		events = append(events, e)
	}

	c.JSON(http.StatusOK, gin.H{"data": events})
}

func GetEvent(c *gin.Context) {
	id := c.Param("id")

	var e models.Event

	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, title, description, date, location, quota, created_at FROM events WHERE id=$1`,
		id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.Date, &e.Location, &e.Quota, &e.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var in models.Event

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var updated models.Event
	err := config.DB.QueryRow(
		context.Background(),
		`UPDATE events
		 SET title=$1, description=$2, date=$3, location=$4, quota=$5
		 WHERE id=$6
		 RETURNING id, title, description, date, location, quota, created_at`,
		in.Title, in.Description, in.Date, in.Location, in.Quota, id,
	).Scan(&updated.ID, &updated.Title, &updated.Description, &updated.Date,
		&updated.Location, &updated.Quota, &updated.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "event updated!",
		"data":    updated,
	})
}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	var deleted models.Event
	err := config.DB.QueryRow(
		context.Background(),
		`DELETE FROM events WHERE id=$1 
		RETURNING id, title`,
		id,
	).Scan(&deleted.ID, &deleted.Title)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event deleted", "data": deleted})
}
