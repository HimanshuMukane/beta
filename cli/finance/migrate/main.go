package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"core/utils"
	"finance/db"
)

func main() {
	runMigrate := flag.Bool("migrate", true, "Execute pending schema migrations")
	runSeed := flag.Bool("seed", false, "Seed default initial company and chart of accounts")
	flag.Parse()

	if err := utils.LoadConfig(); err != nil {
		log.Printf("Note: utils.LoadConfig: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "especa_root_password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "finance_db"
	}

	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", dbPassword, dbHost, dbPort, dbName)
	database, err := db.OpenDB(dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()

	if *runMigrate {
		log.Printf("Executing Finance migrations on %s...", dbName)
		if err := db.RunMigrations(ctx, database, "db/migrations"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations executed successfully.")
	}

	if *runSeed {
		log.Println("Seeding initial company and chart of accounts...")
		companyID, err := db.SeedInitialCompanyAndAccounts(ctx, database)
		if err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		log.Printf("Seeding complete. Primary Company ID: %s", companyID)
	}
}
