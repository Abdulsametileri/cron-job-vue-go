package controllers

import (
	"bytes"
	"fmt"
	"github.com/magiconair/properties/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func writeErrorMsg(err error) string {
	return err.Error() + "\n"
}

func createHttpReq(method string, endpoint string, body *bytes.Buffer) (w *httptest.ResponseRecorder, req *http.Request) {
	if body == nil {
		body = bytes.NewBuffer(make([]byte, 512))
	}
	req = httptest.NewRequest(method, endpoint, body)
	w = httptest.NewRecorder()
	req.Form = url.Values{}
	return
}

func fileUploadRequest() (body *bytes.Buffer, contentType string) {
	file, err := os.Open("../test.png")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	body = new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test")
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}
	return body, writer.FormDataContentType()
}

func TestAlarmController(t *testing.T) {
	userService := &userSvc{}
	awsClient := &awsClient{}
	alarmCtrl := NewAlarmController(userService, awsClient, nil)

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
	t.Run("Getting success when I form params succesfully specified", func(t *testing.T) {
		body, contentType := fileUploadRequest()

		w, req := createHttpReq(http.MethodPost, "/api/create-alarm", body)
		req.Form.Set("token", "sametintokeni")
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

		assert.Equal(t, resp.StatusCode, http.StatusOK)
	})
}
