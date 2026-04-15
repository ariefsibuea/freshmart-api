package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ariefsibuea/freshmart-api/internal/model"
	pkgerr "github.com/ariefsibuea/freshmart-api/internal/pkg/errors"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product model.Product) (model.Product, error) {
	query := `
		INSERT INTO products (name, price, product_type, description, quantity)
		VALUES (?, ?, ?, ?, ?)
	`

	var description sql.NullString
	if product.Description != nil {
		description = sql.NullString{String: *product.Description, Valid: true}
	}

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Price,
		product.ProductType.String(),
		description,
		product.Quantity,
	)
	if err != nil {
		return model.Product{}, err
	}

	createdID, err := result.LastInsertId()
	if err != nil {
		return model.Product{}, err
	}

	return r.Get(ctx, createdID)
}

func (r *productRepository) Fetch(ctx context.Context, filter model.ProductFilter) ([]model.Product, int64, error) {
	offset := (filter.Page - 1) * filter.PageSize

	whereClause := ""
	args := []any{}

	if filter.Name != "" {
		whereClause += " WHERE name LIKE ?"
		args = append(args, "%"+filter.Name+"%")
	}

	if filter.ProductType.IsValid() {
		if whereClause == "" {
			whereClause += " WHERE product_type = ?"
		} else {
			whereClause += " AND product_type = ?"
		}
		args = append(args, filter.ProductType.String())
	}

	countQuery := "SELECT COUNT(*) FROM products" + whereClause
	totalProducts := int64(0)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalProducts); err != nil {
		return nil, 0, fmt.Errorf("count products failed: %w", err)
	}

	orderClause := " ORDER BY "
	switch filter.SortBy {
	case "price":
		orderClause += "price"
	case "name":
		orderClause += "name"
	case "date":
		orderClause += "created_at"
	default:
		orderClause += "created_at"
	}

	// NOTE: add id as unique column when sorting the products in case
	// the remaining sort columns have the same values across products.
	if filter.Order == "asc" {
		orderClause += " ASC, id ASC"
	} else {
		orderClause += " DESC, id DESC"
	}

	limitOffsetClause := " LIMIT ? OFFSET ?"
	args = append(args, filter.PageSize, offset)

	query := `
		SELECT id, name, price, product_type, description, quantity, created_at, updated_at
		FROM products
	` + whereClause + orderClause + limitOffsetClause

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var product model.Product
		var productTypeStr string
		var description sql.NullString

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&productTypeStr,
			&description,
			&product.Quantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		product.ProductType = model.ProductType(productTypeStr)
		if description.Valid {
			product.Description = &description.String
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, totalProducts, nil
}

func (r *productRepository) Get(ctx context.Context, id int64) (model.Product, error) {
	query := `
		SELECT id, name, price, product_type, description, quantity, created_at, updated_at
		FROM products
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var product model.Product
	var productTypeStr string
	var description sql.NullString

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&productTypeStr,
		&description,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Product{}, pkgerr.NotFoundErrorf("product with id %d not found", id)
		}
		return model.Product{}, err
	}

	product.ProductType = model.ProductType(productTypeStr)
	if description.Valid {
		product.Description = &description.String
	}
	return product, nil
}
