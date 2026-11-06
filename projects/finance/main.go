package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"finance/db"
	"finance/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "finance.db"
	}

	database, err := db.OpenDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	ctx := context.Background()
	companyID, err := db.SeedInitialCompanyAndAccounts(ctx, database)
	if err != nil {
		panic(err)
	}

	router := routes.SetupRouter(database, companyID)

	fmt.Printf("Finance Management Server running at http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
