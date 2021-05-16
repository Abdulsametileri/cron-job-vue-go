package controllers

import (
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	"net/http"
)

type TokenController interface {
	ValidateToken(http.ResponseWriter, *http.Request)
}

type tokenController struct {
	bc BaseController
	us userservice.UserService
}

func NewTokenController(bc BaseController, us userservice.UserService) TokenController {
	return &tokenController{
		bc: bc,
		us: us,
	}
}

func (t tokenController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		t.bc.Error(w, http.StatusNotFound, ErrMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		t.bc.Error(w, http.StatusBadRequest, ErrTokenNotFound)
		return
	}

	user, err := t.us.GetUserByToken(token)
	if err != nil || user.TelegramId == 0 {
		t.bc.Error(w, http.StatusBadRequest, ErrTokenDoesNotExist)
		return
	}

	t.bc.Data(w, http.StatusOK, nil, "")
}
