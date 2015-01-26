package main

import (
	"meowtrics/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreEvent_ExpectedData(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	err := StoreEvent(testEvent)
	assert.Nil(t, err, "Error is not nil")
}

func TestStoreEvent_MissingEventId(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	testEvent.EventId = nil
	err := StoreEvent(testEvent)
	assert.Equal(t, InvalidParametersError, err, "Error should be invalid parameters")
}

func TestRetrieveEvent(t *testing.T) {
	eventMap = make(map[string]model.ClientEventData)
	testEvent := generateTestClientEvent()
	err := StoreEvent(testEvent)
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
