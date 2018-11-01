package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/monokrome/progress"

	// SQL dialects
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var err error

	options, _, err := progress.NewOptions("progress")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
	}

	database, err := gorm.Open(options.Storage.Backend, options.Storage.Options)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database connection: %v\n", err)
		os.Exit(1)
	}

	defer database.Close()

	progress.EnsureSchema(database)

	transaction := database.Begin()

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
			fmt.Fprintf(os.Stderr, "Command failed: %v\n", r)
			os.Exit(1)
		}
	}()

	CommandLine(options, transaction).Execute()
	transaction.Commit()
}
