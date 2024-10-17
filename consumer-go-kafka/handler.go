package main

import (
	"errors"
	"fmt"
)

var repository ProductRepository

func handler(product Product, event string) error {
	fmt.Println("received product:", product)
	fmt.Println("received product event:", event)

	switch event {
	case "CREATED", "UPDATED":
		return repository.Insert(&product)
	case "DELETED":
		return repository.Delete(product.ID)
	default:
		return errors.New("unable to process event")
	}
}
