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

		api.POST("/events/:id/join", controllers.JoinEvent)
		api.GET("/tickets", controllers.MyTickets)

		api.GET("/events/:id/quizzes", controllers.GetQuizByEvent)

		api.POST("/events/:id/quizzes/submit", controllers.SubmitQuiz)
		api.GET("/quizzes/submissions/me", controllers.MyQuizSubmissions)
	}

	admin := r.Group("/api/admin")
	admin.Use(middlewares.Auth("admin"))
	{
		admin.GET("/dashboard", controllers.AdminDashboard)

		admin.GET("/users", controllers.GetAllUser)
		admin.GET("/users/:id", controllers.GetUserById)

		admin.GET("/events", controllers.ListEvents)
		admin.POST("/events", controllers.CreateEvent)
		admin.GET("/events/:id", controllers.GetEvent)

		admin.GET("/tickets", controllers.ListTickets)
		admin.PUT("/tickets/:id/approve", controllers.ApproveTicket)
		admin.PUT("/tickets/:id/reject", controllers.RejectTicket)

		api.POST("/events/:id/quizzes", controllers.CreateQuiz)
		admin.GET("/events/:id/quizzes/submissions", controllers.ListQuizSubmissionsByEvent)
	}

	return r

}
