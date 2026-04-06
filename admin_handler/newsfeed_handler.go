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



func AdminUpdateNewsFeed(c *fiber.Ctx) error {

	var updateNewsFeed data.UpdateNewsFeedRequest
	if err := c.BodyParser(&updateNewsFeed); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	id := c.Params("id")

	newsFeedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	sqlStatement := `UPDATE news_feed
	SET title = COALESCE(NULLIF($1, ''), title),
		content = COALESCE(NULLIF($2, ''), content),
		image_url = COALESCE(NULLIF($3, ''), image_url),
		is_published = COALESCE($4, is_published),
		scheduled_at = COALESCE($5, scheduled_at),
		updated_at = NOW(),
		updated_by = $6
		WHERE id = $7 AND deleted_at IS NULL
		RETURNING id,title,content,image_url,is_published,scheduled_at,created_at,created_by,updated_at,updated_by
	`
	row := database.InitiateDataBase().QueryRow(
		c.Context(),
		sqlStatement,
		updateNewsFeed.Title,
		updateNewsFeed.Content,
		updateNewsFeed.ImageUrl,
		updateNewsFeed.IsPublished,
		updateNewsFeed.ScheduledAt,
		userEmail,
		newsFeedID,
	)

	var updated data.NewsFeedResponse
	err = row.Scan(
		&updated.ID,
		&updated.Title,
		&updated.Content,
		&updated.ImageUrl,
		&updated.IsPublished,
		&updated.ScheduledAt,
		&updated.CreatedAt,
		&updated.CreatedBy,
		&updated.UpdatedAt,
		&updated.UpdatedBy,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ticket updated successfully",
		"ticket":  updated,
	})
}

func AdminDeleteNewsFeed(c *fiber.Ctx) error {
	id := c.Params("id")

	newsFeedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user email not found",
		})
	}

	sqlStatement := `
        UPDATE news_feed
        SET deleted_at = NOW(),
            deleted_by = $1
        WHERE id = $2
          AND deleted_at IS NULL
    `

	result, err := database.InitiateDataBase().Exec(
		c.Context(),
		sqlStatement,
		userEmail,
		newsFeedID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "newsfeed not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "newsfeed deleted successfully",
	})
}
