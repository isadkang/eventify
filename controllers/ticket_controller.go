package controllers

import (
	"context"
	"eventify/config"
	"eventify/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JoinEvent(c *gin.Context) {
	userID := c.GetInt("user_id")
	eventID := c.Param("id")
	
	fmt.Println("DEBUG userID:", userID, "eventID:", eventID)
	// cek apakah event ada
	var exists bool
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM events WHERE id=$1)",
		eventID,
	).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	// cek apakah user udh join
	var count int
	err = config.DB.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM tickets WHERE user_id=$1 AND event_id=$2",
		userID, eventID,
	).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket already exists"})
		return
	}

	// insert ticket
	_, err = config.DB.Exec(
		context.Background(),
		"INSERT INTO tickets (user_id, event_id, status) VALUES ($1, $2, 'pending')",
		userID, eventID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ticket created, pending approval"})
}

func MyTickets(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := config.DB.Query(
		context.Background(),
		`SELECT t.id, t.event_id, t.status, e.title 
		FROM tickets t 
		JOIN events e ON e.id = t.event_id
		WHERE t.user_id=$1`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.EventID, &t.Status, &t.EventName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		t.UserID = userID
		tickets = append(tickets, t)
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func ListTickets(c *gin.Context) {
	rows, err := config.DB.Query(
		context.Background(),
		`SELECT t.id, t.user_id, u.name, t.event_id, e.title, t.status 
		FROM tickets t
		JOIN users u ON t.user_id = u.id
		JOIN events e ON t.event_id = e.id`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.UserID, &t.Username, &t.EventID, &t.EventName, &t.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tickets = append(tickets, t)
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func ApproveTicket(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM tickets WHERE id=$1)",
		id,
	).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	_, err = config.DB.Exec(
		context.Background(),
		"UPDATE tickets SET status='approved' WHERE id=$1",
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket approved"})
}

func RejectTicket(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM tickets WHERE id=$1)",
		id,
	).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	_, err = config.DB.Exec(
		context.Background(),
		"UPDATE tickets SET status='rejected' WHERE id=$1",
		id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket rejected"})
}
