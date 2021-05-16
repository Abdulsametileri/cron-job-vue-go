package controllers

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTokenController(t *testing.T) {
	us := &userSvc{}
	tokenCtrl := NewTokenController(us)

	t.Run("Any method except get is not allowed", func(t *testing.T) {
		w, req := createHttpReq(http.MethodPost, "/api/validate-token", nil)
		tokenCtrl.ValidateToken(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusNotFound)
		assert.Equal(t, string(body), writeErrorMsg(ErrMethodNotAllowed))
	})
	t.Run("Token does not exist error", func(t *testing.T) {
		w, req := createHttpReq(http.MethodGet, "/api/validate-token", nil)
		tokenCtrl.ValidateToken(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(body), writeErrorMsg(ErrTokenNotFound))
	})
	t.Run("user does not exist with specified token", func(t *testing.T) {
		w, req := createHttpReq(http.MethodGet, "/api/validate-token/?token=not-exist-token", nil)
		tokenCtrl.ValidateToken(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
		assert.Equal(t, string(body), writeErrorMsg(ErrTokenDoesNotExist))
	})
}
