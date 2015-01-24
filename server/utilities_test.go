package main

import (
	"net/http"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLogger     = log.New()
	testConfigName = "meowtricsConfig"
	testConfigPath = ""
)

func TestLoadAppProperties(t *testing.T) {
	err := LoadAppProperties(testConfigName, testConfigPath, testLogger)
	assert.NoError(t, err, "Error loading config")
	assert.Equal(t, "3003", viper.GetString("appPort"), "App port property doesn't match")
	assert.Equal(t, "10", viper.GetString("appGracefulShutdownTimeinSeconds"), "appGracefulShutdownTimeinSeconds property doesn't match")
}

func TestNotFoundHandler(t *testing.T) {

	test := GenerateHandleTester(t, NotFoundHandler())

	w := test("GET", "")

	assert.Equal(t, http.StatusNotFound, w.Code, "Not found handler")
	//assert.Equal(t, `{"error":{"code":"NOT_FOUND","message":"Nothing to see here","description":""}}`, w.Body.String(), "Not found handler response")

	w = test("POST", "")

	assert.Equal(t, http.StatusNotFound, w.Code, "Not found handler")
	//assert.Equal(t, `{"error":{"code":"NOT_FOUND","message":"Nothing to see here","description":""}}`, w.Body.String(), "Not found handler response")
}

func TestHeartBeatHandler(t *testing.T) {

	test := GenerateHandleTester(t, HeartBeatHandler())

	w := test("GET", "")
	assert.Equal(t, http.StatusOK, w.Code, "HeartBeat response status")
}
