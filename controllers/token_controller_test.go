package controllers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTokenController(t *testing.T) {
	us := &userSvc{}
	bc := NewBaseController()
	tokenCtrl := NewTokenController(bc, us)

	t.Run("Any method except get is not allowed", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/validate-token", nil)
		tokenCtrl.ValidateToken(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		props := &Props{}
		_ = json.Unmarshal(body, props)

		assert.Equal(t, props.Code, http.StatusNotFound)
		assert.Equal(t, props.Message, ErrMethodNotAllowed.Error())
	})
	t.Run("Token does not exist error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodGet, "/api/validate-token", nil)
		tokenCtrl.ValidateToken(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		props := &Props{}
		_ = json.Unmarshal(body, props)

		assert.Equal(t, props.Code, http.StatusBadRequest)
		assert.Equal(t, props.Message, ErrTokenNotFound.Error())
	})
	t.Run("user does not exist with specified token", func(t *testing.T) {
		w, req := createHttpReq(http.MethodGet, "/api/validate-token/?token=not-exist-token", nil)
		tokenCtrl.ValidateToken(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		props := &Props{}
		_ = json.Unmarshal(body, props)

		assert.Equal(t, props.Code, http.StatusBadRequest)
		assert.Equal(t, props.Message, ErrTokenDoesNotExist.Error())
	})
}
