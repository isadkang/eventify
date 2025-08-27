package controllers

import (
	"context"
	"eventify/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminDashboard(c *gin.Context) {
	ctx := context.Background()

	var totalUser, totalEvent, totalSubmissions, totalTickets int
	var passed, failed, approved, pending, rejected int

	config.DB.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&totalUser)
	config.DB.QueryRow(ctx, `SELECT COUNT(*) FROM events`).Scan(&totalEvent)

	config.DB.QueryRow(ctx, `SELECT COUNT(*) FROM quiz_submissions`).Scan(&totalSubmissions)
	config.DB.QueryRow(ctx, `SELECT COUNT(*) FROM quiz_submissions WHERE status='passed'`).Scan(&passed)
	config.DB.QueryRow(ctx, `SELECT COUNT(*) FROM quiz_submissions WHERE status='failed'`).Scan(&failed)
	
	config.DB.QueryRow(ctx, `SELET COUNT(*) FROM tickets`).Scan(&totalTickets)
	config.DB.QueryRow(ctx, `SELET COUNT(*) FROM tickets WHERE status='approved'`).Scan(&approved)
	config.DB.QueryRow(ctx, `SELET COUNT(*) FROM tickets WHERE status='pending'`).Scan(&pending)
	config.DB.QueryRow(ctx, `SELET COUNT(*) FROM tickets WHERE status='rejected'`).Scan(&rejected)


	c.JSON(http.StatusOK, gin.H{
		"users": gin.H{
			"total": totalUser,
		},
		"events": gin.H{
			"total": totalEvent,
		},
		"quiz_submissions": gin.H{
			"total": totalSubmissions,
			"quiz passed": passed,
			"quiz failed": failed,
		},
		"tickets": gin.H{
			"total": totalTickets,
			"ticket approved": approved,
			"ticket pending": pending,
			"ticket rejected": rejected,
		},
	})
}