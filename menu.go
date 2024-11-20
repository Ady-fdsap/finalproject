package main

import (
	"fmt"
	"os"
)

func menu() {
	for {
		fmt.Println("=================================")
		fmt.Println("Select an option:")
		fmt.Println("1. Register Employee")
		fmt.Println("2. Exit")
		var choice int
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			registerEmployee(db)
		case 2:
			fmt.Println("Exiting program")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Choose again.")
		}
	}
}
