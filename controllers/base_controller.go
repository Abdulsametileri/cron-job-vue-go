package controllers

import (
	"encoding/json"
	"net/http"
)

type Props struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type BaseController interface {
	Data(w http.ResponseWriter, code int, data interface{}, message string)
	Error(w http.ResponseWriter, code int, err error)
}

type baseController struct {
}

func NewBaseController() BaseController {
	return &baseController{}
}

func (bc *baseController) Data(w http.ResponseWriter, code int, data interface{}, message string) {
	ret := &Props{
		Code:    code,
		Data:    data,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	jData, _ := json.Marshal(&ret)
	w.Write(jData)
}

func (bc *baseController) Error(w http.ResponseWriter, code int, friendlyErrorForClient error) {
	ret := &Props{
		Code:    code,
		Data:    nil,
		Message: friendlyErrorForClient.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	jData, _ := json.Marshal(&ret)
	w.Write(jData)
}
