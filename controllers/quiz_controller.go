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
	eventID := c.Param("id")

	var quiz models.Quiz

	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, question, options FROM quizzes WHERE id = $1`, eventID,
	).Scan(&quiz.ID, &quiz.Question, &quiz.Options)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": quiz})
}