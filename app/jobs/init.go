package jobs

import (
	"time"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("cron.update_cache", UpdateCache{})
		jobs.In(5*time.Second, UpdateCache{})
	})
}
