package main

import (
	"refrij/controllers"
	"refrij/initializers"
	"refrij/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	//Setup gin app
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.POST("/api/signup", controllers.SignUp)
	r.POST("/api/login", controllers.Login)
	r.PUT("/api/change-password/:user_id", middleware.RequireAuth, controllers.ChangePassword)
	r.PUT("/api/update/:user_id", middleware.RequireAuth, controllers.UpdateProfile)
	r.PUT("/api/user/:user_id", middleware.RequireAuth, controllers.GetUserDetail)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()

}
