package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"llm/internal/clients"
	"llm/internal/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Parse command line arguments
func main() {
	// Get the database connection pool
	pool := clients.GetPool()

	// Make sure to close the connection pool
	defer pool.Close()

	// Load configuration
	config := utils.GetConfig()

	// Migrate subcommand
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	up := migrateCmd.Bool("up", false, "Migrate database all the way up")
	down := migrateCmd.Bool("down", false, "Migrate database all the way down")
	drop := migrateCmd.Bool("drop", false, "Drop database")
	steps := migrateCmd.Int("steps", 0, "Step up or down the migration with n steps")
	force := migrateCmd.Int("force", -1, "Force database to version")

	// Parse arguments for the appropriate subcommand
	if len(os.Args) < 2 {
		fmt.Println("Usage: app <command> [<args>]")
		fmt.Println("Available commands: migrate")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate":
		// Parse flags
		migrateCmd.Parse(os.Args[2:])

		// Construct connection string
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
			config.Postgres.Username,
			config.Postgres.Password,
			config.Postgres.Host,
			config.Postgres.Port,
			config.Postgres.Database,
			config.Postgres.Options,
		)

		// Setup migrations
		migration, err := migrate.New(config.Migrations, dsn)

		// Make sure we can connect
		if err != nil {
			log.Fatal("Cannot connect to database")
		}

		// Get the migration version
		version, dirty, err := migration.Version()

		// Exit on dirty database
		if dirty {
			log.Println("NOTE DATABSE IS DIRTY")
		}

		// Log migration version
		log.Printf("Current DB version %d", version)

		// Drop database
		if *drop {
			Drop(migration)
		}

		// Force migrate database to version
		if *force > -1 {
			Force(migration, *force)
		}

		// Migrate database
		if *steps != 0 {
			Steps(migration, *steps)
		}

		// Migrate to latest version
		if *down {
			Down(migration)
		}

		// Migrate to latest version
		if *up {
			Up(migration)
		}

		// Get the migration version
		version, dirty, err = migration.Version()

		// Exit on dirty database
		if dirty {
			log.Println("NOTE DATABSE IS DIRTY")
		}

		// Log migration version
		log.Printf("New DB version %d", version)

		// Exit ctl
		os.Exit(0)
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}

}

// Log database version
func Version(migration *migrate.Migrate) {

	// Get the migration version
	version, dirty, _ := migration.Version()

	// Exit on dirty database
	if dirty {
		log.Println("NOTE DATABSE IS DIRTY")
	}

	// Log migration version
	log.Printf("New database version %d", version)
}

// Drop database
func Drop(migration *migrate.Migrate) {

	// Log to console
	log.Println("DROPPING DATABASE")

	// Start migration
	err := migration.Drop()

	// Error Check
	if err != nil {
		panic(err)
	}

	// Log to console
	log.Println("Database has been dropped")
}

// Migrate database to latest version
func Up(migration *migrate.Migrate) {

	// Log to console
	log.Println("Migrating database to latest version")

	// Start migration
	err := migration.Up()

	// Error Check
	if err != nil {
		log.Fatal(err)
	}

	// Log to console
	log.Println("Database migration complete")
}

// Migrate database all the way down
func Down(migration *migrate.Migrate) {

	// Log to console
	log.Println("Downgrade database to lowest version")

	// Start migration
	err := migration.Down()

	// Error Check
	if err != nil {
		log.Fatal(err)
	}

	// Log to console
	log.Println("Database migration complete")

}

// Migrate database n steps up or down
func Steps(migration *migrate.Migrate, steps int) {

	// Log to console
	log.Printf("Migrate database with %d steps", steps)

	// Start migration
	err := migration.Steps(steps)

	// Error Check
	if err != nil {
		log.Fatal(err)
	}

	// Log to console
	log.Println("Database migration complete")

}

// Migrate database to the nth version
func Force(migration *migrate.Migrate, version int) {

	// Log to console
	log.Printf("Force migrate database to version %d", version)

	// Start migration
	err := migration.Force(version)

	// Error Check
	if err != nil {
		log.Fatal(err)
	}

	// Log to console
	log.Printf("Successfully forced to version %d", version)

}
