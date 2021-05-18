package utils

import (
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	t.Run("Write job pretty", func(t *testing.T) {
		job := models.Job{
			Tag:            "tag",
			Name:           "name",
			UserTelegramId: 123,
			UserToken:      "token",
			ImageUrl:       "http://etc.",
			RepeatType:     "1",
			Time:           "22:10",
			Status:         models.JobValid,
		}
		fmt.Println(PrettyPrint(job))
	})
}
