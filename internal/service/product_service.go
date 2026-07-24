package service

import (
	"fmt"
	"strings"

	"Anto7304.com/internal/models"
	"Anto7304.com/internal/repository"
	"Anto7304.com/internal/utils"
	"github.com/google/uuid"
)

// product serviuces handle all the business logic
type ProduuctService struct {
	repo models.ProductRepository
}

func NewProductService(repo models.ProductRepository) *ProduuctService {
	return &ProduuctService{
		repo: repo,
	}
}

// create product to add a new product to the catalog
func (s *ProduuctService) CreateProduct(name, description, category, sku string, price, weight float64, stockQuantity int) (*models.Product, error) {
	//sanitize input
	name = utils.SanitizeString(name)
	description = utils.SanitizeString(description)
	category = utils.SanitizeString(category)
	sku = utils.SanitizeString(sku)

	//validate all inputs
	if err := utils.ValidateProductName(name); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := utils.ValidateCategory(category); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := utils.ValidatePrice(price); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := utils.ValidateSKU(sku); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := utils.ValidateStock(stockQuantity); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	//check if sku already exists
	existingProduct, err := s.repo.GetBySKU(sku)

	if err != nil && err != repository.ErrNotFound {
		return nil, fmt.Errorf("failed to check sku uniqueness: %w", err)
	}

	if existingProduct != nil {
		return nil, fmt.Errorf("SKU '%s' already exist for product %s", sku, existingProduct.Name)
	}

	//create the product
	product := &models.Product{
		ID:            uuid.New(),
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      category,
		SKU:           sku,
		Weight:        weight,
	}

	//store in database
	err = s.repo.Create(product)

	if err != nil {
		if err == repository.ErrDuplicateSku {
			return nil, fmt.Errorf("SKU '%s' is already taken. Please use a different SKU", sku)
		}
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil

}

func (s *ProduuctService) GetProduct(id string) (*models.Product, error) {

	//validate id format
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("product id cannot be empty")
	}

	//try to parse uuid to check the format
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %s", id)
	}

	//get from repository
	product, err := s.repo.GetById(id)

	if err != nil {
		if err == repository.ErrNotFound {
			return nil, fmt.Errorf("product with ID '%s' not found", id)
		}
		return nil, fmt.Errorf("failed to get the product: %w", err)
	}

	return product, nil

}

func (s *ProduuctService) GetProductBySKU(sku string) (*models.Product, error) {
	//validate sku
	sku = strings.TrimSpace(sku)

	if sku == "" {
		return nil, fmt.Errorf("SKU cannot be empty")
	}

	//get from repository
	product, err := s.repo.GetBySKU(sku)

	if err != nil {
		if err == repository.ErrNotFound {
			return nil, fmt.Errorf("product with SKU '%s' not found", sku)
		}
		return nil, fmt.Errorf("failed to get the product by SKU: %w", err)
	}

	return product, nil

}

func (s *ProduuctService) ListProducts(page, limit int) ([]models.Product, error) {
	//sanitize to avoid negative numbers to the DB
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	//get from repository
	products, err := s.repo.GetAll(page, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to list the products: %w", err)
	}

	return products, nil

}

func (s *ProduuctService) UpdateProduct(id string, name, description, category, sku *string, price, weight *float64, stockQuantity *int) (*models.Product, error) {

	//get the existing product
	existingProduct, err := s.GetProduct(id)
	if err != nil {
		return nil, err
	}

	//apply updateas if needed
	if name != nil {
		sanitizedName := utils.SanitizeString(*name)
		if err := utils.ValidateProductName(sanitizedName); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		existingProduct.Name = sanitizedName
	}

	if description != nil {
		existingProduct.Description = utils.SanitizeString(*description)
	}

	if category != nil {
		sanitizedCategory := utils.SanitizeString(*category)

		if err := utils.ValidateCategory(sanitizedCategory); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		existingProduct.Category = sanitizedCategory
	}

	if sku != nil {
		sanitizedSKU := strings.TrimSpace(*sku)
		if err := utils.ValidateSKU(sanitizedSKU); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}

		//check if new SKU conflicts with another
		if sanitizedSKU != existingProduct.SKU {
			otherProduct, err := s.repo.GetBySKU(sanitizedSKU)

			if err != nil && err != repository.ErrNotFound {
				return nil, fmt.Errorf("failed to check SKU uniqueness: %w", err)
			}

			if otherProduct != nil && otherProduct.ID != existingProduct.ID {
				return nil, fmt.Errorf("SKU '%s' is already taken by product %s", sanitizedSKU, otherProduct.Name)
			}

			existingProduct.SKU = sanitizedSKU
		}
	}

	if price != nil {
		if err := utils.ValidatePrice(*price); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		existingProduct.Price = *price
	}

	if stockQuantity != nil {
		if err := utils.ValidateStock(*stockQuantity); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		existingProduct.StockQuantity = *stockQuantity
	}

	if weight != nil {
		existingProduct.Weight = *weight
	}
	//save updates

	err = s.repo.Update(existingProduct)

	if err != nil {
		if err == repository.ErrDuplicateSku {
			return nil, fmt.Errorf("SKU '%s' is already taken. Please use different SKU", existingProduct.SKU)
		}
		return nil, fmt.Errorf("failed to update product: %w", err)

	}
	return existingProduct, nil

}

func (s *ProduuctService) DeleteProduct(id string) error {

	//validate
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("Product ID cannot be empty")
	}

	// try to parse the id
	_, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product id: %s", id)
	}

	//delete form the repository

	err = s.repo.Delete(id)

	if err != nil {
		if err == repository.ErrNotFound {
			return fmt.Errorf("Product with ID '%s' not found", id)
		}
		return fmt.Errorf("failed to delete the product: %w", err)
	}

	return nil

}

func (s *ProduuctService) SearchProducts(query string) ([]models.Product, error) {

	//validate
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	//check minimum

	if len(query) < 2 {
		return nil, fmt.Errorf("search query must be atleast 2 characters")
	}

	products, err := s.repo.Search(query)

	if err != nil {
		return nil, fmt.Errorf("failed to search the query: %w", err)
	}

	return products, nil

}

func (s *ProduuctService) GetTotalProducts() (int, error) {
	return 0, fmt.Errorf("not implemented yet")
}
