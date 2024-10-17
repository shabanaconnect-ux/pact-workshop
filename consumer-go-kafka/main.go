package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

var productRepository = &ProductRepository{
	Products: map[string]*Product{
		"10": {
			Name:    "Pizza",
			ID:      "10",
			Type:    "food",
			Version: "1.0.0",
		},
	},
}

func main() {
	go KafkaConsumer()
	go setupRouter()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Shutting down server...")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/product/:id", GetProduct)
	go func() {
		if err := router.Run("localhost:8080"); err != nil {
			fmt.Printf("Failed to run server: %v\n", err)
		}
	}()
	return router
}
