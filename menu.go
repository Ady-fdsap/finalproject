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
		fmt.Println("3. Exit")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			registerEmployee(db)
		case 2:
			err := deleteEmployee(db)
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}
			break

		case 3:
			fmt.Println("Exiting program")
			os.Exit(0)

		default:
			fmt.Println("Invalid choice. Choose again.")
		}
	}
}
