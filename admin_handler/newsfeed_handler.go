package adminhandler

import (
	"fmt"
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"
	"subscription_manager/utils"

	"github.com/gofiber/fiber/v2"
)

func AdminCreateNewsFeed(c *fiber.Ctx) error {
	var req data.NewsFeedRequest

	body := c.Body()
	fmt.Println(string(body))
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := utils.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(
			fmt.Sprintf("Validation error: %s", err),
		)
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	sqlstatement := `
    INSERT INTO news_feed (title, content, image_url, is_published, scheduled_at, created_by, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, title, content, COALESCE(image_url, ''), is_published, scheduled_at, published_at, created_by, created_at, updated_by, updated_at
	`

	var inserted data.NewsFeedResponse

	err := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlstatement,
		req.Title,
		req.Content,
		req.ImageUrl,
		req.IsPublished,
		req.ScheduledAt,
		userEmail,
		userEmail,
	).Scan(
		&inserted.ID,
		&inserted.Title,
		&inserted.Content,
		&inserted.ImageUrl,
		&inserted.IsPublished,
		&inserted.ScheduledAt,
		&inserted.PublishedAt,
		&inserted.CreatedBy,
		&inserted.CreatedAt,
		&inserted.UpdatedBy,
		&inserted.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(inserted)
}

func AdminGetAllNewsFeed(c *fiber.Ctx) error {

	sqlStatement := `SELECT id,title,content,image_url,is_published,scheduled_at,published_at,created_by,created_at,updated_by,updated_at
	FROM news_feed WHERE deleted_at IS NULL`

	rows, err := database.InitiateDataBase().Query(c.Context(), sqlStatement)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	defer rows.Close()

	newsfeed := []data.NewsFeedResponse{}

	for rows.Next() {
		var nf data.NewsFeedResponse
		err := rows.Scan(
			&nf.ID,
			&nf.Title,
			&nf.Content,
			&nf.ImageUrl,
			&nf.IsPublished,
			&nf.ScheduledAt,
			&nf.PublishedAt,
			&nf.CreatedBy,
			&nf.CreatedAt,
			&nf.UpdatedBy,
			&nf.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		newsfeed = append(newsfeed, nf)
	}

	return c.Status(fiber.StatusOK).JSON(newsfeed)
}

func AdminGetNewsFeed(c *fiber.Ctx) error {
	id := c.Params("id")

	newsFeedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	sqlStatement := `SELECT id,title,content,image_url,is_published,scheduled_at,published_at,created_by,created_at,updated_by,updated_at
	FROM news_feed WHERE id =$1 AND deleted_at IS NULL`

	var nf data.NewsFeedResponse

	err = database.InitiateDataBase().QueryRow(c.Context(), sqlStatement, newsFeedID).Scan(
		&nf.ID,
		&nf.Title,
		&nf.Content,
		&nf.ImageUrl,
		&nf.IsPublished,
		&nf.ScheduledAt,
		&nf.PublishedAt,
		&nf.CreatedBy,
		&nf.CreatedAt,
		&nf.UpdatedBy,
		&nf.UpdatedAt,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(nf)

}
