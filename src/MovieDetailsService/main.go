package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	tmdbAccessToken := os.Getenv("TMDB_ACCESS_TOKEN")
	if tmdbAccessToken == "" {
		log.Fatal("TMDB_ACCESS_TOKEN environment variable not set.")
	}
	tmdb := NewTMDbClient(tmdbAccessToken)
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	router.GET("/details/:id", func(c *gin.Context) {
		id := c.Param("id")

		details, err := tmdb.getMovieDetails(id)
		if err != nil {
			log.Printf("Error fetching movie details for ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie details", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, details)
	})

	router.Run(":3300")
}
