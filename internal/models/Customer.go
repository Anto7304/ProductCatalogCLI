package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID uuid.UUID `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	FullName string `json:"full_name" db:"full_name"`
	Phone string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`


}


type CustomerRepository interface{
	CreateUser(customer *Customer)error

	GetById(id string)(*Customer,error)

	GetAll(page,limit int)([]Customer,error)

	Delete(id string)error


}
