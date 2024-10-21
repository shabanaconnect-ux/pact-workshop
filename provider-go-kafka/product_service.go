package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func incrementVersion(v string) string {
	var version int
	if v != "" {
		fmt.Sscanf(v, "v%d", &version)
		version++
	} else {
		version = 1
	}
	return fmt.Sprintf("v%d", version)
}

func createEvent(product Product, eventType string) ProductEvent {
	incrementVersion(product.Version)
	if product.ID == "" {
		product.ID = uuid.New().String()
	}
	return ProductEvent{
		Product: product,
		Event:   eventType,
	}
}

func CreateProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	Save(&product)
	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	if err := Update(&product); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "file not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	if err := Delete(&product); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "file not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
