//for the employee registration api endpoint

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	if employee.ID == "" || employee.FirstName == "" || employee.LastName == "" || employee.Password == "" || employee.Role == "" {
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

	if len(employee.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Encrypt the password
	secretKey := GetSecretKey()
	encryptedPassword, err := EncryptPassword(employee.Password, secretKey)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusBadRequest)
		return
	}

	// Insert the employee into the database
	_, err = api.database.Exec(`
    INSERT INTO employees (id, first_name, last_name, date_added, password, role)
    VALUES ($1, $2, $3, CURRENT_DATE, $4, $5);
`, employee.ID, employee.FirstName, employee.LastName, encryptedPassword, employee.Role)

	if err != nil {
		log.Println(err) // Print the error
		http.Error(w, "Failed to insert employee", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Employee registered successfully")
}
