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

type PostHandleTester func(method string, params string) *httptest.ResponseRecorder

// Given the current test runner and an http.Handler, generate a
// HandleTester which will test its the handler against its input

func GeneratePostHandleTester(t *testing.T, handleFunc http.Handler, contentType string) PostHandleTester {

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

//Util function specifically for setting accept header for testing GET requests

type GetHandleTester func(location string, acceptHeader string) *httptest.ResponseRecorder

func GenerateGetHandleTester(t *testing.T, handleFunc http.Handler) GetHandleTester {

	return func(location string, acceptHeader string) *httptest.ResponseRecorder {

		req, err := http.NewRequest("GET", location, strings.NewReader(""))
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set("Accept", acceptHeader)
		req.Body.Close()
		w := httptest.NewRecorder()
		getSubrouter.ServeHTTP(w, req)
		return w
	}
}

func generateTestClientEvent() model.ClientEventData {
	id := "123"
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

	test := GeneratePostHandleTester(t, NotFoundHandler(), "application/json")

	w := test("GET", "")
	assert.Equal(t, http.StatusNotFound, w.Code, "Not found handler")

	w = test("POST", "")
	assert.Equal(t, http.StatusNotFound, w.Code, "Not found handler")

}

func TestHeartBeatHandler(t *testing.T) {

	test := GeneratePostHandleTester(t, HeartBeatHandler(), "application/json")

	w := test("GET", "")
	assert.Equal(t, http.StatusOK, w.Code, "HeartBeat response status")
}

//----------------------JSON POST tests------------------------------

func TestCreateEventHandler_ValidJsonRequest(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GeneratePostHandleTester(t, CreateEventHandler(), "application/json")
	uploadReq := generateTestClientEventUploadRequest_Valid()
	jsonReq, err := json.Marshal(uploadReq)
	if err != nil {
		panic("Cannot marshal data. Error: " + err.Error())
	}

	w := test("POST", string(jsonReq))
	assert.Equal(t, http.StatusOK, w.Code, "Valid JSON request should be properly posted")

	_, ok := eventMap["123"]
	assert.True(t, ok, "Map should contain event")
}

func TestCreateEventHandler_InvalidJsonRequestWithNoEventId(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GeneratePostHandleTester(t, CreateEventHandler(), "application/json")
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
	test := GeneratePostHandleTester(t, CreateEventHandler(), "application/json")

	w := test("POST", "randomString")
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid data should receive an error")

	assert.Equal(t, 0, len(eventMap), "eventMap should be empty")
}

//--------------------------------Unsupported media type-----------

func TestCreateEventHandler_UnsupportedMediaTypePost(t *testing.T) {
	test := GeneratePostHandleTester(t, CreateEventHandler(), "application/meow")
	w := test("POST", "meowtrics")
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code, "UnsupportedMedia should be the header")
}

func TestRetrieveEventHandler_UnsupportedMediaTypeGet(t *testing.T) {
	testGet := GenerateGetHandleTester(t, RetrieveEventHandler())

	w := testGet("/v1/events/1", "application/meow")
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code, "UnsupportedMedia should be the header")
}

//------------------------------Protobuf POST tests----------------------

func TestCreateEventHandler_ValidProtobufRequest(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GeneratePostHandleTester(t, CreateEventHandler(), "application/x-protobuf")
	uploadReq := generateTestClientEventUploadRequest_Valid()
	protoBytes, err := proto.Marshal(&uploadReq)
	if err != nil {
		panic("Cannot marshal data. Error: " + err.Error())
	}

	w := test("POST", string(protoBytes))
	assert.Equal(t, http.StatusOK, w.Code, "Valid protocol buffer request should be properly posted")

	assert.Equal(t, 1, len(eventMap), "eventMap should have one entry")

	_, ok := eventMap["123"]
	assert.True(t, ok, "Map should contain event")
}

/*
//This test is not required because it is impossible to marshal model into protobuf without required field

func TestCreateEventHandler_InvalidProtobufRequestWithNoEventId(t *testing.T) {

	eventMap = make(map[string]model.ClientEventData)
	test := GeneratePostHandleTester(t, CreateEventHandler(), "application/x-protobuf")
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
	test := GeneratePostHandleTester(t, CreateEventHandler(), "application/x-protobuf")

	w := test("POST", "randomString")
	assert.Equal(t, http.StatusBadRequest, w.Code, "Invalid data should receive an error")

	assert.Equal(t, 0, len(eventMap), "eventMap should be empty")
}

//--------------------------------JSON GET tests----------------------------

func TestRetrieveEventHandler_ValidRouteVariable_JSON(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	StoreEvent(testEvent)

	_, ok := eventMap["123"]
	assert.True(t, ok, "Map should contain event")

	testGet := GenerateGetHandleTester(t, RetrieveEventHandler())
	w := testGet("/v1/events/123", "application/json")

	assert.Equal(t, http.StatusOK, w.Code, "Http status should be 200")

	actualEvent := new(model.ClientEventData)
	err := json.Unmarshal([]byte(w.Body.String()), actualEvent)
	if err != nil {
		panic("Error unmarshalling json response: " + err.Error())
	}

	assert.Equal(t, testEvent.GetData(), actualEvent.GetData(), "Event data should be equal")
}

func TestRetrieveEventHandler_RecordNotFound_JSON(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	StoreEvent(testEvent)

	_, ok := eventMap["123"]
	assert.True(t, ok, "Map should contain event")

	testGet := GenerateGetHandleTester(t, RetrieveEventHandler())
	w := testGet("/v1/events/12", "application/json")

	assert.Equal(t, http.StatusNotFound, w.Code, "Http status should be 404")
}

func TestRetrieveEventHandler_InValidRouteVariable_JSON(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	newEventId := "abc"
	testEvent.EventId = &newEventId
	StoreEvent(testEvent)

	_, ok := eventMap[newEventId]
	assert.True(t, ok, "Map should contain event")

	testGet := GenerateGetHandleTester(t, RetrieveEventHandler())
	w := testGet("/v1/events/abc", "application/json")

	assert.Equal(t, http.StatusNotFound, w.Code, "Http status should be 404")
}

//-------------------------Protobuf GET-----------------

func TestRetrieveEventHandler_ValidRouteVariable_Protobuf(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	StoreEvent(testEvent)

	_, ok := eventMap["123"]
	assert.True(t, ok, "Map should contain event")
	assert.Equal(t, len(eventMap), 1, "eventMap should have only 1 entry")

	testGet := GenerateGetHandleTester(t, RetrieveEventHandler())
	w := testGet("/v1/events/123", "application/x-protobuf")

	assert.Equal(t, http.StatusOK, w.Code, "Http status should be 200")

	actualEvent := new(model.ClientEventData)
	err := proto.Unmarshal([]byte(w.Body.String()), actualEvent)
	if err != nil {
		panic("Error unmarshalling protobuf response: " + err.Error())
	}

	assert.Equal(t, testEvent.GetData(), actualEvent.GetData(), "Event data should be equal")
}

func TestRetrieveEventHandler_RecordNotFound_Protobuf(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	StoreEvent(testEvent)

	_, ok := eventMap["123"]
	assert.True(t, ok, "Map should contain event")

	testGet := GenerateGetHandleTester(t, RetrieveEventHandler())
	w := testGet("/v1/events/12", "application/x-protobuf")

	assert.Equal(t, http.StatusNotFound, w.Code, "Http status should be 404")
}
