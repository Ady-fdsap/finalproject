package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func menu() {
	for {
		fmt.Println("=================================")
		fmt.Println("Select an option:")
		fmt.Println("1. Add employee")
		fmt.Println("2. Restart application")
		fmt.Println("3. Exit")

		reader := bufio.NewReader(os.Stdin)
		option, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		optionInt, err := strconv.Atoi(option[:len(option)-1])
		if err != nil {
			fmt.Println("Invalid option:", option)
			continue
		}

		switch optionInt {
		case 1:
			fmt.Print("Enter employee ID (or 'abort' to cancel): ")
			var id string
			fmt.Scanln(&id)

			if id == "abort" {
				fmt.Println("Employee addition aborted.")
				break
			}

			fmt.Print("Enter employee name (or 'abort' to cancel): ")
			var name string
			fmt.Scanln(&name)

			if name == "abort" {
				fmt.Println("Employee addition aborted.")
				break
			}

			fmt.Print("Enter employee password (or 'abort' to cancel): ")
			var password string
			fmt.Scanln(&password)

			if password == "abort" {
				fmt.Println("Employee addition aborted.")
				break
			}

			err := registerEmployee(id, name, password)
			if err != nil {
				fmt.Println("Error registering employee:", err)
			} else {
				fmt.Println("Employee registered successfully")
			}
		case 2:
			fmt.Println("Restarting application...")
			// Restart the application
			os.Exit(0)

		case 3:
			fmt.Println("Exiting application...")
			// Exit the application
			os.Exit(0)

		default:
			fmt.Println("Invalid option:", option)
		}
	}
}
