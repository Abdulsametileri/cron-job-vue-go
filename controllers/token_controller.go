package controllers

import (
	"github.com/Abdulsametileri/cron-job-vue-go/services/userservice"
	"net/http"
)

type TokenController interface {
	ValidateToken(http.ResponseWriter, *http.Request)
}

type tokenController struct {
	us userservice.UserService
}

func NewTokenController(us userservice.UserService) TokenController {
	return &tokenController{
		us: us,
	}
}

func (t tokenController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusNotFound)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, ErrTokenNotFound.Error(), http.StatusBadRequest)
		return
	}

	user, err := t.us.GetUserByToken(token)
	if err != nil || user.TelegramId == 0 {
		http.Error(w, ErrTokenDoesNotExist.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader()
}
