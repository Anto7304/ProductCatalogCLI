package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"Anto7304.com/internal/models"
	"github.com/google/uuid"
)

// PostgresProductRepository implements models.Product
// this is the PostgreSql version of the repository
type PostgresProductRepository struct {
	db *sql.DB
}

// Create new repository instance through constructor function
// we are reutning to the interface for decoupling and easy to test with mock data
func NewRepositoryProductRepository(db *sql.DB) models.ProductRepository {
	return &PostgresProductRepository{
		db: db,
	}
}

//Create a product

func (r *PostgresProductRepository) Create(product *models.Product) error {
	//generate uuid
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}

	//sql query to insert a new product
	//using parametarize query to avoid sql injection
	query := `
		INSERT INTO products(id,name,description,price,stock_quantity,category,sku,weight)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	`

	//execute the query with the product fields as parameter
	_, err := r.db.Exec(
		query,
		product.ID,   //$1
		product.Name, //$2
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
		product.SKU,
		product.Weight,
	)

	if err != nil {
		//check for duplicate SKU error
		//PosgreSql error code 23505 is "unique viuolation"

		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "sku") {
			return ErrDuplicateSku
		}
		return fmt.Errorf("failed to create product %w", err)
	}
	return nil

}

func (r *PostgresProductRepository) GetBySKU(sku string) (*models.Product, error) {
	//query to select  a product by sku

	query := `
		SELECT id, name, description,price, stock_quantity,category, sku,weight,created_at,updated_at
		FROM products
		WHERE sku = $1 AND deleted_at IS NULL
	`

	//empty struct to hold the query results
	var product models.Product

	//scan allows writing to the empty product struct
	err := r.db.QueryRow(query, sku).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.Category,
		&product.SKU,
		&product.Weight,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		//no products found
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return &product, nil

}

// get product by id

func (r *PostgresProductRepository) GetById(id string) (*models.Product, error) {

	productID, err := uuid.Parse(id)

	if err != nil {
		return nil, ErrInvalid
	}

	query := `
		SELECT id, name, description, price, stock_quantity, category,sku,weight,created_at,updated_at
		FROM products

		WHERE id = $1 AND deleted_at IS NULL

	`

	var product models.Product

	err = r.db.QueryRow(query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.Category,
		&product.SKU,
		&product.Weight,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get by id: %w", err)
	}

	return &product, nil

}

func (r *PostgresProductRepository) GetAll(page, limit int) ([]models.Product, error) {

	//validate pagination parameter

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	// calculate the offset for SQL OFFSET clause
	//page 1: offset 0, page 2: offset limit,page 3 : offset 2* limit

	offset := (page - 1) * limit

	query := `
	
		SELECT id,name, description,price, stock_quantity, category,sku, weight, created_at, updated_at
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2

	
	
	`

	rows, err := r.db.Query(query, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	defer rows.Close() //always close rows when done

	var products []models.Product

	//itearate through the rows
	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Category,
			&product.SKU,
			&product.Weight,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan the product: %w", err)

		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating products: %w", err)
	}

	return products, nil

}

func (r *PostgresProductRepository) Update(product *models.Product) error {

	query := `
		UPDATE products
		SET name  = $1,description =$2,price = $3,stock_quantity = $4,category = $5, sku = $6,weight = $7,updated_at = CURRENT_TIMESTAMP
		WHERE id = $8 AND deleted_at IS NULL
	
	
	`
	result, err := r.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
		product.SKU,
		product.Weight,
		product.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "sku") {
			return ErrDuplicateSku
		}
		return fmt.Errorf("failed to update the product: %w", err)
	}

	//check if any rows were affected

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected : %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil

}

func (r *PostgresProductRepository) Delete(id string) error {
	productID, err := uuid.Parse(id)

	if err != nil {
		return ErrInvalid
	}

	query := `
	
		UPDATE products
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	
	`

	result, err := r.db.Exec(query, productID)

	if err != nil {
		return fmt.Errorf("failed to delete the product : %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil

}

func (r *PostgresProductRepository) Search(query string) ([]models.Product, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	//we add wildcards for partial matching

	// ILIKE is case sensitive LIKE

	searchPattern := "%" + strings.ToLower(query) + "%"

	sql := `
		SELECT id, name, description, price,stock_quantity, category,sku, weight,created_at, updated_at
		FROM products
		WHERE deleted_at IS NULL
		AND (
			LOWER(name) ILIKE $1 OR
			LOWER(description) ILIKE $1 OR
			LOWER(category) ILIKE $1
		)

		ORDER BY name
		LIMIT 50
	
	`
	rows, err := r.db.Query(sql, searchPattern)

	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Category,
			&product.SKU,
			&product.Weight,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return products, nil

}
