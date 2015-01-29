package main

import (
	"meowtrics/model"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var r render.Render

const (
	APPLICATION_PROTOBUF = "application/x-protobuf"
	APPLICATION_JSON     = "application/json"
	APPLICATION_ALL      = "*/*"
)

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
		case APPLICATION_JSON:
			status, errResp = processJsonPost(req, meowtricsLogger)
		case APPLICATION_PROTOBUF:
			status, errResp = processProtobufPost(req, meowtricsLogger)
		default:
			status, errResp = processUnsupportedMediaTypePost(req, meowtricsLogger)
		}
		//Response body to this post request is always returned as JSON because it is human readable in case of errors
		r.JSON(w, status, errResp)
	})
}

func RetrieveEventHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		acceptHeader := req.Header.Get("Accept")
		id := mux.Vars(req)["id"]
		switch acceptHeader {
		case APPLICATION_PROTOBUF:
			status, data := processProtobufGet(id, meowtricsLogger)
			w.Header().Set("Content-Type", APPLICATION_PROTOBUF)
			r.Data(w, status, data)
		case APPLICATION_JSON, APPLICATION_ALL, "":
			status, event := processJsonGet(id, meowtricsLogger)
			r.JSON(w, status, event)
		default:
			status, errResp := processUnsupportedMediaTypeGet(req, meowtricsLogger)
			r.JSON(w, status, errResp)
		}
	})
}
