package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
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

	// Check if there are any attendance records associated with the employee
	var attendanceCount int
	err = db.QueryRow("SELECT COUNT(*) FROM attendance WHERE employee_id = $1", employeeID).Scan(&attendanceCount)
	if err != nil {
		return err
	}

	if attendanceCount > 0 {
		// If there are attendance records, prompt the user to confirm deletion
		fmt.Println("Warning: This employee has attendance records associated with them. Deleting this employee will also delete these records.")
		fmt.Print("Are you sure you want to delete employee " + employeeID + " (" + firstName + " " + lastName + ")? (y/n): ")
		var confirm string
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) != "y" {
			fmt.Println("Deletion cancelled")
			return nil
		}

		// Delete the attendance records
		_, err = db.Exec("DELETE FROM attendance WHERE employee_id = $1", employeeID)
		if err != nil {
			return err
		}
	}

	// Prompt the user to confirm deletion
	fmt.Print("Are you sure you want to delete employee " + employeeID + " (" + firstName + " " + lastName + ")? (y/n): ")
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	// Proceed with deleting the employee
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

	//watermark ni ady
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all employee information from the database
	rows, err := db.Query("SELECT id, role, last_name, first_name FROM employees")
	if err != nil {
		http.Error(w, "Failed to retrieve employee information", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Stream the employee information to the response body
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Role, &employee.LastName, &employee.FirstName)
		if err != nil {
			http.Error(w, "Failed to scan employee row", http.StatusInternalServerError)
			return
		}

		// Marshal the employee information to JSON
		employeeInfo, err := json.Marshal(employee)
		if err != nil {
			http.Error(w, "Failed to marshal employee information", http.StatusInternalServerError)
			return
		}

		// Write the employee information to the response body
		_, err = w.Write(append(employeeInfo, '\n'))
		if err != nil {
			http.Error(w, "Failed to write employee information", http.StatusInternalServerError)
			return
		}

		// Flush the response body to the client
		if fw, ok := w.(http.Flusher); ok {
			fw.Flush()
		}
	}
}

type Employee struct {
	ID        string `json:"ID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Password  string `json:"Password"`
	Role      string `json:"Role"`
}
