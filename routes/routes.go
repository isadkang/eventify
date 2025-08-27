package routes

import (
	"eventify/controllers"
	"eventify/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "🚀 API IS RUNNING"})
	})

	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	api := r.Group("/api")
	api.Use(middlewares.Auth(""))
	{
		api.GET("/me", controllers.Me)
	}

	admin := r.Group("/api/admin")
	admin.Use(middlewares.Auth("admin"))
	{
		admin.GET("/dashboard", controllers.AdminDashboard)

		admin.GET("/users", controllers.GetAllUser)
		admin.GET("/users/:id", controllers.GetUserById)
	}

	return r

}