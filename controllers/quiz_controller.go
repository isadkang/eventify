package controllers

import (
	"context"
	"eventify/config"
	"eventify/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateQuiz(c *gin.Context) {
	eventID := c.Param("id")

	var in models.Quiz

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := config.DB.Exec(
		context.Background(),
		`INSERT INTO quizzes (event_id, question, options, answer_key) 
		VALUES ($1, $2, $3, $4)`,
		eventID, in.Question, in.Options, in.AnswerKey,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "quiz created"})
}

func GetQuizByEvent(c *gin.Context) {
	userID := c.GetInt("user_id")
	eventID := c.Param("id")

	var joined bool
	err := config.DB.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM tickets WHERE user_id=$1 AND event_id=$2)",
		userID, eventID,
	).Scan(&joined)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !joined {
		c.JSON(http.StatusForbidden, gin.H{"error": "you must join this event to access quizzes"})
		return
	}

	rows, err := config.DB.Query(
		context.Background(),
		`SELECT id, event_id, question, options 
		FROM quizzes 
		WHERE event_id = $1`, eventID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var q models.Quiz
		if err := rows.Scan(&q.ID, &q.EventID, &q.Question, &q.Options); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		quizzes = append(quizzes, q)
	}

	if len(quizzes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no quizzes found for this event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": quizzes})
}
