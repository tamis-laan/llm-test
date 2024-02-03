package routes

import (
	"context"
	"errors"
	"fmt"
	"llm/internal/clients"
	"llm/internal/models"
	"llm/internal/utils"
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
			`conversation:{title:"%s",description:"%s"}`,
			request.Title,
			request.Description,
		)

		// Create embedding
		embedding := utils.Embed(template)

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

}
