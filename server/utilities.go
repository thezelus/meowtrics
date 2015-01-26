package main

import (
	"errors"
	"fmt"
	"os"

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
	HeaderNotRecognized      = "HEADER_NOT_RECOGNIZED"
	RecordNotFound           = "REQUESTED_RECORD_NOT_FOUND"
	Fatal                    = "FATAL_OPERATION"
	UnsupportedMedia         = "UNSUPPORTED_MEDIA_TYPE"
)

var (
	InvalidParametersError = errors.New(InvalidRequestParameters)
	RecordNotFoundError    = errors.New(RecordNotFound)
	FatalError             = errors.New(Fatal)
	UnsupportedMediaError  = errors.New(UnsupportedMedia)
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
