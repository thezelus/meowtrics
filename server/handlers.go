package main

import (
	"meowtrics/model"
	"net/http"
	"time"

	"github.com/unrolled/render"
)

var r render.Render

func HeartBeatHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		status := "OK"
		timestamp := time.Now().UTC().String()
		r.JSON(w, http.StatusOK, model.HeartBeat{Status: &status, Timestamp: &timestamp})
	})
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		code := ContentNotFound
		errorMesage := "Nothing to see here"
		response := model.ErrorResponse{Code: &code, ErrorMessage: &errorMesage}
		r.JSON(w, http.StatusNotFound, response)
	})
}
