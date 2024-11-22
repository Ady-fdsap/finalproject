package main

import (
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
	var storedPassword string
	err := db.QueryRow(`
        SELECT password
        FROM employees
        WHERE id = $1;
    `, employeeID).Scan(&storedPassword)

	if err != nil {
		http.Error(w, "false", http.StatusUnauthorized)
		log.Println("Failed login attempt from", ipAddress)
		return
	}

	// Compare the provided password with the stored password
	if password == storedPassword {
		log.Println("Successful login from", ipAddress)
		w.Write([]byte("true"))
	} else {
		log.Println("Failed login attempt from", ipAddress)
		http.Error(w, "false", http.StatusUnauthorized)
	}

	logEntry := RequestLog{
		Timestamp: time.Now(),
		Method:    r.Method,
		Latitude:  0, // Not applicable for login attempts
		Longitude: 0, // Not applicable for login attempts
		IPAddress: ipAddress,
	}

	_, err = db.Exec(`
    INSERT INTO login_attempts (timestamp, ip_address, employee_id, success)
    VALUES ($1, $2, $3, $4);
`, logEntry.Timestamp, logEntry.IPAddress, employeeID, password == storedPassword)

	if err != nil {
		log.Fatal(err)
	}

}
