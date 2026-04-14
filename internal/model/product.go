package model

import (
	"strconv"
	"strings"
	"time"

	pkgerr "github.com/ariefsibuea/freshmart-api/internal/pkg/errors"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
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

var ValidSortBy = map[string]bool{
	"":      true, // default will be used
	"price": true,
	"name":  true,
	"date":  true,
}

var ValidOrder = map[string]bool{
	"":     true, // default will be used
	"asc":  true,
	"desc": true,
}

type Product struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Price       float64     `json:"price"`
	ProductType ProductType `json:"product_type"`
	Description *string     `json:"description,omitempty"`
	Quantity    int         `json:"quantity"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
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

type ProductFilter struct {
	Name        string
	ProductType ProductType
	SortBy      string
	Order       string
	Page        int
	PageSize    int
}

func NewProductFilter(name, productType, page, pageSize, sortBy, order string) (ProductFilter, error) {
	filter := ProductFilter{
		Name:     name,
		SortBy:   sortBy,
		Order:    order,
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
	}

	if productType != "" {
		filter.ProductType = ProductType(productType)
	}

	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			return ProductFilter{}, pkgerr.BadRequestError("invalid query parameter 'page'")
		}
		if p > 0 {
			filter.Page = p
		}
	}

	if pageSize != "" {
		ps, err := strconv.Atoi(pageSize)
		if err != nil {
			return ProductFilter{}, pkgerr.BadRequestError("invalid query parameter 'page_size'")
		}
		if ps > 0 {
			filter.PageSize = ps
		}
		if ps > MaxPageSize {
			filter.PageSize = MaxPageSize
		}
	}

	return filter, nil
}

func (f *ProductFilter) Validate() error {
	if f.ProductType.String() != "" && !f.ProductType.IsValid() {
		return pkgerr.BadRequestErrorf("invalid product_type: must be one of %v", ValidProductTypes)
	}

	if !ValidSortBy[f.SortBy] {
		return pkgerr.BadRequestError("invalid sort_by: must be one of 'price', 'name', 'date'")
	}

	if !ValidOrder[f.Order] {
		return pkgerr.BadRequestError("invalid order: must be one of 'asc', 'desc'")
	}

	return nil
}
