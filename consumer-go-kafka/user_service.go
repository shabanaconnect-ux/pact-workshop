package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Product is the domain object
type Product struct {
	ID      string `json:"id" pact:"example=10"`
	Name    string `json:"name" pact:"example=pizza"`
	Type    string `json:"type" pact:"example=food"`
	Version string `json:"version"`
}

// ProductEvent extends Product and adds event information
type ProductEvent struct {
	Product
	Event string `json:"event" pact:"example=CREATED"`
}

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

// GetProduct fetches a product if authenticated and exists
func GetProduct(c *gin.Context) {
	product, err := productRepository.ByID(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "file not found"})
	} else {
		c.JSON(http.StatusOK, product)
	}
}
