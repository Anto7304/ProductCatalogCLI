package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"Anto7304.com/internal/models"
	"Anto7304.com/internal/repository"
	"Anto7304.com/internal/service"
	_ "github.com/lib/pq"
)

func main() {

	connstr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connstr)

	if err != nil {
		log.Fatal("failed to open database connection: %w", err)
	}

	defer db.Close()

	//test connection

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to postgreSQL!")
	fmt.Printf("conneted to database: %s\n", os.Getenv("DB_NAME"))

	//section 2: initialize layers
	//initialize repository
	repo := repository.NewRepositoryProductRepository(db)

	//initialize service layer

	svc := service.NewProductService(repo)

	//section 3 RUN THE CLI

	runCLI(svc)

}

//cli menuu

func runCLI(svc *service.ProduuctService) {

	//create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	//min loop runs until user select "Exist"
	for {
		//display the menu
		displayMenu()

		//read the user choice
		fmt.Println("Choose option")
		scanner.Scan()

		option := strings.TrimSpace(scanner.Text())

		//routing to the appropriate function
		switch option {
		case "1":
			handleAddProduct(svc, scanner)
		case "2":
			handleListProducts(svc, scanner)
		case "3":
			handleGetProduct(svc, scanner)
		case "4":
			handleUpdateProduuct(svc, scanner)
		case "5":
			handleDeleteProduct(svc, scanner)
		case "6":
			handleSearchProduct(svc, scanner)
		case "7":
			fmt.Println("\n Goodbye! Thankyou for using Product catalog CLI.")
			return
		default:
			fmt.Println("Invalid option. Please choose 1-7")
		}

		fmt.Print("\n Press enter to continue...")
		scanner.Scan()

	}

}

// section 5 display functions
// display menu for the main menu
func displayMenu() {
	fmt.Println()
	fmt.Println("Product Catalog Manager")
	fmt.Println("==========================")
	fmt.Println("1. Add Product")
	fmt.Println("2. List Products")
	fmt.Println("3. Get Product")
	fmt.Println("4. Update Product")
	fmt.Println("5. Delete Product")
	fmt.Println("6. Search Products")
	fmt.Println("7. Exit")
	fmt.Println("==========================")
}

//displayProduct to show a single products details

func displayProduct(product *models.Product) {
	fmt.Println()
	fmt.Println("Product Details")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("ID:          %s\n", product.ID)
	fmt.Printf("Name:        %s\n", product.Name)
	fmt.Printf("Description: %s\n", product.Description)
	fmt.Printf("Price:       KSh %.2f\n", product.Price)
	fmt.Printf("Stock:       %d\n", product.StockQuantity)
	fmt.Printf("Category:    %s\n", product.Category)
	fmt.Printf("SKU:         %s\n", product.SKU)
	fmt.Printf("Weight:      %s\n", product.Weight)
	fmt.Printf("Created:     %s\n", product.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:     %s\n", product.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("-", 60))
}

//displayProducts showing a list of products in a table

func displayProducts(products []models.Product, page, limit int) {
	if len(products) == 0 {
		fmt.Println("\n No products found")
		return
	}

	fmt.Printf("\n Products (page %d, %d per page):\n", page, limit)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-4s %-2s %-15s %-10s %-10s\n", "#", "Name", "SKU", "Price", "Stock")
	fmt.Println(strings.Repeat("-", 60))

	for i, p := range products {
		fmt.Printf("%-4d %-25.25s %-15s KSh %-5.2f %-10d\n",
			i+1+(page-1)*limit,
			p.Name,
			p.SKU,
			p.Price,
			p.StockQuantity,
		)
	}

	fmt.Println(strings.Repeat("-", 60))

}

//input helpers

func readString(scanner *bufio.Scanner, prompt string) string {

	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())

}

func readStringOptional(scanner *bufio.Scanner, prompt string) *string {

	fmt.Println(prompt)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return nil
	}

	return &input

}

func readFloat(scanner *bufio.Scanner, prompt string) (float64, error) {

	fmt.Println(prompt)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" {
		return 0, fmt.Errorf("input cannot be empty")

	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %s", input)
	}

	return value, nil

}

func readFloatOptional(scanner *bufio.Scanner, prompt string) *float64 {

	fmt.Println(prompt)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" {
		return nil
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Printf("Invalid number format: %s\n", input)
		return nil
	}

	return &value

}

func readInt(scanner *bufio.Scanner, prompt string) (int, error) {

	fmt.Print(prompt)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" {
		return 0, fmt.Errorf("inpuut cannot be empty")

	}

	value, err := strconv.Atoi(input)

	if err != nil {
		return 0, fmt.Errorf("invalid number format: %s", input)

	}

	return value, nil

}

func readIntOptional(scanner *bufio.Scanner, prompt string) *int {
	fmt.Print(prompt)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return nil
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("invalid number format: %s\n", input)
	}

	return &value

}

//section 7 :Handler function

