package main

import (
	"fmt"
	"log"
	"os"
)

func menu() {
	for {
		fmt.Println("=================================")
		fmt.Println("Select an option:")
		fmt.Println("1. Register Employee")
		fmt.Println("2. Delete Employee")
		fmt.Println("3. Display Employees")
		fmt.Println("4. Update Geofences")
		fmt.Println("5. Exit")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			err := registerEmployee(db)
			if err != nil {
				fmt.Println(err)
			}

		case 2:
			err := deleteEmployee(db)
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}

		case 3:
			err := displayEmployees(db)
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}

		case 4:
			updateGeofences(db)
			fmt.Println("Geofences updated successfully")

		case 5:
			fmt.Println("Exiting program")
			os.Exit(0)

		default:
			fmt.Println("Invalid choice. Choose again.")
		}
	}
}
