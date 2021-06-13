package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAlarmController(t *testing.T) {
	bc := NewBaseController()
	userService := &userSvc{}
	jobService := &jobSvc{}
	awsClient := &awsClient{}
	telegramClient := &telegramClient{}

	alarmCtrl := NewAlarmController(bc, userService, jobService, awsClient, telegramClient, nil)

	t.Run("Paginate Alarm", func(t *testing.T) {
		t.Run("Is get not allow", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/paginate-alarm", nil)
			alarmCtrl.PaginateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusNotFound)
			assert.Equal(t, res.Message, ErrMethodNotAllowed.Error())
		})
		t.Run("Getting second page with size 5 is empty", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/paginate-alarm?pageNo=2&pageSize=5", nil)
			alarmCtrl.PaginateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusOK)
			assert.Equal(t, res.Message, "")
		})
		t.Run("Getting first page with size 5 is full", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/paginate-alarm?pageNo=1&pageSize=5", nil)
			alarmCtrl.PaginateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusOK)
			assert.Equal(t, res.Message, "")

			val, cvrtOk := res.Data.(map[string]interface{})
			assert.True(t, cvrtOk)

			var jobs []models.Job
			jS, err := json.Marshal(val["jobs"])
			assert.NoError(t, err)
			json.Unmarshal(jS, &jobs)

			assert.EqualValues(t, 100, val["total"])
			assert.Len(t, jobs, 5)
			for i := 0; i < 5; i++ {
				assert.Equal(t, jobs[i].Tag, fmt.Sprintf("tag-%d", i))
				assert.Equal(t, jobs[i].Name, fmt.Sprintf("name-%d", i))
				assert.EqualValues(t, jobs[i].Status, models.JobValid)
			}
		})
	})

	t.Run("CreateAlarm", func(t *testing.T) {
		t.Run("Is Get Not Allowed", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/create-alarm", nil)
			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusNotFound)
			assert.Equal(t, res.Message, ErrMethodNotAllowed.Error())
		})
		t.Run("Getting Token error", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrFieldNotFound("token").Error())
		})
		t.Run("Non empty token but getting empty name error", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
			req.Form.Set("token", "token")

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrFieldNotFound("name").Error())
		})
		t.Run("Non empty {token, name} but getting empty time error", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
			req.Form.Set("token", "token")
			req.Form.Set("name", "name")

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrFieldNotFound("time").Error())
		})
		t.Run("Non empty {token, name, time} but getting empty repeatType error", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
			req.Form.Set("token", "token")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrFieldNotFound("repeat time").Error())
		})
		t.Run("Non empty {token,name,time,repeatType} but getting reading image file err", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
			req.Form.Set("token", "token")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrReadingFile.Error())
		})
		t.Run("Getting token err occured in db", func(t *testing.T) {
			body, contentType := fileUploadRequest()

			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
			req.Form.Set("token", "db-err")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")
			req.Form.Set("fileName", "test")
			req.Form.Set("fileType", "image/png")

			req.Header.Add("Content-Type", contentType)

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrDb.Error())
		})
		t.Run("when job exist db error", func(t *testing.T) {
			body, contentType := fileUploadRequest()

			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
			req.Form.Set("token", "job-already-exist-db-error")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")
			req.Form.Set("fileName", "arbitrary-name")
			req.Form.Set("fileType", "image/png")

			req.Header.Add("Content-Type", contentType)

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrGettingJob.Error())
		})
		t.Run("when job already exist error", func(t *testing.T) {
			body, contentType := fileUploadRequest()

			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
			req.Form.Set("token", "job-already-exist")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")
			req.Form.Set("fileName", "arbitrary-name")
			req.Form.Set("fileType", "image/png")

			req.Header.Add("Content-Type", contentType)

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrJobAlreadyExist.Error())
		})
		t.Run("Getting non exist token error", func(t *testing.T) {
			body, contentType := fileUploadRequest()

			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
			req.Form.Set("token", "not-exist-token")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")
			req.Form.Set("fileName", "test")
			req.Form.Set("fileType", "image/png")

			req.Header.Add("Content-Type", contentType)

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrTokenDoesNotExistInUrl.Error())
		})
		t.Run("Getting s3 upload error", func(t *testing.T) {
			body, contentType := fileUploadRequest()

			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
			req.Form.Set("token", "sametintokeni")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")
			req.Form.Set("fileName", "badFileName")
			req.Form.Set("fileType", "image/png")

			req.Header.Add("Content-Type", contentType)

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrS3Upload.Error())
		})
		t.Run("When job created, job err occured, delete uploaded file in s3 also occured", func(t *testing.T) {
			body, contentType := fileUploadRequest()

			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
			req.Form.Set("token", "sametintokeni")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")
			req.Form.Set("fileName", "error-scenario-with-s3")
			req.Form.Set("fileType", "image/png")

			req.Header.Add("Content-Type", contentType)

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrDeleteFileS3.Error())
		})
		t.Run("When job created, job err occured, delete uploaded file in s3 is success return add job error", func(t *testing.T) {
			body, contentType := fileUploadRequest()

			w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
			req.Form.Set("token", "sametintokeni")
			req.Form.Set("name", "name")
			req.Form.Set("time", "23:14")
			req.Form.Set("repeatType", "5")
			req.Form.Set("fileName", "error-scenario-job")
			req.Form.Set("fileType", "image/png")

			req.Header.Add("Content-Type", contentType)

			alarmCtrl.CreateAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrAddingJob.Error())
		})
	})

	t.Run("ListAlarm", func(t *testing.T) {
		t.Run("Is Post not allowed", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/list-alarm", nil)
			alarmCtrl.ListAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusNotFound)
			assert.Equal(t, res.Message, ErrMethodNotAllowed.Error())
		})
		t.Run("Error when token not specified", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/list-alarm", nil)
			alarmCtrl.ListAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrTokenDoesNotExistInUrl.Error())
		})
		t.Run("Error occured in db when validating the token", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/list-alarm/?token=db-err", nil)
			alarmCtrl.ListAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrDb.Error())
		})
		t.Run("Error user cannot exist given token", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/list-alarm/?token=not-exist-token", nil)
			alarmCtrl.ListAlarm(w, req)

			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrUserDoesNotExist.Error())
		})
		t.Run("Error getting the job list", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/list-alarm/?token=job-list-err", nil)
			alarmCtrl.ListAlarm(w, req)

			res := parseBody(w)
			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrGettingJobList.Error())
		})
		t.Run("Getting job list empty", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/list-alarm/?token=123", nil)
			alarmCtrl.ListAlarm(w, req)

			res := parseBody(w)
			assert.Equal(t, res.Code, http.StatusOK)

			val, _ := res.Data.([]interface{})
			lenItems := len(val)

			assert.Equal(t, lenItems, 0)
		})
		t.Run("Getting the job list item", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/list-alarm/?token=three-job-list-item", nil)
			alarmCtrl.ListAlarm(w, req)

			res := parseBody(w)
			assert.Equal(t, res.Code, http.StatusOK)

			val, _ := res.Data.([]interface{})
			lenItems := len(val)

			assert.Equal(t, lenItems, 3)
		})
	})

	t.Run("DeleteAlarm", func(t *testing.T) {
		t.Run("Is Get Not Allowed", func(t *testing.T) {
			w, req := createHttpReq(http.MethodGet, "/api/delete-alarm", nil)
			alarmCtrl.DeleteAlarm(w, req)
			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusNotFound)
			assert.Equal(t, res.Message, ErrMethodNotAllowed.Error())
		})
		t.Run("Is tag is speficied in url", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/delete-alarm", nil)
			alarmCtrl.DeleteAlarm(w, req)
			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrTagDoesNotExistInUrl.Error())
		})
		t.Run("Error occured when deleting the job in db with specified tag", func(t *testing.T) {
			w, req := createHttpReq(http.MethodPost, "/api/delete-alarm?tag=error-tag", nil)
			alarmCtrl.DeleteAlarm(w, req)
			res := parseBody(w)

			assert.Equal(t, res.Code, http.StatusBadRequest)
			assert.Equal(t, res.Message, ErrJobDelete.Error())
		})
	})
}
