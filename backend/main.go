package main

import (
	"Lpanel/backend/config"
	"Lpanel/backend/handlers"
	"Lpanel/backend/models"
	"Lpanel/backend/utils"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Initialize Database
	config.InitDB()

	// Start Site Runner
	var sites []models.Site
	config.DB.Find(&sites)
	utils.Runner.UpdateSites(sites)

	r := gin.Default()

	// Serve Frontend
	r.Use(static.Serve("/", static.LocalFile("./frontend", true)))

	// API Routes
	api := r.Group("/api")
	{
		sites := api.Group("/sites")
		{
			sites.GET("", handlers.ListSites)
			sites.POST("", handlers.CreateSite)
			sites.DELETE("/:id", handlers.DeleteSite)
			sites.POST("/:id/toggle", handlers.ToggleSite)
		}

		system := api.Group("/system")
		{
			system.GET("/status", handlers.GetSystemStatus)
		}

		files := api.Group("/files")
		{
			files.GET("/list", handlers.ListFiles)
			files.DELETE("", handlers.DeleteFile)
		}
	}

	log.Println("Lpanel server starting on :8080")
	r.Run(":8080")
}
