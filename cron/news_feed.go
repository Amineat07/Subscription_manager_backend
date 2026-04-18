package cron

import (
	"context"
	"log"
	"subscription_manager/data"
	"subscription_manager/database"
	ssehandler "subscription_manager/see_handler"
)

func PublishScheduledNews() {
	log.Println("cronjob running...")

	sqlStatement := `
        UPDATE news_feed
        SET is_published = true,
            published_at = NOW()
        WHERE scheduled_at <= NOW()
        AND is_published = false
        AND deleted_at IS NULL
        RETURNING id, title, content, COALESCE(image_url,''), is_published, scheduled_at, published_at, created_by, created_at, updated_by, updated_at
    `

	rows, err := database.InitiateDataBase().Query(context.Background(), sqlStatement)
	if err != nil {
		log.Println("cronjob error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item data.NewsFeedResponse
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Content,
			&item.ImageUrl,
			&item.IsPublished,
			&item.ScheduledAt,
			&item.PublishedAt,
			&item.CreatedBy,
			&item.CreatedAt,
			&item.UpdatedBy,
			&item.UpdatedAt,
		)
		if err != nil {
			log.Println("cronjob scan error:", err)
			continue
		}

		ssehandler.BroadcastNewsFeed(item)
		log.Printf("cronjob: broadcasted newsfeed id=%d\n", item.ID)
	}

	log.Println("cronjob ran: scheduled news feed published")
}
