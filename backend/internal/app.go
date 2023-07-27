package internal

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	db *sql.DB
}

// NewApp to create and initialize app
func NewApp(path string) (*App, error) {
	app := &App{}
	err := app.connectDatabase(path)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {

	router := gin.Default()

	// err := app.createTable()
	// if err != nil {
	// 	return err
	// }

	// Configure CORS middleware with desired options
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // Replace with false to restrict allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/tasks", app.GetTasks)
	router.POST("/tasks", app.AddTask)
	router.DELETE("/tasks/:id", app.DeleteTask)
	router.PUT("/tasks", app.UpdateTask)

	router.Run(":3000")

	fmt.Println("Server started on port: 3000")
	err := http.ListenAndServe(":3000", nil)
	return err
}
