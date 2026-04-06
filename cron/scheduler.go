package cron

import "github.com/robfig/cron/v3"

func StartCronJobs() {
	c := cron.New()
	c.AddFunc("* * * * *", PublishScheduledNews)
	c.Start()
}
