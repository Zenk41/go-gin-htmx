package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Zenk41/go-gin-htmx/api"
	"github.com/Zenk41/go-gin-htmx/firebase"
	"github.com/Zenk41/go-gin-htmx/handlers"
	"github.com/Zenk41/go-gin-htmx/middlewares"
	"github.com/Zenk41/go-gin-htmx/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type FirebaseConfig struct {
	ProjectID string `json:"project_id"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("apiKey environment variable not set")
	}

	domain := os.Getenv("DOMAIN")
	if apiKey == "" {
		log.Fatal("DOMAIN environment variable not set")
	}

	firebaseServiceAccount := os.Getenv("SERVICE_ACCOUNT_FILE")
	if firebaseServiceAccount == "" {
		log.Fatal("Firebase Project ID variable not set")
	}

	httpAddr := flag.String("addr", "0.0.0.0"+port, "Listen address")
	flag.Parse()

	gin.SetMode(gin.DebugMode) // Ensure gin is in debug mode
	gin.DefaultWriter = os.Stdout

	// Path to your JSON configuration file
	filePath := firebaseServiceAccount

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal the JSON data into the struct
	var config FirebaseConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	fireStoreClient, err := firebase.Firestore(firebaseServiceAccount, config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer fireStoreClient.Close()

	firebaseApi := api.NewFirebaseApi(apiKey)

	userRepo := models.NewUserRepository(fireStoreClient)
	taskRepo := models.NewTaskRepository(fireStoreClient)

	firebaseAuth, err := firebase.Auth(firebaseServiceAccount)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	userHandler := handlers.NewUserHandler(userRepo, apiKey, firebaseApi, domain)
	taskHandler := handlers.NewTaskHandler(taskRepo, userRepo, firebaseAuth)
	pageHandler := handlers.NewPageHandler(userRepo, taskRepo, firebaseApi, firebaseAuth)

	routesInit := handlerList{
		userHandler: userHandler,
		taskHandler: taskHandler,
		pageHandler: pageHandler,
	}

	e := gin.New()

	routesInit.RoutesRegister(e)

	e.SetTrustedProxies(nil)

	if err := e.Run(*httpAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

type handlerList struct {
	userHandler handlers.UserHandler
	taskHandler handlers.TaskHandler
	pageHandler handlers.PageHandler
}

func (hl *handlerList) RoutesRegister(e *gin.Engine) {
	e.Use(gin.Recovery())                 // Add recovery middleware for panic recovery
	e.Use(middlewares.StructuredLogger()) // Apply logging middleware

	e.Static("/public", "./public")

	// Pages
	// auth page
	e.GET("/login", hl.pageHandler.Login)
	e.GET("/register", hl.pageHandler.Register)
	// home page
	e.GET("/", hl.pageHandler.Home)

	// task
	task := e.Group("/task")
	task.POST("", hl.taskHandler.CreateNewTask)
	task.PUT("", hl.taskHandler.EditTaskById)
	task.PUT("/:id/done", hl.taskHandler.DoneTaskById)
	task.DELETE("/:id", hl.taskHandler.DeleteTaskById)
	task.POST("/update", hl.taskHandler.GetTasksByDate)
	task.PUT("/done-all", hl.taskHandler.DoneAllTaskDayByDate)

	// component
	comp := e.Group("/component")
	comp.POST("/task-edit", hl.taskHandler.EditTaskModal)
	comp.POST("/task-delete", hl.taskHandler.DeleteTaskModal)

	// auth
	auth := e.Group("/auth")
	auth.POST("/register", hl.userHandler.Register)
	auth.POST("/login", hl.userHandler.Login)
	auth.POST("/logout", hl.userHandler.Logout)

	// validates
	validates := e.Group("/validate")
	validates.POST("/email", handlers.Make(handlers.ValidateEmailHandler))
	validates.POST("/password", handlers.Make(handlers.ValidatePasswordHandler))
}
