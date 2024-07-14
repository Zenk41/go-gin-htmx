package main

import (
	"flag"
	"log"
	"os"

	"github.com/Zenk41/go-gin-htmx/handlers"
	"github.com/Zenk41/go-gin-htmx/middlewares"
	"github.com/Zenk41/go-gin-htmx/views/home"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	httpAddr := flag.String("addr", "0.0.0.0"+port, "Listen address")
	flag.Parse()

	gin.SetMode(gin.DebugMode) // Ensure gin is in debug mode
	gin.DefaultWriter = os.Stdout

	route := Routes()

	route.SetTrustedProxies(nil)

	if err := route.Run(*httpAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func Routes() *gin.Engine {
	router := gin.New()                        // Create a new gin router without default middleware
	router.Use(gin.Recovery())                 // Add recovery middleware for panic recovery
	router.Use(middlewares.StructuredLogger()) // Azpply logging middleware


	router.Static("/public", "./public")

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"aku": "kamu",
		})
	})

	router.GET("/home", handlers.Make(func(ctx *gin.Context) error {
		return handlers.Render(ctx, home.Index())
	}))

	return router
}
