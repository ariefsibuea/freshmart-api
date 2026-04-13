package model

import (
	"database/sql"
	"strings"
	"time"

	pkgerr "github.com/ariefsibuea/freshmart-api/internal/pkg/errors"
)

type ProductType string

// List of product types based on enums in storage.
const (
	ProductTypeSayuran ProductType = "Sayuran"
	ProductTypeProtein ProductType = "Protein"
	ProductTypeBuah    ProductType = "Buah"
	ProductTypeSnack   ProductType = "Snack"
)

var ValidProductTypes = []ProductType{
	ProductTypeSayuran,
	ProductTypeProtein,
	ProductTypeBuah,
	ProductTypeSnack,
}

func (t ProductType) IsValid() bool {
	switch t {
	case ProductTypeSayuran, ProductTypeProtein, ProductTypeBuah, ProductTypeSnack:
		return true
	}
	return false
}

func (t ProductType) String() string {
	return string(t)
}

type Product struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Price       float64        `json:"price"`
	ProductType ProductType    `json:"product_type"`
	Description sql.NullString `json:"description"`
	Quantity    int            `json:"quantity"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string      `json:"name"`
	Price       float64     `json:"price"`
	ProductType ProductType `json:"product_type"`
	Description *string     `json:"description"`
	Quantity    int         `json:"quantity"`
}

func (r *CreateProductRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return pkgerr.ValidationError("name is required")
	}

	if r.Price <= 0 {
		return pkgerr.ValidationError("price must be greater than 0")
	}

	if !r.ProductType.IsValid() {
		return pkgerr.ValidationErrorf("invalid product_type: must be one of %v", ValidProductTypes)
	}

	if r.Quantity < 0 {
		return pkgerr.ValidationError("quantity must be greater than or equal to 0")
	}

	return nil
}
