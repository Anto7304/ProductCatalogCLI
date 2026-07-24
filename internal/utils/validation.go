package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateProductName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("product name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("product name must be atleast two character")
	}

	if len(name) > 255 {
		return fmt.Errorf("product name cannot exceed 255 character")
	}

	return nil
}

func ValidatePrice(price float64) error {
	if price < 0 {
		return fmt.Errorf("price cannot be negative")
	}

	//check unreasonanble price

	if price > 1000000 {
		return fmt.Errorf("price seems unreasonably high. please double check the price")
	}
	return nil
}

func ValidateStock(quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("stock quantity cannot be negative")
	}

	if quantity > 1000000 {
		return fmt.Errorf("quantity seems unreasonably high. Please double check")
	}
	return nil
}

func ValidateSKU(sku string) error {
	sku = strings.TrimSpace(sku)

	if sku == "" {
		return fmt.Errorf("sku cannot be empty")
	}

	if len(sku) > 100 {
		return fmt.Errorf("sku cannot exceed 100 characters")
	}

	//check format: alphanumeric,hyphens, underscore only

	matched, err := regexp.MatchString(`^[a-zA-Z0-9\-_]+$`, sku)

	if err != nil {
		return fmt.Errorf("error validating the sku format: %w", err)
	}

	if !matched {
		return fmt.Errorf("SKU can only containe letters,numbers, hyphens and underscore")
	}

	return nil

}

func ValidateCategory(category string) error {
	category = strings.TrimSpace(category)

	if category == "" {
		return fmt.Errorf("category cannot be empty")
	}

	if len(category) > 100 {
		return fmt.Errorf("category cannot exceed 100 characters")
	}

	return nil
}

func SanitizeString(input string) string {

	if input == "" {
		return ""
	}

	input = strings.TrimSpace(input)

	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(input, "")

}
