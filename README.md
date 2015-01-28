###What is Meowtrics?###
Meowtrics is a metrics collection server written in Go. For this first version it provides a very basic API for storing and retrieving events. Current version is built using [negroni](https://github.com/codegangsta/negroni), [gorilla mux](http://www.gorillatoolkit.org/pkg/mux), [protocol buffers](https://developers.google.com/protocol-buffers/), and stores data in memory.

###API Docs###

####Heartbeat####

**Request**

For checking the service status, heartbeat calls can be made.

Path - `/heartbeat`

**Response**

The response would be in the following JSON format.

```javascript
{
    "status": "OK",
    "timestamp": "2015-01-28 07:18:14.450707079 +0000 UTC"
}
```

####ClientEventUploadRequest####

Method - `POST`

Path - `/v1/events` 

Content-Type - `application/json` or `application/x-protobuf`

**Request body format for JSON**

```javascript
{
  "request_id": "testRequestId",
  "device_type": "testDeviceAndroid",
  "events": [
    {
      "event_id": "123",
      "event_type": 1,
      "timestamp": 1422409858,
      "data": "testTestTestTestTest"
    }
  ]
}
```
For protobuf, more details can be found in the [metrics.proto](https://github.com/thezelus/meowtrics/blob/master/model/metrics.proto) file in the model package.


**Response**

For successful POST, the response body will be empty. For error cases, please check the error details section below.


####ClientEventData####

**Request**

Method - `GET`

Path - `/v1/events/id` (id should be replaced by the eventId)

Accept header - `application/json` or `application/x-protobuf` or `*/*` or none

Note: This version of the meowtrics server only accepts numerical ids for GET requests, non-numerical ids would get a NOT_FOUND error response.

**Response**

Successful response to `GET` request is returned in following JSON format -

```javascript
{
  "event_id": "123",
  "event_type": 1,
  "timestamp": 1422409858,
  "data": "testTestTestTestTest"
}
```

####Error details####

**Response**

In case of errors, the response body will include the error response in following JSON format.

```javascript
{
    "code": "SAMPLE",
    "error_message": "sample",
    "description": "sample-meow"
}
```

**Error codes**

- NOT_FOUND
- INVALID_REQUEST_PARAMETERS
- MALFORMED_REQUEST
- HEADER_NOT_RECOGNIZED
- REQUESTED_RECORD_NOT_FOUND
- FATAL_OPERATION
- UNSUPPORTED_MEDIA_TYPE


###Notes###
- I like to create a config directory with the name same as the project under '$GOPATH/bin/config/', and this is set as the DefaultDeploymentPath for the config file ('$GOPATH/bin/config/meowtrics/' for this project).
- Viper is configured to check first in the default deployment directory and then in the injected config path.
- Each ClientEventUploadRequest POST can have multiple events, to achieve transcational behavior the datastore will be switched to boltdb in a later version. For now the client events bundle is validated to achieve atomicity, either all of them are stored or an error response is sent back. 
- Partial storage is performed in case of errors from in memory database StoreEvent() method
- Header -->  "Content-Type" ---> "application/json" OR "application/x-protobuf"
- POST calls have no restriction on eventId type (can be string or integers), GET calls only accept numeric values as id

###Done List###
- [X] Implement in memory database access to use for testing. 
- [X] Implement POST and GET handlers. 
- [X] The POST handler(s) should allow both protobuf and JSON encoding as input depending on the Content-Type header. 
- [X] The GET handler(s) should allow both protobuf and JSON as output depending on the Accept header and if no Accept 	header is present, or set to “*/*”, output should be JSON by default. 
- [X] The GET request handlers should only accept numerical id’s and return a 404 otherwise. 
- [X] Modify metrics.proto so that clients can submit arbitrary key/value pairs as part of the ClientEventData object. 
- [X] Add test cases demonstrating successful storage and retrieval plus error cases. 

###TODO List###
- [X] Update documentation to include details about the API endpoints
- [X] Write deployment scripts
- [X] Include more tests
- [ ] Dockerize meowtrics
- [ ] Modify datasource to use [bolt](https://github.com/boltdb/bolt) instead of a map

###Why call it *Meowtrics*?###
Because it is a metrics collection server, but.....

<img src="https://38.media.tumblr.com/21ff9c82d8c0a686a03e6aa12683ddc2/tumblr_mvj9n2YhH11r4sj1co2_500.gif" width="350px" height="225px">

<img src="http://31.media.tumblr.com/902475db2312e77265b1e527261ee0f1/tumblr_mig9ppVJfQ1qjjnt0o1_500.gif" width="350px" height="225px"> 

<img src="http://24.media.tumblr.com/tumblr_m9k621fdMK1ry5v76o7_500.gif" width="350px" height="225px" >   
And [Super Troopers](http://www.imdb.com/title/tt0247745/) is an awesome movie.

###Why not use an existing solution?###
I will try this out for meow, this is more of a learning exercise to use [protobuf](https://github.com/golang/protobuf) in Go. If you want to extend this solution or if you have feature suggestions, you should fork this repository right meow.


(P.S. - Inspiration from a  [test project](https://github.com/kikinteractive/server-metrics-test))