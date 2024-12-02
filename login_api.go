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

	if storedPassword == password {
		_, err := db.Exec(`
            INSERT INTO attendance (employee_id, check_in_time)
            VALUES ($1, $2);
        `, employeeID, time.Now())
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to check in employee", http.StatusInternalServerError)
			return
		}

		log.Println("Check-in record inserted successfully")
		// Create a response string containing the role, first name, and last name
		response := fmt.Sprintf("%s %s %s", role, firstName, lastName)
		w.Write([]byte(response))

	} else {

		log.Printf("Successful login attempt from %s (employee ID: %s, role: %s, first name: %s, last name: %s)\n", ipAddress, employeeID, role, firstName, lastName)

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
		return
	}

}
