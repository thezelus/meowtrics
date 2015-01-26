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

func CreateEventHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contentHeader := req.Header.Get("Content-Type")
		var status int
		var errResp *model.ErrorResponse
		switch contentHeader {
		case "application/json":
			status, errResp = processJsonPost(req, meowtricsLogger)
		case "application/x-protobuf":
			status, errResp = processProtobufPost(req, meowtricsLogger)
		default:
			status, errResp = processUnsupportedMediaTypePost(req, meowtricsLogger)
		}
		//Response body to this post request is always returned as JSON because it is human readable in case of errors
		r.JSON(w, status, errResp)
	})
}
