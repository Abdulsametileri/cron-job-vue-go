package controllers

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAlarmController(t *testing.T) {
	userService := &userSvc{}
	jobService := &jobSvc{}
	awsClient := &awsClient{}
	telegramClient := &telegramClient{}

	alarmCtrl := NewAlarmController(userService, jobService, awsClient, telegramClient, nil)

	t.Run("Is Get Not Allowed", func(t *testing.T) {
		w, req := createHttpReq(http.MethodGet, "/api/create-alarm", nil)
		alarmCtrl.CreateAlarm(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusNotFound)
		assert.Equal(t, string(body), writeErrorMsg(ErrMethodNotAllowed))
	})
	t.Run("Getting Token error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		alarmCtrl.CreateAlarm(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(body), writeErrorMsg(ErrTokenNotFound))
	})
	t.Run("Non empty token but getting empty name error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")

		alarmCtrl.CreateAlarm(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(body), writeErrorMsg(ErrNameNotFound))
	})
	t.Run("Non empty {token, name} but getting empty time error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")
		req.Form.Set("name", "name")

		alarmCtrl.CreateAlarm(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(body), writeErrorMsg(ErrTimeNotFound))
	})
	t.Run("Non empty {token, name, time} but getting empty repeatType error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")

		alarmCtrl.CreateAlarm(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(body), writeErrorMsg(ErrRepeatTypeNotFound))
	})
	t.Run("Non empty {token,name,time,repeatType} but getting reading image file err", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", nil)
		req.Form.Set("token", "token")
		req.Form.Set("name", "name")
		req.Form.Set("time", "23:14")
		req.Form.Set("repeatType", "5")

		alarmCtrl.CreateAlarm(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(body), writeErrorMsg(ErrReadingFile))
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

		resp := w.Result()
		bodyd, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		fmt.Println(string(bodyd))

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(bodyd), writeErrorMsg(ErrDb))
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

		resp := w.Result()
		bodyd, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		fmt.Println(string(bodyd))

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(bodyd), writeErrorMsg(ErrTokenDoesNotExist))
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

		resp := w.Result()
		bodyd, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		fmt.Println(string(bodyd))

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(bodyd), writeErrorMsg(ErrS3Upload))
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

		resp := w.Result()
		bodyd, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(bodyd), writeErrorMsg(ErrGettingJob))
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

		resp := w.Result()
		bodyd, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(bodyd), writeErrorMsg(ErrJobAlreadyExist))
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

		resp := w.Result()
		bodyd, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(bodyd), writeErrorMsg(ErrDeleteFileS3))
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

		resp := w.Result()
		bodyd, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(bodyd), writeErrorMsg(ErrAddingJob))
	})
}
