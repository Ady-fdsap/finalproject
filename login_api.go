package main

import (
	"log"
	"net/http"

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
}
