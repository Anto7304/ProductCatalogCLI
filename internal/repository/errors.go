package repository

import "errors"


var (

	ErrNotFound = errors.New("pruduct not found")

	ErrDuplicateSku = errors.New("a product with this SKU already exists")

	ErrInvalid = errors.New("invalid product ID format")

	
)