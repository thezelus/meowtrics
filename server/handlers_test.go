package main

import (
	"encoding/json"
	"meowtrics/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

//---------------Test Utils ------------------------

type HandleTester func(method string, params string) *httptest.ResponseRecorder

// Given the current test runner and an http.Handler, generate a
// HandleTester which will test its the handler against its input

func GenerateHandleTester(t *testing.T, handleFunc http.Handler, contentType string) HandleTester {

	// Given a method type ("GET", "POST", etc) and
	// parameters, serve the response against the handler
	// and return the ResponseRecorder.

	return func(method string, params string) *httptest.ResponseRecorder {

		req, err := http.NewRequest(method, "", strings.NewReader(params))
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set("Content-Type", contentType)
		req.Body.Close()
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}

func generateTestClientEvent() model.ClientEventData {
	id := "testEvent123"
	eventType := model.ClientEventType_UNKNOWN
	timestamp := time.Now().Unix()
	data := "testTestTestTestTest"
	return model.ClientEventData{EventId: &id, EventType: &eventType, Timestamp: &timestamp, Data: &data}
}

func generateTestClientEventUploadRequest_Valid() model.ClientEventUploadRequest {
	var mockReq model.ClientEventUploadRequest
	requestId := "testRequestId"
	deviceType := "testDeviceAndroid"
	testEvent := generateTestClientEvent()

	mockReq.RequestId = &requestId
	mockReq.DeviceType = &deviceType
	mockReq.Events = append(mockReq.Events, &testEvent)
	return mockReq
}

func generateTestClientEventUploadRequest_Invalid() model.ClientEventUploadRequest {
	var mockReq model.ClientEventUploadRequest
	requestId := "testRequestId"
	deviceType := "testDeviceAndroid"
	testEvent := generateTestClientEvent()
	testEvent.EventId = nil

	mockReq.RequestId = &requestId
	mockReq.DeviceType = &deviceType
	mockReq.Events = append(mockReq.Events, &testEvent)
	return mockReq
}

//--------------------Handler Tests----------------------

func TestNotFoundHandler(t *testing.T) {

	test := GenerateHandleTester(t, NotFoundHandler(), "application/json")

	w := test("GET", "")
	assert.Equal(t, http.StatusNotFound, w.Code, "Not found handler")

	w = test("POST", "")
	assert.Equal(t, http.StatusNotFound, w.Code, "Not found handler")

}

func TestHeartBeatHandler(t *testing.T) {

	test := GenerateHandleTester(t, HeartBeatHandler(), "application/json")

	w := test("GET", "")
	assert.Equal(t, http.StatusOK, w.Code, "HeartBeat response status")
}

//----------------------JSON tests------------------------------

func TestCreateEventHandler_ValidJsonRequest(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GenerateHandleTester(t, CreateEventHandler(), "application/json")
	uploadReq := generateTestClientEventUploadRequest_Valid()
	jsonReq, err := json.Marshal(uploadReq)
	if err != nil {
		panic("Cannot marshal data. Error: " + err.Error())
	}

	w := test("POST", string(jsonReq))
	assert.Equal(t, http.StatusOK, w.Code, "Valid JSON request should be properly posted")

	_, ok := eventMap["testEvent123"]
	assert.True(t, ok, "Map should contain event")
}

func TestCreateEventHandler_InvalidJsonRequestWithNoEventId(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GenerateHandleTester(t, CreateEventHandler(), "application/json")
	uploadReq := generateTestClientEventUploadRequest_Invalid()
	jsonReq, err := json.Marshal(uploadReq)
	if err != nil {
		panic("Cannot marshal data. Error: " + err.Error())
	}

	w := test("POST", string(jsonReq))
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid data without eventid should receive an error")

	assert.Equal(t, 0, len(eventMap), "eventMap should be empty")
}

func TestCreateEventHandler_MalformedJsonData(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GenerateHandleTester(t, CreateEventHandler(), "application/json")

	w := test("POST", "randomString")
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid data should receive an error")

	assert.Equal(t, 0, len(eventMap), "eventMap should be empty")
}

//--------------------------------Unsupported media type-----------

func TestCreateEventHandler_UnsupportedMediaTypePost(t *testing.T) {
	test := GenerateHandleTester(t, CreateEventHandler(), "application/meow")
	w := test("POST", "meowtrics")
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code, "UnsupportedMedia should be the header")
}

//------------------------------Protobuf tests----------------------

func TestCreateEventHandler_ValidProtobufRequest(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GenerateHandleTester(t, CreateEventHandler(), "application/x-protobuf")
	uploadReq := generateTestClientEventUploadRequest_Valid()
	protoBytes, err := proto.Marshal(&uploadReq)
	if err != nil {
		panic("Cannot marshal data. Error: " + err.Error())
	}

	w := test("POST", string(protoBytes))
	assert.Equal(t, http.StatusOK, w.Code, "Valid protocol buffer request should be properly posted")

	assert.Equal(t, 1, len(eventMap), "eventMap should have one entry")

	_, ok := eventMap["testEvent123"]
	assert.True(t, ok, "Map should contain event")
}

/*
//This test is not required because it is impossible to marshal model into protobuf without required field

func TestCreateEventHandler_InvalidProtobufRequestWithNoEventId(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GenerateHandleTester(t, CreateEventHandler(), "application/x-protobuf")
	uploadReq := generateTestClientEventUploadRequest_Invalid()
	protoBytes, err := proto.Marshal(&uploadReq)
	if err != nil {
		panic("Cannot marshal data. Error: " + err.Error())
	}

	w := test("POST", string(protoBytes))
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid data without eventid should receive an error")

	assert.Equal(t, 0, len(eventMap), "eventMap should be empty")
}
*/

func TestCreateEventHandler_MalformedProtobufData(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GenerateHandleTester(t, CreateEventHandler(), "application/x-protobuf")

	w := test("POST", "randomString")
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid data should receive an error")

	assert.Equal(t, 0, len(eventMap), "eventMap should be empty")
}
