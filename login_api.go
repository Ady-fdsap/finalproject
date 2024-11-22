package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func (api *API) handleEmployeeLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("[LOGIN] Received login request with ID:", r.URL.Query().Get("id"))
	log.Println("[LOGIN]: Received login request with password:", r.URL.Query().Get("password"))

	employeeID := r.URL.Query().Get("id")
	password := r.URL.Query().Get("password")

	// Query the database to check if the employee ID and password match
	var storedPassword string
	err := db.QueryRow(`
        SELECT password
        FROM employees
        WHERE id = $1;
    `, employeeID).Scan(&storedPassword)

	log.Println("[DEBUG] Stored password from database:", storedPassword)

	if err != nil {
		log.Println("[LOGIN] Error querying database:", err)
		http.Error(w, "[LOGIN] Failed to query database", http.StatusInternalServerError)
		return
	}

	// Compare the provided password with the stored password
	if password == storedPassword {
		log.Println("[LOGIN] Login successful!")
		w.Write([]byte("true"))
	} else {
		log.Println("[LOGIN] Invalid password")
		w.Write([]byte("false"))
	}
}
