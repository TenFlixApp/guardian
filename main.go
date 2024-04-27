package main

import (
	"log"
	"time"

	"guardian/data"
	"guardian/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	data.ConnectToDB()
	defer data.CloseDB()

	configCors := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "creditential", "Content-Length", "Content-Type", "Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router := gin.Default()
	router.Use(cors.New(configCors))

	router.POST("/register", routes.RegisterRoute)
	router.POST("/login", routes.LoginRoute)
	router.GET("/metrics", routes.GetDashboardStatsRoute)

	router.Run(":8080")
}
