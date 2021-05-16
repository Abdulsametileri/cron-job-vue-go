package jobservice

import (
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/jobrepo"
)

type JobService interface {
	AddJob(job models.Job) error
	GetJobByFields(map[string]interface{}) (models.Job, error)
}

type jobService struct {
	Repo jobrepo.Repo
}

func NewJobService(jRepo jobrepo.Repo) JobService {
	return &jobService{
		Repo: jRepo,
	}
}

func (j jobService) AddJob(job models.Job) error {
	return j.Repo.AddJob(job)
}

func (j jobService) GetJobByFields(m map[string]interface{}) (models.Job, error) {
	return j.Repo.GetJobByFields(m)
}
