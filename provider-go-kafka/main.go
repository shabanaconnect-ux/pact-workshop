package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func main() {
	go setupRouter()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Shutting down server...")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/products", CreateProduct)
	router.PUT("/products/:id", UpdateProduct)
	router.DELETE("/products/:id", DeleteProduct)
	go func() {
		if err := router.Run("localhost:8081"); err != nil {
			fmt.Printf("Failed to run server: %v\n", err)
		}
	}()
	return router
}
