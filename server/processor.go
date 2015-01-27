package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"meowtrics/model"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

//------------------GET-----------------------

func processJsonGet(id string, logger *log.Logger) (int, *model.ClientEventData) {

	event, err := RetrieveEvent(id)
	switch err {
	case nil:
		return http.StatusOK, event
	case RecordNotFoundError:
		logger.WithFields(log.Fields{"method": "processJsonGet", "id": id, "error": RecordNotFoundError.Error()}).Infoln("Record not found")
		return http.StatusNotFound, nil
	case InvalidParametersError:
		logger.WithFields(log.Fields{"method": "processJsonGet", "error": InvalidParametersError.Error()}).Warningln("Invalid id passed through router")
		return http.StatusInternalServerError, nil
	}
	
//This should never be executed
	return http.StatusInternalServerError, nil
}

func processProtobufGet(id string, logger *log.Logger) (int, []byte) {

	event, err := RetrieveEvent(id)
	if err != nil {
		switch err {
		case RecordNotFoundError:
			logger.WithFields(log.Fields{"method": "processProtobufGet", "id": id, "error": RecordNotFoundError.Error()}).Infoln("Record not found")
			return http.StatusNotFound, nil
		case InvalidParametersError:
			logger.WithFields(log.Fields{"method": "processProtobufGet", "error": InvalidParametersError.Error()}).Warningln("    Invalid id passed through router")
			return http.StatusInternalServerError, nil
		}
	}

	protoBytes, err := proto.Marshal(event)
	if err != nil {
		logger.WithFields(log.Fields{"method": "processProtobufGet", "id": id, "error": err.Error()}).Warningln("Error marshaling model to protocol buffer byte array")
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, protoBytes
}

func processUnsupportedMediaTypeGet(req *http.Request, logger *log.Logger) (int, *model.ErrorResponse) {
	logger.WithFields(log.Fields{"method": "processUnsupportedMediaTypeGet", "error": UnsupportedMedia}).Infoln("Accept: " + req.Header.Get("Accept"))

	errCode := UnsupportedMedia
	errMsg := "Accept specifies a media type that is not supported by this resource"
	return http.StatusUnsupportedMediaType, &model.ErrorResponse{Code: &errCode, ErrorMessage: &errMsg}
}

//-----------------POST-----------------------

func processJsonPost(req *http.Request, logger *log.Logger) (int, *model.ErrorResponse) {

	uploadRequest, err := decodeJson(req.Body)
	if err != nil {
		logger.WithFields(log.Fields{"method": "processJsonPost", "error": err.Error()}).Warningln("Error decoding json")

		errCode := MalformedRequest
		errMsg := "Request body contains malformed JSON"
		return http.StatusBadRequest, &model.ErrorResponse{Code: &errCode, ErrorMessage: &errMsg}
	}

	err, errResp := processUploadRequest(*uploadRequest, logger)
	if err != nil {
		switch err {
		case InvalidParametersError:
			return http.StatusBadRequest, errResp
		case FatalError:
			return http.StatusInternalServerError, errResp
		}
	}

	return http.StatusOK, nil
}

func processProtobufPost(req *http.Request, logger *log.Logger) (int, *model.ErrorResponse) {

	uploadRequest, err := decodeProtobuf(req.Body)
	if err != nil {
		logger.WithFields(log.Fields{"method": "processProtobufPost", "error": err.Error()}).Warningln("Error decoding body")

		errCode := MalformedRequest
		errMsg := "Request body contains malformed buffered data"
		return http.StatusBadRequest, &model.ErrorResponse{Code: &errCode, ErrorMessage: &errMsg}
	}

	err, errResp := processUploadRequest(*uploadRequest, logger)
	if err != nil {
		switch err {
		case InvalidParametersError:
			return http.StatusBadRequest, errResp
		case FatalError:
			return http.StatusInternalServerError, errResp
		}
	}

	return http.StatusOK, nil
}

//Can be used for logging in case there's a system in place to ban IP addresses that try to DDOS the service.
func processUnsupportedMediaTypePost(req *http.Request, logger *log.Logger) (int, *model.ErrorResponse) {
	logger.WithFields(log.Fields{"method": "processUnsupportedMediaTypePost", "error": UnsupportedMedia}).Infoln("Content-Type: " + req.Header.Get("Content-Type"))

	errCode := UnsupportedMedia
	errMsg := "Content-Type specifies a media type that is not supported by this resource"
	return http.StatusUnsupportedMediaType, &model.ErrorResponse{Code: &errCode, ErrorMessage: &errMsg}
}

//-----------------------------------------------------

func decodeJson(r io.ReadCloser) (uploadRequest *model.ClientEventUploadRequest, err error) {
	uploadRequest = new(model.ClientEventUploadRequest)
	err = json.NewDecoder(r).Decode(uploadRequest)
	return
}

func decodeProtobuf(r io.ReadCloser) (*model.ClientEventUploadRequest, error) {

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	uploadRequest := new(model.ClientEventUploadRequest)
	err = proto.Unmarshal(data, uploadRequest)
	if err != nil {
		return nil, err
	}

	return uploadRequest, nil
}

//------------------------------------------------------

/*
Each upload request can have multiple events, to achieve transcational behavior the datastore will be switched to boltdb in a later version. For now the client events bundle is validated to achieve atomicity, either all of them are stored or an error response is sent back.

Partial storage is performed in case of errors from in memory database StoreEvent() method
*/
func processUploadRequest(uploadRequest model.ClientEventUploadRequest, logger *log.Logger) (error, *model.ErrorResponse) {
	flag, index := hasValidEventIds(uploadRequest.GetEvents())
	if !flag {
		logger.WithFields(log.Fields{"method": "processUploadRequest", "error": InvalidParametersError.Error(), "requestId": uploadRequest.GetRequestId()}).Warningln("Error validating eventIds in the upload request")

		errCode := InvalidRequestParameters
		errMsg := "Event bundle has an event with invalid eventId"
		errDes := "Event index (count starts from 0): " + strconv.Itoa(index)
		return InvalidParametersError, &model.ErrorResponse{Code: &errCode, ErrorMessage: &errMsg, Description: &errDes}
	}

	for i, event := range uploadRequest.GetEvents() {
		err := StoreEvent(*event)
		if err != nil {
			logger.WithFields(log.Fields{"method": "processUploadRequest", "error": FatalError.Error(), "requestId": uploadRequest.GetRequestId()}).Errorln("Error storing event with index: " + strconv.Itoa(i))

			errCode := Fatal
			errMsg := "Error storing events, aborting"
			return FatalError, &model.ErrorResponse{Code: &errCode, ErrorMessage: &errMsg}
		}
	}

	logger.WithFields(log.Fields{"method": "processUploadRequest", "requestId": uploadRequest.RequestId}).Infoln("Request successfully processed")
	return nil, nil
}

func hasValidEventIds(events []*model.ClientEventData) (bool, int) {
	for i, event := range events {
		if event.GetEventId() == "" {
			return false, i
		}
	}
	return true, -1
}
