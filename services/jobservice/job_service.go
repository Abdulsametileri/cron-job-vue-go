package jobservice

import (
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/jobrepo"
)

type JobService interface {
	ListAllValidJobsByToken(token string) ([]models.Job, error)
	ListAllValidJobs() ([]models.Job, error)
	AddJob(job models.Job) error
	GetJobByFields(map[string]interface{}) (models.Job, error)
	DeleteJobByTag(tag string) error
}

type jobService struct {
	Repo jobrepo.Repo
}

func NewJobService(jRepo jobrepo.Repo) JobService {
	return &jobService{
		Repo: jRepo,
	}
}

func (j jobService) ListAllValidJobs() ([]models.Job, error) {
	return j.Repo.ListAllValidJobs()
}

func (j jobService) ListAllValidJobsByToken(token string) ([]models.Job, error) {
	return j.Repo.ListAllValidJobsByToken(token)
}

func (j jobService) AddJob(job models.Job) error {
	return j.Repo.AddJob(job)
}

func (j jobService) GetJobByFields(m map[string]interface{}) (models.Job, error) {
	return j.Repo.GetJobByFields(m)
}

func (j jobService) DeleteJobByTag(tag string) error {
	return j.Repo.DeleteJobByTag(tag)
}
