package jobservice

import (
	"errors"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/repository/jobrepo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJobService_AddJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobRepo := jobrepo.NewMockRepo(ctrl)

	noJobSpecifiedErr := errors.New("No job specified")
	jobRepo.
		EXPECT().
		AddJob(models.Job{}).
		Return(noJobSpecifiedErr).
		Times(1)

	jobService := NewJobService(jobRepo)

	err := jobService.AddJob(models.Job{})

	require.ErrorIs(t, err, noJobSpecifiedErr)
}
