package main

import (
	"fmt"
	"meowtrics/model"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gopkg.in/tylerb/graceful.v1"
)

var (
	file            *os.File
	meowtricsLogger = log.New()
	logFileName     = "log-meowtrics.log"
	configFileName  = "meowtricsConfig"
	router          *mux.Router
	postSubrouter   *mux.Router
	getSubrouter    *mux.Router
	eventMap        map[string]model.ClientEventData
)

func init() {
	initApp()
	initRouter()
}

func initApp() {

	err := InitializeLogger(file, logFileName, meowtricsLogger, new(log.JSONFormatter))
	if err != nil {
		fmt.Println("Error initializing logger, defaulting to stdout. Error: " + err.Error())
	}

	err = LoadAppProperties(configFileName, DeploymentConfigDefaultPath, meowtricsLogger)
	if err != nil {
		meowtricsLogger.Panicln("Error reading app properties:" + err.Error())
	}

	eventMap = make(map[string]model.ClientEventData)
}

func initRouter() {

	router = mux.NewRouter()
	router.StrictSlash(true)

	postSubrouter = router.PathPrefix("/v1/").Methods("POST").Subrouter()
	postSubrouter.Handle("/events", CreateEventHandler())

	getSubrouter = router.PathPrefix("/v1/").Methods("GET").Subrouter()
	getSubrouter.Handle("/events/{id:[0-9]+}", RetrieveEventHandler())

	router.Handle("/heartbeat", HeartBeatHandler())
	router.NotFoundHandler = NotFoundHandler()
}

func main() {

	n := negroni.Classic()
	n.UseHandler(router)

	appGracefulShutdownTimeinSeconds, err := strconv.Atoi(viper.GetString("appGracefulShutdownTimeinSeconds"))
	if err != nil {
		meowtricsLogger.Errorln("Error reading config: " + err.Error())
	}

	graceful.Run(":"+viper.GetString("appPort"), time.Duration(appGracefulShutdownTimeinSeconds)*time.Second, n)

	defer cleanup(file, meowtricsLogger)
}

func cleanup(file *os.File, logger *log.Logger) {
	logger.Println("Starting clean up")

	logger.Println("Closing file stream")
	file.Close()
}
