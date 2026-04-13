package repository

import (
	"context"
	"database/sql"
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

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Price,
		product.ProductType.String(),
		product.Description,
		product.Quantity,
	)
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return model.Product{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	product.ID = lastID
	return product, nil
}

func (r *productRepository) Fetch(ctx context.Context, filter ProductFilter) ([]model.Product, int64, error) {
	page := filter.Page
	if page < 1 {
		page = DefaultPage
	}

	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	offset := (page - 1) * pageSize

	whereClause := ""
	args := []interface{}{}

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
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
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

	if filter.Order == "asc" {
		orderClause += " ASC"
	} else {
		orderClause += " DESC"
	}

	limitOffsetClause := fmt.Sprintf(" LIMIT ? OFFSET ?")
	args = append(args, pageSize, offset)

	query := `
		SELECT id, name, price, product_type, description, quantity, created_at, updated_at
		FROM products
	` + whereClause + orderClause + limitOffsetClause

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query products: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan product row: %w", err)
		}

		product.ProductType = model.ProductType(productTypeStr)
		if description.Valid {
			product.Description = description.String
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating product rows: %w", err)
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

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&productTypeStr,
		&product.Description,
		&product.Quantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, pkgerr.NotFoundErrorf("product with id %d not found", id)
		}
		return model.Product{}, fmt.Errorf("failed to find product by id: %w", err)
	}

	product.ProductType = model.ProductType(productTypeStr)
	return product, nil
}
