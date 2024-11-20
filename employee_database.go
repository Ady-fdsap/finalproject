package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func registerEmployee(db *sql.DB) error {
	var id, firstName, lastName, password string

	fmt.Print("Enter Employee ID: ")
	reader := bufio.NewReader(os.Stdin)
	id, _ = reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Enter First Name: ")
	fmt.Scanln(&firstName)

	fmt.Print("Enter Last Name: ")
	fmt.Scanln(&lastName)

	fmt.Print("Enter Password: ")
	fmt.Scanln(&password)

	_, err := db.Exec(`
		INSERT INTO employees (id, first_name, last_name, date_added, password)
		VALUES ($1, $2, $3, CURRENT_DATE, $4);
	`, id, firstName, lastName, password)

	if err != nil {
		return fmt.Errorf("failed to register employee: %v", err)
	}

	fmt.Println("Employee registered successfully")
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
