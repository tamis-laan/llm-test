package clients

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"llm/internal/utils"

	"github.com/go-playground/validator/v10"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Singleton migration
var migration *migrate.Migrate

// Return migration
func GetMigration() *migrate.Migrate {
	return migration
}

func init() {

	// Load configuration
	config := utils.GetConfig()

	// Define error
	var err error

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
	migration, err = migrate.New(config.Migrations, dsn)

	// Handle error
	if err != nil {
		panic(err)
	}

}

// Singleton connection pool
var pool *pgxpool.Pool

// Return the connection pool
func GetPool() *pgxpool.Pool {
	return pool
}

// Run when imported
func init() {

	// Log to console
	log.Println("Connecting to PostgreSQL")

	// Get the global c5672onfiguration object
	config := utils.GetConfig()

	// Construct connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
		config.Postgres.Options,
	)

	// Predefined error object
	var err error

	// Create pool NOTE: This can be configured
	pool, err = pgxpool.New(context.Background(), dsn)

	// Error check
	if err != nil {
		log.Fatal(err)
	}

	// Ping server
	err = pool.Ping(context.Background())

	// Error check
	if err != nil {
		log.Fatal(err)
	}

}

// Singleton validator
var validate *validator.Validate

// Return the validator
func GetValidor() *validator.Validate {
	return validate
}

func init() {
	// Create validator
	validate = validator.New()
}

// Http client
var Http *http.Client

// Initialize http request client
func init() {
	// Create the client
	Http = &http.Client{}
}
