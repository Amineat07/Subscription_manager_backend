package sharedhandler

import (
	"strconv"
	"subscription_manager/data"
	"subscription_manager/database"

	"github.com/gofiber/fiber/v2"
)

func GetAllNewsFeed(c *fiber.Ctx) error {

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

func GetNewsFeed(c *fiber.Ctx) error {
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
