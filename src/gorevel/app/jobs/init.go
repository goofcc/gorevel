package jobs

import (
	"time"

	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
)

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("cron.update_cache", UpdateCache{})
		jobs.In(5*time.Second, UpdateCache{})
	})
}
