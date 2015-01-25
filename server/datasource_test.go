package main

import (
	"meowtrics/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateTestClientEvent() model.ClientEventData {
	id := "testEvent123"
	eventType := model.ClientEventType_UNKNOWN
	timestamp := time.Now().Unix()
	data := "testTestTestTestTest"
	return model.ClientEventData{EventId: &id, EventType: &eventType, Timestamp: &timestamp, Data: &data}
}

func TestStoreEvent_ExpectedData(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	err := StoreEvent(testEvent, meowtricsLogger)
	assert.Nil(t, err, "Error is not nil")
}

func TestStoreEvent_MissingEventId(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	testEvent.EventId = nil
	err := StoreEvent(testEvent, meowtricsLogger)
	assert.Equal(t, InvalidParametersError, err, "Error should be invalid parameters")
}

func TestRetrieveEvent(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	err := StoreEvent(testEvent, meowtricsLogger)
	assert.Nil(t, err, "Error is not nil")

	actualEvent, err := RetrieveEvent(testEvent.GetEventId())
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, testEvent.GetData(), actualEvent.GetData(), "Event data should be equal")

	actualEvent, err = RetrieveEvent("")
	assert.Nil(t, actualEvent, "Event should be nil")
	assert.Equal(t, InvalidParametersError, err, "Error should be invalid parameters")

	actualEvent, err = RetrieveEvent("absentEvent")
	assert.Nil(t, actualEvent, "Event should be nil")
	assert.Equal(t, RecordNotFoundError, err, "Error should be record not found")
}
