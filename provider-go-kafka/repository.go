package main

import "errors"

// Save adds or updates a product in the repository
func Save(product *Product) error {
	KakfaProducer(createEvent(*product, "CREATED"))
	return nil
}

// Delete removes a product from the repository by its ID
func Delete(product *Product) error {
	KakfaProducer(createEvent(*product, "DELETED"))
	return nil
}

// Update modifies an existing product in the repository
func Update(product *Product) error {
	KakfaProducer(createEvent(*product, "UPDATED"))
	return nil
}

var (
	// ErrNotFound represents a resource not found (404)
	ErrNotFound = errors.New("not found")

	// ErrUnauthorized represents a Forbidden (403)
	ErrUnauthorized = errors.New("unauthorized")

	// ErrEmpty is returned when input string is empty
	ErrEmpty = errors.New("empty string")
)
