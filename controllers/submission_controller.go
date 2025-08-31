package controllers

import (
	"context"
	"database/sql"
	"eventify/config"
	"eventify/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubmitQuiz(c *gin.Context) {
	userID := c.GetInt("user_id")
	eventID := c.Param("id")

	rows, err := config.DB.Query(
		context.Background(),
		`SELECT id, answer_key FROM quizzes WHERE event_id = $1`, eventID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch quiz"})
		return
	}
	defer rows.Close()

	type Answer struct {
		QuizID int    `json:"quiz_id"`
		Answer string `json:"answer"`
	}
	var answers []Answer
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	correctAnswers := make(map[int]string)
	for rows.Next() {
		var qid int
		var ans string
		if err := rows.Scan(&qid, &ans); err == nil {
			correctAnswers[qid] = ans
		}
	}

	totalQuestions := len(correctAnswers)
	if totalQuestions == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no quizzes for this event"})
		return
	}

	pointPerQuestion := 100 / totalQuestions
	score := 0
	correctCount := 0

	for _, ans := range answers {
		if correctAnswers[ans.QuizID] == ans.Answer {
			score += pointPerQuestion
			correctCount++
		}
	}

	if correctCount == totalQuestions {
		score = 100
	}

	status := "failed"
	if score >= 50 {
		status = "passed"
	}

	_, err = config.DB.Exec(
		context.Background(),
		`INSERT INTO quiz_submissions (user_id, event_id, score, status) 
		 VALUES ($1, $2, $3, $4)`,
		userID, eventID, score, status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save submission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "quiz submitted",
		"score":   score,
		"status":  status,
	})
}


func MyQuizSubmissions(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := config.DB.Query(
		context.Background(),
		`SELECT q.id, q.user_id, u.name, q.event_id, e.title, q.score, q.status, q.submitted_at 
		 FROM quiz_submissions q 
		 JOIN users u ON q.user_id = u.id
		 JOIN events e ON q.event_id = e.id
		 WHERE q.user_id = $1 ORDER BY q.submitted_at DESC`, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var submissions []models.QuizSubmission
	for rows.Next() {
		var sub models.QuizSubmission
		var score sql.NullInt32
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.Username ,&sub.EventID, &sub.EventTitle, &score, &sub.Status, &sub.SubmittedAt); err == nil {
			if score.Valid {
				val := int(score.Int32)
				sub.Score = &val
			}
			submissions = append(submissions, sub)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": submissions})
}

func ListQuizSubmissionsByEvent(c *gin.Context) {
	eventID := c.Param("id")

	rows, err := config.DB.Query(
		context.Background(),
		`SELECT id, user_id, event_id, score, status, submitted_at 
		 FROM quiz_submissions WHERE event_id = $1 ORDER BY submitted_at DESC`, eventID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch submissions"})
		return
	}
	defer rows.Close()

	var submissions []models.QuizSubmission
	for rows.Next() {
		var sub models.QuizSubmission
		var score sql.NullInt32
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.EventID, &score, &sub.Status, &sub.SubmittedAt); err == nil {
			if score.Valid {
				val := int(score.Int32)
				sub.Score = &val
			}
			submissions = append(submissions, sub)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": submissions})
}