func handleAddProduct(svc *service.ProduuctService, scanner *bufio.Scanner) {

	fmt.Println("\n Add new product")
	fmt.Println(strings.Repeat("-", 60))

	fmt.Println("Enter products details (press enter after each):")

	//read inputs

	name := readString(scanner, "Name: ")
	if name == "" {
		fmt.Println("Produuct name is required.")
		return
	}
	description := readString(scanner, "Description: ")

	price, err := readFloat(scanner, "Price (KES): ")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	category := readString(scanner, "Category: ")
	sku := readString(scanner, "SKU: ")
	if sku == "" {
		fmt.Println("SKU is required.")
		return
	}
	weight, err := readFloat(scanner, "Weight: ")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	stock, err := readInt(scanner, "Stock: ")

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	product, err := svc.CreateProduct(name, description, category, sku, price, weight, stock)

	if err != nil {
		fmt.Printf("Failed to create product: %v\n", err)
		return
	}

	fmt.Println("Product created successfully")
	fmt.Printf("Id: %s\n", product.ID)
	fmt.Printf("Name: %s\n", product.Name)
	fmt.Printf("SKU: %s\n", product.SKU)

}

func handleListProducts(svc *service.ProduuctService, scanner *bufio.Scanner) {
	fmt.Println("\nList Products")
	fmt.Println("----------------------------------------")

	// Read pagination
	page, err := readInt(scanner, "Page number (default 1): ")
	if err != nil {
		page = 1
	}

	limit, err := readInt(scanner, "Items per page (default 10, max 50): ")
	if err != nil {
		limit = 10
	}

	// Call the service
	products, err := svc.ListProducts(page, limit)
	if err != nil {
		fmt.Printf("Failed to list products: %v\n", err)
		return
	}

	// Display results
	displayProducts(products, page, limit)
}

func handleGetProduct(svc *service.ProduuctService, scanner *bufio.Scanner) {
	fmt.Println("\nGet Product ")
	fmt.Println(strings.Repeat("-", 60))

	//read id
	id := readString(scanner, "Product ID: ")
	if id == "" {
		fmt.Println("Product ID is required")
		return
	}

	//call the service

	product, err := svc.GetProduct(id)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	displayProduct(product)

}

func handleUpdateProduuct(svc *service.ProduuctService, scanner *bufio.Scanner) {

	fmt.Println("\nUpdate product")
	fmt.Println(strings.Repeat("-", 60))

	//read id

	id := readString(scanner, "Product ID: ")
	if id == "" {
		fmt.Println("Product ID is required")
		return
	}

	//get current product to show current values

	current, err := svc.GetProduct(id)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("\nCurrent Product: \n")
	displayProduct(current)

	fmt.Printf("\nEnter new values (or press Enter to keep current)")
	name := readStringOptional(scanner, fmt.Sprintf("Name [%s]", current.Name))
	description := readStringOptional(scanner, fmt.Sprintf("Description [%s]", current.Description))
	category := readStringOptional(scanner, fmt.Sprintf("Category [%s]", current.Category))
	sku := readStringOptional(scanner, fmt.Sprintf("SKU [%s]", current.SKU))
	price := readFloatOptional(scanner, fmt.Sprintf("Price [%.2f]", current.Price))
	weight := readFloatOptional(scanner, fmt.Sprintf("Weight [%.2f]", current.Weight))
	stock := readIntOptional(scanner, fmt.Sprintf("Stock [%d]", current.StockQuantity))

	//call the service
	product, err := svc.UpdateProduct(id, name, description, category, sku, price, weight, stock)

	if err != nil {
		fmt.Printf("Failed to update product: %v\n", err)
		return
	}

	//success
	fmt.Println("\nProduct updated successfully")
	displayProduct(product)

}

func handleDeleteProduct(svc *service.ProduuctService, scanner *bufio.Scanner) {
	fmt.Println("\n Delete Product")
	fmt.Println(strings.Repeat("-", 60))

	//read id
	id := readString(scanner, "Product ID: ")
	if id == "" {
		fmt.Println("Product ID is required")
		return
	}

	//confirm deletion
	fmt.Println("Are you sure? This action cannot be undone. (Y/N): ")
	scanner.Scan()
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return

	}

	//call the service
	err := svc.DeleteProduct(id)

	if err != nil {
		fmt.Printf("Failed to delete product: %v\n", err)
		return
	}

	//success
	fmt.Println("\nProduct deleted successfully!")

}

func handleSearchProduct(svc *service.ProduuctService, scanner *bufio.Scanner) {
	fmt.Println("Search Products")
	fmt.Println(strings.Repeat("-", 60))

	//read query
	query := readString(scanner, "Search term: ")
	if query == "" {
		fmt.Println("Search term cannot be empty")
		return
	}

	//call the service
	products, err := svc.SearchProducts(query)

	if err != nil {
		fmt.Printf("failed to search the products")
		return
	}

	//duisplay results

	if len(products) == 0 {
		fmt.Printf("\n no products found matching '%s' \n", query)
		return
	}

	fmt.Printf("\n🔍 Search Results for '%s' (%d found):\n", query, len(products))
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-4s %-25s %-15s %-10s\n", "#", "Name", "SKU", "Price")
	fmt.Println(strings.Repeat("-", 60))

	for i, p := range products {
		fmt.Printf("%-4d %-25.25s %-15s KSh %-8.2f\n",
			i+1,
			p.Name,
			p.SKU,
			p.Price,
		)
	}

	fmt.Println(strings.Repeat("-", 60))
}
