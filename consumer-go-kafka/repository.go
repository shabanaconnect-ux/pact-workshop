package main

import "errors"

// ProductRepository is an in-memory product database.
type ProductRepository struct {
	Products map[string]*Product
}

// ByID finds a product by its ID
func (u *ProductRepository) ByID(ID string) (*Product, error) {
	for _, product := range u.Products {
		if product.ID == ID {
			return product, nil
		}
	}
	return nil, ErrNotFound
}
// CreateProduct adds a new product to the repository
func (u *ProductRepository) Insert(product *Product) error {
	if product.ID == "" {
		return ErrEmpty
	}
	if _, exists := u.Products[product.ID]; exists {
		return errors.New("product already exists")
	}
	if u.Products == nil {
		u.Products = make(map[string]*Product)
	}
	u.Products[product.ID] = product
	return nil
}

// DeleteProduct removes a product from the repository by its ID
func (u *ProductRepository) Delete(ID string) error {
	if ID == "" {
		return ErrEmpty
	}
	if _, exists := u.Products[ID]; !exists {
		return ErrNotFound
	}
	delete(u.Products, ID)
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
