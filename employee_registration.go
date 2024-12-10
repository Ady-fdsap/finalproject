//for the employee registration api endpoint

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
)

func (api *API) handleRegisterEmployee(w http.ResponseWriter, r *http.Request) {
	// Get the request body
	var employee Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the employee data
	if employee.ID == "" || employee.FirstName == "" || employee.LastName == "" || employee.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate the password
	if !hasNumber(employee.Password) {
		http.Error(w, "Password must contain at least one number", http.StatusBadRequest)
		return
	}

	if !hasCapitalLetter(employee.Password) {
		http.Error(w, "Password must contain at least one capital letter", http.StatusBadRequest)
		return
	}

	// Check if password meets length requirements
	if len(employee.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Register the employee
	err = api.registerEmployeeAPI(employee)
	if err != nil {
		http.Error(w, "Failed to register employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (api *API) registerEmployeeAPI(employee Employee) error {
	// Check if the ID already exists
	var count int
	err := api.database.QueryRow(`
        SELECT COUNT(*) FROM employees WHERE id = $1
    `, employee.ID).Scan(&count)

	if err != nil {
		return fmt.Errorf("failed to check existing ID: %v", err.Error())
	}

	if count > 0 {
		return fmt.Errorf("employee with ID %v already exists", employee.ID)
	}

	// Encrypt the password
	log.Println("Executing SQL query to insert employee")
	secretKey := GetSecretKey()
	encryptedPassword, err := EncryptPassword(employee.Password, secretKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %v", err.Error())
	}

	// Insert the employee into the database
	_, err = api.database.Exec(`
        INSERT INTO employees (id, first_name, last_name, date_added, password, role)
        VALUES ($1, $2, $3, CURRENT_DATE, $4, $5);
    `, employee.ID, employee.FirstName, employee.LastName, encryptedPassword, employee.Role)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			return pqErr
		}
		return err
	}

	return nil
}
