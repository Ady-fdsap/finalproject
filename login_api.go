package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func (api *API) handleEmployeeLogin(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	employeeID := r.URL.Query().Get("id")
	password := r.URL.Query().Get("password")

	if employeeID == "" || password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Query the database to check if the employee ID and password match
	stmt, err := db.Prepare("SELECT password, role, first_name, last_name FROM employees WHERE id = $1")
	if err != nil {
		http.Error(w, "failed to prepare query", http.StatusInternalServerError)
		return
	}

	var storedPassword, role, firstName, lastName string
	err = stmt.QueryRow(employeeID).Scan(&storedPassword, &role, &firstName, &lastName)
	if err != nil {
		http.Error(w, "failed to retrieve password", http.StatusInternalServerError)
		return
	}

	decryptedPassword, err := DecryptPassword(storedPassword, GetSecretKey())
	if err != nil {
		http.Error(w, "Failed to decrypt password", http.StatusInternalServerError)
		return
	}

	if decryptedPassword != password {
		// Log the failed login attempt
		_, err = db.Exec(`
			INSERT INTO login_attempts (timestamp, ip_address, employee_id, success)
			VALUES ($1, $2, $3, $4);
		`, time.Now(), ipAddress, employeeID, false)
		if err != nil {
			log.Fatal(err)
		}
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	log.Println("Check-in record inserted successfully")
	// Log the successful login attempt
	_, err = db.Exec(`
        INSERT INTO login_attempts (timestamp, ip_address, employee_id, success)
        VALUES ($1, $2, $3, $4);
    `, time.Now(), ipAddress, employeeID, true)
	if err != nil {
		log.Fatal(err)
	}

	response := fmt.Sprintf("%s %s %s", role, firstName, lastName)
	w.Write([]byte(response))
}
