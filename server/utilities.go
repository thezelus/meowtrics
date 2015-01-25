package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DeploymentConfigDefaultPath = "$GOPATH/bin/config/meowtrics/"
)

const (
	ContentNotFound          = "NOT_FOUND"
	InvalidRequestParameters = "INVALID_REQUEST_PARAMETERS"
	MalformedRequest         = "MALFORMED_REQUEST"
	HeaderAbsent             = "HEADER_MISSING"
	RecordNotFound           = "REQUESTED_RECORD_NOT_FOUND"
	FatalError               = "FATAL_ERROR"
)

var (
	InvalidParametersError = errors.New(InvalidRequestParameters)
	RecordNotFoundError    = errors.New(RecordNotFound)
)

func InitializeLogger(file *os.File, logFileName string, logger *log.Logger, format log.Formatter) error {

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Log file cannot be opened")
		file.Close()
		file = os.Stdout
	}

	logger.Out = file
	if format != nil {
		logger.Formatter = format
	}

	return err
}

//Looks for config file in default deployment directory and the injected confiPath directory

func LoadAppProperties(configName string, configPath string, logger *log.Logger) error {
	logger.Println("Loading configuration from " + configName)

	viper.SetConfigName(configName)
	viper.AddConfigPath(DeploymentConfigDefaultPath)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		logger.Errorln("Configuration couldn't be initialized: " + err.Error())
		return err
	}

	return nil
}

//--------Test utils------

type HandleTester func(method string, params string) *httptest.ResponseRecorder

// Given the current test runner and an http.Handler, generate a
// HandleTester which will test its the handler against its input

func GenerateHandleTester(t *testing.T, handleFunc http.Handler) HandleTester {

	// Given a method type ("GET", "POST", etc) and
	// parameters, serve the response against the handler
	// and return the ResponseRecorder.

	return func(method string, params string) *httptest.ResponseRecorder {

		req, err := http.NewRequest(method, "", strings.NewReader(params))
		if err != nil {
			t.Errorf("%v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Body.Close()
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}
