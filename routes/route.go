package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/manlikehenryy/go-gin-project/controllers"
	"github.com/manlikehenryy/go-gin-project/database"
	"github.com/manlikehenryy/go-gin-project/middleware"
)

func Setup(app *gin.Engine) {
	
    controllers.InitDB(database.DB)


	app.POST("/api/register", controllers.Register)
	app.POST("/api/login", controllers.Login)
	app.GET("/api/logout", controllers.Logout)

	app.Use(middleware.IsAuthenticated)

	app.POST("/api/task", controllers.CreateTask)
	app.GET("/api/task", controllers.GetAllTasks)
	app.GET("/api/task/:id", controllers.GetTask)
	app.PUT("/api/task/:id", controllers.UpdateTask)
	app.GET("/api/user-tasks", controllers.UsersTask)
	app.DELETE("/api/task/:id", controllers.DeleteTask)
}