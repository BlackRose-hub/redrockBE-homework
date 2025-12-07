package main

import (
	"course-system/config"
	"course-system/handlers"
	"course-system/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/courses", handlers.GetCourses)
		auth.POST("/enroll", handlers.EnrollCourse)
		auth.GET("/my-courses", handlers.GetMyCourses)
		auth.DELETE("/drop/:course_id", handlers.DropCourse)

		admin := auth.Group("/")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/courses", handlers.CreateCourse)
		}
	}

	log.Println("🚀 服务器启动在 http://localhost:8081")

	r.Run(":8081")
}
