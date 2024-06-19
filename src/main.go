package main

import (
	"forum/logic" // Importing the logic package which contains the core functionalities of the forum

	_ "github.com/mattn/go-sqlite3" // Importing the SQLite3 driver for database operations
)

func main() {
	// Clear the console
	ClearConsole()

	// Initialize the database and load the data
	logic.InitData()

	// Launch the web application
	logic.LaunchApp()
}

// clearConsole clears the console screen
func ClearConsole() {
	print("\033[H\033[2J") // Clear the console screen
}
