package routes

import (
	"context"
	"llm/internal/clients"
	"llm/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
)

// Embed string into vector
func embedActor(input string) pgvector.Vector {
	// Create a 1024-zero vector
	zeroVector := make([]float32, 1024)

	// Convert the zero vector to a pgvector.Vector
	return pgvector.NewVector(zeroVector)
}

func SetRouteActor(router fiber.Router) {

	// Request validator
	validate := clients.GetValidor()

	// Get the database connection pool
	pool := clients.GetPool()

	////////////////////////// /conversation.get
	router.Get("/actor", func(c *fiber.Ctx) error {
		return c.JSON("Finding actor by id")
	})

	///////////////////////// /conversation.post
	router.Post("/actor", func(c *fiber.Ctx) error {

		// Request structure
		type Request struct {
			Name      string `validate:"required,max=64"`
			Embedding pgvector.Vector
		}

		// Create the request
		request := Request{}

		// Parse body into request
		c.BodyParser(&request)

		// Run validation
		err := validate.Struct(request)

		// Error check
		if err != nil {
			return err
		}

		// Log the request
		log.Println(request)

		// Add conversation to database
		rows, err := pool.Query(context.Background(), `
			INSERT INTO actors 
				(name, embedding) 
			VALUES 
				($1, $2)
			RETURNING *`,
			request.Name,
			embedActor(request.Name),
		)

		// Release connection
		defer rows.Close()

		// Error check
		if err != nil {
			return err
		}

		// Collect results
		responses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Actor])

		// Error check
		if err != nil {
			return err
		}

		// log.Println(responses[0].Embedding)

		// Return the result
		return c.JSON(responses[0])
	})

}
