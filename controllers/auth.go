package controllers

import (
	"context"
	"eventify/config"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var in struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), 10)

	_, err := config.DB.Exec(
		context.Background(),
		`INSERT INTO users (name, email, password, role) VALUES ($1,$2,$3,'peserta')`,
		in.Name, in.Email, string(hash),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "register ok"})
}

func Login(c *gin.Context) {
	var in struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id int
	var hashed string
	var role string
	err := config.DB.QueryRow(
		context.Background(),
		`SELECT id, password, role FROM users WHERE email=$1`,
		in.Email,
	).Scan(&id, &hashed, &role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email/password"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(in.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email/password"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	c.JSON(http.StatusOK, gin.H{"message": "login ok", "token": signed})
}
