package cron

import (
	"context"
	"log"
	"subscription_manager/database"
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
	`

	_, err := database.InitiateDataBase().Exec(context.Background(), sqlStatement)

	if err != nil {
		log.Println("cronjob error:", err)
	} else {
		log.Println("cronjob ran: scheduled news feed published")
	}
}
