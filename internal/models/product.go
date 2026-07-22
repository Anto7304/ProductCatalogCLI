package models

import (
	"time"

	"github.com/google/uuid"
)

//represent data fields in database in sql

type Product struct{
	ID uuid.UUID `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price float64 `json:"price" db:"price"`
	StockQuantity int `json:"stock_quantity" db:"stock_quantity"`
	Category string `json:"category" db:"category"`
	SKU string `json:"sku" db:"sku"`
	Weight float64 `json:"weight" db:"weight"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

//product repository defines the operations we can perfom on products
//The interface - says what we can do, not how to do it

type ProductRepository interface{

	Create(product *Product)error

	GetById(id string) (*Product,error)

	GetBySKU(sku string)(*Product,error)

	GetAll(page,limit int)([]Product,error)
	GetByCategory(category string)([]Product,error)

	Update(product *Product)error

	Delete(id string) error

	Search(query string) ([]Product,error)

}