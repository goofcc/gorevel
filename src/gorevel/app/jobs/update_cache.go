package jobs

import (
	"time"

	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"

	c "gorevel/app/controllers"
	"gorevel/app/models"
	"gorevel/app/routes"
)

type UpdateCache struct{}

func (uc UpdateCache) Run() {
	revel.WARN.Print("update cache...")
	for page := 1; page <= models.CachePageSize; page++ {
		cache.Delete("topics" + routes.Topic.Index(page))
		c.GetTopics(page, "", "created", routes.Topic.Index(page))
		time.Sleep(time.Second)

		cache.Delete("topics" + routes.Topic.Hot(page))
		c.GetTopics(page, "", "hits", routes.Topic.Hot(page))
		time.Sleep(time.Second)

		cache.Delete("topics" + routes.Topic.Good(page))
		c.GetTopics(page, "good = true", "created", routes.Topic.Good(page))
		time.Sleep(time.Second)
	}
}
