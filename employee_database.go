package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func registerEmployee(db *sql.DB) error {
	var id, firstName, lastName, password, role string

	fmt.Print("Enter Role [Admin] [Employee] [Intern] (or 'abort' to cancel): ")
	fmt.Scanln(&role)
	role = strings.ToLower(role)
	if role == "abort" {
		fmt.Println("Registration cancelled")
		return nil
	}

	fmt.Print("Enter Employee ID (or 'abort' to cancel): ")
	reader := bufio.NewReader(os.Stdin)
	id, _ = reader.ReadString('\n')
	id = strings.TrimSpace(id)
	if strings.ToLower(id) == "abort" {
		fmt.Println("Registration cancelled")
		return nil
	}

	fmt.Print("Enter First Name (or 'abort' to cancel): ")
	fmt.Scanln(&firstName)
	if strings.ToLower(firstName) == "abort" {
		fmt.Println("Registration cancelled")
		return nil
	}

	fmt.Print("Enter Last Name (or 'abort' to cancel): ")
	fmt.Scanln(&lastName)
	if strings.ToLower(lastName) == "abort" {
		fmt.Println("Registration cancelled")
		return nil
	}

	fmt.Print("Enter Password (or 'abort' to cancel): ")
	fmt.Scanln(&password)
	if strings.ToLower(password) == "abort" {
		fmt.Println("Registration cancelled")
		return nil
	}

	// Encrypt the password
	secretKey := GetSecretKey()
	encryptedPassword, err := EncryptPassword(password, secretKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %v", err)
	}

	// Check if password meets requirements
	if len(password) < 8 {
		fmt.Println("Password must be at least 8 characters long")
		return registerEmployee(db)
	}
	if !hasCapitalLetter(password) {
		fmt.Println("Password must contain at least one capital letter")
		return registerEmployee(db)
	}
	if !hasNumber(password) {
		fmt.Println("Password must contain at least one number")
		return registerEmployee(db)
	}

	_, err = db.Exec(`
    INSERT INTO employees (id, first_name, last_name, date_added, password, role)
    VALUES ($1, $2, $3, CURRENT_DATE, $4, $5);
`, id, firstName, lastName, encryptedPassword, role)

	if err != nil {
		return fmt.Errorf("failed to register employee: %v", err)
	}
	return nil
}

func deleteEmployee(db *sql.DB) error {
	fmt.Print("Enter Employee ID to delete: ")
	reader := bufio.NewReader(os.Stdin)
	employeeID, _ := reader.ReadString('\n')
	employeeID = strings.TrimSpace(employeeID) // remove any leading or trailing spaces

	var firstName, lastName string
	err := db.QueryRow(`
		SELECT first_name, last_name
		FROM employees
		WHERE id = $1;
	`, employeeID).Scan(&firstName, &lastName)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Employee not found")
			return nil
		}
		return fmt.Errorf("failed to retrieve employee details: %v", err)
	}

	fmt.Print("Are you sure you want to delete employee " + employeeID + " (" + firstName + " " + lastName + ")? (y/n): ")
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	result, err := db.Exec(`
		DELETE FROM employees
		WHERE id = $1;
	`, employeeID)

	if err != nil {
		return fmt.Errorf("failed to delete employee: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	fmt.Println("Employee deleted successfully")
	return nil
}

func displayEmployees(db *sql.DB) error {
	var offset int
	for {
		rows, err := db.Query(`
			SELECT id, first_name, last_name
			FROM employees
			ORDER BY last_name
			LIMIT 10 OFFSET $1;
		`, offset)
		if err != nil {
			return fmt.Errorf("failed to retrieve employees: %v", err)
		}
		defer rows.Close()

		fmt.Println("Employees (sorted by last name):")
		var count int
		for rows.Next() {
			var id, firstName, lastName string
			err := rows.Scan(&id, &firstName, &lastName)
			if err != nil {
				return fmt.Errorf("failed to scan employee row: %v", err)
			}
			fmt.Printf("%s %s (%s)\n", firstName, lastName, id)
			count++
		}

		if count == 0 {
			fmt.Println("No more employees to display.")
			break
		}

		fmt.Println("Options:")
		fmt.Println("1. Next Page")
		fmt.Println("2. Exit")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			offset += 10
		case 2:
			return nil
		default:
			fmt.Println("Invalid choice. Choose again.")
		}
	}

	return nil
}

func (api *API) handleGetEmployeeInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Parse the request body
	var requestBody struct {
		Keyword string `json:"keyword"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if the specific word is present in the request body
	if !strings.Contains(requestBody.Keyword, "strawberry shortcake") {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve all employee information from the database
	rows, err := db.Query("SELECT id, role, last_name, first_name FROM employees")
	if err != nil {
		http.Error(w, "Failed to retrieve employee information", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Role, &employee.LastName, &employee.FirstName)
		if err != nil {
			http.Error(w, "Failed to scan employee row", http.StatusInternalServerError)
			return
		}
		employees = append(employees, employee)
	}

	// Marshal the employee information to JSON
	employeeInfo, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, "Failed to marshal employee information", http.StatusInternalServerError)
		return
	}

	// Send the employee information as a response body
	w.Header().Set("Content-Type", "application/json")
	w.Write(employeeInfo)
}

type Employee struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}
