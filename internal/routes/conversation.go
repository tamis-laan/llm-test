package routes

import (
	"context"
	"errors"
	"fmt"
	"llm/internal/clients"
	"llm/internal/models"
	"llm/internal/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetRouteConversation(router fiber.Router) {

	// Request validator
	validate := clients.GetValidor()

	// Get the database connection pool
	pool := clients.GetPool()

	////////////////////////// /conversation.get
	router.Get("/conversation", func(c *fiber.Ctx) error {
		// Get id parameter
		id, err := strconv.ParseUint(c.Query("id"), 0, 64)

		// Error check
		if err != nil {
			return errors.New("id must be of type unsigned int")
		}

		// Get conversation from database
		rows, err := pool.Query(context.Background(), `
			SELECT * FROM conversations
			WHERE id=$1`,
			id,
		)

		// Release connection
		defer rows.Close()

		// Error check
		if err != nil {
			return err
		}

		// Collect results
		responses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Conversation])

		// Error check
		if err != nil {
			return err
		}

		return c.JSON(responses[0])
	})

	///////////////////////// /conversation.post
	router.Post("/conversation", func(c *fiber.Ctx) error {

		// Request structure
		type Request struct {
			Title       string `json:"title"       validate:"required,max=32"`
			Description string `json:"description" validate:"required,max=512"`
		}

		// Create the request
		request := Request{}

		// Parse body into request
		err := c.BodyParser(&request)

		// Error check
		if err != nil {
			return err
		}

		// Run validation
		err = validate.Struct(request)

		// Error check
		if err != nil {
			return err
		}

		// Create embedding template
		template := fmt.Sprintf(
			// `conversation:{title:"%s",description:"%s"}`,
			`"%s" "%s"`,
			request.Title,
			request.Description,
		)

		// Create embedding
		embedding, err := utils.Embed(template)

		// Error check
		if err != nil {
			return err
		}

		// Add conversation to database
		rows, err := pool.Query(context.Background(), `
			INSERT INTO conversations 
				(title, description, embedding) 
			VALUES 
				($1, $2, $3)
			RETURNING *`,
			request.Title,
			request.Description,
			embedding,
		)

		// Release connection
		defer rows.Close()

		// Error check
		if err != nil {
			return err
		}

		// Collect results
		responses, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Conversation])

		// Error check
		if err != nil {
			return err
		}

		// Return the result
		return c.JSON(responses[0])
	})

	////////////////////////// /query.get
	router.Get("/query", func(c *fiber.Ctx) error {

		// Get id parameter
		query := c.Query("q")

		log.Println(query)

		// Embed the query
		queryEmbedding, err := utils.Embed(query)

		// Error check
		if err != nil {
			return err
		}

		log.Println(queryEmbedding.Slice()[0])

		// Similarity search
		rows, err := pool.Query(context.Background(), `
			SELECT (embedding <-> $1) AS score, * FROM conversations ORDER BY score`,
			queryEmbedding,
		)

		// Release connection
		defer rows.Close()

		// Error check
		if err != nil {
			return err
		}

		type Output struct {
			models.Conversation
			Score float64 `db:"score"`
		}

		// Collect results
		responses, err := pgx.CollectRows(rows, pgx.RowToStructByName[Output])

		// Error check
		if err != nil {
			return err
		}

		// Return found
		return c.JSON(responses)
	})

}
