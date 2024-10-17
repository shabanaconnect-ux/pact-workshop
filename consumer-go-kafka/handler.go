package main

import (
	"errors"
	"fmt"
)

func handler(product Product, event string) error {
	fmt.Println("received product:", product)
	fmt.Println("received product event:", event)

	switch event {
	case "CREATED", "UPDATED":
		return productRepository.Insert(&product)
	case "DELETED":
		return productRepository.Delete(product.ID)
	default:
		return errors.New("unable to process event")
	}
}
