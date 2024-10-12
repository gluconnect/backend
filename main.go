package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stanleymw/glucose/endpoints"
)

func main() {
	db, err := bolt.Open("gluconnect.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Use(cors.Default())

	router.POST("/register", endpoints.Register(db))

	authorized := router.Group("/")
	authorized.Use(endpoints.Auth(db))
	{
		authorized.POST("/add_reading", endpoints.AddReading(db))
		authorized.POST("/get_readings", endpoints.GetReadings(db))

		// returns 200 if authorized
		authorized.POST("/verify", endpoints.Verify())
	}

	router.Run("0.0.0.0:24816")
}
