package main

import (
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
