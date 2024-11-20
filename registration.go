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
