package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/products", CreateProduct)
	router.PUT("/products/:id", UpdateProduct)
	router.DELETE("/products/:id", DeleteProduct)
	router.Run("localhost:8080")
}
