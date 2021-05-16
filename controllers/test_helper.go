package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

func writeErrorMsg(err error) string {
	return err.Error() + "\n"
}

func parseBody(w *httptest.ResponseRecorder) Props {
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	props := Props{}
	_ = json.Unmarshal(body, &props)

	return props
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
	file, err := os.Open("./test.png")
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
