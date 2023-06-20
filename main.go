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
	initializers.SyncDatabase()
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
	r.GET("/api/get-user-detail/:user_id", middleware.RequireAuth, controllers.GetUserDetail)

	r.GET("/api/getRefrigeratorDetail/:refrigerator_id", middleware.RequireAuth, controllers.GetRefrigeratorDetail)
	r.GET("/api/getRefiregators/:user_id", middleware.RequireAuth, controllers.GetRefiregators)
	r.POST("/api/createRefrigerator", middleware.RequireAuth, controllers.CreateRefrigerator)
	r.GET("/api/updateRefrigerator/:refrigerator_id", middleware.RequireAuth, controllers.UpdateRefrigerator)
	r.DELETE("/api/deleteRefrigerator/:refrigerator_id", middleware.RequireAuth, controllers.DeleteRefrigerator)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()

}
